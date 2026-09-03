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

