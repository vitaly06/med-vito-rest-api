package httpserver

import (
	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

// RegisterPaymentRoutes — как Nest PaymentController (/payment).
func RegisterPaymentRoutes(app fiber.Router, pay *service.PaymentService, auth *service.AuthService) {
	g := app.Group("/payment")
	sess := authmw.RequireSession(auth)

	g.Post("/create", sess, func(c *fiber.Ctx) error {
		var body struct {
			Amount      float64 `json:"amount"`
			Description *string `json:"description"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		desc := ""
		if body.Description != nil {
			desc = *body.Description
		}
		me := authmw.UserFromLocals(c)
		out, err := pay.CreatePayment(c.UserContext(), me.ID, body.Amount, desc)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	// Webhook Т-Банк — без сессии; тело как есть для проверки подписи.
	g.Post("/notification", func(c *fiber.Ctx) error {
		out, err := pay.HandleNotification(c.UserContext(), c.Body())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(out)
	})

	g.Get("/history", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := pay.GetUserPayments(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/check-status", sess, func(c *fiber.Ctx) error {
		var body struct {
			PaymentID string `json:"paymentId"`
		}
		if err := c.BodyParser(&body); err != nil || body.PaymentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Нужен paymentId"})
		}
		out, err := pay.CheckPaymentStatus(c.UserContext(), body.PaymentID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
