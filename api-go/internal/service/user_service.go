package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"med-vito/api-go/internal/config"
	mailpkg "med-vito/api-go/internal/pkg/mail"
	"med-vito/api-go/internal/pkg/s3client"
	"med-vito/api-go/internal/pkg/sms"
	"med-vito/api-go/internal/rbac"
	"med-vito/api-go/internal/repository"
)

const verifyEmailPrefix = "verify-email:"

type fiberMap = map[string]any

var profileTypeLabels = map[string]string{
	"IP":         "Индивидуальный предприниматель",
	"OOO":        "Юридическое лицо",
	"INDIVIDUAL": "Физическое лицо",
}

// UserService — логика user.controller / user.service (Nest).
type UserService struct {
	cfg   config.Config
	users *repository.UserPG
	rdb   *redis.Client
	s3    *s3client.Client // nil — без загрузки фото
	sms   *sms.Sender
}

func NewUserService(cfg config.Config, users *repository.UserPG, rdb *redis.Client, s3 *s3client.Client) *UserService {
	return &UserService{
		cfg:   cfg,
		users: users,
		rdb:   rdb,
		s3:    s3,
		sms:   sms.NewSender(cfg.MTSBearer, cfg.NotisendAPIKey, cfg.NotisendProject),
	}
}

func (s *UserService) sendMail(to, subject, html string) {
	if strings.TrimSpace(s.cfg.SMTPHost) == "" || strings.TrimSpace(to) == "" {
		return
	}
	_ = mailpkg.SendHTMLSmart(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPFrom, to, subject, html, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure)
}

func (s *UserService) FindAllAdmin(ctx context.Context) ([]fiberMap, error) {
	rows, err := s.users.ListUsersAdmin(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]fiberMap, 0, len(rows))
	for _, r := range rows {
		roleName := "USER"
		if r.RoleName != nil {
			rn := strings.TrimSpace(*r.RoleName)
			if rn != "" {
				if strings.EqualFold(rn, "default") {
					roleName = "USER"
				} else {
					roleName = rn
				}
			}
		}
		out = append(out, fiberMap{
			"id":              r.ID,
			"createdAt":       r.CreatedAt,
			"isBanned":        r.IsBanned,
			"fullName":        r.FullName,
			"email":           r.Email,
			"phoneNumber":     r.PhoneNumber,
			"rating":          r.Rating,
			"profileType":     r.ProfileType,
			"photo":           r.Photo,
			"balance":         r.Balance,
			"bonusBalance":    r.BonusBalance,
			"role":            roleName,
			"isEmailVerified": r.IsEmailVerified,
			"isPhoneVerified": r.IsPhoneVerified,
			"products":        r.ProductsCount,
			"productStats": fiberMap{
				"total":      r.ProductsCount,
				"drafts":     r.DraftsCount,
				"moderation": r.ModerationCount,
				"active":     r.ActiveCount,
				"hidden":     r.HiddenCount,
				"denied":     r.DeniedCount,
				"paid":       r.PaidCount,
				"free":       r.FreeCount,
			},
			"adsLimit": fiberMap{
				"total":     int32(repository.FreeAdsMonthlyLimit),
				"used":      r.UsedFreeAds,
				"remaining": maxInt32(0, int32(repository.FreeAdsMonthlyLimit)-r.UsedFreeAds),
				"costPerAd": repository.AdListingCost,
			},
			"reviewsCount": r.ReviewsCount,
		})
	}
	return out, nil
}

