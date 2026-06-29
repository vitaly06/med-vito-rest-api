// @title Med Vito REST API (Go)
// @version 1.0
// @description REST API (Go/Fiber). Документация и Try it out: GET /docs/index.html. Для цепочки slug в path используй %2F (например elektronika%2Ftelefony).
// @description **Сессия в Swagger:** кнопка Authorize, схема SessionId - вставь значение cookie session_id (после POST /auth/sign-in в Try it out или из DevTools). WithCredentials уже включён.
// @description **Socket.IO:** тот же хост и CORS что у REST; путь `/socket.io`; авторизация cookie `session_id`. Неймспейсы `/chat` и `/support`.
// @description Клиент `socket.io`: надёжнее **v2.x / v3.x**; движок Go (`googollee/go-socket.io` ~v1.7) часто **не коннектится** с **socket.io-client v4** из‑за протокола handshake.
// @description **/chat** — исходящие: `joinChat` {chatId}, `sendMessage` {chatId, content}, `markAsRead` {chatId}. Входящие: `newMessage`, `messagesRead`, `newChatMessage`.
// @description **/support** — исходящие: `joinTicket`, `leaveTicket`, `sendSupportMessage` {ticketId, message:{text}}, `supportTyping` {ticketId, isTyping}, `updateTicketStatus` {ticketId, status} (модератор). Входящие: `newSupportMessage`, `moderatorResponse`, `supportUserTyping`, `ticketStatusUpdated`; по неймспейсу `newSupportTicket`; в комнате тикета `ticketAssigned`.
// @host localhost:3000
// @BasePath /
// @tag.name websocket
// @tag.description Socket.IO на `/socket.io` (см. блок **Socket.IO** в описании API выше). В Swagger нет Try it out для WS.
// @tag.name product-draft
// @tag.description Черновики объявлений: create-draft, my-drafts, publish-draft; cookie session_id.
// @securityDefinitions.apikey SessionId
// @in cookie
// @name session_id
package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	docs "med-vito/api-go/docs"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/dbmigrate"
	"med-vito/api-go/internal/httpserver"
	"med-vito/api-go/internal/pkg/s3client"
	"med-vito/api-go/internal/repository"
	"med-vito/api-go/internal/service"
)

