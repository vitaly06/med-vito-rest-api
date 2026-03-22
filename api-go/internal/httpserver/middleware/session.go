package middleware

import (
	"github.com/gofiber/fiber/v2"

	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/service"
)

const UserLocalsKey = "user"

// RequireSession — как SessionAuthGuard: кука session_id, пользователь в Locals.
func RequireSession(auth *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sid := c.Cookies("session_id")
		if sid == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": fiber.StatusUnauthorized,
				"message":    "Сессия не найдена",
			})
		}
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": fiber.StatusUnauthorized,
				"message":    "Сессия недействительна или истекла",
			})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Ваш аккаунт заблокирован. Пожалуйста, обратитесь в поддержку",
			})
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
}

// RequireAdmin — как AdminSessionAuthGuard: роль admin.
func RequireAdmin(auth *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sid := c.Cookies("session_id")
		if sid == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": fiber.StatusUnauthorized,
				"message":    "Сессия не найдена",
			})
		}
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": fiber.StatusUnauthorized,
				"message":    "Сессия недействительна или истекла",
			})
		}
		if u.RoleName == nil || *u.RoleName != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Доступ запрещён",
			})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Ваш аккаунт заблокирован. Пожалуйста, обратитесь в поддержку",
			})
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
}

// OptionalSession — как OptionalSessionAuthGuard.
func OptionalSession(auth *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sid := c.Cookies("session_id")
		if sid == "" {
			return c.Next()
		}
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil || u == nil {
			return c.Next()
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
}

// UserFromLocals — после RequireSession / RequireAdmin.
func UserFromLocals(c *fiber.Ctx) *domain.UserEntity {
	v := c.Locals(UserLocalsKey)
	if v == nil {
		return nil
	}
	u, _ := v.(*domain.UserEntity)
	return u
}
