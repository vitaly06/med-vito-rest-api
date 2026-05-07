package httpserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func RegisterDealRoutes(app fiber.Router, deals *service.DealService, auth *service.AuthService) {
	g := app.Group("/deals")
	sess := authmw.RequireSession(auth)

	g.Post("/", sess, func(c *fiber.Ctx) error {
		var body service.CreateDealRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.CreateDeal(c.UserContext(), me.ID, body)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	g.Get("/my-purchases", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := deals.MyPurchases(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/my-sales", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := deals.MySales(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/my", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := deals.MyAllDeals(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/:id", sess, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.GetDeal(c.UserContext(), me.ID, id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/:id/cdek-qr", sess, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.GetDeal(c.UserContext(), me.ID, id)
		if err != nil {
			return writeAppError(c, err)
		}
		cdekMap, _ := out["cdek"].(map[string]any)
		if cdekMap == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"statusCode": 404, "message": "Данные CDEK не найдены"})
		}
		qrCodeData, _ := cdekMap["qrCodeData"].(*string)
		if qrCodeData == nil || *qrCodeData == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"statusCode": 404, "message": "QR-код еще не сформирован"})
		}
		return c.JSON(fiber.Map{
			"qrCodeData":  *qrCodeData,
			"trackNumber": cdekMap["trackNumber"],
			"trackingUrl": cdekMap["trackingUrl"],
			"orderUuid":   cdekMap["orderUuid"],
		})
	})

	g.Post("/:id/pay", sess, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.PayDeal(c.UserContext(), me.ID, id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/:id/mark-shipped", sess, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		var body service.MarkShippedRequest
		if len(c.Body()) > 0 {
			if err := c.BodyParser(&body); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
			}
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.MarkShipped(c.UserContext(), me.ID, id, body)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/:id/confirm-delivery", sess, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.ConfirmDelivery(c.UserContext(), me.ID, id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/:id/open-dispute", sess, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		var body struct {
			Reason string `json:"reason"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.OpenDispute(c.UserContext(), me.ID, id, body.Reason)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/:id/cancel", sess, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.CancelDeal(c.UserContext(), me.ID, id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}

func parseDealID(c *fiber.Ctx) (int32, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return 0, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
	}
	return int32(id), nil
}