func main() {
	loadEnvFiles()
	cfg := config.Load()
	applySwaggerRuntimeInfo(cfg.BaseURL)
	if cfg.DatabaseURL == "" {
		log.Fatal("нужен DATABASE_URL в окружении (PostgreSQL, как в Prisma schema)")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("postgres ping: %v", err)
	}
	if cfg.AutoMigrate {
		if err := dbmigrate.Apply(ctx, pool); err != nil {
			log.Fatalf("auto-migrate: %v", err)
		}
		log.Printf("auto-migrate: applied successfully")
	}

	opt := &redis.Options{Addr: cfg.RedisAddr}
	if cfg.RedisPassword != "" {
		opt.Password = cfg.RedisPassword
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer rdb.Close()

	logRepo := repository.NewLogPG(pool)
	logSvc := service.NewLogService(logRepo)
	kbRepo := repository.NewKnowledgeBasePG(pool)
	kbSvc := service.NewKnowledgeBaseService(kbRepo)
	catRepo := repository.NewCategoryPG(pool)
	catSvc := service.NewCategoryService(catRepo)
	userRepo := repository.NewUserPG(pool)
	authSvc := service.NewAuthService(cfg, rdb, userRepo)

	var s3c *s3client.Client
	if sc, err := s3client.New(cfg.S3Endpoint, cfg.S3Region, cfg.S3Bucket, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3PublicBase); err == nil {
		s3c = sc
	} else {
		log.Printf("S3: %v (загрузка аватара отключена)", err)
	}
	userSvc := service.NewUserService(cfg, userRepo, rdb, s3c)
	prodRepo := repository.NewProductPG(pool)
	statRepo := repository.NewStatisticsPG(pool)
	statSvc := service.NewStatisticsService(statRepo)
	supRepo := repository.NewSupportPG(pool)
	supSvc := service.NewSupportService(supRepo)
	prodSvc := service.NewProductService(prodRepo, s3c, userSvc, statSvc, supSvc)
	moderationSvc := service.NewModerationService(cfg, prodRepo, supSvc)
	moderationAdminSvc := service.NewModerationAdminService(prodRepo)
	revRepo := repository.NewReviewPG(pool)
	chatRepo := repository.NewChatPG(pool)
	revSvc := service.NewReviewService(revRepo, chatRepo, cfg)
	chatSvc := service.NewChatService(chatRepo)
	payRepo := repository.NewPaymentPG(pool)
	dealRepo := repository.NewDealPG(pool)
	paySvc := service.NewPaymentService(cfg, payRepo, dealRepo)
	cdekSvc := service.NewCDEKService(cfg)
	reservationRepo := repository.NewReservationPG(pool)
	dealSvc := service.NewDealService(cfg, dealRepo, logRepo, chatRepo, paySvc, reservationRepo, cdekSvc, userRepo)
	dealSvc.StartPayoutWorker(ctx)
	reservationSvc := service.NewReservationService(cfg, reservationRepo, userRepo, chatRepo)
	reservationSvc.StartWorker(ctx)
	promoRepo := repository.NewPromotionPG(pool)
	promoSvc := service.NewPromotionService(promoRepo)
	addrSvc := service.NewAddressService(cfg.DadataAPIKey)
	banRepo := repository.NewBannerPG(pool)
	banSvc := service.NewBannerService(banRepo, s3c)

	raw := os.Getenv("CORS_ORIGIN")
	var origins string
	if raw != "" {
		parts := strings.Split(raw, ",")
		for i := range parts {
			v := strings.TrimSpace(parts[i])
			v = strings.Trim(v, `"'`)
			v = strings.TrimRight(v, "/")
			parts[i] = v
		}
		origins = strings.Join(parts, ",")
	}

	app := httpserver.NewApp(origins, httpserver.AppDeps{
		Config:      cfg,
		Log:         logSvc,
		Knowledge:   kbSvc,
		Category:    catSvc,
		Auth:        authSvc,
		User:        userSvc,
		Product:     prodSvc,
		Moderation:  moderationAdminSvc,
		Review:      revSvc,
		Chat:        chatSvc,
		Payment:     paySvc,
		Promotion:   promoSvc,
		Statistics:  statSvc,
		Support:     supSvc,
		Address:     addrSvc,
		Banner:      banSvc,
		CDEK:        cdekSvc,
		Deal:        dealSvc,
		Reservation: reservationSvc,
	})
	moderationSvc.Start(ctx)
	revSvc.StartAIModerationWorker(ctx)

	addr := ":" + cfg.Port
	log.Printf("Fiber %s — Socket.IO /socket.io namespaces /chat, /support; … /docs", addr)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("fiber shutdown: %v", err)
		}
	}()
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}

func loadEnvFiles() {
	seen := map[string]struct{}{}
	candidates := make([]string, 0, 16)

	addChain := func(base string) {
		if strings.TrimSpace(base) == "" {
			return
		}
		base = filepath.Clean(base)
		dir := base
		for i := 0; i < 6; i++ {
			candidates = append(candidates, filepath.Join(dir, ".env"))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	if wd, err := os.Getwd(); err == nil {
		addChain(wd)
	}
	if exe, err := os.Executable(); err == nil {
		addChain(filepath.Dir(exe))
	}

	paths := make([]string, 0, len(candidates))
	for _, p := range candidates {
		cp := filepath.Clean(p)
		if _, ok := seen[cp]; ok {
			continue
		}
		if _, err := os.Stat(cp); err != nil {
			continue
		}
		seen[cp] = struct{}{}
		paths = append(paths, cp)
	}
	if len(paths) == 0 {
		return
	}
	// Overload: последний файл побеждает. Сначала родители, в конце — .env у cwd, чтобы api-go/.env перекрывал корень репы.
	slices.Reverse(paths)
	_ = godotenv.Overload(paths...)
}

func applySwaggerRuntimeInfo(baseURL string) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("swagger base url parse: %v", err)
		return
	}
	if u.Host != "" {
		docs.SwaggerInfo.Host = u.Host
	}
	if u.Path != "" {
		p := strings.TrimSpace(u.Path)
		if p != "" && p != "/" {
			docs.SwaggerInfo.BasePath = p
		}
	}
	if u.Scheme != "" {
		docs.SwaggerInfo.Schemes = []string{u.Scheme}
	}
}
