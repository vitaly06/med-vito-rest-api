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

var ErrBannerNotFound = errors.New("banner: not found")

type BannerPG struct {
	pool *pgxpool.Pool
}

func NewBannerPG(pool *pgxpool.Pool) *BannerPG {
	return &BannerPG{pool: pool}
}

// BannerRow — строка "Banner".
type BannerRow struct {
	ID             int32     `json:"id"`
	Name           string    `json:"name"`
	PhotoURL       string    `json:"photoUrl"`
	Place          string    `json:"place"`
	ModerateState  string    `json:"moderateState,omitempty"`
	NavigateToURL  string    `json:"navigateToUrl"`
	UserID         int32     `json:"userId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (r *BannerPG) Insert(ctx context.Context, name, photoURL, place, moderate, navigate string, userID int32) (*BannerRow, error) {
	var b BannerRow
	err := r.pool.QueryRow(ctx, `
		INSERT INTO "Banner" (name, "photoUrl", place, "moderateState", "navigateToUrl", "userId", "createdAt", "updatedAt")
		VALUES ($1, $2, $3::"BannerPlace", $4::"BannerModerate", $5, $6, NOW(), NOW())
		RETURNING id, name, "photoUrl", place::text, "moderateState"::text, "navigateToUrl", "userId", "createdAt", "updatedAt"`,
		name, photoURL, place, moderate, navigate, userID,
	).Scan(&b.ID, &b.Name, &b.PhotoURL, &b.Place, &b.ModerateState, &b.NavigateToURL, &b.UserID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BannerPG) FindAllApproved(ctx context.Context) ([]BannerRow, error) {
	return r.queryBanners(ctx, `WHERE "moderateState" = 'APPROVED'::"BannerModerate" ORDER BY "createdAt" DESC`, nil)
}

func (r *BannerPG) FindByPlaceApproved(ctx context.Context, place string) ([]BannerRow, error) {
	return r.queryBanners(ctx, `WHERE place = $1::"BannerPlace" AND "moderateState" = 'APPROVED'::"BannerModerate" ORDER BY "createdAt" DESC`, []any{place})
}

func (r *BannerPG) queryBanners(ctx context.Context, where string, args []any) ([]BannerRow, error) {
	q := `SELECT id, name, "photoUrl", place::text, "moderateState"::text, "navigateToUrl", "userId", "createdAt", "updatedAt" FROM "Banner" ` + where
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBannerRows(rows)
}

func scanBannerRows(rows pgx.Rows) ([]BannerRow, error) {
	var out []BannerRow
	for rows.Next() {
		var b BannerRow
		if err := rows.Scan(&b.ID, &b.Name, &b.PhotoURL, &b.Place, &b.ModerateState, &b.NavigateToURL, &b.UserID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BannerPG) FindRandomApproved(ctx context.Context, limit int) ([]BannerRow, error) {
	if limit < 1 {
		limit = 5
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, "photoUrl", place::text, "moderateState"::text, "navigateToUrl", "userId", "createdAt", "updatedAt"
		FROM "Banner"
		WHERE "moderateState" = 'APPROVED'::"BannerModerate"
		ORDER BY RANDOM()
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBannerRows(rows)
}

func (r *BannerPG) GetByID(ctx context.Context, id int32) (*BannerRow, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, "photoUrl", place::text, "moderateState"::text, "navigateToUrl", "userId", "createdAt", "updatedAt"
		FROM "Banner" WHERE id = $1`, id)
	var b BannerRow
	err := row.Scan(&b.ID, &b.Name, &b.PhotoURL, &b.Place, &b.ModerateState, &b.NavigateToURL, &b.UserID, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrBannerNotFound
	}
	return &b, err
}

func (r *BannerPG) Update(ctx context.Context, id int32, photoURL *string, place, navigateToURL, name *string) (*BannerRow, error) {
	sets := []string{`"updatedAt" = NOW()`}
	args := []any{id}
	n := 2
	if photoURL != nil {
		sets = append(sets, fmt.Sprintf(`"photoUrl" = $%d`, n))
		args = append(args, *photoURL)
		n++
	}
	if place != nil && *place != "" {
		sets = append(sets, fmt.Sprintf(`place = $%d::"BannerPlace"`, n))
		args = append(args, *place)
		n++
	}
	if navigateToURL != nil {
		sets = append(sets, fmt.Sprintf(`"navigateToUrl" = $%d`, n))
		args = append(args, *navigateToURL)
		n++
	}
	if name != nil && *name != "" {
		sets = append(sets, fmt.Sprintf(`name = $%d`, n))
		args = append(args, *name)
		n++
	}
	if len(sets) == 1 {
		return r.GetByID(ctx, id)
	}
	q := fmt.Sprintf(`UPDATE "Banner" SET %s WHERE id = $1`, strings.Join(sets, ", "))
	cmd, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, ErrBannerNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *BannerPG) Delete(ctx context.Context, id int32) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM "Banner" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrBannerNotFound
	}
	return nil
}

