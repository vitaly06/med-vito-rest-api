package service

import (
	"context"
	"errors"
	"fmt"

	"med-vito/api-go/internal/repository"
)

type PromotionService struct {
	repo    *repository.PromotionPG
	support *SupportService
}

func NewPromotionService(repo *repository.PromotionPG, support *SupportService) *PromotionService {
	return &PromotionService{repo: repo, support: support}
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
	var ins *repository.PromotionInsufficientError
	if errors.As(err, &ins) {
		return nil, &AppError{400, fmt.Sprintf("Недостаточно средств. Требуется: %g₽, доступно: %g₽", ins.Required, ins.Available)}
	}
	if err != nil {
		return nil, err
	}
	if s.support != nil {
		msg := fmt.Sprintf("🚀 Продвижение объявления успешно активировано!\nСписано: %d ₽\nСрок действия: %d дней (до %s).",
			res.TotalPrice, res.Days, res.EndDate.Format("02.01.2006 15:04"))
		go s.support.NotifyUserBilling(context.Background(), userID, msg)
	}
	return map[string]any{
		"message": "Продвижение успешно активировано",
		"promotion": map[string]any{
			"id":         res.ID,
			"days":       res.Days,
			"totalPrice": res.TotalPrice,
			"startDate":  res.StartDate,
			"endDate":    res.EndDate,
		},
	}, nil
}

