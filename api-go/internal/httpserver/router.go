package httpserver

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/service"
)

// AppDeps вЂ” Р·Р°РІРёСЃРёРјРѕСЃС‚Рё HTTP-СЃР»РѕСЏ (РїРѕ РјРµСЂРµ РјРёРіСЂР°С†РёРё РїРѕРїРѕР»РЅСЏРµС‚СЃСЏ).
type AppDeps struct {
	Config      config.Config
	Log         *service.LogService
	Knowledge   *service.KnowledgeBaseService
	Category    *service.CategoryService
	Auth        *service.AuthService
	User        *service.UserService
	Product     *service.ProductService
	Moderation  *service.ModerationAdminService
	Review      *service.ReviewService
	Chat        *service.ChatService
	Payment     *service.PaymentService
	Promotion   *service.PromotionService
	Statistics  *service.StatisticsService
	Support     *service.SupportService
	Address     *service.AddressService
	Banner      *service.BannerService
	CDEK        *service.CDEKService
	Deal        *service.DealService
	Reservation *service.ReservationService
}

// NewApp СЃРѕР±РёСЂР°РµС‚ Fiber: middleware + РјР°СЂС€СЂСѓС‚С‹ (handlers = Р±С‹РІС€РёРµ controllers).
func NewApp(corsOrigins string, deps AppDeps) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		BodyLimit:             80 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusRequestEntityTooLarge {
				return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
					"statusCode": fiber.StatusRequestEntityTooLarge,
					"message":    "Файлы слишком большие. Допустимо до 10 МБ на одно фото и до 80 МБ на одно объявление.",
				})
			}
			return fiber.DefaultErrorHandler(c, err)
		},
	})
	app.Use(recover.New())
	if corsOrigins != "" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     corsOrigins,
			AllowCredentials: true,
		}))
	}

	app.Get("/health", HealthHandler)
	RegisterLogRoutes(app, deps.Log)
	RegisterKnowledgeBaseRoutes(app, deps.Knowledge, deps.Auth)
	RegisterAuthRoutes(app, deps.Config, deps.Auth)
	RegisterCategoryRoutes(app, deps.Category, deps.Auth)
	RegisterSubcategoryRoutes(app, deps.Category, deps.Auth)
	RegisterSubcategoryTypeRoutes(app, deps.Category, deps.Auth)
	RegisterTypeFieldRoutes(app, deps.Category, deps.Auth)
	RegisterUserRoutes(app, deps.User, deps.Auth)
	RegisterProductRoutes(app, deps.Product, deps.Auth)
	RegisterModerationRoutes(app, deps.Moderation, deps.Auth)
	RegisterReviewRoutes(app, deps.Review, deps.Auth)
	RegisterChatRoutes(app, deps.Chat, deps.Auth)
	RegisterPaymentRoutes(app, deps.Payment, deps.Auth)
	RegisterPromotionRoutes(app, deps.Promotion, deps.Auth)
	RegisterStatisticsRoutes(app, deps.Statistics, deps.Auth)
	RegisterSupportRoutes(app, deps.Support, deps.Auth)
	RegisterAddressRoutes(app, deps.Address)
	RegisterBannerRoutes(app, deps.Banner, deps.Auth)
	RegisterCDEKRoutes(app, deps.CDEK)
	RegisterDealRoutes(app, deps.Deal, deps.Auth)
	RegisterReservationRoutes(app, deps.Reservation, deps.Auth)
	RegisterSocketIO(app, corsOrigins, deps.Auth, deps.Chat, deps.Support)
	RegisterChatWS(app, deps.Auth, deps.Chat)
	RegisterSupportWS(app, deps.Auth, deps.Support)

	// OpenAPI 2 + Swagger UI (Try it out). doc.json РёР· РїР°РєРµС‚Р° docs (swag init).
	app.Get("/docs/*", swagger.New(swagger.Config{
		Title:                  "Med Vito API (Go)",
		WithCredentials:        true,
		TryItOutEnabled:        true,
		PersistAuthorization:   true,
		DisplayRequestDuration: true,
	}))

	return app
}

