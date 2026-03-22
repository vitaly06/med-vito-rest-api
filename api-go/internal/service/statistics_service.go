package service

import (
	"context"
	"strings"

	"med-vito/api-go/internal/repository"
)

type StatisticsService struct {
	repo *repository.StatisticsPG
}

func NewStatisticsService(repo *repository.StatisticsPG) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func normalizeStatsPeriod(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", nil
	}
	switch p {
	case "day", "week", "month", "quarter", "half-year", "year":
		return p, nil
	default:
		return "", &AppError{400, "Некорректный period"}
	}
}

func (s *StatisticsService) GetUserStatistics(ctx context.Context, userID int32, period string, categoryID *int32, region *string, productID *int32) (map[string]any, error) {
	p, err := normalizeStatsPeriod(period)
	if err != nil {
		return nil, err
	}
	v, ph, fav, pOut, err := s.repo.UserAnalyticTotals(ctx, userID, p, categoryID, region, productID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"period":          pOut,
		"totalViews":      v,
		"totalPhoneViews": ph,
		"totalFavorites":  fav,
	}, nil
}

func (s *StatisticsService) GetProductsAnalytic(ctx context.Context, userID int32) ([]repository.ProductAnalyticRow, error) {
	return s.repo.ProductsAnalytic(ctx, userID)
}