func (s *UserService) ToggleBanned(ctx context.Context, id int32) (fiberMap, error) {
	was, err := s.users.GetBannedState(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Пользователь не найден"}
		}
		return nil, err
	}
	if err := s.users.SetBanned(ctx, id, !was); err != nil {
		return nil, err
	}
	if was {
		return fiberMap{"message": "Пользователь успешно разбанен"}, nil
	}
	return fiberMap{"message": "Пользователь успешно забанен"}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) (fiberMap, error) {
	email, _, _ := s.users.GetUserEmailPhone(ctx, id)
	if err := s.users.DeleteUserByID(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Пользователь не найден"}
		}
		return nil, err
	}
	if email != "" {
		go s.sendMail(email, "Удаление аккаунта - Торгуй Сам", "<p>Ваш аккаунт был удален администратором.</p>")
	}
	return fiberMap{"message": "Пользователь успешно удалён"}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, id int32) (fiberMap, error) {
	row, err := s.users.GetUserInfoRow(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{400, "Пользователь не найден"}
		}
		return nil, err
	}
	now := time.Now().UTC()
	last := row.LastAdLimitReset.UTC()
	used := row.UsedFreeAds
	if now.Month() != last.Month() || now.Year() != last.Year() {
		used = 0
		if err := s.users.ResetMonthlyFreeAds(ctx, id, now); err != nil {
			return nil, err
		}
	}
	avg, count, err := s.users.ApprovedReviewsStats(ctx, id)
	if err != nil {
		return nil, err
	}
	rating := 0.0
	if count > 0 {
		rating = avg
	}
	label := profileTypeLabels[row.ProfileType]
	if label == "" {
		label = row.ProfileType
	}
	totalLimit := int32(repository.FreeAdsMonthlyLimit)
	remaining := totalLimit - used
	if remaining < 0 {
		remaining = 0
	}
	return fiberMap{
		"id":           row.ID,
		"fullName":     row.FullName,
		"profileType":  label,
		"phoneNumber":  row.PhoneNumber,
		"balance":      row.Balance,
		"bonusBalance": row.BonusBalance,
		"photo":        row.Photo,
		"rating":       rating,
		"reviewsCount": count,
		"adsLimit": fiberMap{
			"total":     totalLimit,
			"used":      used,
			"remaining": remaining,
			"costPerAd": 50,
		},
	}, nil
}

func (s *UserService) GetRemainingFreeAds(ctx context.Context, userID int32) (fiberMap, error) {
	st, err := s.users.GetFreeAdsState(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Пользователь не найден"}
		}
		return nil, err
	}
	now := time.Now().UTC()
	reset := st.LastAdLimitReset.UTC()
	used := st.UsedFreeAds
	if now.Month() != reset.Month() || now.Year() != reset.Year() {
		if err := s.users.ResetMonthlyFreeAds(ctx, userID, now); err != nil {
			return nil, err
		}
		used = 0
	}
	totalLimit := int32(repository.FreeAdsMonthlyLimit)
	rem := totalLimit - used
	if rem < 0 {
		rem = 0
	}
	next := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	return fiberMap{
		"total":     totalLimit,
		"used":      used,
		"remaining": rem,
		"costPerAd": 50,
		"nextReset": next,
	}, nil
}

func (s *UserService) GetPhoneNumber(ctx context.Context, sellerID, viewerID int32) (fiberMap, error) {
	ph, err := s.users.GetPhoneForUser(ctx, sellerID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Продавец с данным id не найден"}
		}
		return nil, err
	}
	if strings.TrimSpace(ph) == "" {
		return nil, &AppError{400, "У продавца отсутствует номер телефона"}
	}
	ok, err := s.users.PhoneViewExists(ctx, viewerID, sellerID)
	if err != nil {
		return nil, err
	}
	if !ok {
		if err := s.users.InsertPhoneView(ctx, viewerID, sellerID); err != nil {
			return nil, err
		}
	}
	return fiberMap{"phoneNumber": ph}, nil
}

