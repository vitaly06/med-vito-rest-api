package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewPG struct {
	pool *pgxpool.Pool
}

func NewReviewPG(pool *pgxpool.Pool) *ReviewPG {
	return &ReviewPG{pool: pool}
}

// IsPgUniqueViolation — нарушение уникального ограничения (дубликат отзыва).
func IsPgUniqueViolation(err error) bool {
	var pe *pgconn.PgError
	return errors.As(err, &pe) && pe.Code == "23505"
}

func (r *ReviewPG) UserExists(ctx context.Context, id int32) (bool, error) {
	var one int
	err := r.pool.QueryRow(ctx, `SELECT 1 FROM "User" WHERE id = $1`, id).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *ReviewPG) InsertReview(ctx context.Context, reviewedByID, reviewedUserID int32, rating float64, text *string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO "Review" ("reviewedById", "reviewedUserId", rating, text, "moderateState", "createdAt")
		VALUES ($1, $2, $3, $4, 'MODERATE'::"ReviewModerate", NOW())`,
		reviewedByID, reviewedUserID, rating, text)
	return err
}

type ApprovedReviewRow struct {
	Rating    float64
	Text      *string
	CreatedAt time.Time
	AuthorID  int32
	FullName  string
}

func (r *ReviewPG) ListApprovedForUser(ctx context.Context, userID int32) ([]ApprovedReviewRow, error) {
	const q = `
		SELECT r.rating, r.text, r."createdAt", u.id, u."fullName"
		FROM "Review" r
		JOIN "User" u ON u.id = r."reviewedById"
		WHERE r."reviewedUserId" = $1 AND r."moderateState" = 'APPROVED'::"ReviewModerate"
		ORDER BY r."createdAt" DESC`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ApprovedReviewRow
	for rows.Next() {
		var row ApprovedReviewRow
		if err := rows.Scan(&row.Rating, &row.Text, &row.CreatedAt, &row.AuthorID, &row.FullName); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

type ModerateReviewRow struct {
	ID             int32
	Rating         float64
	Text           *string
	CreatedAt      time.Time
	ReviewedByID   int32
	ReviewedByName string
	ReviewedUserID int32
	ReviewedName   string
}

func (r *ReviewPG) ListModerateQueue(ctx context.Context) ([]ModerateReviewRow, error) {
	const q = `
		SELECT r.id, r.rating, r.text, r."createdAt",
			b.id, b."fullName", s.id, s."fullName"
		FROM "Review" r
		JOIN "User" b ON b.id = r."reviewedById"
		JOIN "User" s ON s.id = r."reviewedUserId"
		WHERE r."moderateState" = 'MODERATE'::"ReviewModerate"
		ORDER BY r."createdAt" DESC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ModerateReviewRow
	for rows.Next() {
		var row ModerateReviewRow
		if err := rows.Scan(&row.ID, &row.Rating, &row.Text, &row.CreatedAt,
			&row.ReviewedByID, &row.ReviewedByName, &row.ReviewedUserID, &row.ReviewedName); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *ReviewPG) SetReviewModeration(ctx context.Context, reviewID int32, state string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE "Review" SET "moderateState" = $2::"ReviewModerate" WHERE id = $1`,
		reviewID, state)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
