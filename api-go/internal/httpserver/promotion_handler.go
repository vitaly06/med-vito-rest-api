package httpserver

import (
	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

// RegisterPromotionRoutes — как Nest PromotionController.
func RegisterPromotionRoutes(app fiber.Router, promo *service.PromotionService, auth *service.AuthService) {
	g := app.Group("/promotion")

	g.Get("/all-promotions", func(c *fiber.Ctx) error {
		out, err := promo.AllPromotions(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	sess := authmw.RequireSession(auth)
	g.Post("/add-promotion", sess, func(c *fiber.Ctx) error {
		var body struct {
			ProductID   int32 `json:"productId"`
			PromotionID int32 `json:"promotionId"`
			Days        int32 `json:"days"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		if body.ProductID < 1 || body.PromotionID < 1 || body.Days < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "productId, promotionId и days должны быть положительными"})
		}
		me := authmw.UserFromLocals(c)
		out, err := promo.AddPromotion(c.UserContext(), me.ID, body.ProductID, body.PromotionID, body.Days)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})
}