func (r *BannerPG) SetModerateState(ctx context.Context, id int32, state string) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE "Banner" SET "moderateState" = $2::"BannerModerate", "updatedAt" = NOW() WHERE id = $1`, id, state)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrBannerNotFound
	}
	return nil
}

// BannerModerateRow — очередь модерации с автором.
type BannerModerateRow struct {
	ID            int32     `json:"id"`
	Name          string    `json:"name"`
	PhotoURL      string    `json:"photoUrl"`
	Place         string    `json:"place"`
	NavigateToURL string    `json:"navigateToUrl"`
	CreatedAt     time.Time `json:"createdAt"`
	User          struct {
		ID       int32   `json:"id"`
		FullName string  `json:"fullName"`
		Photo    *string `json:"photo,omitempty"`
	} `json:"user"`
}

func (r *BannerPG) ListModerateQueue(ctx context.Context) ([]BannerModerateRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT b.id, b.name, b."photoUrl", b.place::text, b."navigateToUrl", b."createdAt",
			u.id, u."fullName", u.photo
		FROM "Banner" b
		JOIN "User" u ON u.id = b."userId"
		WHERE b."moderateState" = 'MODERATE'::"BannerModerate"
		ORDER BY b."createdAt" DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []BannerModerateRow
	for rows.Next() {
		var row BannerModerateRow
		if err := rows.Scan(&row.ID, &row.Name, &row.PhotoURL, &row.Place, &row.NavigateToURL, &row.CreatedAt,
			&row.User.ID, &row.User.FullName, &row.User.Photo); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

// BannerViewRow — ответ registerView.
type BannerViewRow struct {
	ID        int32     `json:"id"`
	BannerID  int32     `json:"bannerId"`
	UserID    *int32    `json:"userId"`
	IPAddress *string   `json:"ipAddress"`
	ViewedAt  time.Time `json:"viewedAt"`
}

func (r *BannerPG) FindRecentView(ctx context.Context, bannerID int32, since time.Time, userID *int32, ip string) (*BannerViewRow, error) {
	var q strings.Builder
	args := []any{bannerID, since}
	q.WriteString(`SELECT id, "bannerId", "userId", "ipAddress", "viewedAt" FROM "BannerView" WHERE "bannerId" = $1 AND "viewedAt" >= $2`)
	n := 3
	if userID != nil && strings.TrimSpace(ip) != "" {
		q.WriteString(fmt.Sprintf(` AND ("userId" = $%d OR "ipAddress" = $%d)`, n, n+1))
		args = append(args, *userID, strings.TrimSpace(ip))
	} else if userID != nil {
		q.WriteString(fmt.Sprintf(` AND "userId" = $%d`, n))
		args = append(args, *userID)
	} else if strings.TrimSpace(ip) != "" {
		q.WriteString(fmt.Sprintf(` AND "ipAddress" = $%d`, n))
		args = append(args, strings.TrimSpace(ip))
	}
	q.WriteString(` ORDER BY "viewedAt" DESC LIMIT 1`)
	row := r.pool.QueryRow(ctx, q.String(), args...)
	var v BannerViewRow
	err := row.Scan(&v.ID, &v.BannerID, &v.UserID, &v.IPAddress, &v.ViewedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *BannerPG) InsertView(ctx context.Context, bannerID int32, userID *int32, ip *string) (*BannerViewRow, error) {
	var v BannerViewRow
	err := r.pool.QueryRow(ctx, `
		INSERT INTO "BannerView" ("bannerId", "userId", "ipAddress", "viewedAt")
		VALUES ($1, $2, $3, NOW())
		RETURNING id, "bannerId", "userId", "ipAddress", "viewedAt"`,
		bannerID, userID, ip,
	).Scan(&v.ID, &v.BannerID, &v.UserID, &v.IPAddress, &v.ViewedAt)
	return &v, err
}

func (r *BannerPG) CountViews(ctx context.Context, bannerID int32) (int64, error) {
	var n int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::bigint FROM "BannerView" WHERE "bannerId" = $1`, bannerID).Scan(&n)
	return n, err
}

func (r *BannerPG) CountDistinctViewUsers(ctx context.Context, bannerID int32) (int64, error) {
	var n int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(DISTINCT "userId")::bigint FROM "BannerView" WHERE "bannerId" = $1 AND "userId" IS NOT NULL`, bannerID).Scan(&n)
	return n, err
}

type BannerViewsByDay struct {
	Date  string `json:"date"`
	Views int64  `json:"views"`
}

func (r *BannerPG) ViewsByDay(ctx context.Context, bannerID int32, since time.Time) ([]BannerViewsByDay, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT ("viewedAt"::date)::text, COUNT(*)::bigint
		FROM "BannerView"
		WHERE "bannerId" = $1 AND "viewedAt" >= $2
		GROUP BY "viewedAt"::date
		ORDER BY ("viewedAt"::date) DESC`, bannerID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []BannerViewsByDay
	for rows.Next() {
		var d BannerViewsByDay
		if err := rows.Scan(&d.Date, &d.Views); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// UserBannerStat — как getUserBannersStats map в Nest.
type UserBannerStat struct {
	ID         int32     `json:"id"`
	Name       string    `json:"name"`
	Place      string    `json:"place"`
	TotalViews int64     `json:"totalViews"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (r *BannerPG) UserBannersStats(ctx context.Context, userID int32) ([]UserBannerStat, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT b.id, b.name, b.place::text, b."createdAt",
			(SELECT COUNT(*)::bigint FROM "BannerView" v WHERE v."bannerId" = b.id) AS vc
		FROM "Banner" b
		WHERE b."userId" = $1
		ORDER BY b."createdAt" DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []UserBannerStat
	for rows.Next() {
		var s UserBannerStat
		if err := rows.Scan(&s.ID, &s.Name, &s.Place, &s.CreatedAt, &s.TotalViews); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
