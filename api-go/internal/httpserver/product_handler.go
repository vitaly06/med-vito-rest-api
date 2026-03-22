package httpserver

import (
	"io"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func collectImageFiles(c *fiber.Ctx, field string, max int) ([]service.UploadedFile, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	fhs := form.File[field]
	if len(fhs) > max {
		fhs = fhs[:max]
	}
	var out []service.UploadedFile
	for _, fh := range fhs {
		f, err := fh.Open()
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(f)
		_ = f.Close()
		if err != nil {
			return nil, err
		}
		out = append(out, service.UploadedFile{
			Name:        fh.Filename,
			ContentType: fh.Header.Get("Content-Type"),
			Body:        b,
		})
	}
	return out, nil
}

func productSearchUseExpanded(c *fiber.Ctx) bool {
	keys := []string{"search", "categorySlug", "subCategorySlug", "typeSlug", "minPrice", "maxPrice", "state", "region", "profileType", "fieldValues", "sortBy", "page", "limit"}
	for _, k := range keys {
		if strings.TrimSpace(c.Query(k)) != "" {
			return true
		}
	}
	return false
}

func parseProductSearchQuery(c *fiber.Ctx) service.ProductSearchQuery {
	q := service.ProductSearchQuery{SortBy: c.Query("sortBy")}
	if v := strings.TrimSpace(c.Query("search")); v != "" {
		q.Search = &v
	}
	if v := strings.TrimSpace(c.Query("categorySlug")); v != "" {
		q.CategorySlug = &v
	}
	if v := strings.TrimSpace(c.Query("subCategorySlug")); v != "" {
		q.SubCategorySlug = &v
	}
	if v := strings.TrimSpace(c.Query("typeSlug")); v != "" {
		q.TypeSlug = &v
	}
	if v := strings.TrimSpace(c.Query("minPrice")); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			x := int32(n)
			q.MinPrice = &x
		}
	}
	if v := strings.TrimSpace(c.Query("maxPrice")); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			x := int32(n)
			q.MaxPrice = &x
		}
	}
	if v := strings.TrimSpace(c.Query("state")); v != "" {
		q.State = &v
	}
	if v := strings.TrimSpace(c.Query("region")); v != "" {
		q.Region = &v
	}
	if v := strings.TrimSpace(c.Query("profileType")); v != "" {
		q.ProfileType = &v
	}
	if v := strings.TrimSpace(c.Query("fieldValues")); v != "" {
		q.FieldValuesJSON = &v
	}
	if v := strings.TrimSpace(c.Query("page")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			q.Page = n
		}
	}
	if v := strings.TrimSpace(c.Query("limit")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			q.Limit = n
		}
	}
	return q
}

// RegisterProductRoutes — как Nest ProductController; статические пути до :id.
func RegisterProductRoutes(app fiber.Router, p *service.ProductService, auth *service.AuthService) {
	g := app.Group("/product")
	sess := authmw.RequireSession(auth)
	adm := authmw.RequireAdmin(auth)
	opt := authmw.OptionalSession(auth)

	g.Post("/create", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		files, err := collectImageFiles(c, "images", 8)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Файлы: " + err.Error()})
		}
		out, err := p.CreateProduct(c.UserContext(), me.ID,
			c.FormValue("name"), c.FormValue("price"), c.FormValue("state"),
			c.FormValue("description"), c.FormValue("address"),
			c.FormValue("categoryId"), c.FormValue("subcategoryId"), c.FormValue("typeId"),
			c.FormValue("fieldValues"), c.FormValue("videoUrl"), files)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	g.Get("/available-filters", func(c *fiber.Ctx) error {
		cs := nonemptyPtr(c.Query("categorySlug"))
		ss := nonemptyPtr(c.Query("subCategorySlug"))
		ts := nonemptyPtr(c.Query("typeSlug"))
		out, err := p.AvailableFilters(c.UserContext(), cs, ss, ts)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/all-products", opt, func(c *fiber.Ctx) error {
		var viewer *int32
		if u := authmw.UserFromLocals(c); u != nil {
			viewer = &u.ID
		}
		use := productSearchUseExpanded(c)
		out, err := p.FindAll(c.UserContext(), viewer, parseProductSearchQuery(c), use)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/random-products", opt, func(c *fiber.Ctx) error {
		var viewer *int32
		if u := authmw.UserFromLocals(c); u != nil {
			viewer = &u.ID
		}
		out, err := p.RandomProducts(c.UserContext(), viewer)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/user-products/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := p.ProductsByUserID(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/add-to-favorites/:id", sess, func(c *fiber.Ctx) error {
		pid, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		me := authmw.UserFromLocals(c)
		out, err := p.AddFavorite(c.UserContext(), me.ID, int32(pid))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Delete("/remove-from-favorites/:id", sess, func(c *fiber.Ctx) error {
		pid, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		me := authmw.UserFromLocals(c)
		out, err := p.RemoveFavorite(c.UserContext(), me.ID, int32(pid))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/my-favorites", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := p.MyFavorites(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/product-card/:id", opt, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		var viewer *int32
		if u := authmw.UserFromLocals(c); u != nil {
			viewer = &u.ID
		}
		out, err := p.GetProductCard(c.UserContext(), int32(id), viewer)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/toggle-product/:id", sess, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		me := authmw.UserFromLocals(c)
		out, err := p.ToggleProduct(c.UserContext(), int32(id), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/moderate-product/:id", adm, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		st := c.Query("status")
		reason := c.Query("reason")
		if err := p.ModerateProduct(c.UserContext(), int32(id), st, reason); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{})
	})

	g.Get("/all-products-to-moderate", adm, func(c *fiber.Ctx) error {
		out, err := p.AllProductsToModerate(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/promoted-products", adm, func(c *fiber.Ctx) error {
		out, err := p.AllPromotedProducts(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/toggle-promotion/:promotionId", adm, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("promotionId"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный promotionId"})
		}
		out, err := p.TogglePromotion(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Delete("/:id", sess, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		me := authmw.UserFromLocals(c)
		out, err := p.DeleteProduct(c.UserContext(), int32(id), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Patch("/:id", sess, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		me := authmw.UserFromLocals(c)
		files, err := collectImageFiles(c, "images", 8)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Файлы: " + err.Error()})
		}
		out, err := p.UpdateProduct(c.UserContext(), int32(id), me.ID,
			c.FormValue("name"), c.FormValue("price"), c.FormValue("state"),
			c.FormValue("description"), c.FormValue("address"), c.FormValue("videoUrl"),
			c.FormValue("fieldValues"), files)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}

func nonemptyPtr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}
