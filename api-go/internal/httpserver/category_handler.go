package httpserver

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/encoding/charmap"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func parseCategoryPayload(c *fiber.Ctx) (string, *string, error) {
	var body struct {
		Name string  `json:"name"`
		Slug *string `json:"slug"`
	}
	if err := c.BodyParser(&body); err != nil {
		return "", nil, err
	}

	name := strings.TrimSpace(body.Name)
	slug := body.Slug
	if slug != nil {
		s := strings.TrimSpace(*slug)
		if s == "" {
			slug = nil
		} else {
			slug = &s
		}
	}

	// Совместимость со старым/альтернативным фронтом.
	if name == "" || slug == nil {
		var raw map[string]any
		if err := c.BodyParser(&raw); err == nil {
			if name == "" {
				for _, key := range []string{"name", "title", "categoryName"} {
					if v, ok := raw[key]; ok {
						if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
							name = strings.TrimSpace(s)
							break
						}
					}
				}
			}
			if slug == nil {
				for _, key := range []string{"slug", "url", "categorySlug"} {
					if v, ok := raw[key]; ok {
						if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
							t := strings.TrimSpace(s)
							slug = &t
							break
						}
					}
				}
			}
		}
	}

	return name, slug, nil
}

func writeAppError(c *fiber.Ctx, err error) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	var ae *service.AppError
	if errors.As(err, &ae) {
		return c.Status(ae.Status).JSON(fiber.Map{
			"statusCode": ae.Status,
			"message":    normalizeResponseMessage(ae.Message),
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"statusCode": fiber.StatusInternalServerError,
		"message":    normalizeResponseMessage(err.Error()),
	})
}

func normalizeResponseMessage(msg string) string {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return msg
	}
	if !strings.ContainsAny(msg, "РСЃ") {
		return msg
	}
	decoded, err := charmap.Windows1251.NewEncoder().Bytes([]byte(msg))
	if err != nil {
		return msg
	}
	fixed := strings.TrimSpace(string(decoded))
	if fixed == "" || !utf8.ValidString(fixed) {
		return msg
	}
	return fixed
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

	createCategory := func(c *fiber.Ctx) error {
		name, slug, err := parseCategoryPayload(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректное тело запроса",
			})
		}
		out, err := cat.CreateCategory(c.UserContext(), name, slug)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	}
	adm.Post("/create-category", createCategory)
	adm.Post("/create", createCategory) // Алиас для фронта.

	updateCategory := func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректный id",
			})
		}
		name, slug, err := parseCategoryPayload(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Некорректное тело запроса",
			})
		}
		out, err := cat.UpdateCategory(c.UserContext(), int32(id), name, slug)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	}
	adm.Put("/update-category/:id", updateCategory)
	adm.Patch("/update-category/:id", updateCategory) // Часто фронт шлёт PATCH.
	adm.Patch("/:id", updateCategory)                 // Алиас под REST-стиль.

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
