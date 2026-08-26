package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticsPG struct {
	pool *pgxpool.Pool
}

func NewStatisticsPG(pool *pgxpool.Pool) *StatisticsPG {
	return &StatisticsPG{pool: pool}
}

// PeriodStart — граница периода в локальном времени (как new Date() в Nest).
func PeriodStart(period string) (since time.Time, ok bool) {
	switch period {
	case "", "all-time":
		return time.Time{}, false
	case "day":
		now := time.Now().In(time.Local)
		t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return t, true
	case "week":
		return time.Now().In(time.Local).AddDate(0, 0, -7), true
	case "month":
		now := time.Now().In(time.Local)
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), true
	case "quarter":
		now := time.Now().In(time.Local)
		q0 := (int(now.Month()) - 1) / 3 * 3
		return time.Date(now.Year(), time.Month(q0+1), 1, 0, 0, 0, 0, now.Location()), true
	case "half-year":
		now := time.Now().In(time.Local)
		m := time.January
		if now.Month() >= time.July {
			m = time.July
		}
		return time.Date(now.Year(), m, 1, 0, 0, 0, 0, now.Location()), true
	case "year":
		now := time.Now().In(time.Local)
		return time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location()), true
	default:
		return time.Time{}, false
	}
}

func appendProductFilters(sb *strings.Builder, args *[]any, userID int32, categoryID *int32, region *string, productID *int32) int {
	*args = append(*args, userID)
	sb.WriteString(`WHERE p."userId" = $1`)
	idx := 2
	if categoryID != nil {
		fmt.Fprintf(sb, ` AND sc."categoryId" = $%d`, idx)
		*args = append(*args, *categoryID)
		idx++
	}
	if region != nil && strings.TrimSpace(*region) != "" {
		fmt.Fprintf(sb, ` AND p.address ILIKE $%d`, idx)
		*args = append(*args, "%"+strings.TrimSpace(*region)+"%")
		idx++
	}
	if productID != nil {
		fmt.Fprintf(sb, ` AND p.id = $%d`, idx)
		*args = append(*args, *productID)
		idx++
	}
	return idx
}

// UserAnalyticTotals — COUNT просмотров товаров, показов телефона, избранного (логика фильтров как в Nest).
func (r *StatisticsPG) UserAnalyticTotals(ctx context.Context, userID int32, period string, categoryID *int32, region *string, productID *int32) (totalViews, totalPhoneViews, totalFavorites int64, periodOut string, err error) {
	periodOut = period
	if periodOut == "" {
		periodOut = "all-time"
	}
	since, hasSince := PeriodStart(period)

	var args []any
	var sb strings.Builder
	sb.WriteString(`
		SELECT COUNT(*)::bigint FROM "ProductView" pv
		INNER JOIN "Product" p ON p.id = pv."productId"
		LEFT JOIN "SubCategory" sc ON sc.id = p."subCategoryId" `)
	next := appendProductFilters(&sb, &args, userID, categoryID, region, productID)
	if hasSince {
		fmt.Fprintf(&sb, ` AND pv."viewedAt" >= $%d`, next)
		args = append(args, since)
	}
	err = r.pool.QueryRow(ctx, sb.String(), args...).Scan(&totalViews)
	if err != nil {
		return 0, 0, 0, periodOut, err
	}

	args = args[:0]
	sb.Reset()
	sb.WriteString(`SELECT COUNT(*)::bigint FROM "PhoneNumberView" WHERE "viewedUserId" = $1`)
	args = append(args, userID)
	if hasSince {
		sb.WriteString(` AND "viewedAt" >= $2`)
		args = append(args, since)
	}
	err = r.pool.QueryRow(ctx, sb.String(), args...).Scan(&totalPhoneViews)
	if err != nil {
		return 0, 0, 0, periodOut, err
	}

	args = args[:0]
	sb.Reset()
	sb.WriteString(`
		SELECT COUNT(*)::bigint FROM "FavoriteAction" fa
		INNER JOIN "Product" p ON p.id = fa."productId"
		LEFT JOIN "SubCategory" sc ON sc.id = p."subCategoryId" `)
	next = appendProductFilters(&sb, &args, userID, categoryID, region, productID)
	if hasSince {
		fmt.Fprintf(&sb, ` AND fa."addedAt" >= $%d`, next)
		args = append(args, since)
	}
	err = r.pool.QueryRow(ctx, sb.String(), args...).Scan(&totalFavorites)
	if err != nil {
		return 0, 0, 0, periodOut, err
	}
	return totalViews, totalPhoneViews, totalFavorites, periodOut, nil
}

