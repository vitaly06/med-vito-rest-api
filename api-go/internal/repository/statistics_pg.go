package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

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
}

func (r *StatisticsPG) ProductsAnalytic(ctx context.Context, userID int32) ([]ProductAnalyticRow, error) {
	var phoneTotal int64
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::bigint FROM "PhoneNumberView" WHERE "viewedUserId" = $1`, userID).Scan(&phoneTotal); err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT p.id, p.name, p.price, p.images,
			(SELECT COUNT(*)::bigint FROM "ProductView" pv WHERE pv."productId" = p.id) AS vcount,
			(SELECT COUNT(*)::bigint FROM "_UserFavorites" uf WHERE uf."A" = p.id) AS fcount
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
		if err := rows.Scan(&id, &name, &price, &images, &vcount, &fcount); err != nil {
			return nil, err
		}
		var img *string
		if len(images) > 0 && images[0] != "" {
			img = &images[0]
		}
		out = append(out, ProductAnalyticRow{
			ID:               id,
			Image:            img,
			Name:             name,
			Price:            price,
			Views:            vcount,
			FavoritedBy:      fcount,
			PhoneNumberViews: phoneTotal,
		})
	}
	return out, rows.Err()
}
