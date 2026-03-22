package httpserver

import (
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

func firstFormValue(form map[string][]string, key string) string {
	if form == nil {
		return ""
	}
	v := form[key]
	if len(v) == 0 {
		return ""
	}
	return strings.TrimSpace(v[0])
}

// RegisterBannerRoutes — как Nest BannerController; статические пути до :id.
func RegisterBannerRoutes(app fiber.Router, ban *service.BannerService, auth *service.AuthService) {
	g := app.Group("/banner")
	sess := authmw.RequireSession(auth)
	adm := authmw.RequireAdmin(auth)
	opt := authmw.OptionalSession(auth)

	g.Post("/", sess, func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Нужен multipart/form-data"})
		}
		defer func() { _ = form.RemoveAll() }()
		fhs := form.File["image"]
		if len(fhs) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Изображение баннера обязательно"})
		}
		fh := fhs[0]
		file, err := fh.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Файл не читается"})
		}
		body, err := io.ReadAll(file)
		_ = file.Close()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Файл не читается"})
		}
		up := service.UploadedFile{Name: fh.Filename, ContentType: fh.Header.Get("Content-Type"), Body: body}
		place := firstFormValue(form.Value, "place")
		nav := firstFormValue(form.Value, "navigateToUrl")
		name := firstFormValue(form.Value, "name")
		me := authmw.UserFromLocals(c)
		isAdmin := me.RoleName != nil && *me.RoleName == "admin"
		out, err := ban.Create(c.UserContext(), me.ID, isAdmin, up, place, nav, name)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	g.Get("/random", func(c *fiber.Ctx) error {
		out, err := ban.FindRandom(c.UserContext(), 5)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/", func(c *fiber.Ctx) error {
		var place *string
		if v := strings.TrimSpace(c.Query("place")); v != "" {
			place = &v
		}
		out, err := ban.FindAll(c.UserContext(), place)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/my-stats/all", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := ban.GetUserBannersStats(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/all-banners-to-moderate", adm, func(c *fiber.Ctx) error {
		out, err := ban.AllBannersToModerate(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/moderate/:id", adm, func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil || id < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		st := strings.TrimSpace(c.Query("status"))
		if st == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Нужен query status"})
		}
		out, err := ban.ModerateBanner(c.UserContext(), int32(id), st)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/:id/view", opt, func(c *fiber.Ctx) error {
		bid, err := c.ParamsInt("id")
		if err != nil || bid < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		var uid *int32
		if u := authmw.UserFromLocals(c); u != nil {
			x := u.ID
			uid = &x
		}
		ip := c.IP()
		out, err := ban.RegisterView(c.UserContext(), int32(bid), uid, ip)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(out)
	})

	g.Get("/:id/stats", sess, func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil || id < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := ban.GetBannerStats(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/:id", adm, func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil || id < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Нужен multipart/form-data"})
		}
		defer func() { _ = form.RemoveAll() }()
		var file *service.UploadedFile
		if fhs := form.File["image"]; len(fhs) > 0 {
			fh := fhs[0]
			f, err := fh.Open()
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Файл не читается"})
			}
			b, err := io.ReadAll(f)
			_ = f.Close()
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Файл не читается"})
			}
			file = &service.UploadedFile{Name: fh.Filename, ContentType: fh.Header.Get("Content-Type"), Body: b}
		}
		var place, nav, name *string
		if v := firstFormValue(form.Value, "place"); v != "" {
			place = &v
		}
		if v, ok := form.Value["navigateToUrl"]; ok && len(v) > 0 {
			x := strings.TrimSpace(v[0])
			nav = &x
		}
		if v := firstFormValue(form.Value, "name"); v != "" {
			name = &v
		}
		out, err := ban.Update(c.UserContext(), int32(id), file, place, nav, name)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Delete("/:id", adm, func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil || id < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := ban.Remove(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil || id < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := ban.FindOne(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
