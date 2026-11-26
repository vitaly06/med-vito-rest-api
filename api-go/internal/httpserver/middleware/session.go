package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/rbac"
	"med-vito/api-go/internal/service"
)

const UserLocalsKey = "user"

func GetSessionIDs(c *fiber.Ctx) []string {
	var ids []string
	seen := make(map[string]struct{})
	add := func(v string) {
		v = strings.TrimSpace(v)
		if v == "" {
			return
		}
		if _, ok := seen[v]; ok {
			return
		}
		seen[v] = struct{}{}
		ids = append(ids, v)
	}

	// Парсим сырой заголовок Cookie, чтобы при наличии дублей cookie (host-only + .torguisam.ru) проверить все
	if rawCookie := c.Get("Cookie"); rawCookie != "" {
		for _, part := range strings.Split(rawCookie, ";") {
			kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
			if len(kv) == 2 {
				k := strings.TrimSpace(kv[0])
				if k == "session_id" || k == "ws_session_id" {
					add(kv[1])
				}
			}
		}
	}
	add(c.Cookies("session_id"))
	add(c.Cookies("ws_session_id"))
	authHdr := strings.TrimSpace(c.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(authHdr), "bearer ") {
		add(authHdr[7:])
	}
	return ids
}

func GetSessionID(c *fiber.Ctx) string {
	ids := GetSessionIDs(c)
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}

func resolveSessionUser(c *fiber.Ctx, auth *service.AuthService) (*domain.UserEntity, string, error) {
	ids := GetSessionIDs(c)
	if len(ids) == 0 {
		return nil, "", nil
	}
	var lastErr error
	for _, sid := range ids {
		u, err := auth.UserFromSession(c.UserContext(), sid)
		if err != nil {
			lastErr = err
			continue
		}
		if u != nil {
			return u, sid, nil
		}
	}
	return nil, ids[0], lastErr
}

func RequireSession(auth *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ids := GetSessionIDs(c)
		if len(ids) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Сессия не найдена"})
		}
		u, sid, err := resolveSessionUser(c, auth)
		if err != nil && u == nil {
			if strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no session") {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Сессия недействительна или истекла"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"statusCode": fiber.StatusInternalServerError, "message": err.Error()})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Сессия недействительна или истекла"})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": "Ваш аккаунт заблокирован", "code": "USER_BANNED"})
		}
		if auth.SessionRequiresVKOnboarding(c.UserContext(), sid, u) && !allowDuringVKOnboarding(c.Path()) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Требуется завершить привязку и подтверждение email и телефона",
				"code":       "VK_ONBOARDING_REQUIRED",
			})
		}
		if auth.SessionRequiresYandexOnboarding(c.UserContext(), sid, u) && !allowDuringYandexOnboarding(c.Path()) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Требуется подтвердить номер телефона по SMS",
				"code":       "YANDEX_ONBOARDING_REQUIRED",
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

func allowDuringYandexOnboarding(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	if path == "/auth/me" || path == "/auth/logout" {
		return true
	}
	return strings.HasPrefix(path, "/auth/yandex/onboarding/")
}

func RequireRoleLevel(auth *service.AuthService, minLevel int, denyMessage string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ids := GetSessionIDs(c)
		if len(ids) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Сессия не найдена"})
		}
		u, _, err := resolveSessionUser(c, auth)
		if err != nil && u == nil {
			if strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no session") {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Сессия недействительна или истекла"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"statusCode": fiber.StatusInternalServerError, "message": err.Error()})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Сессия недействительна или истекла"})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": "Ваш аккаунт заблокирован", "code": "USER_BANNED"})
		}
		if !rbac.HasMinRole(u.RoleName, minLevel) {
			if denyMessage == "" {
				denyMessage = "Доступ запрещен"
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": denyMessage})
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

func RequirePermission(auth *service.AuthService, permission, denyMessage string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ids := GetSessionIDs(c)
		if len(ids) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Сессия не найдена"})
		}
		u, _, err := resolveSessionUser(c, auth)
		if err != nil && u == nil {
			if strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no session") {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Сессия недействительна или истекла"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"statusCode": fiber.StatusInternalServerError, "message": err.Error()})
		}
		if u == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"statusCode": fiber.StatusUnauthorized, "message": "Сессия недействительна или истекла"})
		}
		if u.IsBanned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": "Ваш аккаунт заблокирован", "code": "USER_BANNED"})
		}
		if !rbac.HasPermission(u.RoleName, permission) {
			if denyMessage == "" {
				denyMessage = "Недостаточно прав"
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"statusCode": fiber.StatusForbidden, "message": denyMessage})
		}
		c.Locals(UserLocalsKey, u)
		return c.Next()
	}
}

func OptionalSession(auth *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u, _, err := resolveSessionUser(c, auth)
		if err != nil || u == nil || u.IsBanned {
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
