package httpserver

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/service"
)

func supportIsModerator(u *domain.UserEntity) bool {
	if u == nil || u.RoleName == nil {
		return false
	}
	r := *u.RoleName
	return r == "moderator" || r == "admin"
}

func supportIsAdmin(u *domain.UserEntity) bool {
	return u != nil && u.RoleName != nil && *u.RoleName == "admin"
}

func parseSupportPageLimit(c *fiber.Ctx) (page, limit int) {
	page, limit = 1, 10
	if v := strings.TrimSpace(c.Query("page")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := strings.TrimSpace(c.Query("limit")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	return page, limit
}

func queryStringPtr(c *fiber.Ctx, key string) *string {
	v := strings.TrimSpace(c.Query(key))
	if v == "" {
		return nil
	}
	return &v
}

// RegisterSupportRoutes — как Nest SupportController; порядок маршрутов: my, all, stats, затем :id/...
func RegisterSupportRoutes(app fiber.Router, sup *service.SupportService, auth *service.AuthService) {
	g := app.Group("/support")
	sess := authmw.RequireSession(auth)

	g.Post("/tickets", sess, func(c *fiber.Ctx) error {
		var body struct {
			Theme     string  `json:"theme"`
			Subject   string  `json:"subject"`
			Message   string  `json:"message"`
			Priority  *string `json:"priority"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		pr := ""
		if body.Priority != nil {
			pr = *body.Priority
		}
		me := authmw.UserFromLocals(c)
		out, err := sup.CreateTicket(c.UserContext(), me.ID, body.Theme, body.Subject, body.Message, pr)
		if err != nil {
			return writeAppError(c, err)
		}
		NotifySupportNewTicket(out)
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	g.Get("/tickets/my", sess, func(c *fiber.Ctx) error {
		page, limit := parseSupportPageLimit(c)
		me := authmw.UserFromLocals(c)
		out, err := sup.GetUserTickets(c.UserContext(), me.ID, page, limit, queryStringPtr(c, "status"), queryStringPtr(c, "priority"), queryStringPtr(c, "theme"), queryStringPtr(c, "search"))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/tickets/all", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		if !supportIsModerator(me) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Доступ запрещен. Требуется роль модератора или администратора",
			})
		}
		page, limit := parseSupportPageLimit(c)
		out, err := sup.GetAllTickets(c.UserContext(), page, limit, queryStringPtr(c, "status"), queryStringPtr(c, "priority"), queryStringPtr(c, "theme"), queryStringPtr(c, "search"))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/stats", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		if !supportIsAdmin(me) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Доступ запрещен. Требуется роль администратора",
			})
		}
		out, err := sup.TicketStats(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/tickets/:id/messages", sess, func(c *fiber.Ctx) error {
		tid, err := c.ParamsInt("id")
		if err != nil || tid < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		var body struct {
			Text string `json:"text"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		out, err := sup.SendMessage(c.UserContext(), int32(tid), me.ID, body.Text, supportIsModerator(me))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/tickets/:id/assign", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		if !supportIsModerator(me) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Доступ запрещен. Требуется роль модератора или администратора",
			})
		}
		tid, err := c.ParamsInt("id")
		if err != nil || tid < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := sup.AssignTicket(c.UserContext(), int32(tid), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		NotifySupportTicketAssigned(int32(tid), me.ID)
		return c.JSON(out)
	})

	g.Put("/tickets/:id", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		if !supportIsModerator(me) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": fiber.StatusForbidden,
				"message":    "Доступ запрещен. Требуется роль модератора или администратора",
			})
		}
		tid, err := c.ParamsInt("id")
		if err != nil || tid < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		var body struct {
			Status   *string `json:"status"`
			Priority *string `json:"priority"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		out, err := sup.UpdateTicket(c.UserContext(), int32(tid), me.ID, body.Status, body.Priority)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/tickets/:id", sess, func(c *fiber.Ctx) error {
		tid, err := c.ParamsInt("id")
		if err != nil || tid < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		me := authmw.UserFromLocals(c)
		out, err := sup.GetTicket(c.UserContext(), int32(tid), me.ID, supportIsModerator(me))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
