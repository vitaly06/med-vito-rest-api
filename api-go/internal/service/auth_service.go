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
	sessionTTL       = 30 * 24 * time.Hour
	verifyPhoneTTL   = time.Hour
	forgotPassTTL    = time.Hour
	bcryptCost       = 10
	sessionKeyPrefix = "session:"
	verifyKeyPrefix  = "verify-phone:"
	forgotKeyPrefix  = "forgot-password:"
	vkidPKCEPrefix   = "vkid-pkce:"
	vkidPKCETTL      = 10 * time.Minute
)

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
	UserID      int32  `json:"userId"`
	Email       string `json:"email"`
	ProfileType string `json:"profileType"`
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
	sp := sessionPayload{UserID: u.ID, Email: u.Email, ProfileType: u.ProfileType}
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
			return nil, "", &AppError{401, "Неверные данные для входа"}
		}
		return nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", &AppError{401, "Неверные данные для входа"}
	}
	sid := generateSessionID()
	sp := sessionPayload{UserID: u.ID, Email: u.Email, ProfileType: u.ProfileType}
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

func (s *AuthService) SignInWithVK(ctx context.Context, code, state string) (*signInResponse, string, error) {
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
	var tokenPayload struct {
		AccessToken string `json:"access_token"`
		UserID      int64  `json:"user_id"`
		Email       string `json:"email"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(tokenBody, &tokenPayload); err != nil {
		return nil, "", err
	}
	if strings.TrimSpace(tokenPayload.Error) != "" {
		return nil, "", &AppError{401, "VK OAuth token error: " + tokenPayload.Error}
	}
	if strings.TrimSpace(tokenPayload.AccessToken) == "" {
		return nil, "", &AppError{401, "VK OAuth: пустой access_token"}
	}

	userQ := url.Values{}
	if !s.cfg.VkIDEnabled {
		userQ.Set("access_token", tokenPayload.AccessToken)
		userQ.Set("v", "5.131")
		userQ.Set("fields", "photo_200,screen_name")
	}
	userURL := s.cfg.VkOAuthUserInfoURL
	if enc := userQ.Encode(); enc != "" {
		userURL += "?" + enc
	}
	userReq, err := http.NewRequestWithContext(ctx, http.MethodGet, userURL, nil)
	if err != nil {
		return nil, "", err
	}
	userReq.Header.Set("Authorization", "Bearer "+tokenPayload.AccessToken)
	userRes, err := s.client.Do(userReq)
	if err != nil {
		return nil, "", err
	}
	defer userRes.Body.Close()
	userBody, err := io.ReadAll(userRes.Body)
	if err != nil {
		return nil, "", err
	}
	if userRes.StatusCode < 200 || userRes.StatusCode >= 300 {
		return nil, "", &AppError{401, "VK OAuth userinfo error: " + truncateForErr(string(userBody))}
	}
	vkID, fullName, emailFromProfile, err := parseVKUserInfo(userBody, tokenPayload.UserID)
	if err != nil {
		return nil, "", err
	}
	email := strings.ToLower(strings.TrimSpace(tokenPayload.Email))
	if email == "" {
		email = strings.ToLower(strings.TrimSpace(emailFromProfile))
	}

	user, err := s.findOrCreateOAuthUser(ctx, "vk", vkID, email, "", fullName)
	if err != nil {
		return nil, "", err
	}

	sid := generateSessionID()
	sp := sessionPayload{UserID: user.ID, Email: user.Email, ProfileType: user.ProfileType}
	b, _ := json.Marshal(sp)
	if err := s.rdb.Set(ctx, sessionKeyPrefix+sid, b, sessionTTL).Err(); err != nil {
		return nil, "", err
	}

	var photo *string
	if user.Photo != nil && *user.Photo != "" {
		p := s.cfg.BaseURL + *user.Photo
		photo = &p
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
	if strings.TrimSpace(s.cfg.SMTPHost) == "" {
		return &AppError{500, "SMTP не настроен (SMTP_HOST)"}
	}
	u, err := s.users.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &AppError{400, "Пользователя с такой почтой не существует"}
		}
		return err
	}
	code := s.generateVerifyCode()
	payload := map[string]string{"id": fmt.Sprintf("%d", u.ID), "code": code}
	b, _ := json.Marshal(payload)
	if err := s.rdb.Set(ctx, forgotKeyPrefix+code, b, forgotPassTTL).Err(); err != nil {
		return err
	}
	htmlBody, err := mailpkg.ForgotPasswordHTML(code)
	if err != nil {
		return &AppError{500, "Не удалось сформировать письмо"}
	}
	if err := mailpkg.SendHTMLSmart(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword,
		s.cfg.SMTPFrom, email, "Код восстановления пароля - Торгуй Сам", htmlBody, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure); err != nil {
		return &AppError{400, "Ошибка отправки письма: " + err.Error()}
	}
	return nil
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
	if email != "" {
		if u, err := s.users.FindUserByEmail(ctx, email); err == nil {
			return u, nil
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}
	if phone != "" {
		if u, err := s.users.FindUserByLogin(ctx, phone); err == nil {
			return u, nil
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	providerSlug := strings.ToLower(strings.TrimSpace(provider))
	if providerSlug == "" {
		providerSlug = "oauth"
	}
	if externalID == "" {
		externalID = generateSessionID()[:12]
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
				continue
			}
			return nil, err
		}
		return s.users.FindUserByID(ctx, uid)
	}
	return nil, &AppError{500, "Не удалось создать пользователя MAX"}
}

func parseVKUserInfo(body []byte, fallbackID int64) (id, fullName, email string, err error) {
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
			fullName = "VK User"
		}
		return id, fullName, "", nil
	}

	var vkid struct {
		Sub       string `json:"sub"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Name      string `json:"name"`
	}
	if e := json.Unmarshal(body, &vkid); e == nil {
		id = strings.TrimSpace(vkid.Sub)
		if id == "" && fallbackID > 0 {
			id = strconv.FormatInt(fallbackID, 10)
		}
		fullName = strings.TrimSpace(vkid.Name)
		if fullName == "" {
			fullName = strings.TrimSpace(strings.TrimSpace(vkid.FirstName) + " " + strings.TrimSpace(vkid.LastName))
		}
		if fullName == "" {
			fullName = "VK User"
		}
		return id, fullName, strings.TrimSpace(vkid.Email), nil
	}

	if fallbackID > 0 {
		return strconv.FormatInt(fallbackID, 10), "VK User", "", nil
	}
	return "", "", "", &AppError{401, "VK OAuth userinfo parse error"}
}
