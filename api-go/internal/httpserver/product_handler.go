package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

const maxProductImages = 15

func scalarToString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func multipartFirst(form *multipart.Form, keys ...string) string {
	for _, k := range keys {
		if form == nil || form.Value == nil {
			continue
		}
		if vs := form.Value[k]; len(vs) > 0 && strings.TrimSpace(vs[0]) != "" {
			return vs[0]
		}
	}
	return ""
}

func readImageFilesFromMultipart(form *multipart.Form, field string, max int) ([]service.UploadedFile, error) {
	if form == nil {
		return nil, nil
	}
	fhs := form.File[field]
	if len(fhs) > max {
		return nil, fmt.Errorf("РјР°РєСЃРёРјСѓРј %d РёР·РѕР±СЂР°Р¶РµРЅРёР№", max)
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

// parseCreateDraftInputs вЂ” multipart | JSON | urlencoded; Р±РµР· MultipartForm РґР»СЏ РЅРµ-multipart (РёРЅР°С‡Рµ 400).
func parseCreateDraftInputs(c *fiber.Ctx) (
	name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, fieldJSON, videoStr string,
	files []service.UploadedFile,
	err error,
) {
	ct := strings.ToLower(c.Get("Content-Type"))
	if strings.HasPrefix(ct, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err != nil {
			return "", "", "", "", "", "", "", "", "", "", "", nil, err
		}
		name = multipartFirst(form, "name")
		priceStr = multipartFirst(form, "price")
		quantityStr = multipartFirst(form, "quantity")
		state = multipartFirst(form, "state")
		description = multipartFirst(form, "description")
		address = multipartFirst(form, "address")
		categoryStr = multipartFirst(form, "categoryId")
		subStr = multipartFirst(form, "subcategoryId", "subCategoryId")
		typeStr = multipartFirst(form, "typeId", "typeld")
		fieldJSON = multipartFirst(form, "fieldValues")
		videoStr = multipartFirst(form, "videoUrl")
		files, err = readImageFilesFromMultipart(form, "images", maxProductImages)
		return name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, fieldJSON, videoStr, files, err
	}
	if strings.Contains(ct, "application/json") {
		raw := bytes.TrimSpace(c.Body())
		if len(raw) == 0 {
			return "", "", "", "", "", "", "", "", "", "", "", nil, nil
		}
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			return "", "", "", "", "", "", "", "", "", "", "", nil, err
		}
		fv := ""
		if fvRaw, ok := body["fieldValues"]; ok {
			if b, err := json.Marshal(fvRaw); err == nil {
				fv = strings.TrimSpace(string(b))
			}
		}
		sub := scalarToString(body["subcategoryId"])
		if sub == "" {
			sub = scalarToString(body["subCategoryId"])
		}
		return scalarToString(body["name"]),
			scalarToString(body["price"]),
			scalarToString(body["quantity"]),
			scalarToString(body["state"]),
			scalarToString(body["description"]),
			scalarToString(body["address"]),
			scalarToString(body["categoryId"]),
			sub,
			func() string {
				t := scalarToString(body["typeId"])
				if t == "" {
					t = scalarToString(body["typeld"])
				}
				return t
			}(),
			fv,
			scalarToString(body["videoUrl"]),
			nil,
			nil
	}
	name = c.FormValue("name")
	priceStr = c.FormValue("price")
	quantityStr = c.FormValue("quantity")
	state = c.FormValue("state")
	description = c.FormValue("description")
	address = c.FormValue("address")
	categoryStr = c.FormValue("categoryId")
	subStr = c.FormValue("subcategoryId")
	if subStr == "" {
		subStr = c.FormValue("subCategoryId")
	}
	typeStr = c.FormValue("typeId")
	if typeStr == "" {
		typeStr = c.FormValue("typeld")
	}
	fieldJSON = c.FormValue("fieldValues")
	videoStr = c.FormValue("videoUrl")
	return name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, fieldJSON, videoStr, nil, nil
}

// parseUpdateProductInputs вЂ” PATCH С‚РѕРІР°СЂР°: РїРѕРґРґРµСЂР¶РёРІР°РµС‚ РєР°С‚РµРіРѕСЂРёРё Рё typeId.
func parseUpdateProductInputs(c *fiber.Ctx) (
	name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, videoStr, fieldJSON string,
	files []service.UploadedFile,
	err error,
) {
	ct := strings.ToLower(c.Get("Content-Type"))
	if strings.HasPrefix(ct, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err != nil {
			return "", "", "", "", "", "", "", "", "", "", "", nil, err
		}
		name = multipartFirst(form, "name")
		priceStr = multipartFirst(form, "price")
		quantityStr = multipartFirst(form, "quantity")
		state = multipartFirst(form, "state")
		description = multipartFirst(form, "description")
		address = multipartFirst(form, "address")
		categoryStr = multipartFirst(form, "categoryId")
		subStr = multipartFirst(form, "subcategoryId", "subCategoryId")
		typeStr = multipartFirst(form, "typeId", "typeld")
		videoStr = multipartFirst(form, "videoUrl")
		fieldJSON = multipartFirst(form, "fieldValues")
		files, err = readImageFilesFromMultipart(form, "images", maxProductImages)
		return name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, videoStr, fieldJSON, files, err
	}
	if strings.Contains(ct, "application/json") {
		raw := bytes.TrimSpace(c.Body())
		if len(raw) == 0 {
			return "", "", "", "", "", "", "", "", "", "", "", nil, nil
		}
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			return "", "", "", "", "", "", "", "", "", "", "", nil, err
		}
		fv := ""
		if fvRaw, ok := body["fieldValues"]; ok {
			if b, err := json.Marshal(fvRaw); err == nil {
				fv = strings.TrimSpace(string(b))
			}
		}
		sub := scalarToString(body["subcategoryId"])
		if sub == "" {
			sub = scalarToString(body["subCategoryId"])
		}
		return scalarToString(body["name"]),
			scalarToString(body["price"]),
			scalarToString(body["quantity"]),
			scalarToString(body["state"]),
			scalarToString(body["description"]),
			scalarToString(body["address"]),
			scalarToString(body["categoryId"]),
			sub,
			func() string {
				t := scalarToString(body["typeId"])
				if t == "" {
					t = scalarToString(body["typeld"])
				}
				return t
			}(),
			scalarToString(body["videoUrl"]),
			fv,
			nil,
			nil
	}
	name = c.FormValue("name")
	priceStr = c.FormValue("price")
	quantityStr = c.FormValue("quantity")
	state = c.FormValue("state")
	description = c.FormValue("description")
	address = c.FormValue("address")
	categoryStr = c.FormValue("categoryId")
	subStr = c.FormValue("subcategoryId")
	if subStr == "" {
		subStr = c.FormValue("subCategoryId")
	}
	typeStr = c.FormValue("typeId")
	if typeStr == "" {
		typeStr = c.FormValue("typeld")
	}
	videoStr = c.FormValue("videoUrl")
	fieldJSON = c.FormValue("fieldValues")
	return name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, videoStr, fieldJSON, nil, nil
}

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
	keys := []string{"search", "categorySlug", "subCategorySlug", "typeSlug", "minPrice", "maxPrice", "minRating", "maxRating", "state", "region", "profileType", "hasSecureDeal", "fieldValues", "sortBy", "page", "limit"}
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
	if v := strings.TrimSpace(c.Query("minRating")); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			q.MinRating = &n
		}
	}
	if v := strings.TrimSpace(c.Query("maxRating")); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			q.MaxRating = &n
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
	if v := strings.TrimSpace(c.Query("hasSecureDeal")); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			q.HasSecureDeal = &b
		}
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