// UpdateSettingsMultipart — поля из multipart + опциональный файл photo.
func (s *UserService) UpdateSettingsMultipart(ctx context.Context, userID int32, fullName, phone, email, isAnswersRaw, profileType string, photoName, photoCT string, photoBody []byte) (fiberMap, error) {
	patch := repository.UserSettingsPatch{}
	if fullName != "" {
		patch.FullName = &fullName
	}
	if phone != "" {
		_, curPhone, err := s.users.GetUserEmailPhone(ctx, userID)
		if err == nil {
			normNew, errNew := normalizePhone(phone)
			normCur, _ := normalizePhone(curPhone)
			if errNew == nil && normNew != normCur {
				// Изменение номера телефона требует отдельного подтверждения через SMS!
				// Напрямую через обычные настройки не перезаписываем.
			} else if errNew == nil {
				patch.PhoneNumber = &normNew
			} else {
				patch.PhoneNumber = &phone
			}
		} else {
			patch.PhoneNumber = &phone
		}
	}
	if strings.TrimSpace(email) != "" {
		cleanEmail := strings.ToLower(strings.TrimSpace(email))
		patch.Email = &cleanEmail
	}
	if isAnswersRaw != "" {
		v := isAnswersRaw == "true" || isAnswersRaw == "1"
		patch.IsAnswersCall = &v
	}
	if profileType != "" {
		pt := strings.TrimSpace(profileType)
		patch.ProfileType = &pt
	}
	if len(photoBody) > 0 {
		if s.s3 == nil {
			return nil, &AppError{500, "S3 не настроен — загрузка фото недоступна"}
		}
		old, err := s.users.GetUserPhoto(ctx, userID)
		if err != nil && err != repository.ErrNotFound {
			return nil, err
		}
		if old != nil && *old != "" {
			_ = s.s3.DeleteByURL(ctx, *old)
		}
		url, err := s.s3.Upload(ctx, "users", photoName, photoCT, photoBody)
		if err != nil {
			return nil, &AppError{400, "Ошибка загрузки файла в S3: " + err.Error()}
		}
		patch.Photo = &url
	}
	if err := s.users.UpdateUserSettings(ctx, userID, patch); err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Пользователь не найден"}
		}
		return nil, err
	}
	// Ответ как в Nest: развёрнутый dto + photo при загрузке
	out := fiberMap{}
	if patch.FullName != nil {
		out["fullName"] = *patch.FullName
	}
	if patch.PhoneNumber != nil {
		out["phoneNumber"] = *patch.PhoneNumber
	}
	if patch.Email != nil {
		out["email"] = *patch.Email
	}
	if patch.IsAnswersCall != nil {
		out["isAnswersCall"] = *patch.IsAnswersCall
	}
	if patch.ProfileType != nil {
		out["profileType"] = *patch.ProfileType
	}
	if patch.Photo != nil {
		out["photo"] = *patch.Photo
	}
	return out, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, userID int32) (fiberMap, error) {
	u, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Пользователь не найден"}
		}
		return nil, err
	}
	if s.cfg.SMTPHost == "" {
		return nil, &AppError{500, "SMTP не настроен (SMTP_HOST)"}
	}
	code := generateVerifyCode6()
	payload, _ := json.Marshal(fiberMap{"id": strconv.Itoa(int(u.ID)), "code": code})
	if err := s.rdb.Set(ctx, verifyEmailPrefix+code, string(payload), time.Hour).Err(); err != nil {
		return nil, err
	}
	subj := "Код подтверждения почты - Торгуй Сам"
	htmlBody, err := mailpkg.VerifyEmailHTML(code)
	if err != nil {
		return nil, &AppError{500, "Не удалось сформировать письмо"}
	}
	if err := mailpkg.SendHTMLSmart(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword,
		s.cfg.SMTPFrom, u.Email, subj, htmlBody, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure); err != nil {
		return nil, &AppError{400, "Ошибка отправки письма: " + err.Error()}
	}
	return fiberMap{"message": "Письмо с кодом подтверждения отправлено на почту"}, nil
}

type verifyEmailCache struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

