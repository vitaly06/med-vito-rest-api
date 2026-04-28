package httpserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func RegisterReservationRoutes(app fiber.Router, reservations *service.ReservationService, auth *service.AuthService) {
	g := app.Group("/reservations")
	sess := authmw.RequireSession(auth)

	g.Post("/", sess, func(c *fiber.Ctx) error {
		var body service.CreateReservationRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		out, err := reservations.Create(c.UserContext(), me.ID, body)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	g.Get("/my", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := reservations.MyList(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/product/:productId", func(c *fiber.Ctx) error {
		productID, err := strconv.ParseInt(c.Params("productId"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный productId"})
		}
		out, err := reservations.ProductReservationInfo(c.UserContext(), int32(productID))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/:id/cancel-by-buyer", sess, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		me := authmw.UserFromLocals(c)
		out, err := reservations.CancelByBuyer(c.UserContext(), me.ID, id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/:id/cancel-by-seller", sess, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		var body struct {
			Reason string `json:"reason"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		out, err := reservations.CancelBySeller(c.UserContext(), me.ID, id, body.Reason)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/:id/extend", sess, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		me := authmw.UserFromLocals(c)
		out, err := reservations.Extend(c.UserContext(), me.ID, id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/product-settings", sess, func(c *fiber.Ctx) error {
		var body service.UpdateProductReservationSettingsRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		if err := reservations.UpdateProductSettings(c.UserContext(), me.ID, body); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Настройки резервирования обновлены"})
	})
}
