package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mathrand "math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/domain"
	mailpkg "med-vito/api-go/internal/pkg/mail"
	"med-vito/api-go/internal/repository"
)

const (
	sessionTTL          = 30 * 24 * time.Hour
	verifyPhoneTTL      = time.Hour
	forgotPassTTL       = time.Hour
	bcryptCost          = 10
	sessionKeyPrefix    = "session:"
	verifyKeyPrefix     = "verify-phone:"
	forgotKeyPrefix     = "forgot-password:"
	vkidPKCEPrefix      = "vkid-pkce:"
	vkidPKCETTL         = 10 * time.Minute
	vkIDPublicInfoURL   = "https://id.vk.ru/oauth2/public_info"
	vkVerifyEmailPrefix = "vk-verify-email:"
	vkVerifyPhonePrefix = "vk-verify-phone:"
)

type vkPublicInfo struct {
	UserID      string
	FirstName   string
	LastInitial string
	Avatar      string
	EmailMasked string
	PhoneMasked string
}

type oauthProfile struct {
	Provider   string
	ExternalID string
	Email      string
	Phone      string
	FullName   string
	Avatar     string
}

type AuthService struct {
	cfg    config.Config
	rdb    *redis.Client
	users  *repository.UserPG
	client *http.Client
}

func NewAuthService(cfg config.Config, rdb *redis.Client, users *repository.UserPG) *AuthService {
	return &AuthService{
		cfg:    cfg,
		rdb:    rdb,
		users:  users,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *AuthService) defaultUserRoleID(ctx context.Context) (int32, error) {
	if id, err := s.users.RoleIDByName(ctx, "USER"); err == nil {
		return id, nil
	}
	return s.users.RoleIDByName(ctx, "default")
}

type sessionPayload struct {
	UserID              int32  `json:"userId"`
	Email               string `json:"email"`
	ProfileType         string `json:"profileType"`
	AuthProvider        string `json:"authProvider,omitempty"`
	RequireVKOnboarding bool   `json:"requireVkOnboarding,omitempty"`
}

type signUpCache struct {
	Data struct {
		FullName    string `json:"fullName"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		Password    string `json:"password"`
	} `json:"data"`
	Code string `json:"code"`
}

type signInResponse struct {
	Message string `json:"message"`
	User    struct {
		ID          int32   `json:"id"`
		Email       string  `json:"email"`
		FullName    string  `json:"fullName"`
		PhoneNumber string  `json:"phoneNumber"`
		ProfileType string  `json:"profileType"`
		Photo       *string `json:"photo"`
	} `json:"user"`
}

func (s *AuthService) SignUp(ctx context.Context, where string, fullName, email, phone, password string) error {
	ok, err := s.users.FindUserByEmailOrPhone(ctx, email, phone)
	if err != nil {
		return err
	}
	if ok {
		return &AppError{400, "Данный пользователь уже существует"}
	}
	code := s.generateVerifyCode()
	fmt.Println(code)
	switch strings.ToLower(strings.TrimSpace(where)) {
	case "telegram":
		if s.cfg.NotisendAPIKey == "" {
			return &AppError{500, "NOTISEND_API_KEY не задан"}
		}
		u := fmt.Sprintf(
			"https://sms.notisend.ru/api/message/send?project=%s&message=%s&recipients=%s&apikey=%s",
			url.QueryEscape(s.cfg.NotisendProject),
			url.QueryEscape("Код подтверждения: "+code),
			url.QueryEscape(phone),
			url.QueryEscape(s.cfg.NotisendAPIKey),
		)
		var resp map[string]any
		if err := s.httpGETJSON(ctx, u, &resp); err != nil {
			return err
		}
		if st, _ := resp["status"].(string); st == "error" {
			return &AppError{500, "Не удалось отправить сообщение"}
		}
	case "sms":
		if s.cfg.MTSBearer == "" {
			return &AppError{500, "MTS_TOKEN не задан"}
		}
		body := map[string]any{
			"submits": []any{
				map[string]any{"msid": phone, "message": "Код подтверждения: " + code},
			},
			"naming": "Torguisamru",
		}
		if err := s.httpPostJSON(ctx,
			"https://api.mts.ru/client-omni-adapter_production/1.0.2/mcom/messageManagement/messages",
			body,
			s.cfg.MTSBearer,
			nil,
		); err != nil {
			return err
		}
	default:
		return &AppError{400, "Where должен быть telegram или sms"}
	}
	payload := signUpCache{Code: code}
	payload.Data.FullName = fullName
	payload.Data.Email = email
	payload.Data.PhoneNumber = phone
	payload.Data.Password = password
	b, _ := json.Marshal(payload)
	if err := s.rdb.Set(ctx, verifyKeyPrefix+code, b, verifyPhoneTTL).Err(); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) VerifyMobileCode(ctx context.Context, code string) (*signInResponse, string, error) {
	raw, err := s.rdb.Get(ctx, verifyKeyPrefix+code).Bytes()
	if err == redis.Nil || len(raw) == 0 {
		return nil, "", &AppError{400, "Код подтверждения не найден или истек"}
	}
	if err != nil {
		return nil, "", err
	}
	var cached signUpCache
	if err := json.Unmarshal(raw, &cached); err != nil {
		return nil, "", err
	}
	if cached.Code != code {
		return nil, "", &AppError{400, "Неверный код подтверждения"}
	}
	roleID, err := s.defaultUserRoleID(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, "", &AppError{404, "Роль USER/default не найдена"}
		}
		return nil, "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(cached.Data.Password), bcryptCost)
	if err != nil {
		return nil, "", err
	}
	uid, err := s.users.GenerateUniqueUserID(ctx)
	if err != nil {
		return nil, "", &AppError{500, err.Error()}
	}
	if err := s.users.InsertUser(ctx, uid, cached.Data.FullName, cached.Data.Email, cached.Data.PhoneNumber, string(hash), roleID); err != nil {
		return nil, "", err
	}
	_ = s.rdb.Del(ctx, verifyKeyPrefix+code)

	u, err := s.users.FindUserByID(ctx, uid)
	if err != nil {
		return nil, "", err
	}
	sid := generateSessionID()
	sp := sessionPayload{UserID: u.ID, Email: u.Email, ProfileType: u.ProfileType, AuthProvider: "password"}
	b, _ := json.Marshal(sp)
	if err := s.rdb.Set(ctx, sessionKeyPrefix+sid, b, sessionTTL).Err(); err != nil {
		return nil, "", err
	}

	var photo *string
	if u.Photo != nil && *u.Photo != "" {
		p := s.cfg.BaseURL + *u.Photo
		photo = &p
	}
	out := &signInResponse{Message: "Вы успешно зарегистрировались!"}
	out.User.ID = u.ID
	out.User.Email = u.Email
	out.User.FullName = u.FullName
	out.User.PhoneNumber = u.PhoneNumber
	out.User.ProfileType = u.ProfileType
	out.User.Photo = photo
	return out, sid, nil
}

func (s *AuthService) SignIn(ctx context.Context, login, password string) (*signInResponse, string, error) {
	u, err := s.users.FindUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, "", &AppError{401, "Пользователь не существует"}
		}
		return nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", &AppError{401, "Неверный пароль"}
	}
	sid := generateSessionID()
	sp := sessionPayload{UserID: u.ID, Email: u.Email, ProfileType: u.ProfileType, AuthProvider: "password"}
	b, _ := json.Marshal(sp)
	if err := s.rdb.Set(ctx, sessionKeyPrefix+sid, b, sessionTTL).Err(); err != nil {
		return nil, "", err
	}
	var photo *string
	if u.Photo != nil && *u.Photo != "" {
		p := s.cfg.BaseURL + *u.Photo
		photo = &p
	}
	out := &signInResponse{Message: "Вы успешно авторизовались!"}
	out.User.ID = u.ID
	out.User.Email = u.Email
	out.User.FullName = u.FullName
	out.User.PhoneNumber = u.PhoneNumber
	out.User.ProfileType = u.ProfileType
	out.User.Photo = photo
	return out, sid, nil
}

func (s *AuthService) TIDAuthURL(state string) (string, error) {
	if strings.TrimSpace(s.cfg.TIDClientID) == "" {
		return "", &AppError{500, "T-ID не настроен: TID_CLIENT_ID"}
	}
	if strings.TrimSpace(s.cfg.TIDRedirectURI) == "" {
		return "", &AppError{500, "T-ID не настроен: TID_REDIRECT_URI"}
	}
	q := url.Values{}
	q.Set("client_id", s.cfg.TIDClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", s.cfg.TIDRedirectURI)
	if strings.TrimSpace(s.cfg.TIDScope) != "" {
		q.Set("scope", s.cfg.TIDScope)
	}
	if strings.TrimSpace(state) == "" {
		state = generateSessionID()
	}
	q.Set("state", state)
	return strings.TrimRight(s.cfg.TIDAuthorizeURL, "?") + "?" + q.Encode(), nil
}

func (s *AuthService) SignInWithTID(ctx context.Context, code, state string) (*signInResponse, string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, "", &AppError{400, "Нужен code"}
	}
	if strings.TrimSpace(s.cfg.TIDClientID) == "" ||
		strings.TrimSpace(s.cfg.TIDClientSecret) == "" ||
		strings.TrimSpace(s.cfg.TIDRedirectURI) == "" ||
		strings.TrimSpace(s.cfg.TIDTokenURL) == "" {
		return nil, "", &AppError{500, "T-ID не настроен в .env"}
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", s.cfg.TIDClientID)
	form.Set("client_secret", s.cfg.TIDClientSecret)
	form.Set("redirect_uri", s.cfg.TIDRedirectURI)
	form.Set("code", code)
	if strings.TrimSpace(state) != "" {
		form.Set("state", strings.TrimSpace(state))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.TIDTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := s.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, "", &AppError{401, "T-ID token error: " + truncateForErr(string(body))}
	}

	accessToken, idToken, err := parseTIDTokenResponse(body)
	if err != nil {
		return nil, "", err
	}
	profile := oauthProfile{Provider: "tbank"}
	claims := parseOAuthJWTClaims(idToken)
	fillOAuthProfileFromClaims(&profile, claims)
	if accessToken != "" && strings.TrimSpace(s.cfg.TIDUserInfoURL) != "" {
		userinfo, err := s.fetchOAuthUserInfo(ctx, s.cfg.TIDUserInfoURL, accessToken)
		if err == nil {
			fillOAuthProfileFromClaims(&profile, userinfo)
		}
	}
	profile.Provider = "tbank"
	profile.ExternalID = strings.TrimSpace(profile.ExternalID)
	if profile.ExternalID == "" {
		return nil, "", &AppError{401, "T-ID не вернул идентификатор пользователя"}
	}

	user, err := s.findOrCreateTrustedOAuthUser(ctx, profile)
	if err != nil {
		return nil, "", err
	}

	sid := generateSessionID()
	sp := sessionPayload{
		UserID:       user.ID,
		Email:        user.Email,
		ProfileType:  user.ProfileType,
		AuthProvider: "tbank",
	}
	b, _ := json.Marshal(sp)
	if err := s.rdb.Set(ctx, sessionKeyPrefix+sid, b, sessionTTL).Err(); err != nil {
		return nil, "", err
	}

	var photo *string
	if user.Photo != nil && *user.Photo != "" {
		p := s.cfg.BaseURL + *user.Photo
		photo = &p
	} else if strings.TrimSpace(profile.Avatar) != "" {
		a := strings.TrimSpace(profile.Avatar)
		photo = &a
	}
	out := &signInResponse{Message: "Вы успешно авторизовались через T-ID!"}
	out.User.ID = user.ID
	out.User.Email = user.Email
	out.User.FullName = user.FullName
	out.User.PhoneNumber = user.PhoneNumber
	out.User.ProfileType = user.ProfileType
	out.User.Photo = photo
	return out, sid, nil
}

func (s *AuthService) VKAuthURL(state string) (string, error) {
	if strings.TrimSpace(s.cfg.VkOAuthClientID) == "" {
		return "", &AppError{500, "VK OAuth не настроен: VK_OAUTH_CLIENT_ID"}
	}
	if strings.TrimSpace(s.cfg.VkOAuthRedirectURI) == "" {
		return "", &AppError{500, "VK OAuth не настроен: VK_OAUTH_REDIRECT_URI"}
	}
	q := url.Values{}
	q.Set("client_id", s.cfg.VkOAuthClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", s.cfg.VkOAuthRedirectURI)
	if strings.TrimSpace(state) == "" {
		state = generateSessionID()
	}
	if s.cfg.VkIDEnabled {
		verifier := generatePKCEVerifier()
		if err := s.rdb.Set(context.Background(), vkidPKCEPrefix+state, verifier, vkidPKCETTL).Err(); err != nil {
			return "", err
		}
		q.Set("code_challenge", pkceS256(verifier))
		q.Set("code_challenge_method", "S256")
	}
	if strings.TrimSpace(s.cfg.VkOAuthScope) != "" {
		q.Set("scope", s.cfg.VkOAuthScope)
	}
	q.Set("state", state)
	return s.cfg.VkOAuthAuthorizeURL + "?" + q.Encode(), nil
}

func (s *AuthService) SignInWithVK(ctx context.Context, code, state, deviceID string) (*signInResponse, string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, "", &AppError{400, "Нужен code"}
	}
	if strings.TrimSpace(s.cfg.VkOAuthClientID) == "" ||
		strings.TrimSpace(s.cfg.VkOAuthClientSecret) == "" ||
		strings.TrimSpace(s.cfg.VkOAuthRedirectURI) == "" ||
		strings.TrimSpace(s.cfg.VkOAuthTokenURL) == "" ||
		strings.TrimSpace(s.cfg.VkOAuthUserInfoURL) == "" {
		return nil, "", &AppError{500, "VK OAuth не настроен в .env"}
	}

	tokenQ := url.Values{}
	tokenQ.Set("grant_type", "authorization_code")
	tokenQ.Set("client_id", s.cfg.VkOAuthClientID)
	tokenQ.Set("redirect_uri", s.cfg.VkOAuthRedirectURI)
	tokenQ.Set("code", code)
	if s.cfg.VkIDEnabled {
		state = strings.TrimSpace(state)
		if state == "" {
			return nil, "", &AppError{400, "Нужен state для VK ID"}
		}
		verifier, err := s.rdb.Get(ctx, vkidPKCEPrefix+state).Result()
		if err == redis.Nil || strings.TrimSpace(verifier) == "" {
			return nil, "", &AppError{401, "VK ID: истекла сессия авторизации"}
		}
		if err != nil {
			return nil, "", err
		}
		tokenQ.Set("code_verifier", verifier)
		if state != "" {
			tokenQ.Set("state", state)
		}
		deviceID = strings.TrimSpace(deviceID)
		if deviceID != "" {
			tokenQ.Set("device_id", deviceID)
		}
		_ = s.rdb.Del(ctx, vkidPKCEPrefix+state)
	}
	if strings.TrimSpace(s.cfg.VkOAuthClientSecret) != "" {
		tokenQ.Set("client_secret", s.cfg.VkOAuthClientSecret)
	}

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.VkOAuthTokenURL, strings.NewReader(tokenQ.Encode()))
	if err != nil {
		return nil, "", err
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenRes, err := s.client.Do(tokenReq)
	if err != nil {
		return nil, "", err
	}
	defer tokenRes.Body.Close()
	tokenBody, err := io.ReadAll(tokenRes.Body)
	if err != nil {
		return nil, "", err
	}
	if tokenRes.StatusCode < 200 || tokenRes.StatusCode >= 300 {
		return nil, "", &AppError{401, "VK OAuth token error: " + truncateForErr(string(tokenBody))}
	}
	accessToken, tokenEmail, tokenUserID, idToken, err := parseVKTokenResponse(tokenBody)
	if err != nil {
		return nil, "", err
	}
	if accessToken == "" {
		return nil, "", &AppError{401, "VK OAuth: пустой access_token"}
	}

	public := parseVKIDTokenPublicClaims(idToken)
	if fetched, err := s.fetchVKPublicInfo(ctx, idToken); err == nil {
		public = mergeVKPublicInfo(public, fetched)
	}

	userBody, err := s.fetchVKUserProfile(ctx, accessToken)
	if err != nil {
		return nil, "", err
	}
	vkID, fullName, emailFromProfile, err := parseVKUserInfo(userBody, tokenUserID)
	if err != nil {
		return nil, "", err
	}
	if strings.TrimSpace(vkID) == "" {
		vkID = strings.TrimSpace(public.UserID)
	}
	if strings.TrimSpace(fullName) == "" {
		fullName = strings.TrimSpace(strings.TrimSpace(public.FirstName + " " + public.LastInitial))
	}
	email := strings.ToLower(strings.TrimSpace(tokenEmail))
	if email == "" {
		email = strings.ToLower(strings.TrimSpace(emailFromProfile))
	}

	user, err := s.findOrCreateOAuthUser(ctx, "vk", vkID, email, "", fullName)
	if err != nil {
		return nil, "", err
	}
	_ = s.users.UpsertOAuthIdentity(ctx, "vk", strings.TrimSpace(vkID), user.ID)
	user = s.applyVKProfileEmail(ctx, user, email)
	user = s.applyVKProfileName(ctx, user, fullName)
	_ = s.users.SetEmailVerified(ctx, user.ID, false)
	_ = s.users.SetPhoneVerified(ctx, user.ID, false)

	sid := generateSessionID()
	sp := sessionPayload{
		UserID:              user.ID,
		Email:               user.Email,
		ProfileType:         user.ProfileType,
		AuthProvider:        "vk",
		RequireVKOnboarding: IsVKOnboardingRequired(user, true),
	}
	b, _ := json.Marshal(sp)
	if err := s.rdb.Set(ctx, sessionKeyPrefix+sid, b, sessionTTL).Err(); err != nil {
		return nil, "", err
	}

	var photo *string
	if user.Photo != nil && *user.Photo != "" {
		p := s.cfg.BaseURL + *user.Photo
		photo = &p
	} else if strings.TrimSpace(public.Avatar) != "" {
		a := strings.TrimSpace(public.Avatar)
		photo = &a
	}
	out := &signInResponse{Message: "Вы успешно авторизовались через VK!"}
	out.User.ID = user.ID
	out.User.Email = user.Email
	out.User.FullName = user.FullName
	out.User.PhoneNumber = user.PhoneNumber
	out.User.ProfileType = user.ProfileType
	out.User.Photo = photo
	return out, sid, nil
}

func (s *AuthService) CreateSessionForUserID(ctx context.Context, userID int32) (string, error) {
	u, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	sid := generateSessionID()
	sp := sessionPayload{UserID: u.ID, Email: u.Email, ProfileType: u.ProfileType, AuthProvider: "password"}
	b, _ := json.Marshal(sp)
	if err := s.rdb.Set(ctx, sessionKeyPrefix+sid, b, sessionTTL).Err(); err != nil {
		return "", err
	}
	return sid, nil
}

func isLegacyVKIdentity(u *domain.UserEntity) bool {
	if u == nil {
		return false
	}
	email := strings.ToLower(strings.TrimSpace(u.Email))
	phone := strings.ToUpper(strings.TrimSpace(u.PhoneNumber))
	return strings.HasSuffix(email, "@oauth.local") || strings.HasPrefix(phone, "VK_")
}

func IsVKOnboardingRequired(u *domain.UserEntity, isVKUser bool) bool {
	if u == nil || !isVKUser {
		return false
	}
	email := strings.ToLower(strings.TrimSpace(u.Email))
	phone := strings.TrimSpace(u.PhoneNumber)
	if email == "" || phone == "" {
		return true
	}
	if strings.HasSuffix(email, "@oauth.local") {
		return true
	}
	if strings.HasPrefix(strings.ToUpper(phone), "VK_") {
		return true
	}
	return !u.IsEmailVerified || !u.IsPhoneVerified
}

func (s *AuthService) VKOnboardingStatus(ctx context.Context, userID int32) (map[string]any, error) {
	u, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	hasVKIdentity, err := s.users.HasOAuthProviderForUser(ctx, userID, "vk")
	if err != nil {
		return nil, err
	}
	isVKUser := hasVKIdentity || isLegacyVKIdentity(u)
	return map[string]any{
		"required":        IsVKOnboardingRequired(u, isVKUser),
		"isEmailVerified": u.IsEmailVerified,
		"isPhoneVerified": u.IsPhoneVerified,
		"email":           u.Email,
		"phoneNumber":     u.PhoneNumber,
	}, nil
}

func (s *AuthService) IsVKOnboardingRequiredForUser(ctx context.Context, u *domain.UserEntity) bool {
	if u == nil {
		return false
	}
	hasVKIdentity, err := s.users.HasOAuthProviderForUser(ctx, u.ID, "vk")
	if err != nil {
		return false
	}
	isVKUser := hasVKIdentity || isLegacyVKIdentity(u)
	return IsVKOnboardingRequired(u, isVKUser)
}

func (s *AuthService) VKOnboardingStartEmail(ctx context.Context, userID int32, email string) (int32, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return 0, &AppError{400, "Нужно указать email"}
	}
	current, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	targetUserID := current.ID
	if otherID, err := s.users.FindUserIDByEmail(ctx, email); err != nil {
		return 0, err
	} else if otherID != nil && *otherID != current.ID {
		targetUserID = *otherID
		target, err := s.users.FindUserByID(ctx, targetUserID)
		if err != nil {
			return 0, err
		}
		// Merge VK profile data (except email/phone) into existing account.
		patch := repository.UserSettingsPatch{}
		if strings.TrimSpace(current.FullName) != "" {
			name := strings.TrimSpace(current.FullName)
			patch.FullName = &name
		}
		if current.Photo != nil && strings.TrimSpace(*current.Photo) != "" {
			ph := strings.TrimSpace(*current.Photo)
			patch.Photo = &ph
		}
		if patch.FullName != nil || patch.Photo != nil {
			_ = s.users.UpdateUserSettings(ctx, target.ID, patch)
		}
		if provider, externalID, ok := inferOAuthIdentityFromUser(current); ok {
			_ = s.users.UpsertOAuthIdentity(ctx, provider, externalID, target.ID)
		}
	}

	if err := s.users.SetEmail(ctx, targetUserID, email); err != nil {
		return 0, err
	}
	if err := s.users.SetEmailVerified(ctx, targetUserID, false); err != nil {
		return 0, err
	}
	if err := s.sendVKEmailCode(ctx, targetUserID, email); err != nil {
		return 0, err
	}
	return targetUserID, nil
}

func (s *AuthService) VKOnboardingVerifyEmailCode(ctx context.Context, userID int32, code string) error {
	code = strings.TrimSpace(code)
	raw, err := s.rdb.Get(ctx, vkVerifyEmailPrefix+code).Result()
	if err == redis.Nil || raw == "" {
		return &AppError{400, "Неверный код подтверждения"}
	}
	if err != nil {
		return err
	}
	var cached struct {
		UserID int32  `json:"userId"`
		Email  string `json:"email"`
		Code   string `json:"code"`
	}
	if err := json.Unmarshal([]byte(raw), &cached); err != nil {
		return err
	}
	if cached.Code != code || cached.UserID != userID {
		return &AppError{400, "Неверный код подтверждения"}
	}
	if err := s.users.SetEmail(ctx, cached.UserID, cached.Email); err != nil {
		return err
	}
	if err := s.users.SetEmailVerified(ctx, cached.UserID, true); err != nil {
		return err
	}
	_ = s.rdb.Del(ctx, vkVerifyEmailPrefix+code).Err()
	return nil
}

func (s *AuthService) VKOnboardingStartPhone(ctx context.Context, userID int32, phone string) error {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return &AppError{400, "Нужно указать номер телефона"}
	}
	if otherID, err := s.users.FindUserIDByPhone(ctx, phone); err != nil {
		return err
	} else if otherID != nil && *otherID != userID {
		return &AppError{400, "Пользователь с таким номером телефона уже существует"}
	}
	code := s.generateVerifyCode()
	payload, _ := json.Marshal(map[string]any{"userId": userID, "phone": phone, "code": code})
	if err := s.rdb.Set(ctx, vkVerifyPhonePrefix+code, payload, verifyPhoneTTL).Err(); err != nil {
		return err
	}
	return s.sendVKPhoneCode(ctx, phone, code)
}

func (s *AuthService) VKOnboardingVerifyPhoneCode(ctx context.Context, userID int32, code string) error {
	code = strings.TrimSpace(code)
	raw, err := s.rdb.Get(ctx, vkVerifyPhonePrefix+code).Bytes()
	if err == redis.Nil || len(raw) == 0 {
		return &AppError{400, "Неверный код подтверждения"}
	}
	if err != nil {
		return err
	}
	var cached struct {
		UserID int32  `json:"userId"`
		Phone  string `json:"phone"`
		Code   string `json:"code"`
	}
	if err := json.Unmarshal(raw, &cached); err != nil {
		return err
	}
	if cached.Code != code || cached.UserID != userID {
		return &AppError{400, "Неверный код подтверждения"}
	}
	if err := s.users.SetPhone(ctx, cached.UserID, cached.Phone); err != nil {
		return err
	}
	if err := s.users.SetPhoneVerified(ctx, cached.UserID, true); err != nil {
		return err
	}
	_ = s.rdb.Del(ctx, vkVerifyPhonePrefix+code).Err()
	return nil
}

func (s *AuthService) sendVKEmailCode(ctx context.Context, userID int32, email string) error {
	if strings.TrimSpace(s.cfg.SMTPHost) == "" {
		return &AppError{500, "SMTP не настроен (SMTP_HOST)"}
	}
	code := s.generateVerifyCode()
	payload, _ := json.Marshal(map[string]any{"userId": userID, "email": email, "code": code})
	if err := s.rdb.Set(ctx, vkVerifyEmailPrefix+code, payload, time.Hour).Err(); err != nil {
		return err
	}
	htmlBody, err := mailpkg.VerifyEmailHTML(code)
	if err != nil {
		return &AppError{500, "Не удалось сформировать письмо"}
	}
	fromAddr := strings.TrimSpace(s.cfg.SMTPFrom)
	if fromAddr == "" {
		fromAddr = strings.TrimSpace(s.cfg.SMTPUser)
	}
	if fromAddr == "" {
		return &AppError{500, "SMTP_FROM/SMTP_USER не задан"}
	}
	return mailpkg.SendHTMLSmart(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword,
		fromAddr, email, "Код подтверждения почты - Торгуй Сам", htmlBody, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure)
}

func (s *AuthService) sendVKPhoneCode(ctx context.Context, phone, code string) error {
	if s.cfg.MTSBearer != "" {
		body := map[string]any{
			"submits": []any{map[string]any{"msid": phone, "message": "Код подтверждения: " + code}},
			"naming":  "Torguisamru",
		}
		return s.httpPostJSON(ctx,
			"https://api.mts.ru/client-omni-adapter_production/1.0.2/mcom/messageManagement/messages",
			body,
			s.cfg.MTSBearer,
			nil,
		)
	}
	if s.cfg.NotisendAPIKey != "" {
		u := fmt.Sprintf(
			"https://sms.notisend.ru/api/message/send?project=%s&message=%s&recipients=%s&apikey=%s",
			url.QueryEscape(s.cfg.NotisendProject),
			url.QueryEscape("Код подтверждения: "+code),
			url.QueryEscape(phone),
			url.QueryEscape(s.cfg.NotisendAPIKey),
		)
		var resp map[string]any
		if err := s.httpGETJSON(ctx, u, &resp); err != nil {
			return err
		}
		return nil
	}
	return &AppError{500, "Не настроен SMS провайдер"}
}

func inferOAuthIdentityFromUser(u *domain.UserEntity) (provider, externalID string, ok bool) {
	if u == nil {
		return "", "", false
	}
	email := strings.ToLower(strings.TrimSpace(u.Email))
	phone := strings.ToUpper(strings.TrimSpace(u.PhoneNumber))
	if strings.HasPrefix(email, "vk_") && strings.HasSuffix(email, "@oauth.local") {
		ext := strings.TrimSuffix(strings.TrimPrefix(email, "vk_"), "@oauth.local")
		if ext != "" {
			return "vk", ext, true
		}
	}
	if strings.HasPrefix(phone, "VK_") {
		ext := strings.TrimPrefix(phone, "VK_")
		if i := strings.Index(ext, "_"); i > 0 {
			ext = ext[:i]
		}
		if ext != "" {
			return "vk", strings.ToLower(ext), true
		}
	}
	return "", "", false
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	if sessionID != "" {
		_ = s.rdb.Del(ctx, sessionKeyPrefix+sessionID).Err()
	}
	return nil
}

func (s *AuthService) UserFromSession(ctx context.Context, sessionID string) (*domain.UserEntity, error) {
	if sessionID == "" {
		return nil, nil
	}
	raw, err := s.rdb.Get(ctx, sessionKeyPrefix+sessionID).Bytes()
	if err == redis.Nil || len(raw) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var sp sessionPayload
	if err := json.Unmarshal(raw, &sp); err != nil {
		return nil, nil
	}
	u, err := s.users.FindUserByID(ctx, sp.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (s *AuthService) SessionRequiresVKOnboarding(ctx context.Context, sessionID string, u *domain.UserEntity) bool {
	if sessionID == "" || u == nil {
		return false
	}
	raw, err := s.rdb.Get(ctx, sessionKeyPrefix+sessionID).Bytes()
	if err != nil || len(raw) == 0 {
		return false
	}
	var sp sessionPayload
	if err := json.Unmarshal(raw, &sp); err != nil {
		return false
	}
	if strings.ToLower(strings.TrimSpace(sp.AuthProvider)) != "vk" {
		return false
	}
	if !sp.RequireVKOnboarding {
		return false
	}
	return IsVKOnboardingRequired(u, true)
}

func parseSessionIDCookie(cookieHeader string) string {
	for _, part := range strings.Split(cookieHeader, ";") {
		part = strings.TrimSpace(part)
		if after, ok := strings.CutPrefix(part, "session_id="); ok {
			return after
		}
	}
	return ""
}

// SocketUserFromCookie — handshake Socket.IO: cookie session_id → пользователь с ролью (чат + support WS).
func (s *AuthService) SocketUserFromCookie(ctx context.Context, cookieHeader string) (*domain.UserEntity, error) {
	sid := parseSessionIDCookie(cookieHeader)
	if sid == "" {
		return nil, fmt.Errorf("no session")
	}
	u, err := s.UserFromSession(ctx, sid)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("invalid session")
	}
	return u, nil
}

func (s *AuthService) Me(ctx context.Context, userID int32) (*domain.MeResponse, error) {
	u, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Пользователь не найден"}
		}
		return nil, err
	}
	var photo *string
	if u.Photo != nil && *u.Photo != "" {
		p := s.cfg.BaseURL + *u.Photo
		photo = &p
	}
	return &domain.MeResponse{
		ID: u.ID, Email: u.Email, FullName: u.FullName, PhoneNumber: u.PhoneNumber,
		ProfileType: u.ProfileType, Photo: photo, Rating: u.Rating, IsAnswersCall: u.IsAnswersCall,
		Role: u.RoleName,
	}, nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	return s.ForgotPasswordBy(ctx, "email", email, "")
}

func (s *AuthService) ForgotPasswordBy(ctx context.Context, where, email, phone string) error {
	where = strings.ToLower(strings.TrimSpace(where))
	if where == "" {
		where = "email"
	}

	var (
		u   *domain.UserEntity
		err error
	)
	switch where {
	case "email":
		email = strings.TrimSpace(email)
		if email == "" {
			return &AppError{400, "Нужно указать email"}
		}
		u, err = s.users.FindUserByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return &AppError{400, "Пользователя с такой почтой не существует"}
			}
			return err
		}
	case "sms":
		phone = strings.TrimSpace(phone)
		if phone == "" {
			return &AppError{400, "Нужно указать номер телефона"}
		}
		u, err = s.users.FindUserByLogin(ctx, phone)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return &AppError{400, "Пользователя с таким номером не существует"}
			}
			return err
		}
	default:
		return &AppError{400, "where должен быть email или sms"}
	}

	code, err := s.reserveUniqueForgotCode(ctx, u.ID)
	if err != nil {
		return err
	}

	if where == "sms" {
		if s.cfg.MTSBearer == "" {
			return &AppError{500, "MTS_TOKEN не задан"}
		}
		body := map[string]any{
			"submits": []any{
				map[string]any{"msid": u.PhoneNumber, "message": "Код восстановления пароля: " + code},
			},
			"naming": "Torguisamru",
		}
		if err := s.httpPostJSON(ctx,
			"https://api.mts.ru/client-omni-adapter_production/1.0.2/mcom/messageManagement/messages",
			body,
			s.cfg.MTSBearer,
			nil,
		); err != nil {
			return &AppError{400, "Ошибка отправки SMS: " + err.Error()}
		}
		return nil
	}

	if strings.TrimSpace(s.cfg.SMTPHost) == "" {
		return &AppError{500, "SMTP не настроен (SMTP_HOST)"}
	}
	htmlBody, err := mailpkg.ForgotPasswordHTML(code)
	if err != nil {
		return &AppError{500, "Не удалось сформировать письмо"}
	}
	fromAddr := strings.TrimSpace(s.cfg.SMTPFrom)
	if fromAddr == "" {
		fromAddr = strings.TrimSpace(s.cfg.SMTPUser)
	}
	if fromAddr == "" {
		return &AppError{500, "SMTP_FROM/SMTP_USER не задан"}
	}
	if err := mailpkg.SendHTMLSmart(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword,
		fromAddr, u.Email, "Код восстановления пароля - Торгуй Сам", htmlBody, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure); err != nil {
		return &AppError{400, "Ошибка отправки письма: " + err.Error()}
	}
	return nil
}

func (s *AuthService) reserveUniqueForgotCode(ctx context.Context, userID int32) (string, error) {
	for i := 0; i < 20; i++ {
		code := s.generateVerifyCode()
		payload := map[string]string{"id": fmt.Sprintf("%d", userID), "code": code}
		b, _ := json.Marshal(payload)
		ok, err := s.rdb.SetNX(ctx, forgotKeyPrefix+code, b, forgotPassTTL).Result()
		if err != nil {
			return "", err
		}
		if ok {
			return code, nil
		}
	}
	return "", &AppError{500, "Не удалось создать уникальный код восстановления"}
}

func (s *AuthService) VerifyForgotCode(ctx context.Context, code string) (int32, error) {
	raw, err := s.rdb.Get(ctx, forgotKeyPrefix+code).Bytes()
	if err == redis.Nil || len(raw) == 0 {
		return 0, &AppError{400, "Неверный код подтверждения"}
	}
	if err != nil {
		return 0, err
	}
	var cached struct {
		ID   string `json:"id"`
		Code string `json:"code"`
	}
	if err := json.Unmarshal(raw, &cached); err != nil {
		return 0, err
	}
	if cached.Code != code {
		return 0, &AppError{400, "Неверный код подтверждения"}
	}
	id64, err := strconv.ParseInt(cached.ID, 10, 32)
	if err != nil {
		return 0, &AppError{400, "Неверный код подтверждения"}
	}
	uid := int32(id64)
	_, err = s.users.FindUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, &AppError{404, "Такого пользователя не существует"}
		}
		return 0, err
	}
	if err := s.users.SetResetVerified(ctx, uid, true); err != nil {
		return 0, err
	}
	_ = s.rdb.Del(ctx, forgotKeyPrefix+code)
	return uid, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int32, password string) error {
	u, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &AppError{404, "Пользователь не найден"}
		}
		return err
	}
	if !u.IsResetVerified {
		return &AppError{403, "Требуется подтверждение сброса пароля"}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	return s.users.UpdatePassword(ctx, userID, string(hash))
}

func (s *AuthService) generateVerifyCode() string {
	return fmt.Sprintf("%d", 100000+mathrand.Intn(900000))
}

func generateSessionID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func generatePKCEVerifier() string {
	b := make([]byte, 48)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func pkceS256(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func (s *AuthService) httpGETJSON(ctx context.Context, u string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, truncateForErr(string(body)))
	}
	return json.Unmarshal(body, out)
}

func (s *AuthService) httpPostJSON(ctx context.Context, urlStr string, body any, bearer string, out any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, truncateForErr(string(b)))
	}
	if out != nil {
		return json.Unmarshal(b, out)
	}
	return nil
}

// truncateForErr — не раздуваем лог/ответ целым HTML от провайдера.
func truncateForErr(s string) string {
	const max = 512
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

func (s *AuthService) findOrCreateOAuthUser(ctx context.Context, provider, externalID, email, phone, fullName string) (*domain.UserEntity, error) {
	providerSlug := strings.ToLower(strings.TrimSpace(provider))
	if providerSlug == "" {
		providerSlug = "oauth"
	}
	externalID = strings.TrimSpace(externalID)
	if externalID == "" {
		externalID = generateSessionID()[:12]
	}

	// 1) Сначала ищем уже созданный OAuth-аккаунт по provider+externalID.
	if existingID, err := s.users.FindOAuthUserIDByProviderExternalID(ctx, providerSlug, externalID); err == nil && existingID != nil {
		u, err := s.users.FindUserByID(ctx, *existingID)
		if err == nil {
			_ = s.users.UpsertOAuthIdentity(ctx, providerSlug, externalID, u.ID)
			return u, nil
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	email = strings.ToLower(strings.TrimSpace(email))
	phone = strings.TrimSpace(phone)
	fullName = strings.TrimSpace(fullName)

	// 2) Если VK вернул email и он уже занят другим аккаунтом — не создаём дубль.
	if email != "" {
		if _, err := s.users.FindUserByEmail(ctx, email); err == nil {
			return nil, &AppError{400, "Почта уже занята"}
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}
	if phone != "" {
		if u, err := s.users.FindUserByLogin(ctx, phone); err == nil {
			_ = s.users.UpsertOAuthIdentity(ctx, providerSlug, externalID, u.ID)
			return u, nil
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	if email == "" {
		email = providerSlug + "_" + externalID + "@oauth.local"
	}
	if phone == "" {
		phone = strings.ToUpper(providerSlug) + "_" + externalID
	}
	if fullName == "" {
		fullName = strings.ToUpper(providerSlug) + " User"
	}

	roleID, err := s.defaultUserRoleID(ctx)
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(generateSessionID()), bcryptCost)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 5; i++ {
		uid, err := s.users.GenerateUniqueUserID(ctx)
		if err != nil {
			return nil, err
		}
		candidateEmail := email
		candidatePhone := phone
		if i > 0 {
			candidateEmail = fmt.Sprintf("%s_%s_%d@oauth.local", providerSlug, externalID, i)
			candidatePhone = fmt.Sprintf("%s_%s_%d", strings.ToUpper(providerSlug), externalID, i)
		}
		if err := s.users.InsertUser(ctx, uid, fullName, candidateEmail, candidatePhone, string(hash), roleID); err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
				// Для VK иногда параллельные запросы могут создать гонку.
				// Пробуем найти уже созданный аккаунт по externalID и вернуть его.
				if existingID, findErr := s.users.FindOAuthUserIDByProviderExternalID(ctx, providerSlug, externalID); findErr == nil && existingID != nil {
					if u, getErr := s.users.FindUserByID(ctx, *existingID); getErr == nil {
						_ = s.users.UpsertOAuthIdentity(ctx, providerSlug, externalID, u.ID)
						return u, nil
					}
				}
				if strings.ToLower(providerSlug) == "vk" && email != "" && !strings.HasSuffix(email, "@oauth.local") {
					return nil, &AppError{400, "Почта уже занята"}
				}
				continue
			}
			return nil, err
		}
		u, err := s.users.FindUserByID(ctx, uid)
		if err != nil {
			return nil, err
		}
		_ = s.users.UpsertOAuthIdentity(ctx, providerSlug, externalID, u.ID)
		return u, nil
	}
	return nil, &AppError{500, "Не удалось создать пользователя MAX"}
}

func (s *AuthService) findOrCreateTrustedOAuthUser(ctx context.Context, p oauthProfile) (*domain.UserEntity, error) {
	provider := strings.ToLower(strings.TrimSpace(p.Provider))
	if provider == "" {
		provider = "oauth"
	}
	externalID := strings.TrimSpace(p.ExternalID)
	if externalID == "" {
		return nil, &AppError{400, "Пустой внешний идентификатор OAuth"}
	}
	if existingID, err := s.users.FindOAuthUserIDByProviderExternalID(ctx, provider, externalID); err == nil && existingID != nil {
		u, err := s.users.FindUserByID(ctx, *existingID)
		if err == nil {
			u = s.applyTrustedOAuthProfile(ctx, u, p)
			return u, nil
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	email := strings.ToLower(strings.TrimSpace(p.Email))
	if email != "" {
		if u, err := s.users.FindUserByEmail(ctx, email); err == nil {
			_ = s.users.UpsertOAuthIdentity(ctx, provider, externalID, u.ID)
			u = s.applyTrustedOAuthProfile(ctx, u, p)
			return u, nil
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	phone := strings.TrimSpace(p.Phone)
	if phone != "" {
		if u, err := s.users.FindUserByLogin(ctx, phone); err == nil {
			_ = s.users.UpsertOAuthIdentity(ctx, provider, externalID, u.ID)
			u = s.applyTrustedOAuthProfile(ctx, u, p)
			return u, nil
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	u, err := s.findOrCreateOAuthUser(ctx, provider, externalID, email, phone, p.FullName)
	if err != nil {
		return nil, err
	}
	u = s.applyTrustedOAuthProfile(ctx, u, p)
	return u, nil
}

func (s *AuthService) applyTrustedOAuthProfile(ctx context.Context, u *domain.UserEntity, p oauthProfile) *domain.UserEntity {
	if u == nil {
		return nil
	}
	patch := repository.UserSettingsPatch{}
	adminPatch := repository.AdminUserPatch{}
	changed := false
	adminChanged := false

	fullName := strings.TrimSpace(p.FullName)
	if fullName != "" && (strings.TrimSpace(u.FullName) == "" || strings.Contains(strings.ToLower(strings.TrimSpace(u.FullName)), " user")) {
		patch.FullName = &fullName
		changed = true
	}
	email := strings.ToLower(strings.TrimSpace(p.Email))
	if email != "" && (strings.TrimSpace(u.Email) == "" || strings.HasSuffix(strings.ToLower(strings.TrimSpace(u.Email)), "@oauth.local")) {
		adminPatch.Email = &email
		adminChanged = true
	}
	phone := strings.TrimSpace(p.Phone)
	if phone != "" && isDealPhoneSynthetic(u.PhoneNumber) {
		patch.PhoneNumber = &phone
		changed = true
	}
	if changed {
		_ = s.users.UpdateUserSettings(ctx, u.ID, patch)
	}
	if adminChanged {
		_ = s.users.UpdateUserAdmin(ctx, u.ID, adminPatch)
	}
	if email != "" {
		_ = s.users.SetEmailVerified(ctx, u.ID, true)
	}
	if phone != "" && !isDealPhoneSynthetic(phone) {
		_ = s.users.SetPhoneVerified(ctx, u.ID, true)
	}
	refreshed, err := s.users.FindUserByID(ctx, u.ID)
	if err == nil {
		return refreshed
	}
	return u
}

func parseTIDTokenResponse(body []byte) (accessToken, idToken string, err error) {
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", "", err
	}
	if e, _ := raw["error"].(string); strings.TrimSpace(e) != "" {
		desc, _ := raw["error_description"].(string)
		msg := e
		if strings.TrimSpace(desc) != "" {
			msg += ": " + desc
		}
		return "", "", &AppError{401, "T-ID token error: " + msg}
	}
	accessToken, _ = raw["access_token"].(string)
	idToken, _ = raw["id_token"].(string)
	return strings.TrimSpace(accessToken), strings.TrimSpace(idToken), nil
}

func parseOAuthJWTClaims(token string) map[string]any {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil
	}
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil
	}
	return claims
}

func (s *AuthService) fetchOAuthUserInfo(ctx context.Context, endpoint, accessToken string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(accessToken))
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("userinfo status=%d: %s", res.StatusCode, truncateForErr(string(body)))
	}
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func fillOAuthProfileFromClaims(profile *oauthProfile, claims map[string]any) {
	if profile == nil || claims == nil {
		return
	}
	if profile.ExternalID == "" {
		profile.ExternalID = pickVKString(claims, "sub", "user_id", "id", "customer_id")
	}
	if profile.Email == "" {
		profile.Email = strings.ToLower(strings.TrimSpace(pickVKString(claims, "email", "mail")))
	}
	if profile.Phone == "" {
		profile.Phone = strings.TrimSpace(pickVKString(claims, "phone_number", "phone", "msisdn"))
	}
	if profile.FullName == "" {
		name := strings.TrimSpace(pickVKString(claims, "name", "full_name"))
		if name == "" {
			first := strings.TrimSpace(pickVKString(claims, "given_name", "first_name"))
			last := strings.TrimSpace(pickVKString(claims, "family_name", "last_name"))
			name = strings.TrimSpace(first + " " + last)
		}
		profile.FullName = name
	}
	if profile.Avatar == "" {
		profile.Avatar = strings.TrimSpace(pickVKString(claims, "picture", "avatar", "photo"))
	}
}

const vkLegacyUsersGetURL = "https://api.vk.com/method/users.get"

func parseVKTokenResponse(body []byte) (accessToken, email string, userID int64, idToken string, err error) {
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", "", 0, "", err
	}
	if e, _ := raw["error"].(string); strings.TrimSpace(e) != "" {
		return "", "", 0, "", &AppError{401, "VK OAuth token error: " + e}
	}
	accessToken, _ = raw["access_token"].(string)
	email, _ = raw["email"].(string)
	idToken, _ = raw["id_token"].(string)
	userID = vkAnyToInt64(raw["user_id"])
	return strings.TrimSpace(accessToken), strings.TrimSpace(email), userID, strings.TrimSpace(idToken), nil
}

func parseVKIDTokenPublicClaims(idToken string) vkPublicInfo {
	idToken = strings.TrimSpace(idToken)
	if idToken == "" {
		return vkPublicInfo{}
	}
	parts := strings.Split(idToken, ".")
	if len(parts) < 2 {
		return vkPublicInfo{}
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return vkPublicInfo{}
	}
	var m map[string]any
	if err := json.Unmarshal(payload, &m); err != nil {
		return vkPublicInfo{}
	}
	out := vkPublicInfo{
		UserID:      pickVKString(m, "sub", "user_id", "id"),
		FirstName:   pickVKString(m, "first_name", "given_name"),
		LastInitial: pickVKString(m, "last_name_initial", "last_initial"),
		Avatar:      pickVKString(m, "picture", "avatar", "photo"),
		EmailMasked: pickVKString(m, "email", "email_masked"),
		PhoneMasked: pickVKString(m, "phone", "phone_masked"),
	}
	return out
}

func (s *AuthService) fetchVKPublicInfo(ctx context.Context, idToken string) (vkPublicInfo, error) {
	idToken = strings.TrimSpace(idToken)
	if idToken == "" {
		return vkPublicInfo{}, fmt.Errorf("empty id_token")
	}
	form := url.Values{}
	form.Set("id_token", idToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, vkIDPublicInfoURL, strings.NewReader(form.Encode()))
	if err != nil {
		return vkPublicInfo{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := s.client.Do(req)
	if err != nil {
		return vkPublicInfo{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return vkPublicInfo{}, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return vkPublicInfo{}, fmt.Errorf("public_info status=%d: %s", res.StatusCode, truncateForErr(string(body)))
	}

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return vkPublicInfo{}, err
	}
	return vkPublicInfo{
		UserID:      pickVKString(raw, "user_id", "sub", "id", "response.user_id", "response.sub", "response.id"),
		FirstName:   pickVKString(raw, "first_name", "response.first_name"),
		LastInitial: pickVKString(raw, "last_name_initial", "last_initial", "response.last_name_initial", "response.last_initial"),
		Avatar:      pickVKString(raw, "avatar", "picture", "photo", "response.avatar", "response.picture", "response.photo"),
		EmailMasked: pickVKString(raw, "email", "email_masked", "response.email", "response.email_masked"),
		PhoneMasked: pickVKString(raw, "phone", "phone_masked", "response.phone", "response.phone_masked"),
	}, nil
}

func mergeVKPublicInfo(a, b vkPublicInfo) vkPublicInfo {
	out := a
	if strings.TrimSpace(out.UserID) == "" {
		out.UserID = b.UserID
	}
	if strings.TrimSpace(out.FirstName) == "" {
		out.FirstName = b.FirstName
	}
	if strings.TrimSpace(out.LastInitial) == "" {
		out.LastInitial = b.LastInitial
	}
	if strings.TrimSpace(out.Avatar) == "" {
		out.Avatar = b.Avatar
	}
	if strings.TrimSpace(out.EmailMasked) == "" {
		out.EmailMasked = b.EmailMasked
	}
	if strings.TrimSpace(out.PhoneMasked) == "" {
		out.PhoneMasked = b.PhoneMasked
	}
	return out
}

func (s *AuthService) fetchVKUserProfile(ctx context.Context, accessToken string) ([]byte, error) {
	if s.cfg.VkIDEnabled {
		body, err := s.postVKIDUserInfo(ctx, accessToken)
		if err != nil {
			return nil, err
		}
		_, name, _, parseErr := parseVKUserInfo(body, 0)
		if parseErr == nil && strings.TrimSpace(name) != "" && !isVKPlaceholderFullName(name) {
			return body, nil
		}
		legacyBody, legErr := s.getVKLegacyUsersGet(ctx, accessToken)
		if legErr == nil {
			_, legName, _, _ := parseVKUserInfo(legacyBody, 0)
			if strings.TrimSpace(legName) != "" && !isVKPlaceholderFullName(legName) {
				return legacyBody, nil
			}
		}
		return body, nil
	}
	return s.getVKLegacyUsersGet(ctx, accessToken)
}

func (s *AuthService) postVKIDUserInfo(ctx context.Context, accessToken string) ([]byte, error) {
	form := url.Values{}
	form.Set("client_id", s.cfg.VkOAuthClientID)
	form.Set("access_token", accessToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.VkOAuthUserInfoURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return s.doVKHTTP(req)
}

func (s *AuthService) getVKLegacyUsersGet(ctx context.Context, accessToken string) ([]byte, error) {
	q := url.Values{}
	q.Set("access_token", accessToken)
	q.Set("v", "5.131")
	q.Set("fields", "photo_200,screen_name")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, vkLegacyUsersGetURL+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	return s.doVKHTTP(req)
}

func (s *AuthService) doVKHTTP(req *http.Request) ([]byte, error) {
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, &AppError{401, "VK OAuth userinfo error: " + truncateForErr(string(body))}
	}
	return body, nil
}

func isVKPlaceholderFullName(name string) bool {
	name = strings.TrimSpace(name)
	return name == "" ||
		strings.HasPrefix(name, "Пользователь VK") ||
		strings.EqualFold(name, "VK USER")
}

func (s *AuthService) applyVKProfileName(ctx context.Context, user *domain.UserEntity, fullName string) *domain.UserEntity {
	fullName = strings.TrimSpace(fullName)
	currentName := strings.TrimSpace(user.FullName)
	if fullName == "" || isVKPlaceholderFullName(fullName) || strings.EqualFold(fullName, currentName) {
		return user
	}
	if err := s.users.UpdateUserSettings(ctx, user.ID, repository.UserSettingsPatch{FullName: &fullName}); err != nil {
		return user
	}
	user.FullName = fullName
	return user
}

func (s *AuthService) applyVKProfileEmail(ctx context.Context, user *domain.UserEntity, email string) *domain.UserEntity {
	email = strings.ToLower(strings.TrimSpace(email))
	if user == nil || email == "" {
		return user
	}
	currentEmail := strings.ToLower(strings.TrimSpace(user.Email))
	if currentEmail == email {
		return user
	}

	otherID, err := s.users.FindUserIDByEmail(ctx, email)
	if err != nil {
		return user
	}
	if otherID != nil && *otherID != user.ID {
		// Email already belongs to another account. Keep current user as is.
		return user
	}

	// Auto-bind VK email when account still has placeholder oauth.local email.
	if strings.HasSuffix(currentEmail, "@oauth.local") || currentEmail == "" {
		if err := s.users.SetEmail(ctx, user.ID, email); err != nil {
			return user
		}
		user.Email = email
	}
	return user
}

func joinVKFullName(last, first, middle string) string {
	return strings.TrimSpace(strings.Join([]string{
		strings.TrimSpace(last),
		strings.TrimSpace(first),
		strings.TrimSpace(middle),
	}, " "))
}

func vkExternalIDFromAny(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case int:
		return strconv.Itoa(t)
	default:
		return ""
	}
}

func vkAnyToInt64(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case int64:
		return t
	case int:
		return int64(t)
	case string:
		n, _ := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		return n
	default:
		return 0
	}
}

func parseVKUserInfo(body []byte, fallbackID int64) (id, fullName, email string, err error) {
	vkFallbackName := func(externalID string, numericID int64) string {
		if strings.TrimSpace(externalID) != "" {
			return "Пользователь VK #" + strings.TrimSpace(externalID)
		}
		if numericID > 0 {
			return "Пользователь VK #" + strconv.FormatInt(numericID, 10)
		}
		return "Пользователь VK"
	}

	var legacy struct {
		Response []struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"response"`
		Error any `json:"error"`
	}
	if e := json.Unmarshal(body, &legacy); e == nil && legacy.Error == nil && len(legacy.Response) > 0 {
		id = strconv.FormatInt(legacy.Response[0].ID, 10)
		first := strings.TrimSpace(legacy.Response[0].FirstName)
		last := strings.TrimSpace(legacy.Response[0].LastName)
		fullName = strings.TrimSpace(first + " " + last)
		if fullName == "" {
			fullName = vkFallbackName(id, legacy.Response[0].ID)
		}
		return id, fullName, "", nil
	}

	var wrapped struct {
		User struct {
			UserID    any    `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
		} `json:"user"`
	}
	if e := json.Unmarshal(body, &wrapped); e == nil && wrapped.User.UserID != nil {
		id = vkExternalIDFromAny(wrapped.User.UserID)
		if id == "" && fallbackID > 0 {
			id = strconv.FormatInt(fallbackID, 10)
		}
		fullName = joinVKFullName(wrapped.User.LastName, wrapped.User.FirstName, "")
		email = strings.TrimSpace(wrapped.User.Email)
		if fullName != "" {
			return id, fullName, email, nil
		}
	}

	var vkid struct {
		Sub       string `json:"sub"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		MidName   string `json:"middle_name"`
		Patronym  string `json:"patronymic"`
		Name      string `json:"name"`
	}
	if e := json.Unmarshal(body, &vkid); e == nil {
		id = strings.TrimSpace(vkid.Sub)
		if id == "" && fallbackID > 0 {
			id = strconv.FormatInt(fallbackID, 10)
		}
		middle := strings.TrimSpace(vkid.MidName)
		if middle == "" {
			middle = strings.TrimSpace(vkid.Patronym)
		}
		fullName = joinVKFullName(vkid.LastName, vkid.FirstName, middle)
		if fullName == "" {
			fullName = strings.TrimSpace(vkid.Name)
		}
		email = strings.TrimSpace(vkid.Email)
		if fullName != "" {
			return id, fullName, email, nil
		}
	}

	var generic map[string]any
	if e := json.Unmarshal(body, &generic); e == nil {
		if id == "" {
			id = pickVKString(generic,
				"sub",
				"user_id",
				"id",
				"user.sub",
				"user.user_id",
				"user.id",
				"user_info.sub",
				"user_info.user_id",
				"user_info.id",
			)
			if id == "" && fallbackID > 0 {
				id = strconv.FormatInt(fallbackID, 10)
			}
		}
		if email == "" {
			email = pickVKString(generic, "email", "user.email", "user_info.email")
		}
		last := pickVKString(generic, "last_name", "user.last_name", "user_info.last_name")
		first := pickVKString(generic, "first_name", "user.first_name", "user_info.first_name")
		middle := pickVKString(generic, "middle_name", "patronymic", "user.middle_name", "user.patronymic", "user_info.middle_name", "user_info.patronymic")
		fullName = joinVKFullName(last, first, middle)
		if fullName == "" {
			fullName = pickVKString(generic, "name", "user.name", "user_info.name")
		}
		if fullName != "" {
			return id, fullName, strings.TrimSpace(email), nil
		}
	}

	if id == "" && fallbackID > 0 {
		id = strconv.FormatInt(fallbackID, 10)
	}
	if id != "" && fullName == "" {
		return id, vkFallbackName(id, fallbackID), "", nil
	}
	return "", "", "", &AppError{401, "VK OAuth userinfo parse error"}
}

func pickVKString(payload map[string]any, paths ...string) string {
	for _, path := range paths {
		raw, ok := pickVKPath(payload, path)
		if !ok {
			continue
		}
		switch v := raw.(type) {
		case string:
			if s := strings.TrimSpace(v); s != "" {
				return s
			}
		case float64:
			return strconv.FormatInt(int64(v), 10)
		case int64:
			return strconv.FormatInt(v, 10)
		case int:
			return strconv.Itoa(v)
		}
	}
	return ""
}

func pickVKPath(payload map[string]any, path string) (any, bool) {
	cur := any(payload)
	parts := strings.Split(path, ".")
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		next, ok := m[p]
		if !ok {
			return nil, false
		}
		cur = next
	}
	return cur, true
}

// ──────────────────────────────────────────────
// Яндекс OAuth
// ──────────────────────────────────────────────

const (
	yandexAuthorizeURL = "https://oauth.yandex.ru/authorize"
	yandexTokenURL     = "https://oauth.yandex.ru/token"
	yandexUserInfoURL  = "https://login.yandex.ru/info"
	yandexPKCEPrefix   = "yandex-pkce:"
	yandexPKCETTL      = 10 * time.Minute
)

// YandexAuthURL генерирует URL для редиректа на Яндекс OAuth.
func (s *AuthService) YandexAuthURL(state string) (string, error) {
	if strings.TrimSpace(s.cfg.YandexClientID) == "" {
		return "", &AppError{500, "Яндекс OAuth не настроен: YANDEX_OAUTH_CLIENT_ID"}
	}
	if strings.TrimSpace(s.cfg.YandexRedirectURI) == "" {
		return "", &AppError{500, "Яндекс OAuth не настроен: YANDEX_OAUTH_REDIRECT_URI"}
	}
	if strings.TrimSpace(state) == "" {
		state = generateSessionID()
	}
	// Сохраняем state в Redis для проверки при callback
	if err := s.rdb.Set(context.Background(), yandexPKCEPrefix+state, "1", yandexPKCETTL).Err(); err != nil {
		return "", err
	}
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", s.cfg.YandexClientID)
	q.Set("redirect_uri", s.cfg.YandexRedirectURI)
	q.Set("state", state)
	q.Set("scope", "login:email login:info login:avatar")
	return yandexAuthorizeURL + "?" + q.Encode(), nil
}

// SignInWithYandex обменивает code на токен, получает профиль и создаёт/находит пользователя.
func (s *AuthService) SignInWithYandex(ctx context.Context, code, state string) (*signInResponse, string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, "", &AppError{400, "Нужен code"}
	}
	if strings.TrimSpace(s.cfg.YandexClientID) == "" || strings.TrimSpace(s.cfg.YandexClientSecret) == "" {
		return nil, "", &AppError{500, "Яндекс OAuth не настроен в .env"}
	}

	// Проверяем state
	state = strings.TrimSpace(state)
	if state != "" {
		val, err := s.rdb.Get(ctx, yandexPKCEPrefix+state).Result()
		if err == redis.Nil || val == "" {
			return nil, "", &AppError{401, "Яндекс OAuth: истекла сессия авторизации"}
		}
		if err != nil {
			return nil, "", err
		}
		_ = s.rdb.Del(ctx, yandexPKCEPrefix+state)
	}

	// Обмен code → access_token
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", s.cfg.YandexClientID)
	form.Set("client_secret", s.cfg.YandexClientSecret)
	form.Set("redirect_uri", s.cfg.YandexRedirectURI)

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, yandexTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, "", err
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenRes, err := s.client.Do(tokenReq)
	if err != nil {
		return nil, "", err
	}
	defer tokenRes.Body.Close()
	tokenBody, err := io.ReadAll(tokenRes.Body)
	if err != nil {
		return nil, "", err
	}
	if tokenRes.StatusCode < 200 || tokenRes.StatusCode >= 300 {
		return nil, "", &AppError{401, "Яндекс OAuth token error: " + truncateForErr(string(tokenBody))}
	}

	var tokenData struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.Unmarshal(tokenBody, &tokenData); err != nil {
		return nil, "", err
	}
	if tokenData.Error != "" {
		return nil, "", &AppError{401, "Яндекс OAuth: " + tokenData.ErrorDesc}
	}
	if tokenData.AccessToken == "" {
		return nil, "", &AppError{401, "Яндекс OAuth: пустой access_token"}
	}

	// Получаем профиль пользователя
	profile, err := s.fetchYandexUserInfo(ctx, tokenData.AccessToken)
	if err != nil {
		return nil, "", err
	}

	user, err := s.findOrCreateTrustedOAuthUser(ctx, *profile)
	if err != nil {
		return nil, "", err
	}
	_ = s.users.UpsertOAuthIdentity(ctx, "yandex", profile.ExternalID, user.ID)

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
	} else if strings.TrimSpace(profile.Avatar) != "" {
		a := strings.TrimSpace(profile.Avatar)
		photo = &a
	}

	out := &signInResponse{Message: "Вы успешно авторизовались через Яндекс!"}
	out.User.ID = user.ID
	out.User.Email = user.Email
	out.User.FullName = user.FullName
	out.User.PhoneNumber = user.PhoneNumber
	out.User.ProfileType = user.ProfileType
	out.User.Photo = photo
	return out, sid, nil
}

