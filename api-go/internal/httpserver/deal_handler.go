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
	adm := authmw.RequireAdmin(auth)

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
		out, err := deals.GetDealCDEKQR(c.UserContext(), me.ID, id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
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

	g.Post("/:id/sync-payment", sess, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.SyncDealPayment(c.UserContext(), me.ID, id)
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

	app.Get("/admin/deals/list", adm, func(c *fiber.Ctx) error {
		out, err := deals.AdminListDeals(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	app.Get("/admin/deals/:id", adm, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		out, err := deals.AdminGetDeal(c.UserContext(), id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	app.Patch("/admin/deals/:id/status", adm, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		var body struct {
			Status string `json:"status"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		out, err := deals.AdminSetStatus(c.UserContext(), me.ID, id, body.Status)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	app.Get("/admin/deals/:id/logs", adm, func(c *fiber.Ctx) error {
		id, err := parseDealID(c)
		if err != nil {
			return err
		}
		out, err := deals.AdminDealLogs(c.UserContext(), id)
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
