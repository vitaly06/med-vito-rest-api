package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"med-vito/api-go/internal/config"
	mailpkg "med-vito/api-go/internal/pkg/mail"
	"med-vito/api-go/internal/repository"
)

// ProductExpiryWorker — фоновый воркер, деактивирующий объявления старше 30 дней.
type ProductExpiryWorker struct {
	cfg          config.Config
	repo         *repository.ProductPG
	support      *SupportService
	httpClient   *http.Client
	pollInterval time.Duration
}

func NewProductExpiryWorker(cfg config.Config, repo *repository.ProductPG, support *SupportService) *ProductExpiryWorker {
	return &ProductExpiryWorker{
		cfg:          cfg,
		repo:         repo,
		support:      support,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
		pollInterval: 1 * time.Hour,
	}
}

// StartExpiryWorker запускает воркер в фоновой горутине.
func (w *ProductExpiryWorker) StartExpiryWorker(ctx context.Context) {
	go func() {
		log.Println("product expiry worker: started")
		// Первый прогон при запуске сервера.
		w.runOnce(ctx)
		ticker := time.NewTicker(w.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("product expiry worker: stopped")
				return
			case <-ticker.C:
				w.runOnce(ctx)
			}
		}
	}()
}

func (w *ProductExpiryWorker) runOnce(ctx context.Context) {
	w.notifyExpiringSoon(ctx)
	w.deactivateExpired(ctx)
}

func (w *ProductExpiryWorker) deactivateExpired(ctx context.Context) {
	expired, err := w.repo.ExpireProducts(ctx)
	if err != nil {
		log.Printf("product expiry worker: ExpireProducts error: %v", err)
		return
	}
	if len(expired) == 0 {
		return
	}
	log.Printf("product expiry worker: deactivated %d expired product(s)", len(expired))
	for _, p := range expired {
		log.Printf("product expiry worker: deactivated product #%d «%s» (userId=%d)", p.ID, p.Name, p.UserID)
		msg := fmt.Sprintf(
			"⏰ Ваше объявление «%s» (ID %d) автоматически снято с публикации — истёк 30-дневный срок размещения.\n\n"+
				"Чтобы снова опубликовать объявление, перейдите в личный кабинет → «Мои объявления» → «Скрытые» и нажмите «Опубликовать снова»"+
				" (или оформите продвижение для продления срока).",
			p.Name, p.ID,
		)
		if w.support != nil {
			go w.support.NotifyUserBilling(context.Background(), p.UserID, msg)
		}
	}
}

func (w *ProductExpiryWorker) notifyExpiringSoon(ctx context.Context) {
	expiring, err := w.repo.ProductsExpiringIn3Days(ctx)
	if err != nil {
		log.Printf("product expiry worker: ProductsExpiringIn3Days error: %v", err)
		return
	}
	for _, p := range expiring {
		log.Printf("product expiry worker: expiring soon product #%d «%s» (userId=%d)", p.ID, p.Name, p.UserID)
		msg := fmt.Sprintf(
			"⚠️ Ваше объявление «%s» (ID %d) будет автоматически снято с публикации через 3 дня — истекает 30-дневный срок.\n\n"+
				"Чтобы продлить размещение, оформите продвижение или опубликуйте объявление заново после снятия.",
			p.Name, p.ID,
		)
		if w.support != nil {
			go w.support.NotifyUserBilling(context.Background(), p.UserID, msg)
		}
		if phone := strings.TrimSpace(p.Phone); phone != "" {
			go w.sendSMS(phone, fmt.Sprintf(
				"TorguiSam.ru: Объявление «%s» будет снято через 3 дня. Продлите в личном кабинете.",
				productExpiryTruncate(p.Name, 40),
			))
		}
		if email := strings.TrimSpace(p.Email); email != "" && !strings.HasSuffix(email, "@oauth.local") {
			go w.sendMail(email,
				"Объявление будет снято через 3 дня — TorguiSam.ru",
				fmt.Sprintf(
					"<p>Здравствуйте!</p>"+
						"<p>Ваше объявление <b>«%s»</b> (ID %d) будет автоматически снято с публикации через 3 дня — истекает 30-дневный срок размещения.</p>"+
						"<p>Чтобы продлить размещение, <a href=\"https://torguisam.ru/profile/my-products\">перейдите в личный кабинет</a>"+
						" и оформите продвижение или опубликуйте объявление заново после снятия.</p>"+
						"<p>С уважением, команда TorguiSam.ru</p>",
					p.Name, p.ID,
				),
			)
		}
	}
}

func (w *ProductExpiryWorker) sendSMS(phone, message string) {
	phone = strings.TrimSpace(phone)
	message = strings.TrimSpace(message)
	if phone == "" || message == "" {
		return
	}
	if bearer := strings.TrimSpace(w.cfg.MTSBearer); bearer != "" {
		body := map[string]any{
			"submits": []any{map[string]any{"msid": phone, "message": message}},
			"naming":  "Torguisamru",
		}
		b, _ := json.Marshal(body)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
			"https://api.mts.ru/client-omni-adapter_production/1.0.2/mcom/messageManagement/messages",
			bytes.NewReader(b))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+bearer)
			resp, err := w.httpClient.Do(req)
			if err != nil {
				log.Printf("product expiry worker: sms(mts) failed phone=%s err=%v", phone, err)
			} else {
				_ = resp.Body.Close()
			}
		}
		return
	}
	if apiKey := strings.TrimSpace(w.cfg.NotisendAPIKey); apiKey != "" {
		u := "https://sms.notisend.ru/api/message/send?project=" + url.QueryEscape(w.cfg.NotisendProject) +
			"&message=" + url.QueryEscape(message) +
			"&recipients=" + url.QueryEscape(phone) +
			"&apikey=" + url.QueryEscape(apiKey)
		resp, err := w.httpClient.Get(u)
		if err != nil {
			log.Printf("product expiry worker: sms(notisend) failed phone=%s err=%v", phone, err)
			return
		}
		_ = resp.Body.Close()
	}
}

func (w *ProductExpiryWorker) sendMail(toEmail, subject, html string) {
	toEmail = strings.TrimSpace(toEmail)
	if toEmail == "" || strings.TrimSpace(w.cfg.SMTPHost) == "" {
		return
	}
	if err := mailpkg.SendHTMLSmart(
		w.cfg.SMTPHost, w.cfg.SMTPPort, w.cfg.SMTPUser, w.cfg.SMTPPassword,
		w.cfg.SMTPFrom, toEmail, subject, html, w.cfg.SMTPSecure, w.cfg.SMTPTLSInsecure,
	); err != nil {
		log.Printf("product expiry worker: mail failed to=%s err=%v", toEmail, err)
	}
}

func productExpiryTruncate(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "…"
}
