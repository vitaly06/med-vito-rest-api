package httpserver

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"med-vito/api-go/internal/service"
)

// RegisterAddressRoutes — как Nest AddressController (без сессии).
func RegisterAddressRoutes(app fiber.Router, addr *service.AddressService) {
	g := app.Group("/address")

	g.Get("/suggestions", func(c *fiber.Ctx) error {
		q := strings.TrimSpace(c.Query("query"))
		if q == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Параметр query обязателен"})
		}
		limit := 5
		if v := strings.TrimSpace(c.Query("limit")); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				limit = n
			}
		}
		out, err := addr.GetSuggestions(c.UserContext(), q, limit)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/validate", func(c *fiber.Ctx) error {
		var body struct {
			Address        string `json:"address"`
			AddressDetails any    `json:"addressDetails"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		if strings.TrimSpace(body.Address) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Поле address обязательно"})
		}
		if !addr.ValidateAddress(c.UserContext(), body.Address) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": 400,
				"message":    "Указанный адрес не найден. Пожалуйста, выберите адрес из предложенных вариантов.",
			})
		}
		return c.JSON(fiber.Map{"valid": true, "message": "Адрес корректен"})
	})
}
