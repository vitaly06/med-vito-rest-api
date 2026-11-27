package httpserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func setSessionCookies(c *fiber.Ctx, cfg config.Config, sid string, maxAgeSec int) {
	var expires time.Time
	if maxAgeSec > 0 {
		expires = time.Now().Add(time.Duration(maxAgeSec) * time.Second)
	} else {
		expires = time.Now().Add(-24 * time.Hour)
		maxAgeSec = -1
	}

	isHTTPS := cfg.Production ||
		strings.EqualFold(c.Get("X-Forwarded-Proto"), "https") ||
		strings.EqualFold(c.Protocol(), "https") ||
		strings.Contains(strings.ToLower(cfg.BaseURL), "https") ||
		!strings.Contains(strings.ToLower(c.Hostname()), "localhost")

	domain := ""
	host := strings.ToLower(c.Hostname())
	if strings.Contains(host, "torguisam.ru") || strings.Contains(strings.ToLower(cfg.BaseURL), "torguisam.ru") {
		domain = ".torguisam.ru"
	}

	buildCookie := func(name, value, dom string, maxAge int, httpOnly bool) string {
		exp := expires
		if maxAge < 0 {
			exp = time.Now().Add(-24 * time.Hour)
		}
		parts := []string{
			fmt.Sprintf("%s=%s", name, value),
			"Path=/",
			fmt.Sprintf("Max-Age=%d", maxAge),
			fmt.Sprintf("Expires=%s", exp.UTC().Format(http.TimeFormat)),
			"SameSite=Lax",
		}
		if dom != "" {
			parts = append(parts, "Domain="+dom)
		}
		if httpOnly {
			parts = append(parts, "HttpOnly")
		}
		if isHTTPS {
			parts = append(parts, "Secure")
		}
		return strings.Join(parts, "; ")
	}

	if domain != "" {
		// Удаляем возможные устаревшие host-only cookie (без Domain), чтобы они не перекрывали .torguisam.ru
		c.Response().Header.Add("Set-Cookie", buildCookie("session_id", "", "", -1, true))
		c.Response().Header.Add("Set-Cookie", buildCookie("ws_session_id", "", "", -1, false))
	}
	c.Response().Header.Add("Set-Cookie", buildCookie("session_id", sid, domain, maxAgeSec, true))
	c.Response().Header.Add("Set-Cookie", buildCookie("ws_session_id", sid, domain, maxAgeSec, false))
}

func clearSessionCookies(c *fiber.Ctx, cfg config.Config) {
	setSessionCookies(c, cfg, "", -1)
}