// fetchYandexUserInfo запрашивает профиль пользователя у Яндекс.
func (s *AuthService) fetchYandexUserInfo(ctx context.Context, accessToken string) (*oauthProfile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, yandexUserInfoURL+"?format=json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "OAuth "+accessToken)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, &AppError{401, "Яндекс OAuth userinfo error: " + truncateForErr(string(body))}
	}

	var info struct {
		ID          string `json:"id"`
		Login       string `json:"login"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		DisplayName string `json:"display_name"`
		RealName    string `json:"real_name"`
		DefaultEmail string `json:"default_email"`
		Emails      []string `json:"emails"`
		DefaultAvatarID string `json:"default_avatar_id"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	if strings.TrimSpace(info.ID) == "" {
		return nil, &AppError{401, "Яндекс не вернул идентификатор пользователя"}
	}

	// Имя
	fullName := strings.TrimSpace(strings.TrimSpace(info.FirstName) + " " + strings.TrimSpace(info.LastName))
	if fullName == "" {
		fullName = strings.TrimSpace(info.DisplayName)
	}
	if fullName == "" {
		fullName = strings.TrimSpace(info.RealName)
	}
	if fullName == "" {
		fullName = strings.TrimSpace(info.Login)
	}

	// Email
	email := strings.ToLower(strings.TrimSpace(info.DefaultEmail))
	if email == "" && len(info.Emails) > 0 {
		email = strings.ToLower(strings.TrimSpace(info.Emails[0]))
	}

	// Аватар
	var avatar string
	if strings.TrimSpace(info.DefaultAvatarID) != "" {
		avatar = "https://avatars.yandex.net/get-yapic/" + info.DefaultAvatarID + "/islands-200"
	}

	return &oauthProfile{
		Provider:   "yandex",
		ExternalID: info.ID,
		Email:      email,
		FullName:   fullName,
		Avatar:     avatar,
	}, nil
}