// RegisterProductRoutes вЂ” РєР°Рє Nest ProductController; СЃС‚Р°С‚РёС‡РµСЃРєРёРµ РїСѓС‚Рё РґРѕ :id.
func RegisterProductRoutes(app fiber.Router, p *service.ProductService, auth *service.AuthService) {
	g := app.Group("/product")
	sess := authmw.RequireSession(auth)
	adm := authmw.RequireAdmin(auth)
	opt := authmw.OptionalSession(auth)

	g.Post("/create", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, fieldJSON, videoStr, files, err := parseCreateDraftInputs(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РўРµР»Рѕ: " + err.Error()})
		}

		isComplete := strings.TrimSpace(name) != "" &&
			strings.TrimSpace(address) != "" &&
			strings.TrimSpace(state) != "" &&
			strings.TrimSpace(priceStr) != "" &&
			strings.TrimSpace(categoryStr) != "" &&
			strings.TrimSpace(subStr) != ""

		if isComplete {
			out, err := p.CreateProduct(c.UserContext(), me.ID,
				name, priceStr, quantityStr, state, description, address,
				categoryStr, subStr, typeStr, fieldJSON, videoStr, files)
			if err == nil {
				return c.Status(fiber.StatusCreated).JSON(out)
			}
			return writeAppError(c, err)
		}

		// Неполная форма — сохраняем как черновик.
		out, err := p.CreateDraft(c.UserContext(), me.ID,
			name, priceStr, quantityStr, state, description, address,
			categoryStr, subStr, typeStr, fieldJSON, videoStr, files)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	g.Post("/create-draft", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, fieldJSON, videoStr, files, err := parseCreateDraftInputs(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Р§РµСЂРЅРѕРІРёРє: " + err.Error()})
		}
		out, err := p.CreateDraft(c.UserContext(), me.ID,
			name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, fieldJSON, videoStr, files)
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

	g.Get("/recommended", opt, func(c *fiber.Ctx) error {
		subID, err := strconv.ParseInt(c.Query("subcategoryId"), 10, 32)
		if err != nil || subID < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ subcategoryId"})
		}
		limit := 20
		if v := strings.TrimSpace(c.Query("limit")); v != "" {
			n, err := strconv.Atoi(v)
			if err == nil && n > 0 {
				limit = n
			}
		}
		var viewer *int32
		if u := authmw.UserFromLocals(c); u != nil {
			viewer = &u.ID
		}
		out, err := p.RecommendedBySubcategory(c.UserContext(), viewer, int32(subID), limit)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/user-products/:id", opt, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ id"})
		}
		var viewer *int32
		if u := authmw.UserFromLocals(c); u != nil {
			viewer = &u.ID
		}
		out, err := p.ProductsByUserID(c.UserContext(), viewer, int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/my-drafts", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := p.MyDrafts(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/publish-draft/:id", sess, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ id"})
		}
		me := authmw.UserFromLocals(c)
		out, err := p.PublishDraft(c.UserContext(), int32(id), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/add-to-favorites/:id", sess, func(c *fiber.Ctx) error {
		pid, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ id"})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ id"})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ id"})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ id"})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ id"})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ promotionId"})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ id"})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РќРµРєРѕСЂСЂРµРєС‚РЅС‹Р№ id"})
		}
		me := authmw.UserFromLocals(c)
		name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, videoStr, fieldJSON, files, err := parseUpdateProductInputs(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "РўРѕРІР°СЂ: " + err.Error()})
		}
		out, err := p.UpdateProduct(c.UserContext(), int32(id), me.ID,
			name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, videoStr, fieldJSON, files)
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
