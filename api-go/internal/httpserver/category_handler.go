package httpserver

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func writeAppError(c *fiber.Ctx, err error) error {
	var ae *service.AppError
	if errors.As(err, &ae) {
		return c.Status(ae.Status).JSON(fiber.Map{
			"statusCode": ae.Status,
			"message":    ae.Message,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"statusCode": fiber.StatusInternalServerError,
		"message":    err.Error(),
	})
}

// RegisterCategoryRoutes — те же пути что у Nest CategoryController (админ — сессия + роль admin).
func RegisterCategoryRoutes(app fiber.Router, cat *service.CategoryService, auth *service.AuthService) {
	g := app.Group("/category")

	g.Get("/find-all", func(c *fiber.Ctx) error {
		out, err := cat.FindAllCategories(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/find-by-id/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректный id",
			})
		}
		out, err := cat.FindByID(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/slug/:slug", func(c *fiber.Ctx) error {
		out, err := cat.FindBySlug(c.UserContext(), c.Params("slug"))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	// Wildcard: /category/path/a/b/c → Params("*") == "a/b/c"
	g.Get("/path/*", func(c *fiber.Ctx) error {
		slugPath := c.Params("*")
		out, err := cat.FindBySlugPath(c.UserContext(), slugPath)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	adm := g.Group("", authmw.RequireAdmin(auth))

	adm.Post("/create-category", func(c *fiber.Ctx) error {
		var body struct {
			Name string  `json:"name"`
			Slug *string `json:"slug"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректное тело запроса",
			})
		}
		out, err := cat.CreateCategory(c.UserContext(), body.Name, body.Slug)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	adm.Put("/update-category/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректный id",
			})
		}
		var body struct {
			Name string  `json:"name"`
			Slug *string `json:"slug"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректное тело запроса",
			})
		}
		out, err := cat.UpdateCategory(c.UserContext(), int32(id), body.Name, body.Slug)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	adm.Delete("/delete-category/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректный id",
			})
		}
		out, err := cat.DeleteCategory(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
