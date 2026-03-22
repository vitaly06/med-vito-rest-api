package httpserver

import (
	"github.com/gofiber/fiber/v2"

	"med-vito/api-go/internal/service"
)

// RegisterLogRoutes — аналог LogController Nest: GET /log/find-all
func RegisterLogRoutes(app fiber.Router, logSvc *service.LogService) {
	app.Get("/log/find-all", func(c *fiber.Ctx) error {
		items, err := logSvc.FindAll(c.UserContext())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}
		return c.JSON(items)
	})
}