// RegisterAuthRoutes — пути как Nest AuthController.
func RegisterAuthRoutes(app fiber.Router, cfg config.Config, auth *service.AuthService) {
	g := app.Group("/auth")

	g.Post("/sign-up", func(c *fiber.Ctx) error {
		where := c.Query("where")
		var body struct {
			FullName    string `json:"fullName"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phoneNumber"`
			Password    string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		if err := auth.SignUp(c.UserContext(), where, body.FullName, body.Email, body.PhoneNumber, body.Password); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Код подтверждения успешно отправлен"})
	})

	g.Post("/verify-mobile-code", func(c *fiber.Ctx) error {
		code := c.Query("code")
		out, sid, err := auth.VerifyMobileCode(c.UserContext(), code)
		if err != nil {
			return writeAppError(c, err)
		}
		setSessionCookies(c, cfg, sid, 30*24*60*60)
		return c.JSON(out)
	})

	g.Post("/sign-in", func(c *fiber.Ctx) error {
		var body struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		out, sid, err := auth.SignIn(c.UserContext(), body.Login, body.Password)
		if err != nil {
			return writeAppError(c, err)
		}
		setSessionCookies(c, cfg, sid, 30*24*60*60)
		return c.JSON(out)
	})

	g.Get("/vk/url", func(c *fiber.Ctx) error {
		authURL, err := auth.VKAuthURL(c.Query("state"))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"url": authURL})
	})

	g.Get("/t-id/url", func(c *fiber.Ctx) error {
		authURL, err := auth.TIDAuthURL(c.Query("state"))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"url": authURL})
	})

	g.Post("/vk/sign-in", func(c *fiber.Ctx) error {
		var body struct {
			Code     string `json:"code"`
			State    string `json:"state"`
			DeviceID string `json:"device_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		out, sid, err := auth.SignInWithVK(c.UserContext(), body.Code, body.State, body.DeviceID)
		if err != nil {
			return writeAppError(c, err)
		}
		setSessionCookies(c, cfg, sid, 30*24*60*60)
		return c.JSON(out)
	})

	g.Post("/t-id/sign-in", func(c *fiber.Ctx) error {
		var body struct {
			Code  string `json:"code"`
			State string `json:"state"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		out, sid, err := auth.SignInWithTID(c.UserContext(), body.Code, body.State)
		if err != nil {
			return writeAppError(c, err)
		}
		setSessionCookies(c, cfg, sid, 30*24*60*60)
		return c.JSON(out)
	})

	g.Get("/yandex/url", func(c *fiber.Ctx) error {
		authURL, err := auth.YandexAuthURL(c.Query("state"))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"url": authURL})
	})

	g.Post("/yandex/sign-in", func(c *fiber.Ctx) error {
		var body struct {
			Code  string `json:"code"`
			State string `json:"state"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		res, err := auth.SignInWithYandex(c.UserContext(), body.Code, body.State)
		if err != nil {
			return writeAppError(c, err)
		}
		// Нового пользователя ещё нет: сессия не выдаётся, пока телефон не подтверждён по SMS.
		if res.RequirePhoneRegistration {
			return c.JSON(fiber.Map{
				"message":                  "Для завершения регистрации подтвердите номер телефона по SMS",
				"requirePhoneRegistration": true,
				"registrationTicket":       res.RegistrationTicket,
			})
		}
		setSessionCookies(c, cfg, res.SessionID, 30*24*60*60)
		return c.JSON(res.SignIn)
	})

	// Регистрация через Яндекс: шаги ниже идут без сессии — по тикету из /yandex/sign-in.
	g.Get("/yandex/register/status", func(c *fiber.Ctx) error {
		out, err := auth.YandexRegistrationStatus(c.UserContext(), c.Query("ticket"))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/yandex/register/send-code", func(c *fiber.Ctx) error {
		var body struct {
			RegistrationTicket string `json:"registrationTicket"`
			PhoneNumber        string `json:"phoneNumber"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		phone, err := auth.YandexRegistrationSendCode(c.UserContext(), body.RegistrationTicket, body.PhoneNumber)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Код подтверждения отправлен по SMS", "phoneNumber": phone})
	})

	g.Post("/yandex/register/verify-code", func(c *fiber.Ctx) error {
		var body struct {
			RegistrationTicket string `json:"registrationTicket"`
			Code               string `json:"code"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		code := body.Code
		if code == "" {
			code = c.Query("code")
		}
		out, sid, err := auth.YandexRegistrationVerifyCode(c.UserContext(), body.RegistrationTicket, code)
		if err != nil {
			return writeAppError(c, err)
		}
		setSessionCookies(c, cfg, sid, 30*24*60*60)
		return c.JSON(out)
	})

	g.Get("/yandex/onboarding/status", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
		out, err := auth.YandexOnboardingStatus(c.UserContext(), u.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/yandex/onboarding/start-phone", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
		var body struct {
			PhoneNumber string `json:"phoneNumber"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		if err := auth.YandexOnboardingStartPhone(c.UserContext(), u.ID, body.PhoneNumber); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Код подтверждения отправлен на телефон"})
	})

	g.Post("/yandex/onboarding/verify-phone", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
		code := c.Query("code")
		if err := auth.YandexOnboardingVerifyPhoneCode(c.UserContext(), u.ID, code); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Телефон успешно подтверждён"})
	})

	g.Get("/vk/onboarding/status", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
		out, err := auth.VKOnboardingStatus(c.UserContext(), u.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/vk/onboarding/start-email", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
		var body struct {
			Email string `json:"email"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		targetUserID, err := auth.VKOnboardingStartEmail(c.UserContext(), u.ID, body.Email)
		if err != nil {
			return writeAppError(c, err)
		}
		if targetUserID != u.ID {
			sid, err := auth.CreateSessionForUserID(c.UserContext(), targetUserID)
			if err != nil {
				return writeAppError(c, err)
			}
			setSessionCookies(c, cfg, sid, 30*24*60*60)
		}
		return c.JSON(fiber.Map{"message": "Код подтверждения отправлен на почту"})
	})

	g.Post("/vk/onboarding/verify-email", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
		code := c.Query("code")
		if err := auth.VKOnboardingVerifyEmailCode(c.UserContext(), u.ID, code); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Почта подтверждена"})
	})

	g.Post("/vk/onboarding/start-phone", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
		var body struct {
			PhoneNumber string `json:"phoneNumber"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		if err := auth.VKOnboardingStartPhone(c.UserContext(), u.ID, body.PhoneNumber); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Код подтверждения отправлен на телефон"})
	})

	g.Post("/vk/onboarding/verify-phone", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
		code := c.Query("code")
		if err := auth.VKOnboardingVerifyPhoneCode(c.UserContext(), u.ID, code); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Телефон подтвержден"})
	})

	g.Get("/me", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
		sid := middleware.GetSessionID(c)
		if sid != "" {
			setSessionCookies(c, cfg, sid, 30*24*60*60)
		}
		me, err := auth.Me(c.UserContext(), u.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(me)
	})

	g.Get("/isAdmin", middleware.RequireAdmin(auth), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"isAdmin": true})
	})

	g.Post("/logout", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		sid := c.Cookies("session_id")
		if err := auth.Logout(c.UserContext(), sid); err != nil {
			return writeAppError(c, err)
		}
		clearSessionCookies(c, cfg)
		return c.JSON(fiber.Map{"message": "Вы успешно вышли из системы!"})
	})

	g.Post("/forgot-password", func(c *fiber.Ctx) error {
		var body struct {
			Email       string `json:"email"`
			PhoneNumber string `json:"phoneNumber"`
			Where       string `json:"where"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		if err := auth.ForgotPasswordBy(c.UserContext(), body.Where, body.Email, body.PhoneNumber); err != nil {
			return writeAppError(c, err)
		}
		if body.Where == "sms" {
			return c.JSON(fiber.Map{"message": "Код восстановления отправлен на номер телефона"})
		}
		return c.JSON(fiber.Map{"message": "Письмо с кодом подтверждения отправлено на почту"})
	})

	g.Post("/verify-code", func(c *fiber.Ctx) error {
		code := c.Query("code")
		uid, err := auth.VerifyForgotCode(c.UserContext(), code)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"userId": uid})
	})

	g.Post("/change-password", func(c *fiber.Ctx) error {
		var body struct {
			UserID   int    `json:"userId"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		if err := auth.ChangePassword(c.UserContext(), int32(body.UserID), body.Password); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Пароль успешно изменён"})
	})
}
