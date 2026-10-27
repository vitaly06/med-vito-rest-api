package httpserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/rbac"
	"med-vito/api-go/internal/service"
)

// RegisterReviewRoutes — как Nest ReviewController.
func RegisterReviewRoutes(app fiber.Router, rev *service.ReviewService, auth *service.AuthService) {
	g := app.Group("/review")
	sess := authmw.RequireSession(auth)
	mod := authmw.RequireModerator(auth)

	g.Post("/send-review", sess, func(c *fiber.Ctx) error {
		var body struct {
			Text           *string `json:"text"`
			Rating         float64 `json:"rating"`
			ReviewedUserID int32   `json:"reviewedUserId"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		out, err := rev.SendReview(c.UserContext(), me.ID, body.ReviewedUserID, body.Rating, body.Text)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/:id", sess, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		var body struct {
			Text   *string `json:"text"`
			Rating float64 `json:"rating"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		isAdmin := rbac.HasMinRole(me.RoleName, 70)
		out, err := rev.UpdateReview(c.UserContext(), me.ID, int32(id), body.Rating, body.Text, isAdmin)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Delete("/:id", sess, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		me := authmw.UserFromLocals(c)
		isAdmin := rbac.HasMinRole(me.RoleName, 70)
		out, err := rev.DeleteReview(c.UserContext(), int32(id), me.ID, isAdmin)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/all-reviews", mod, func(c *fiber.Ctx) error {
		var userID *int32
		if uStr := c.Query("userId"); uStr != "" {
			if u, err := strconv.ParseInt(uStr, 10, 32); err == nil && u > 0 {
				uid := int32(u)
				userID = &uid
			}
		}
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "50"))
		out, err := rev.AllReviews(c.UserContext(), userID, page, limit)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/all-reviews-to-moderate", mod, func(c *fiber.Ctx) error {
		out, err := rev.AllReviewsToModerate(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/user-reviews/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		out, err := rev.GetUserReviews(c.UserContext(), int32(id))
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/moderate-review/:id", mod, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		st := c.Query("status")
		out, err := rev.ModerateReview(c.UserContext(), int32(id), st)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Post("/appeals", sess, func(c *fiber.Ctx) error {
		var body struct {
			ReviewID int32  `json:"reviewId"`
			Reason   string `json:"reason"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		if err := rev.CreateAppeal(c.UserContext(), me.ID, body.ReviewID, body.Reason); err != nil {
			return writeAppError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Апелляция отправлена"})
	})

	g.Get("/my-appeals", sess, func(c *fiber.Ctx) error {
		me := authmw.UserFromLocals(c)
		out, err := rev.MyAppeals(c.UserContext(), me.ID)
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Get("/all-appeals", mod, func(c *fiber.Ctx) error {
		out, err := rev.AllAppeals(c.UserContext())
		if err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(out)
	})

	g.Put("/resolve-appeal/:id", mod, func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil || id < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректный id"})
		}
		var body struct {
			Status string  `json:"status"`
			Note   *string `json:"note"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"statusCode": 400, "message": "Некорректное тело"})
		}
		me := authmw.UserFromLocals(c)
		if err := rev.ResolveAppeal(c.UserContext(), me.ID, id, body.Status, body.Note); err != nil {
			return writeAppError(c, err)
		}
		return c.JSON(fiber.Map{"message": "Апелляция обработана"})
	})
}
