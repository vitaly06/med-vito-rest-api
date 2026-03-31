package httpserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func RegisterKnowledgeBaseRoutes(app fiber.Router, kb *service.KnowledgeBaseService, auth *service.AuthService) {
	g := app.Group("/knowledge-base")

	g.Get("/", func(c *fiber.Ctx) error {
		out, err := kb.ListArticles(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректный id",
			})
		}
		out, err := kb.GetArticleByID(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	adm := g.Group("", authmw.RequireAdmin(auth))

	adm.Post("/", func(c *fiber.Ctx) error {
		var body struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректное тело запроса",
			})
		}
		out, err := kb.CreateArticle(c.UserContext(), body.Title, body.Content)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	adm.Put("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректный id",
			})
		}
		var body struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректное тело запроса",
			})
		}
		out, err := kb.UpdateArticle(c.UserContext(), int32(id), body.Title, body.Content)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	adm.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректный id",
			})
		}
		out, err := kb.DeleteArticle(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
