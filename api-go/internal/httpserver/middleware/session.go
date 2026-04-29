package middleware

import (
	"github.com/gofiber/fiber/v2"

	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/rbac"
	"med-vito/api-go/internal/service"
)

const UserLocalsKey = "user"

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
				"message":    "Ваш аккаунт заблокирован",
			})
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
}

func RequireRoleLevel(auth *service.AuthService, minLevel int, denyMessage string) fiber.Handler {
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
				"message":    "Ваш аккаунт заблокирован",
			})
		}
		if !rbac.HasMinRole(u.RoleName, minLevel) {
			if denyMessage == "" {
				denyMessage = "Доступ запрещен"
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    denyMessage,
			})
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
}

func RequireAdmin(auth *service.AuthService) fiber.Handler {
	return RequireRoleLevel(auth, 90, "Доступ разрешен только для администраторов")
}

func RequireModerator(auth *service.AuthService) fiber.Handler {
	return RequireRoleLevel(auth, 70, "Доступ разрешен только для модераторов и выше")
}

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

func UserFromLocals(c *fiber.Ctx) *domain.UserEntity {
	v := c.Locals(UserLocalsKey)
	if v == nil {
		return nil
	}
	u, _ := v.(*domain.UserEntity)
	return u
}

