package httpserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

// RegisterSubcategoryRoutes — как Nest SubcategoryController (/subcategory).
func RegisterSubcategoryRoutes(app fiber.Router, cat *service.CategoryService, auth *service.AuthService) {
	g := app.Group("/subcategory")
	pub := g.Group("")
	pub.Get("/find-all", func(c *fiber.Ctx) error {
		out, err := cat.AdminListSubcategories(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
	pub.Get("/find-by-id/:id", func(c *fiber.Ctx) error {
		id, err := parseID32(c.Params("id"))
		if err != nil {
			return badRequestID(c)
		}
		out, err := cat.AdminFindSubcategoryByID(c.UserContext(), id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	adm := g.Group("", authmw.RequireAdmin(auth))
	adm.Post("/create-subcategory", func(c *fiber.Ctx) error {
		var body struct {
			Name       string `json:"name"`
			CategoryID int    `json:"categoryId"`
		}
		if err := c.BodyParser(&body); err != nil {
			return badRequestBody(c)
		}
		if body.CategoryID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Id категории обязателен для заполнения",
			})
		}
		out, err := cat.AdminCreateSubcategory(c.UserContext(), body.Name, int32(body.CategoryID))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})
	adm.Patch("/update-subcategory/:id", func(c *fiber.Ctx) error {
		id, err := parseID32(c.Params("id"))
		if err != nil {
			return badRequestID(c)
		}
		var body struct {
			Name       *string `json:"name"`
			CategoryID *int    `json:"categoryId"`
		}
		if err := c.BodyParser(&body); err != nil {
			return badRequestBody(c)
		}
		var catID *int32
		if body.CategoryID != nil {
			if *body.CategoryID <= 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"statusCode": fiber.StatusBadRequest,
					"message":    "Id категории должен быть положительным числом",
				})
			}
			v := int32(*body.CategoryID)
			catID = &v
		}
		out, err := cat.AdminUpdateSubcategory(c.UserContext(), id, body.Name, catID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
	adm.Delete("/delete-subcategory/:id", func(c *fiber.Ctx) error {
		id, err := parseID32(c.Params("id"))
		if err != nil {
			return badRequestID(c)
		}
		out, err := cat.AdminDeleteSubcategory(c.UserContext(), id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}

// RegisterSubcategoryTypeRoutes — как Nest SubcategoryTypeController; сначала статика, потом :id.
func RegisterSubcategoryTypeRoutes(app fiber.Router, cat *service.CategoryService, auth *service.AuthService) {
	g := app.Group("/subcategory-type")
	pub := g.Group("")
	pub.Get("/find-all", func(c *fiber.Ctx) error {
		out, err := cat.AdminListSubcategoryTypes(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
	pub.Get("/find-by-id/:id", func(c *fiber.Ctx) error {
		id, err := parseID32(c.Params("id"))
		if err != nil {
			return badRequestID(c)
		}
		out, err := cat.AdminFindSubcategoryTypeByID(c.UserContext(), id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	adm := g.Group("", authmw.RequireAdmin(auth))
	adm.Post("/create", func(c *fiber.Ctx) error {
		var body struct {
			Name          string `json:"name"`
			SubcategoryID int    `json:"subcategoryId"`
		}
		if err := c.BodyParser(&body); err != nil {
			return badRequestBody(c)
		}
		if body.SubcategoryID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Id подкатегории обязателен для заполнения",
			})
		}
		out, err := cat.AdminCreateSubcategoryType(c.UserContext(), body.Name, int32(body.SubcategoryID))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})
	adm.Patch("/:id", func(c *fiber.Ctx) error {
		id, err := parseID32(c.Params("id"))
		if err != nil {
			return badRequestID(c)
		}
		var body struct {
			Name          *string `json:"name"`
			SubcategoryID *int    `json:"subcategoryId"`
		}
		if err := c.BodyParser(&body); err != nil {
			return badRequestBody(c)
		}
		var subID *int32
		if body.SubcategoryID != nil {
			if *body.SubcategoryID <= 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"statusCode": fiber.StatusBadRequest,
					"message":    "Id подкатегории должен быть положительным числом",
				})
			}
			v := int32(*body.SubcategoryID)
			subID = &v
		}
		out, err := cat.AdminUpdateSubcategoryType(c.UserContext(), id, body.Name, subID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
	adm.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := parseID32(c.Params("id"))
		if err != nil {
			return badRequestID(c)
		}
		out, err := cat.AdminDeleteSubcategoryType(c.UserContext(), id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}

// RegisterTypeFieldRoutes — как Nest TypeFieldController.
func RegisterTypeFieldRoutes(app fiber.Router, cat *service.CategoryService, auth *service.AuthService) {
	g := app.Group("/type-field")
	pub := g.Group("")
	pub.Get("/find-all", func(c *fiber.Ctx) error {
		out, err := cat.AdminListTypeFields(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
	pub.Get("/find-by-id/:id", func(c *fiber.Ctx) error {
		id, err := parseID32(c.Params("id"))
		if err != nil {
			return badRequestID(c)
		}
		out, err := cat.AdminFindTypeFieldByID(c.UserContext(), id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	adm := g.Group("", authmw.RequireAdmin(auth))
	adm.Post("/create", func(c *fiber.Ctx) error {
		var body struct {
			Name   string `json:"name"`
			TypeID int    `json:"typeId"`
		}
		if err := c.BodyParser(&body); err != nil {
			return badRequestBody(c)
		}
		if body.TypeID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"statusCode": fiber.StatusBadRequest,
				"message":    "Id типа подкатегории обязателен для заполнения",
			})
		}
		out, err := cat.AdminCreateTypeField(c.UserContext(), body.Name, int32(body.TypeID))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})
	adm.Patch("/:id", func(c *fiber.Ctx) error {
		id, err := parseID32(c.Params("id"))
		if err != nil {
			return badRequestID(c)
		}
		var body struct {
			Name   *string `json:"name"`
			TypeID *int    `json:"typeId"`
		}
		if err := c.BodyParser(&body); err != nil {
			return badRequestBody(c)
		}
		var tid *int32
		if body.TypeID != nil {
			if *body.TypeID <= 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"statusCode": fiber.StatusBadRequest,
					"message":    "Id типа подкатегории должен быть положительным числом",
				})
			}
			v := int32(*body.TypeID)
			tid = &v
		}
		out, err := cat.AdminUpdateTypeField(c.UserContext(), id, body.Name, tid)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
	adm.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := parseID32(c.Params("id"))
		if err != nil {
			return badRequestID(c)
		}
		out, err := cat.AdminDeleteTypeField(c.UserContext(), id)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}

func parseID32(s string) (int32, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil || n <= 0 {
		return 0, strconv.ErrSyntax
	}
	return int32(n), nil
}

func badRequestID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"statusCode": fiber.StatusBadRequest,
		"message":    "Некорректный id",
	})
}

func badRequestBody(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"statusCode": fiber.StatusBadRequest,
		"message":    "Некорректное тело запроса",
	})
}
