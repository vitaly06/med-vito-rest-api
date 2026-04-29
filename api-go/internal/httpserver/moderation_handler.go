package httpserver

import (
	"strconv"
	"strings"

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

	g.Get("/summary", adm, func(c *fiber.Ctx) error {
		days := 30
		if v := strings.TrimSpace(c.Query("days")); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				days = n
			}
		}
		out, err := moderation.Summary(c.UserContext(), days)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/audit-logs", adm, func(c *fiber.Ctx) error {
		limit := 100
		if v := strings.TrimSpace(c.Query("limit")); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				limit = n
			}
		}
		out, err := moderation.AuditLogs(c.UserContext(), limit)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/appeals", adm, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := moderation.Appeals(c.UserContext(), me, false)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/appeals/:id/review", adm, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil || id < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		var body struct {
			Status  string  `json:"status"`
			Comment *string `json:"comment"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		if err := moderation.ReviewAppeal(c.UserContext(), me, id, body.Status, body.Comment); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Апелляция рассмотрена"})
	})

	pub := app.Group("/moderation")
	sess := authmw.RequireSession(auth)
	pub.Post("/appeals", sess, func(c *fiber.Ctx) error {
		var body struct {
			ProductID int32  `json:"productId"`
			Reason    string `json:"reason"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		if err := moderation.AddAppeal(c.UserContext(), me.ID, body.ProductID, body.Reason); err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Апелляция отправлена"})
	})
	pub.Get("/appeals/my", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := moderation.Appeals(c.UserContext(), me, true)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