func (s *UserService) VerifyEmailCode(ctx context.Context, code string) (fiberMap, error) {
	if code == "" {
		return nil, &AppError{400, "Неверный код подтверждения"}
	}
	raw, err := s.rdb.Get(ctx, verifyEmailPrefix+code).Result()
	if err == redis.Nil || raw == "" {
		return nil, &AppError{400, "Неверный код подтверждения"}
	}
	if err != nil {
		return nil, err
	}
	var cached verifyEmailCache
	if err := json.Unmarshal([]byte(raw), &cached); err != nil || cached.Code != code {
		return nil, &AppError{400, "Неверный код подтверждения"}
	}
	uid64, err := strconv.ParseInt(cached.ID, 10, 32)
	if err != nil {
		return nil, &AppError{400, "Неверный код подтверждения"}
	}
	if _, err := s.users.FindUserByID(ctx, int32(uid64)); err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Такого пользователя не существует"}
		}
		return nil, err
	}
	if err := s.users.SetEmailVerified(ctx, int32(uid64), true); err != nil {
		return nil, err
	}
	_ = s.rdb.Del(ctx, verifyEmailPrefix+code).Err()
	return fiberMap{"message": "Почта успешно подтверждена"}, nil
}

func (s *UserService) SetBalance(ctx context.Context, userID int32, balanceStr string) (fiberMap, error) {
	if _, err := s.users.FindUserByID(ctx, userID); err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Пользователь не найден"}
		}
		return nil, err
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(balanceStr), 64)
	if err != nil {
		return nil, &AppError{400, "Неудалось преобразовать баланс в float"}
	}
	if v < 0 {
		return nil, &AppError{400, "Нельзя установить отрицательный баланс"}
	}
	if err := s.users.SetBonusBalance(ctx, userID, v); err != nil {
		return nil, err
	}
	return fiberMap{"message": "Баланс успешно обновлён"}, nil
}

// AdminUpdateUser — частичное обновление (как updateUser в Nest).
func (s *UserService) AdminUpdateUser(ctx context.Context, userID int32, fullName, email, phone, profileType *string, bonus *float64) (fiberMap, error) {
	curEmail, curPhone, err := s.users.GetUserEmailPhone(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Пользователь не найден"}
		}
		return nil, err
	}
	if email != nil && *email != curEmail {
		other, err := s.users.FindUserIDByEmail(ctx, *email)
		if err != nil {
			return nil, err
		}
		if other != nil && *other != userID {
			return nil, &AppError{400, "Пользователь с таким email уже существует"}
		}
	}
	if phone != nil && *phone != curPhone {
		other, err := s.users.FindUserIDByPhone(ctx, *phone)
		if err != nil {
			return nil, err
		}
		if other != nil && *other != userID {
			return nil, &AppError{400, "Пользователь с таким номером телефона уже существует"}
		}
	}
	p := repository.AdminUserPatch{
		FullName: fullName, Email: email, PhoneNumber: phone, ProfileType: profileType, BonusBalance: bonus,
	}
	if err := s.users.UpdateUserAdmin(ctx, userID, p); err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Пользователь не найден"}
		}
		return nil, err
	}
	if phone != nil && !isDealPhoneSynthetic(*phone) {
		_ = s.users.SetPhoneVerified(ctx, userID, true)
	}
	targetEmail := curEmail
	if email != nil && strings.TrimSpace(*email) != "" {
		targetEmail = strings.TrimSpace(*email)
	}
	if targetEmail != "" {
		if bonus != nil {
			go s.sendMail(targetEmail, "Обновление профиля - Торгуй Сам", fmt.Sprintf("<p>Администратор обновил ваш профиль. Бонусный баланс: <b>%.0f ₽</b>.</p>", *bonus))
		} else {
			go s.sendMail(targetEmail, "Обновление профиля - Торгуй Сам", "<p>Администратор обновил данные вашего профиля.</p>")
		}
	}
	return fiberMap{"message": "Данные пользователя успешно обновлены"}, nil
}

func (s *UserService) AdminSetUserRole(ctx context.Context, userID int32, role string) (fiberMap, error) {
	role = strings.ToUpper(strings.TrimSpace(role))
	switch role {
	case rbac.RoleUser, rbac.RoleUserVerified, rbac.RoleSeniorModerator, rbac.RoleAdmin, rbac.RoleSuperAdmin:
	default:
		return nil, &AppError{400, "Недопустимая роль"}
	}
	if err := s.users.SetUserRoleByName(ctx, userID, role); err != nil {
		if err == repository.ErrNotFound {
			return nil, &AppError{404, "Пользователь или роль не найдены"}
		}
		return nil, err
	}
	email, _, _ := s.users.GetUserEmailPhone(ctx, userID)
	if email != "" {
		go s.sendMail(email, "Изменение роли - Торгуй Сам", fmt.Sprintf("<p>Администратор изменил вашу роль в системе на «<b>%s</b>».</p>", role))
	}
	return fiberMap{"message": "Роль пользователя обновлена", "role": role}, nil
}

