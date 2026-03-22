package httpserver

import "github.com/gofiber/fiber/v2"

// HealthHandler — аналог простого публичного GET без БД (пока нет репозитория).
func HealthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"service": "med-vito-api-go",
	})
}
