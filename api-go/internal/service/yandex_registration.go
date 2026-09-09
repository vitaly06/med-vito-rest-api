package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/repository"
)

const (
	// Незавершённая регистрация через Яндекс живёт только в Redis: пока телефон
	// не подтверждён по SMS, пользователя в базе нет.
	yandexPendingPrefix   = "yandex:pending:"
	yandexPendingTTL      = 30 * time.Minute
	yandexPendingMaxTries = 5
)

// yandexPendingRegistration — профиль Яндекса + состояние SMS-подтверждения до создания аккаунта.
type yandexPendingRegistration struct {
	ExternalID string `json:"externalId"`
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	Avatar     string `json:"avatar"`
	Phone      string `json:"phone"`
	Code       string `json:"code"`
	Attempts   int    `json:"attempts"`
}

// normalizeRuPhone приводит номер к виду +7XXXXXXXXXX — так же, как его шлёт форма регистрации.
func normalizeRuPhone(phone string) (string, error) {
	var b strings.Builder
	for _, r := range phone {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	ds := b.String()
	switch len(ds) {
	case 10:
		if ds[0] == '9' {
			return "+7" + ds, nil
		}
	case 11:
		if ds[0] == '8' || ds[0] == '7' {
			return "+7" + ds[1:], nil
		}
	}
	return "", &AppError{400, "Некорректный номер телефона"}
}

func (s *AuthService) pendingKey(ticket string) string {
	return yandexPendingPrefix + strings.TrimSpace(ticket)
}

// createYandexPendingRegistration кладёт профиль Яндекса в Redis и возвращает тикет регистрации.
func (s *AuthService) createYandexPendingRegistration(ctx context.Context, p oauthProfile) (string, error) {
	pending := yandexPendingRegistration{
		ExternalID: strings.TrimSpace(p.ExternalID),
		Email:      strings.ToLower(strings.TrimSpace(p.Email)),
		FullName:   strings.TrimSpace(p.FullName),
		Avatar:     strings.TrimSpace(p.Avatar),
	}
	ticket := generateSessionID()
	b, _ := json.Marshal(pending)
	if err := s.rdb.Set(ctx, s.pendingKey(ticket), b, yandexPendingTTL).Err(); err != nil {
		return "", err
	}
	return ticket, nil
}

func (s *AuthService) loadYandexPendingRegistration(ctx context.Context, ticket string) (*yandexPendingRegistration, error) {
	ticket = strings.TrimSpace(ticket)
	if ticket == "" {
		return nil, &AppError{400, "Не передан тикет регистрации"}
	}
	raw, err := s.rdb.Get(ctx, s.pendingKey(ticket)).Bytes()
	if errors.Is(err, redis.Nil) || len(raw) == 0 {
		return nil, &AppError{401, "Регистрация через Яндекс устарела, войдите заново"}
	}
	if err != nil {
		return nil, err
	}
	var pending yandexPendingRegistration
	if err := json.Unmarshal(raw, &pending); err != nil {
		return nil, err
	}
	return &pending, nil
}

func (s *AuthService) saveYandexPendingRegistration(ctx context.Context, ticket string, pending *yandexPendingRegistration) error {
	b, _ := json.Marshal(pending)
	return s.rdb.Set(ctx, s.pendingKey(ticket), b, yandexPendingTTL).Err()
}

// YandexRegistrationStatus — состояние незавершённой регистрации (для перезагрузки страницы).
func (s *AuthService) YandexRegistrationStatus(ctx context.Context, ticket string) (map[string]any, error) {
	pending, err := s.loadYandexPendingRegistration(ctx, ticket)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"pending":     true,
		"email":       pending.Email,
		"fullName":    pending.FullName,
		"phoneNumber": pending.Phone,
		"codeSent":    pending.Code != "",
	}, nil
}

// YandexRegistrationSendCode запоминает телефон в незавершённой регистрации и отправляет SMS-код.
func (s *AuthService) YandexRegistrationSendCode(ctx context.Context, ticket, phone string) (string, error) {
	pending, err := s.loadYandexPendingRegistration(ctx, ticket)
	if err != nil {
		return "", err
	}
	normalized, err := normalizeRuPhone(phone)
	if err != nil {
		return "", err
	}

	code := s.generateVerifyCode()
	pending.Phone = normalized
	pending.Code = code
	pending.Attempts = 0
	if err := s.saveYandexPendingRegistration(ctx, ticket, pending); err != nil {
		return "", err
	}
	if err := s.sendVKPhoneCode(ctx, normalized, code); err != nil {
		return "", err
	}
	return normalized, nil
}

