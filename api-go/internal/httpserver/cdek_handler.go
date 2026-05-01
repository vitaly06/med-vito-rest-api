package httpserver

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"med-vito/api-go/internal/service"
)

func RegisterCDEKRoutes(app fiber.Router, cdek *service.CDEKService) {
	g := app.Group("/cdek")

	g.Get("/cities", func(c *fiber.Ctx) error {
		limit := 20
		if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil {
				limit = n
			}
		}
		out, err := cdek.Cities(c.UserContext(), c.Query("city"), limit)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/delivery-points", func(c *fiber.Ctx) error {
		cityCode, err := strconv.Atoi(c.Query("cityCode"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќСѓР¶РµРЅ cityCode"})
		}
		out, err := cdek.DeliveryPoints(c.UserContext(), cityCode)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/calculate", func(c *fiber.Ctx) error {
		var body service.CDEKCalculateRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅРѕРµ С‚РµР»Рѕ"})
		}
		out, err := cdek.Calculate(c.UserContext(), body)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}

