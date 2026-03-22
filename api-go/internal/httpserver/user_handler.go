package httpserver

import (
	"io"
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

// RegisterUserRoutes — пути как Nest UserController.
// Статические сегменты (update-settings, verify-email, …) объявлены выше :id, чтобы Fiber не принял их за id.
func RegisterUserRoutes(app fiber.Router, u *service.UserService, auth *service.AuthService) {
	g := app.Group("/user")
	sess := authmw.RequireSession(auth)
	adm := authmw.RequireAdmin(auth)

	// Список юзеров только для admin (в Nest был открыт — у нас закрыто осознанно).
	g.Get("/find-all", adm, func(c *fiber.Ctx) error {
		out, err := u.FindAllAdmin(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/info/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := u.GetUserInfo(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/remaining-free-ads", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := u.GetRemainingFreeAds(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/show-number/:userId", sess, func(c *fiber.Ctx) error {
		sellerID, err := strconv.ParseInt(c.Params("userId"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный userId"})
		}
		me := authmw.UserFromLocals(c)
		out, err := u.GetPhoneNumber(c.UserContext(), int32(sellerID), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Patch("/update-settings", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		fullName := c.FormValue("fullName")
		phone := c.FormValue("phoneNumber")
		isAns := c.FormValue("isAnswersCall")
		prof := c.FormValue("profileType")
		var photoData []byte
		var photoName, photoCT string
		fh, err := c.FormFile("photo")
		if err == nil && fh != nil {
			photoName = fh.Filename
			photoCT = fh.Header.Get("Content-Type")
			f, err := fh.Open()
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Не удалось прочитать файл"})
			}
			photoData, err = io.ReadAll(f)
			_ = f.Close()
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Не удалось прочитать файл"})
			}
		}
		out, err := u.UpdateSettingsMultipart(c.UserContext(), me.ID, fullName, phone, isAns, prof, photoName, photoCT, photoData)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/verify-email", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := u.VerifyEmail(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/verify-code", func(c *fiber.Ctx) error {
		code := c.Query("code")
		out, err := u.VerifyEmailCode(c.UserContext(), code)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/set-balance/:userId", adm, func(c *fiber.Ctx) error {
		uid, err := strconv.ParseInt(c.Params("userId"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный userId"})
		}
		bal := c.Query("balance")
		out, err := u.SetBalance(c.UserContext(), int32(uid), bal)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/toggle-banned/:id", adm, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := u.ToggleBanned(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Patch("/:id", adm, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		var body struct {
			FullName     *string  `json:"fullName"`
			Email        *string  `json:"email"`
			PhoneNumber  *string  `json:"phoneNumber"`
			ProfileType  *string  `json:"profileType"`
			BonusBalance *float64 `json:"bonusBalance"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		out, err := u.AdminUpdateUser(c.UserContext(), int32(id), body.FullName, body.Email, body.PhoneNumber, body.ProfileType, body.BonusBalance)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Delete("/:id", adm, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := u.DeleteUser(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})
}
