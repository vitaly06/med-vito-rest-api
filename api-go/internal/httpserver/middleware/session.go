package middleware

import (
	"strings"

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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "РЎРµСЃСЃРёСЏ РЅРµ РЅР°Р№РґРµРЅР°"})
		}
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"statusCode": fiber.StatusInternalServerError, "message": err.Error()})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "РЎРµСЃСЃРёСЏ РЅРµРґРµР№СЃС‚РІРёС‚РµР»СЊРЅР° РёР»Рё РёСЃС‚РµРєР»Р°"})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": "Р’Р°С€ Р°РєРєР°СѓРЅС‚ Р·Р°Р±Р»РѕРєРёСЂРѕРІР°РЅ"})
		}
		if service.IsVKOnboardingRequired(u) && !allowDuringVKOnboarding(c.Path()) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Требуется завершить привязку и подтверждение email и телефона",
				"code":       "VK_ONBOARDING_REQUIRED",
			})
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
}

func allowDuringVKOnboarding(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	if path == "/auth/me" || path == "/auth/logout" {
		return true
	}
	return strings.HasPrefix(path, "/auth/vk/onboarding/")
}

func RequireRoleLevel(auth *service.AuthService, minLevel int, denyMessage string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sid := c.Cookies("session_id")
		if sid == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "РЎРµСЃСЃРёСЏ РЅРµ РЅР°Р№РґРµРЅР°"})
		}
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"statusCode": fiber.StatusInternalServerError, "message": err.Error()})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "РЎРµСЃСЃРёСЏ РЅРµРґРµР№СЃС‚РІРёС‚РµР»СЊРЅР° РёР»Рё РёСЃС‚РµРєР»Р°"})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": "Р’Р°С€ Р°РєРєР°СѓРЅС‚ Р·Р°Р±Р»РѕРєРёСЂРѕРІР°РЅ"})
		}
		if !rbac.HasMinRole(u.RoleName, minLevel) {
			if denyMessage == "" {
				denyMessage = "Р”РѕСЃС‚СѓРї Р·Р°РїСЂРµС‰РµРЅ"
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": denyMessage})
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
}

func RequireAdmin(auth *service.AuthService) fiber.Handler {
	return RequireRoleLevel(auth, 90, "Р”РѕСЃС‚СѓРї СЂР°Р·СЂРµС€РµРЅ С‚РѕР»СЊРєРѕ РґР»СЏ Р°РґРјРёРЅРёСЃС‚СЂР°С‚РѕСЂРѕРІ")
}

func RequireModerator(auth *service.AuthService) fiber.Handler {
	return RequireRoleLevel(auth, 70, "Р”РѕСЃС‚СѓРї СЂР°Р·СЂРµС€РµРЅ С‚РѕР»СЊРєРѕ РґР»СЏ РјРѕРґРµСЂР°С‚РѕСЂРѕРІ Рё РІС‹С€Рµ")
}

func RequirePermission(auth *service.AuthService, permission, denyMessage string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sid := c.Cookies("session_id")
		if sid == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "РЎРµСЃСЃРёСЏ РЅРµ РЅР°Р№РґРµРЅР°"})
		}
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"statusCode": fiber.StatusInternalServerError, "message": err.Error()})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "РЎРµСЃСЃРёСЏ РЅРµРґРµР№СЃС‚РІРёС‚РµР»СЊРЅР° РёР»Рё РёСЃС‚РµРєР»Р°"})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": "Р’Р°С€ Р°РєРєР°СѓРЅС‚ Р·Р°Р±Р»РѕРєРёСЂРѕРІР°РЅ"})
		}
		if !rbac.HasPermission(u.RoleName, permission) {
			if denyMessage == "" {
				denyMessage = "РќРµРґРѕСЃС‚Р°С‚РѕС‡РЅРѕ РїСЂР°РІ"
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": denyMessage})
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
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
