package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"med-vito/api-go/internal/config"
	mailpkg "med-vito/api-go/internal/pkg/mail"
	"med-vito/api-go/internal/repository"
)

type ReservationService struct {
	cfg  config.Config
	repo *repository.ReservationPG
}

func NewReservationService(cfg config.Config, repo *repository.ReservationPG, users *repository.UserPG) *ReservationService {
	_ = users
	return &ReservationService{cfg: cfg, repo: repo}
}

type CreateReservationRequest struct {
	ProductID int32   `json:"productId"`
	Hours     *int32  `json:"hours"`
	Note      *string `json:"note"`
}

type UpdateProductReservationSettingsRequest struct {
	ProductID         int32 `json:"productId"`
	AllowReservations bool  `json:"allowReservations"`
}

func (s *ReservationService) Create(ctx context.Context, buyerID int32, req CreateReservationRequest) (map[string]any, error) {
	if req.ProductID <= 0 {
		return nil, &AppError{400, "Нужен productId"}
	}
	pr, err := s.repo.GetProductReserveSettings(ctx, req.ProductID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "Товар не найден"}
	}
	if err != nil {
		return nil, err
	}
	if pr.SellerID == buyerID {
		return nil, &AppError{400, "Нельзя резервировать свой товар"}
	}
	if pr.ModerateState != "APPROVED" || pr.IsHide {
		return nil, &AppError{400, "Товар недоступен для резерва"}
	}
	if !pr.AllowReservations {
		return nil, &AppError{400, "Продавец отключил резервирование для этого объявления"}
	}
	verified, buyerName, err := s.repo.IsUserVerified(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, &AppError{403, "Резерв доступен только верифицированным пользователям"}
	}

	now := time.Now()
	blocked, until, err := s.repo.HasReserveBlock(ctx, buyerID, now)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, &AppError{403, "Функция резерва временно заблокирована до " + until.Format(time.RFC3339)}
	}
	activeCount, err := s.repo.ActiveReservationsCount(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	maxActive := s.cfg.ReservationMaxActive
	if maxActive <= 0 {
		maxActive = 5
	}
	if activeCount >= maxActive {
		return nil, &AppError{400, "Превышен лимит активных резервов"}
	}
	dayCount, err := s.repo.DailyReservationsCount(ctx, buyerID, now.Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}
	dailyLimit := s.cfg.ReservationDailyLimit
	if dailyLimit <= 0 {
		dailyLimit = 4
	}
	if dayCount >= dailyLimit {
		return nil, &AppError{400, "Превышен суточный лимит резервирований"}
	}
	existing, err := s.repo.ActiveReservationByProduct(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &AppError{400, "Товар уже зарезервирован"}
	}

	hours := int32(24)
	note := normalizeStringPtr(req.Note)
	row, err := s.repo.CreateReservation(ctx, req.ProductID, buyerID, pr.SellerID, hours, note, now)
	if err != nil {
		return nil, err
	}
	go s.sendReserveEmail(pr.SellerEmail, pr.ProductName, buyerName, row.ExpiresAt)
	return s.format(row), nil
}

func (s *ReservationService) MyList(ctx context.Context, userID int32) ([]map[string]any, error) {
	rows, err := s.repo.ListForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for i := range rows {
		out = append(out, s.format(&rows[i]))
	}
	return out, nil
}

func (s *ReservationService) ProductReservationInfo(ctx context.Context, productID int32) (map[string]any, error) {
	row, err := s.repo.PublicActiveReservationByProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return map[string]any{"isReserved": false}, nil
	}
	return map[string]any{
		"isReserved":  true,
		"reservation": s.format(row),
	}, nil
}

func (s *ReservationService) CancelByBuyer(ctx context.Context, userID int32, reservationID int64) (map[string]any, error) {
	if err := s.repo.CancelByBuyer(ctx, reservationID, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Активный резерв не найден"}
		}
		return nil, err
	}
	row, err := s.repo.GetByID(ctx, reservationID)
	if err != nil {
		return nil, err
	}
	return s.format(row), nil
}

func (s *ReservationService) CancelBySeller(ctx context.Context, userID int32, reservationID int64, reason string) (map[string]any, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, &AppError{400, "Нужно указать причину отмены"}
	}
	if err := s.repo.CancelBySeller(ctx, reservationID, userID, strings.TrimSpace(reason)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Активный резерв не найден"}
		}
		return nil, err
	}
	row, err := s.repo.GetByID(ctx, reservationID)
	if err != nil {
		return nil, err
	}
	return s.format(row), nil
}

