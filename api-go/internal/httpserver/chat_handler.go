package httpserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

// RegisterChatRoutes — как Nest ChatController; :id/messages раньше :id, чтобы Fiber не съел "messages" как id.
func RegisterChatRoutes(app fiber.Router, chat *service.ChatService, auth *service.AuthService) {
	g := app.Group("/chat")
	sess := authmw.RequireSession(auth)

	g.Post("/start", sess, func(c *fiber.Ctx) error {
		var body struct {
			ProductID int32 `json:"productId"`
		}
		if err := c.BodyParser(&body); err != nil || body.ProductID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело (нужен productId)"})
		}
		me := authmw.UserFromLocals(c)
		out, err := chat.StartChat(c.UserContext(), body.ProductID, me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	g.Get("", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := chat.GetUserChats(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/:id/messages", sess, func(c *fiber.Ctx) error {
		chatID, err := c.ParamsInt("id")
		if err != nil || chatID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id чата"})
		}
		page := 1
		if p := c.Query("page"); p != "" {
			if n, e := strconv.Atoi(p); e == nil && n > 0 {
				page = n
			}
		}
		limit := 50
		if l := c.Query("limit"); l != "" {
			if n, e := strconv.Atoi(l); e == nil && n > 0 {
				limit = n
			}
		}
		me := authmw.UserFromLocals(c)
		out, err := chat.GetChatMessages(c.UserContext(), int32(chatID), me.ID, page, limit)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/:id", sess, func(c *fiber.Ctx) error {
		chatID, err := c.ParamsInt("id")
		if err != nil || chatID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id чата"})
		}
		me := authmw.UserFromLocals(c)
		out, err := chat.GetChatInfo(c.UserContext(), int32(chatID), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