func generateVerifyCode6() string {
	n := 100000 + rand.Intn(900000)
	return strconv.Itoa(n)
}

func maxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

const changePhoneKeyPrefix = "change-phone:"

type changePhoneCache struct {
	UserID int32  `json:"userId"`
	Phone  string `json:"phone"`
	Code   string `json:"code"`
}

func normalizePhone(raw string) (string, error) {
	digits := ""
	for _, r := range raw {
		if r >= '0' && r <= '9' {
			digits += string(r)
		}
	}
	if len(digits) == 10 {
		digits = "7" + digits
	} else if len(digits) == 11 {
		if digits[0] == '8' {
			digits = "7" + digits[1:]
		}
	}
	if len(digits) != 11 || digits[0] != '7' {
		return "", errors.New("некорректный номер телефона (требуется 11 цифр, формат РФ +7)")
	}
	return "+" + digits, nil
}

func (s *UserService) RequestPhoneChange(ctx context.Context, userID int32, rawPhone string) (fiberMap, error) {
	newPhone, err := normalizePhone(rawPhone)
	if err != nil {
		return nil, &AppError{400, err.Error()}
	}

	_, curPhone, err := s.users.GetUserEmailPhone(ctx, userID)
	if err != nil {
		return nil, err
	}
	normCur, _ := normalizePhone(curPhone)
	if newPhone == normCur {
		return nil, &AppError{400, "Новый номер совпадает с текущим"}
	}

	existingUser, err := s.users.FindUserByPhone(ctx, newPhone)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return nil, &AppError{400, "Этот номер телефона уже привязан к другому аккаунту"}
	}

	code := generateVerifyCode6()
	cache := changePhoneCache{
		UserID: userID,
		Phone:  newPhone,
		Code:   code,
	}
	b, _ := json.Marshal(cache)
	if err := s.rdb.Set(ctx, changePhoneKeyPrefix+strconv.Itoa(int(userID)), b, 15*time.Minute).Err(); err != nil {
		return nil, err
	}

	if s.sms != nil {
		s.sms.SendAsync(newPhone, "TorguiSam.ru: Код подтверждения смены номера: "+code)
	}

	return fiberMap{
		"message": "Код подтверждения отправлен по SMS на номер " + newPhone,
		"phone":   newPhone,
	}, nil
}

func (s *UserService) ConfirmPhoneChange(ctx context.Context, userID int32, code string) (fiberMap, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, &AppError{400, "Введите код подтверждения"}
	}

	key := changePhoneKeyPrefix + strconv.Itoa(int(userID))
	raw, err := s.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) || raw == "" {
		return nil, &AppError{400, "Срок действия кода истёк или смена номера не запрашивалась"}
	}
	if err != nil {
		return nil, err
	}

	var cached changePhoneCache
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return nil, &AppError{400, "Ошибка чтения данных кода"}
	}

	if cached.Code != code {
		return nil, &AppError{400, "Неверный код подтверждения"}
	}

	if err := s.users.UpdatePhoneNumber(ctx, userID, cached.Phone); err != nil {
		return nil, err
	}

	_ = s.rdb.Del(ctx, key).Err()

	email, _, _ := s.users.GetUserEmailPhone(ctx, userID)
	if email != "" {
		go s.sendMail(email, "Смена номера телефона - Торгуй Сам", fmt.Sprintf("<p>Номер телефона вашего аккаунта успешно изменён на <b>%s</b>.</p>", cached.Phone))
	}

	return fiberMap{
		"message":     "Номер телефона успешно подтверждён и обновлён",
		"phoneNumber": cached.Phone,
	}, nil
}