// ProductAnalyticRow — элемент GET products-analytic (как Nest).
type ProductAnalyticRow struct {
	ID               int32    `json:"id"`
	Image            *string  `json:"image"`
	Name             string   `json:"name"`
	Price            int32    `json:"price"`
	Views            int64    `json:"views"`
	FavoritedBy      int64    `json:"favoritedBy"`
	PhoneNumberViews int64    `json:"phoneNumberViews"`
	IsPaid           bool     `json:"isPaid"`
	PromotionLevel   int32    `json:"promotionLevel"`
	PaidViews        int64    `json:"paidViews"`
	Clicks           int64    `json:"clicks"`
	Geo              map[string]int64 `json:"geo"`
}

func (r *StatisticsPG) ProductsAnalytic(ctx context.Context, userID int32) ([]ProductAnalyticRow, error) {
	var phoneTotal int64
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::bigint FROM "PhoneNumberView" WHERE "viewedUserId" = $1`, userID).Scan(&phoneTotal); err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.name, p.price, p.images,
			(SELECT COUNT(*)::bigint FROM "ProductView" pv WHERE pv."productId" = p.id) AS vcount,
			(SELECT COUNT(*)::bigint FROM "_UserFavorites" uf WHERE uf."A" = p.id) AS fcount,
			COALESCE((
				SELECT MAX(pr."pricePerDay")
				FROM "ProductPromotion" pp
				JOIN "Promotion" pr ON pr.id = pp."promotionId"
				WHERE pp."productId" = p.id
				  AND pp."isActive" = true
				  AND pp."isPaid" = true
				  AND pp."endDate" >= NOW()
			), 0)::int AS promotion_level
		FROM "Product" p
		WHERE p."userId" = $1
		ORDER BY p."createdAt" DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ProductAnalyticRow
	for rows.Next() {
		var id int32
		var name string
		var price int32
		var images []string
		var vcount, fcount int64
		var promoLevel int32
		if err := rows.Scan(&id, &name, &price, &images, &vcount, &fcount, &promoLevel); err != nil {
			return nil, err
		}
		var img *string
		if len(images) > 0 && images[0] != "" {
			img = &images[0]
		}
		isPaid := promoLevel > 0
		paidViews := int64(0)
		if isPaid {
			paidViews = vcount
		}
		out = append(out, ProductAnalyticRow{
			ID:               id,
			Image:            img,
			Name:             name,
			Price:            price,
			Views:            vcount,
			FavoritedBy:      fcount,
			PhoneNumberViews: phoneTotal,
			IsPaid:           isPaid,
			PromotionLevel:   promoLevel,
			PaidViews:        paidViews,
			Clicks:           vcount,
			Geo:              map[string]int64{},
		})
	}
	return out, rows.Err()
}

func (r *StatisticsPG) InsertSearchQueryStat(ctx context.Context, userID *int32, query, region, categorySlug, subCategorySlug, typeSlug *string, resultsCount int) {
	if query == nil || strings.TrimSpace(*query) == "" {
		return
	}
	_, _ = r.pool.Exec(ctx, `
		INSERT INTO "SearchQueryStat"
		(query, "region", "categorySlug", "subCategorySlug", "typeSlug", "resultsCount", "userId", "createdAt")
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())`,
		strings.TrimSpace(*query), region, categorySlug, subCategorySlug, typeSlug, resultsCount, userID)
}

func (r *StatisticsPG) UserHasPaidAccess(ctx context.Context, userID int32) (bool, error) {
	var one int
	err := r.pool.QueryRow(ctx, `
		SELECT 1
		FROM "ProductPromotion" pp
		JOIN "Product" p ON p.id = pp."productId"
		WHERE p."userId" = $1
		  AND pp."isActive" = true
		  AND pp."isPaid" = true
		  AND pp."endDate" >= NOW()
		LIMIT 1`, userID).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type SearchQueryStatRow struct {
	Query       string `json:"query"`
	Searches    int64  `json:"searches"`
	AvgResults  int64  `json:"avgResults"`
	LastSearched string `json:"lastSearched"`
}

func (r *StatisticsPG) TopSearchQueries(ctx context.Context, days int) ([]SearchQueryStatRow, error) {
	if days <= 0 {
		days = 30
	}
	rows, err := r.pool.Query(ctx, `
		SELECT query, COUNT(*)::bigint, COALESCE(AVG("resultsCount"),0)::bigint, MAX("createdAt")::text
		FROM "SearchQueryStat"
		WHERE "createdAt" >= NOW() - ($1::text || ' days')::interval
		GROUP BY query
		ORDER BY COUNT(*) DESC, MAX("createdAt") DESC
		LIMIT 100`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SearchQueryStatRow
	for rows.Next() {
		var row SearchQueryStatRow
		if err := rows.Scan(&row.Query, &row.Searches, &row.AvgResults, &row.LastSearched); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *StatisticsPG) AdsTypeDashboard(ctx context.Context, userID int32) (map[string]any, error) {
	var freeCount, paidCount int64
	if err := r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE COALESCE((
				SELECT MAX(pr."pricePerDay")
				FROM "ProductPromotion" pp
				JOIN "Promotion" pr ON pr.id = pp."promotionId"
				WHERE pp."productId" = p.id AND pp."isActive" = true AND pp."isPaid" = true AND pp."endDate" >= NOW()
			),0) = 0)::bigint AS free_count,
			COUNT(*) FILTER (WHERE COALESCE((
				SELECT MAX(pr."pricePerDay")
				FROM "ProductPromotion" pp
				JOIN "Promotion" pr ON pr.id = pp."promotionId"
				WHERE pp."productId" = p.id AND pp."isActive" = true AND pp."isPaid" = true AND pp."endDate" >= NOW()
			),0) > 0)::bigint AS paid_count
		FROM "Product" p
		WHERE p."userId" = $1`, userID).Scan(&freeCount, &paidCount); err != nil {
		return nil, err
	}
	return map[string]any{"free": freeCount, "paid": paidCount, "total": freeCount + paidCount}, nil
}

