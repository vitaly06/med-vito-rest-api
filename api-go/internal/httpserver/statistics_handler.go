package httpserver

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

// RegisterStatisticsRoutes — как Nest StatisticsController (сессия).
func RegisterStatisticsRoutes(app fiber.Router, stat *service.StatisticsService, auth *service.AuthService) {
	g := app.Group("/statistics")
	sess := authmw.RequireSession(auth)

	g.Get("/analytic", sess, func(c *fiber.Ctx) error {
		var categoryID *int32
		if v := strings.TrimSpace(c.Query("categoryId")); v != "" {
			n, err := strconv.ParseInt(v, 10, 32)
			if err != nil || n < 1 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный categoryId"})
			}
			x := int32(n)
			categoryID = &x
		}
		var productID *int32
		if v := strings.TrimSpace(c.Query("productId")); v != "" {
			n, err := strconv.ParseInt(v, 10, 32)
			if err != nil || n < 1 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный productId"})
			}
			x := int32(n)
			productID = &x
		}
		var region *string
		if v := strings.TrimSpace(c.Query("region")); v != "" {
			region = &v
		}
		period := c.Query("period")
		me := authmw.UserFromLocals(c)
		out, err := stat.GetUserStatistics(c.UserContext(), me.ID, period, categoryID, region, productID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/products-analytic", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := stat.GetProductsAnalytic(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
