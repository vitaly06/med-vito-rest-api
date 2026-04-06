package httpserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func RegisterModerationRoutes(app fiber.Router, moderation *service.ModerationAdminService, auth *service.AuthService) {
	g := app.Group("/admin/moderation")
	adm := authmw.RequireAdmin(auth)

	g.Get("/products", adm, func(c *fiber.Ctx) error {
		page := 1
		if raw := c.Query("page"); raw != "" {
			if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
				page = parsed
			}
		}
		out, err := moderation.ListProducts(c.UserContext(), c.Query("filter"), page)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/products/:id", adm, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := moderation.GetProduct(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