func (s *ReservationService) Cancel(ctx context.Context, userID int32, reservationID int64, sellerReason *string) (map[string]any, error) {
	row, err := s.repo.FindParticipantReservation(ctx, reservationID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Резерв не найден"}
		}
		return nil, err
	}
	if row.Status != "ACTIVE" {
		return nil, &AppError{400, "Отменить можно только активный резерв"}
	}
	if row.BuyerID == userID {
		return s.CancelByBuyer(ctx, userID, reservationID)
	}
	reason := "Отмена продавцом"
	if sellerReason != nil && strings.TrimSpace(*sellerReason) != "" {
		reason = strings.TrimSpace(*sellerReason)
	}
	if row.SellerID == userID {
		return s.CancelBySeller(ctx, userID, reservationID, reason)
	}
	return nil, &AppError{403, "Нет доступа к этому резерву"}
}

func (s *ReservationService) Extend(ctx context.Context, userID int32, reservationID int64) (map[string]any, error) {
	_ = ctx
	_ = userID
	_ = reservationID
	return nil, &AppError{400, "Продление резерва отключено: срок всегда 24 часа"}
}

func (s *ReservationService) StartWorker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				n, err := s.repo.ExpireDue(ctx, time.Now())
				if err != nil {
					log.Printf("reservation worker: %v", err)
					continue
				}
				if n > 0 {
					log.Printf("reservation worker: expired %d reservation(s)", n)
				}
				_ = s.applyPenaltyForViolations(ctx)
			}
		}
	}()
}

func (s *ReservationService) UpdateProductSettings(ctx context.Context, sellerID int32, req UpdateProductReservationSettingsRequest) error {
	if req.ProductID <= 0 {
		return &AppError{400, "Нужен productId"}
	}
	if err := s.repo.UpdateProductReserveSettings(ctx, req.ProductID, sellerID, req.AllowReservations, 24); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &AppError{404, "Товар не найден или не принадлежит продавцу"}
		}
		return err
	}
	return nil
}

func (s *ReservationService) applyPenaltyForViolations(ctx context.Context) error {
	buyers, err := s.repo.RecentExpiredBuyers(ctx, time.Now().Add(-48*time.Hour))
	if err != nil {
		return err
	}
	for _, buyerID := range buyers {
		statuses, err := s.repo.LastBuyerStatuses(ctx, buyerID, 3)
		if err != nil {
			return err
		}
		if len(statuses) < 3 {
			continue
		}
		if statuses[0] == "EXPIRED" && statuses[1] == "EXPIRED" && statuses[2] == "EXPIRED" {
			days := s.cfg.ReservationBlockDays
			if days <= 0 {
				days = 7
			}
			_ = s.repo.SetReserveBlock(ctx, buyerID, time.Now().AddDate(0, 0, days), "3 неоплаченных резерва подряд")
		}
	}
	return nil
}

func (s *ReservationService) format(r *repository.ReservationRow) map[string]any {
	return map[string]any{
		"id":           r.ID,
		"productId":    r.ProductID,
		"productName":  r.ProductName,
		"buyerId":      r.BuyerID,
		"buyerName":    r.BuyerName,
		"sellerId":     r.SellerID,
		"sellerName":   r.SellerName,
		"status":       localizeReservationStatus(r.Status),
		"statusCode":   r.Status,
		"hours":        r.Hours,
		"note":         r.Note,
		"cancelReason": r.CancelReason,
		"extendedOnce": r.ExtendedOnce,
		"expiresAt":    r.ExpiresAt,
		"createdAt":    r.CreatedAt,
		"cancelledAt":  r.CancelledAt,
	}
}

func (s *ReservationService) sendReserveEmail(toEmail, productName, buyerName string, expiresAt time.Time) {
	if strings.TrimSpace(toEmail) == "" || strings.TrimSpace(s.cfg.SMTPHost) == "" {
		return
	}
	body := "<p>Товар <b>" + productName + "</b> зарезервирован покупателем " + buyerName + " до " + expiresAt.Format("2006-01-02 15:04") + ".</p>"
	_ = mailpkg.SendHTMLSmart(
		s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword,
		s.cfg.SMTPFrom, toEmail, "Новый резерв товара", body, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure,
	)
}

