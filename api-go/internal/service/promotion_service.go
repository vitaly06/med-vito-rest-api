package service

import (
	"context"
	"errors"
	"fmt"

	"med-vito/api-go/internal/repository"
)

type PromotionService struct {
	repo *repository.PromotionPG
}

func NewPromotionService(repo *repository.PromotionPG) *PromotionService {
	return &PromotionService{repo: repo}
}

func (s *PromotionService) AllPromotions(ctx context.Context) ([]repository.PromotionTariff, error) {
	return s.repo.ListTariffs(ctx)
}

func (s *PromotionService) AddPromotion(ctx context.Context, userID, productID, promotionID, days int32) (map[string]any, error) {
	res, err := s.repo.AddProductPromotion(ctx, userID, productID, promotionID, days)
	if errors.Is(err, repository.ErrPromoProductNotFound) {
		return nil, &AppError{404, "Товар не найден"}
	}
	if errors.Is(err, repository.ErrPromoTariffNotFound) {
		return nil, &AppError{404, "Тип продвижения не найден"}
	}
	if errors.Is(err, repository.ErrPromoNotOwner) {
		return nil, &AppError{403, "Вы не можете продвигать чужой товар"}
	}
	var ins *repository.PromotionInsufficientError
	if errors.As(err, &ins) {
		return nil, &AppError{400, fmt.Sprintf("На балансе недостаточно средств. Требуется: %g₽, доступно: %g₽", ins.Required, ins.Available)}
	}
	if err != nil {
		return nil, err
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
