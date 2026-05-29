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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Р РЋР ВµРЎРѓРЎРѓР С‘РЎРЏ Р Р…Р Вµ Р Р…Р В°Р в„–Р Т‘Р ВµР Р…Р В°"})
		}
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"statusCode": fiber.StatusInternalServerError, "message": err.Error()})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Р РЋР ВµРЎРѓРЎРѓР С‘РЎРЏ Р Р…Р ВµР Т‘Р ВµР в„–РЎРѓРЎвЂљР Р†Р С‘РЎвЂљР ВµР В»РЎРЉР Р…Р В° Р С‘Р В»Р С‘ Р С‘РЎРѓРЎвЂљР ВµР С”Р В»Р В°"})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": "Р вЂ™Р В°РЎв‚¬ Р В°Р С”Р С”Р В°РЎС“Р Р…РЎвЂљ Р В·Р В°Р В±Р В»Р С•Р С”Р С‘РЎР‚Р С•Р Р†Р В°Р Р…"})
		}
		if auth.IsVKOnboardingRequiredForUser(c.UserContext(), u) && !allowDuringVKOnboarding(c.Path()) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "РўСЂРµР±СѓРµС‚СЃСЏ Р·Р°РІРµСЂС€РёС‚СЊ РїСЂРёРІСЏР·РєСѓ Рё РїРѕРґС‚РІРµСЂР¶РґРµРЅРёРµ email Рё С‚РµР»РµС„РѕРЅР°",
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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Р РЋР ВµРЎРѓРЎРѓР С‘РЎРЏ Р Р…Р Вµ Р Р…Р В°Р в„–Р Т‘Р ВµР Р…Р В°"})
		}
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"statusCode": fiber.StatusInternalServerError, "message": err.Error()})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Р РЋР ВµРЎРѓРЎРѓР С‘РЎРЏ Р Р…Р ВµР Т‘Р ВµР в„–РЎРѓРЎвЂљР Р†Р С‘РЎвЂљР ВµР В»РЎРЉР Р…Р В° Р С‘Р В»Р С‘ Р С‘РЎРѓРЎвЂљР ВµР С”Р В»Р В°"})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": "Р вЂ™Р В°РЎв‚¬ Р В°Р С”Р С”Р В°РЎС“Р Р…РЎвЂљ Р В·Р В°Р В±Р В»Р С•Р С”Р С‘РЎР‚Р С•Р Р†Р В°Р Р…"})
		}
		if !rbac.HasMinRole(u.RoleName, minLevel) {
			if denyMessage == "" {
				denyMessage = "Р вЂќР С•РЎРѓРЎвЂљРЎС“Р С— Р В·Р В°Р С—РЎР‚Р ВµРЎвЂ°Р ВµР Р…"
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": denyMessage})
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
}

func RequireAdmin(auth *service.AuthService) fiber.Handler {
	return RequireRoleLevel(auth, 90, "Р вЂќР С•РЎРѓРЎвЂљРЎС“Р С— РЎР‚Р В°Р В·РЎР‚Р ВµРЎв‚¬Р ВµР Р… РЎвЂљР С•Р В»РЎРЉР С”Р С• Р Т‘Р В»РЎРЏ Р В°Р Т‘Р СР С‘Р Р…Р С‘РЎРѓРЎвЂљРЎР‚Р В°РЎвЂљР С•РЎР‚Р С•Р Р†")
}

func RequireModerator(auth *service.AuthService) fiber.Handler {
	return RequireRoleLevel(auth, 70, "Р вЂќР С•РЎРѓРЎвЂљРЎС“Р С— РЎР‚Р В°Р В·РЎР‚Р ВµРЎв‚¬Р ВµР Р… РЎвЂљР С•Р В»РЎРЉР С”Р С• Р Т‘Р В»РЎРЏ Р СР С•Р Т‘Р ВµРЎР‚Р В°РЎвЂљР С•РЎР‚Р С•Р Р† Р С‘ Р Р†РЎвЂ№РЎв‚¬Р Вµ")
}

func RequirePermission(auth *service.AuthService, permission, denyMessage string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sid := c.Cookies("session_id")
		if sid == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Р РЋР ВµРЎРѓРЎРѓР С‘РЎРЏ Р Р…Р Вµ Р Р…Р В°Р в„–Р Т‘Р ВµР Р…Р В°"})
		}
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"statusCode": fiber.StatusInternalServerError, "message": err.Error()})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Р РЋР ВµРЎРѓРЎРѓР С‘РЎРЏ Р Р…Р ВµР Т‘Р ВµР в„–РЎРѓРЎвЂљР Р†Р С‘РЎвЂљР ВµР В»РЎРЉР Р…Р В° Р С‘Р В»Р С‘ Р С‘РЎРѓРЎвЂљР ВµР С”Р В»Р В°"})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": "Р вЂ™Р В°РЎв‚¬ Р В°Р С”Р С”Р В°РЎС“Р Р…РЎвЂљ Р В·Р В°Р В±Р В»Р С•Р С”Р С‘РЎР‚Р С•Р Р†Р В°Р Р…"})
		}
		if !rbac.HasPermission(u.RoleName, permission) {
			if denyMessage == "" {
				denyMessage = "Р СњР ВµР Т‘Р С•РЎРѓРЎвЂљР В°РЎвЂљР С•РЎвЂЎР Р…Р С• Р С—РЎР‚Р В°Р Р†"
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
