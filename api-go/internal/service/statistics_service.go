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

func (s *StatisticsService) TrackSearch(ctx context.Context, userID *int32, query, region, categorySlug, subCategorySlug, typeSlug *string, resultsCount int) {
	s.repo.InsertSearchQueryStat(ctx, userID, query, region, categorySlug, subCategorySlug, typeSlug, resultsCount)
}

func (s *StatisticsService) SearchQueriesInsights(ctx context.Context, userID int32, days int) (map[string]any, error) {
	ok, err := s.repo.UserHasPaidAccess(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &AppError{403, "Доступ к статистике поисковых запросов доступен только при активном платном размещении"}
	}
	rows, err := s.repo.TopSearchQueries(ctx, days)
	if err != nil {
		return nil, err
	}
	return map[string]any{"days": days, "items": rows}, nil
}

func (s *StatisticsService) CabinetDashboard(ctx context.Context, userID int32, days int) (map[string]any, error) {
	types, err := s.repo.AdsTypeDashboard(ctx, userID)
	if err != nil {
		return nil, err
	}
	funnel, err := s.repo.TariffFunnel(ctx, userID, days)
	if err != nil {
		return nil, err
	}
	revenue, err := s.repo.RevenueByTypeAndCategory(ctx, days)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"days":              days,
		"adsTypes":          types,
		"tariffFunnel":      funnel,
		"revenueBreakdown":  revenue,
		"tariffClickHeatmap": map[string]any{
			"view_to_select": funnel["tariff_select"],
			"select_to_pay":  funnel["payment"],
			"pay_to_publish": funnel["publication"],
		},
	}, nil
}

func (s *StatisticsService) SystemPlatformDashboard(ctx context.Context, days int) (map[string]any, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.SystemOverallStats(ctx, days)
}
