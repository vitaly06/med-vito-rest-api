package httpserver

import (
	"github.com/gofiber/fiber/v2"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func sessionCookie(cfg config.Config, value string, maxAgeSec int) *fiber.Cookie {
	ss := "Lax"
	if cfg.Production {
		ss = "Strict"
	}
	return &fiber.Cookie{
		Name:     "session_id",
		Value:    value,
		Path:     "/",
		HTTPOnly: true,
		Secure:   cfg.Production,
		SameSite: ss,
		MaxAge:   maxAgeSec,
	}
}

func clearSessionCookie(cfg config.Config) *fiber.Cookie {
	c := sessionCookie(cfg, "", -1)
	c.MaxAge = -1
	return c
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
		if err := auth.VerifyMobileCode(c.UserContext(), code); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Вы успешно зарегистрировались!"})
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
		c.Cookie(sessionCookie(cfg, sid, 30*24*60*60))
		return c.JSON(out)
	})

	g.Get("/me", middleware.RequireSession(auth), func(c *fiber.Ctx) error {
		u := middleware.UserFromLocals(c)
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
		c.Cookie(clearSessionCookie(cfg))
		return c.JSON(fiber.Map{"message": "Вы успешно вышли из системы!"})
	})

	g.Post("/forgot-password", func(c *fiber.Ctx) error {
		var body struct {
			Email string `json:"email"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		if err := auth.ForgotPassword(c.UserContext(), body.Email); err != nil {
			return writeAppError(c, err)
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
