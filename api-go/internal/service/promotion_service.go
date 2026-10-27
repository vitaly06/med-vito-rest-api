package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"med-vito/api-go/internal/config"
	mailpkg "med-vito/api-go/internal/pkg/mail"
	"med-vito/api-go/internal/pkg/sms"
	"med-vito/api-go/internal/repository"
)

type PromotionService struct {
	repo    *repository.PromotionPG
	support *SupportService
	cfg     config.Config
	sms     *sms.Sender
}

func NewPromotionService(repo *repository.PromotionPG, support *SupportService, cfg config.Config) *PromotionService {
	return &PromotionService{
		repo:    repo,
		support: support,
		cfg:     cfg,
		sms:     sms.NewSender(cfg.MTSBearer, cfg.NotisendAPIKey, cfg.NotisendProject),
	}
}

// StartExpiryWorker запускает периодический воркер проверки истечения платных услуг (7 дней, 3 дня, 1 день).
func (s *PromotionService) StartExpiryWorker(ctx context.Context) {
	go func() {
		log.Println("promotion expiry worker: started")
		_ = s.repo.EnsureNotificationTable(ctx)
		// Первый прогон при старте
		s.runExpiryCheck(ctx)
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("promotion expiry worker: stopped")
				return
			case <-ticker.C:
				s.runExpiryCheck(ctx)
			}
		}
	}()
}

func (s *PromotionService) runExpiryCheck(ctx context.Context) {
	intervals := []struct {
		days int
		text string
	}{
		{days: 7, text: "7 дней"},
		{days: 3, text: "3 дня"},
		{days: 1, text: "1 день"},
	}

	for _, item := range intervals {
		list, err := s.repo.FindPromotionsExpiringInDays(ctx, item.days)
		if err != nil {
			log.Printf("promotion expiry worker error for %d days: %v", item.days, err)
			continue
		}
		for _, promo := range list {
			prodName := promo.ProductName
			if prodName == "" {
				prodName = fmt.Sprintf("ID %d", promo.ProductID)
			}
			smsText := fmt.Sprintf("TorguiSam.ru: До окончания платной услуги продвижения объявления «%s» осталось: %s. Продлите в профиле.", prodName, item.text)

			if promo.Phone != "" && s.sms != nil {
				s.sms.SendAsync(promo.Phone, smsText)
			}

			if s.support != nil && promo.UserID > 0 {
				chatMsg := fmt.Sprintf("⏳ До окончания платного продвижения для объявления «%s» осталось: %s (до %s).\nЧтобы продлить услугу, перейдите в профиль.",
					prodName, item.text, promo.EndDate.Format("02.01.2006 15:04"))
				go s.support.NotifyUserBilling(context.Background(), promo.UserID, chatMsg)
			}

			if promo.Email != "" && strings.TrimSpace(s.cfg.SMTPHost) != "" && !strings.HasSuffix(promo.Email, "@oauth.local") {
				subject := fmt.Sprintf("Платное продвижение истекает через %s — TorguiSam.ru", item.text)
				body := fmt.Sprintf("<p>Здравствуйте!</p><p>Срок действия платного продвижения для вашего объявления <b>«%s»</b> истекает через <b>%s</b> (%s).</p><p>Продлить услугу можно в <a href=\"https://torguisam.ru/profile/my-products\">личном кабинете</a>.</p>",
					prodName, item.text, promo.EndDate.Format("02.01.2006 15:04"))
				go func(to, subj, b string) {
					_ = mailpkg.SendHTMLSmart(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPFrom, to, subj, b, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure)
				}(promo.Email, subject, body)
			}

			_ = s.repo.RecordPromotionNotified(ctx, int(promo.ID), item.days)
		}
	}

	// Деактивация завершившихся продвижений
	deactivated, err := s.repo.DeactivateExpiredPromotions(ctx)
	if err != nil {
		log.Printf("promotion expiry worker deactivation error: %v", err)
	} else if len(deactivated) > 0 {
		log.Printf("promotion expiry worker deactivated %d promotion(s)", len(deactivated))
		for _, d := range deactivated {
			if s.support != nil && d.UserID > 0 {
				msg := fmt.Sprintf("Срок действия платного продвижения для объявления «%s» истёк.", d.ProductName)
				go s.support.NotifyUserBilling(context.Background(), d.UserID, msg)
			}
		}
	}
}


func (s *PromotionService) AllPromotions(ctx context.Context) ([]repository.PromotionTariff, error) {
	return s.repo.ListTariffs(ctx)
}

func (s *PromotionService) TrackTariffView(ctx context.Context, userID *int32) {
	s.repo.InsertFunnelEvent(ctx, userID, "tariff_view", nil, nil)
}

func (s *PromotionService) AddPromotion(ctx context.Context, userID, productID, promotionID, days int32) (map[string]any, error) {
	s.repo.InsertFunnelEvent(ctx, &userID, "tariff_select", &promotionID, &productID)
	res, err := s.repo.AddProductPromotion(ctx, userID, productID, promotionID, days)
	if errors.Is(err, repository.ErrPromoProductNotFound) {
		return nil, &AppError{404, "Товар не найден"}
	}
	if errors.Is(err, repository.ErrPromoTariffNotFound) {
		return nil, &AppError{404, "Тариф продвижения не найден"}
	}
	if errors.Is(err, repository.ErrPromoNotOwner) {
		return nil, &AppError{403, "Нельзя продвигать чужой товар"}
	}
	if errors.Is(err, repository.ErrPromoDowngrade) {
		return nil, &AppError{400, "Понижение типа объявления не предусмотрено: доступен только апгрейд"}
	}
	if errors.Is(err, repository.ErrPromoNotApproved) {
		return nil, &AppError{400, "Продвижение доступно только для одобренных объявлений. Дождитесь прохождения модерации"}
	}
	var ins *repository.PromotionInsufficientError
	if errors.As(err, &ins) {
		return nil, &AppError{400, fmt.Sprintf("Недостаточно средств. Требуется: %g₽, доступно: %g₽", ins.Required, ins.Available)}
	}
	if err != nil {
		return nil, err
	}
	if s.support != nil {
		name := res.ProductName
		if name == "" {
			name = fmt.Sprintf("ID %d", productID)
		}
		tariff := res.TariffName
		if tariff == "" {
			tariff = "Платное продвижение"
		}
		msg := fmt.Sprintf("🚀 Продвижение объявления «%s» успешно активировано!\nТариф: %s\nСписано: %d ₽\nСрок действия: %d дн. (до %s).",
			name, tariff, res.TotalPrice, res.Days, res.EndDate.Format("02.01.2006 15:04"))
		go s.support.NotifyUserBilling(context.Background(), userID, msg)
	}
	return map[string]any{
		"message": "Продвижение успешно активировано",
		"promotion": map[string]any{
			"id":          res.ID,
			"productName": res.ProductName,
			"tariffName":  res.TariffName,
			"days":        res.Days,
			"totalPrice":  res.TotalPrice,
			"startDate":   res.StartDate,
			"endDate":     res.EndDate,
		},
	}, nil
}