// YandexRegistrationVerifyCode проверяет SMS-код и только после этого создаёт аккаунт и сессию.
func (s *AuthService) YandexRegistrationVerifyCode(ctx context.Context, ticket, code string) (*signInResponse, string, error) {
	pending, err := s.loadYandexPendingRegistration(ctx, ticket)
	if err != nil {
		return nil, "", err
	}
	code = strings.TrimSpace(code)
	if pending.Code == "" || pending.Phone == "" {
		return nil, "", &AppError{400, "Сначала запросите код по SMS"}
	}
	if code == "" {
		return nil, "", &AppError{400, "Введите код из SMS"}
	}
	if pending.Code != code {
		pending.Attempts++
		if pending.Attempts >= yandexPendingMaxTries {
			_ = s.rdb.Del(ctx, s.pendingKey(ticket)).Err()
			return nil, "", &AppError{429, "Слишком много попыток, войдите через Яндекс заново"}
		}
		_ = s.saveYandexPendingRegistration(ctx, ticket, pending)
		return nil, "", &AppError{400, "Неверный код подтверждения"}
	}

	user, err := s.completeYandexRegistration(ctx, pending)
	if err != nil {
		return nil, "", err
	}
	_ = s.rdb.Del(ctx, s.pendingKey(ticket)).Err()

	sid := generateSessionID()
	sp := sessionPayload{
		UserID:       user.ID,
		Email:        user.Email,
		ProfileType:  user.ProfileType,
		AuthProvider: "yandex",
	}
	b, _ := json.Marshal(sp)
	if err := s.rdb.Set(ctx, sessionKeyPrefix+sid, b, sessionTTL).Err(); err != nil {
		return nil, "", err
	}

	var photo *string
	if user.Photo != nil && *user.Photo != "" {
		p := s.cfg.BaseURL + *user.Photo
		photo = &p
	} else if pending.Avatar != "" {
		a := pending.Avatar
		photo = &a
	}

	out := &signInResponse{Message: "Регистрация через Яндекс завершена!"}
	out.User.ID = user.ID
	out.User.Email = user.Email
	out.User.FullName = user.FullName
	out.User.PhoneNumber = user.PhoneNumber
	out.User.ProfileType = user.ProfileType
	out.User.Photo = photo
	return out, sid, nil
}

// completeYandexRegistration создаёт пользователя с уже подтверждённым телефоном
// либо привязывает Яндекс к аккаунту, которому этот номер уже принадлежит.
func (s *AuthService) completeYandexRegistration(ctx context.Context, pending *yandexPendingRegistration) (*domain.UserEntity, error) {
	profile := oauthProfile{
		Provider:   "yandex",
		ExternalID: pending.ExternalID,
		Email:      pending.Email,
		FullName:   pending.FullName,
		Avatar:     pending.Avatar,
		Phone:      pending.Phone,
	}

	// За время подтверждения аккаунт мог появиться (второе окно, привязка по email).
	if existing, err := s.findExistingYandexUser(ctx, profile); err != nil {
		return nil, err
	} else if existing != nil {
		if err := s.users.SetPhone(ctx, existing.ID, pending.Phone); err == nil {
			_ = s.users.SetPhoneVerified(ctx, existing.ID, true)
		}
		return s.users.FindUserByID(ctx, existing.ID)
	}

	// Номер уже принадлежит аккаунту — привязываем к нему вход через Яндекс.
	if otherID, err := s.users.FindUserIDByPhone(ctx, pending.Phone); err == nil && otherID != nil {
		_ = s.users.UpsertOAuthIdentity(ctx, "yandex", pending.ExternalID, *otherID)
		_ = s.users.SetPhoneVerified(ctx, *otherID, true)
		return s.users.FindUserByID(ctx, *otherID)
	}

	email := pending.Email
	if email != "" {
		if _, err := s.users.FindUserByEmail(ctx, email); err == nil {
			// Почта занята другим аккаунтом — сохраняем плейсхолдер, вход остаётся по телефону.
			email = ""
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}
	if email == "" {
		email = "yandex_" + pending.ExternalID + "@oauth.local"
	}
	fullName := pending.FullName
	if fullName == "" {
		fullName = "YANDEX User"
	}

	roleID, err := s.defaultUserRoleID(ctx)
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(generateSessionID()), bcryptCost)
	if err != nil {
		return nil, err
	}
	uid, err := s.users.GenerateUniqueUserID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.users.InsertUser(ctx, uid, fullName, email, pending.Phone, string(hash), roleID); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return nil, &AppError{400, "Аккаунт с такой почтой или номером уже существует"}
		}
		return nil, err
	}
	_ = s.users.UpsertOAuthIdentity(ctx, "yandex", pending.ExternalID, uid)
	_ = s.users.SetPhoneVerified(ctx, uid, true)
	if pending.Email != "" && !strings.HasSuffix(email, "@oauth.local") {
		_ = s.users.SetEmailVerified(ctx, uid, true)
	}
	return s.users.FindUserByID(ctx, uid)
}

// findExistingYandexUser ищет уже существующий аккаунт по Яндекс-идентификатору или почте.
// Возвращает (nil, nil), если аккаунта ещё нет — тогда нужна регистрация с подтверждением телефона.
func (s *AuthService) findExistingYandexUser(ctx context.Context, p oauthProfile) (*domain.UserEntity, error) {
	externalID := strings.TrimSpace(p.ExternalID)
	if externalID == "" {
		return nil, &AppError{400, "Пустой внешний идентификатор OAuth"}
	}

	existingID, err := s.users.FindOAuthUserIDByProviderExternalID(ctx, "yandex", externalID)
	if err != nil {
		return nil, err
	}
	if existingID != nil {
		u, err := s.users.FindUserByID(ctx, *existingID)
		if err == nil {
			_ = s.users.UpsertOAuthIdentity(ctx, "yandex", externalID, u.ID)
			return s.applyTrustedOAuthProfile(ctx, u, p), nil
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	email := strings.ToLower(strings.TrimSpace(p.Email))
	if email != "" {
		u, err := s.users.FindUserByEmail(ctx, email)
		if err == nil {
			_ = s.users.UpsertOAuthIdentity(ctx, "yandex", externalID, u.ID)
			return s.applyTrustedOAuthProfile(ctx, u, p), nil
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	return nil, nil
}
