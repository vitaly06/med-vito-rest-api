package httpserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/service"
)

// AppDeps — зависимости HTTP-слоя (по мере миграции пополняется).
type AppDeps struct {
	Config     config.Config
	Log        *service.LogService
	Knowledge  *service.KnowledgeBaseService
	Category   *service.CategoryService
	Auth       *service.AuthService
	User       *service.UserService
	Product    *service.ProductService
	Moderation *service.ModerationAdminService
	Review     *service.ReviewService
	Chat       *service.ChatService
	Payment    *service.PaymentService
	Promotion  *service.PromotionService
	Statistics *service.StatisticsService
	Support    *service.SupportService
	Address    *service.AddressService
	Banner     *service.BannerService
}

// NewApp собирает Fiber: middleware + маршруты (handlers = бывшие controllers).
func NewApp(corsOrigins string, deps AppDeps) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
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
	RegisterSocketIO(app, corsOrigins, deps.Auth, deps.Chat, deps.Support)

	// OpenAPI 2 + Swagger UI (Try it out). doc.json из пакета docs (swag init).
	app.Get("/docs/*", swagger.New(swagger.Config{
		Title:                  "Med Vito API (Go)",
		WithCredentials:        true,
		TryItOutEnabled:        true,
		PersistAuthorization:   true,
		DisplayRequestDuration: true,
	}))

	return app
}
