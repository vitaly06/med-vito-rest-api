package httpserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	authmw "med-vito/api-go/internal/httpserver/middleware"
	"med-vito/api-go/internal/service"
)

// RegisterReviewRoutes — как Nest ReviewController.
func RegisterReviewRoutes(app fiber.Router, rev *service.ReviewService, auth *service.AuthService) {
	g := app.Group("/review")
	sess := authmw.RequireSession(auth)
	adm := authmw.RequireAdmin(auth)

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

	g.Get("/all-reviews-to-moderate", adm, func(c *fiber.Ctx) error {
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

	g.Put("/moderate-review/:id", adm, func(c *fiber.Ctx) error {
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
}