func (r *StatisticsPG) TariffFunnel(ctx context.Context, userID int32, days int) (map[string]int64, error) {
	if days <= 0 {
		days = 30
	}
	rows, err := r.pool.Query(ctx, `
		SELECT step, COUNT(*)::bigint
		FROM "TariffFunnelEvent"
		WHERE ("userId" = $1 OR "userId" IS NULL)
		  AND "createdAt" >= NOW() - ($2::text || ' days')::interval
		GROUP BY step`, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int64{
		"tariff_view":   0,
		"tariff_select": 0,
		"payment":       0,
		"publication":   0,
	}
	for rows.Next() {
		var step string
		var cnt int64
		if err := rows.Scan(&step, &cnt); err != nil {
			return nil, err
		}
		out[step] = cnt
	}
	return out, rows.Err()
}

func (r *StatisticsPG) RevenueByTypeAndCategory(ctx context.Context, days int) ([]map[string]any, error) {
	if days <= 0 {
		days = 30
	}
	rows, err := r.pool.Query(ctx, `
		SELECT pr.name, c.name, SUM(pp."totalPrice")::bigint
		FROM "ProductPromotion" pp
		JOIN "Promotion" pr ON pr.id = pp."promotionId"
		JOIN "Product" p ON p.id = pp."productId"
		JOIN "Category" c ON c.id = p."categoryId"
		WHERE pp."isPaid" = true
		  AND pp."createdAt" >= NOW() - ($1::text || ' days')::interval
		GROUP BY pr.name, c.name
		ORDER BY SUM(pp."totalPrice") DESC`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]any
	for rows.Next() {
		var promo, cat string
		var sum int64
		if err := rows.Scan(&promo, &cat, &sum); err != nil {
			return nil, err
		}
		out = append(out, map[string]any{"promotionType": promo, "category": cat, "revenue": sum})
	}
	return out, rows.Err()
}

func (r *StatisticsPG) SystemOverallStats(ctx context.Context, days int) (map[string]any, error) {
	var totalProducts, activeProducts, paidProducts, freeProducts int64
	var totalUsers, totalDeals int64
	var totalRevenue float64

	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM "Product"`).Scan(&totalProducts)
	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM "Product" WHERE "moderateState" = 'APPROVED' AND "isHide" = false`).Scan(&activeProducts)
	_ = r.pool.QueryRow(ctx, `
		SELECT COUNT(DISTINCT "productId")
		FROM "ProductPromotion"
		WHERE "isActive" = true AND "isPaid" = true AND "endDate" >= NOW()
	`).Scan(&paidProducts)

	freeProducts = activeProducts - paidProducts
	if freeProducts < 0 {
		freeProducts = 0
	}

	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM "User"`).Scan(&totalUsers)
	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM "Deal"`).Scan(&totalDeals)
	_ = r.pool.QueryRow(ctx, `SELECT COALESCE(SUM("totalPrice"), 0) FROM "ProductPromotion" WHERE "isPaid" = true`).Scan(&totalRevenue)

	return map[string]any{
		"totalProducts":  totalProducts,
		"activeProducts": activeProducts,
		"paidProducts":   paidProducts,
		"freeProducts":   freeProducts,
		"totalUsers":     totalUsers,
		"totalDeals":     totalDeals,
		"totalRevenue":   totalRevenue,
	}, nil
}
