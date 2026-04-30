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

func (r *ReviewPG) HasDealOrReservation(ctx context.Context, authorID, sellerID int32) (bool, error) {
	var one int
	err := r.pool.QueryRow(ctx, `
		SELECT 1
		FROM "ProductDeal" d
		WHERE d."buyerId" = $1
		  AND d."sellerId" = $2
		  AND d.status IN ('PAID'::"DealStatus",'SHIPPED'::"DealStatus",'DELIVERED'::"DealStatus",'COMPLETED'::"DealStatus")
		LIMIT 1`, authorID, sellerID).Scan(&one)
	if err == nil {
		return true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}
	err = r.pool.QueryRow(ctx, `
		SELECT 1
		FROM "ProductReservation" r
		WHERE r."buyerId" = $1
		  AND r."sellerId" = $2
		  AND r.status IN ('DEAL_CREATED'::"ReservationStatus",'COMPLETED'::"ReservationStatus",'ACTIVE'::"ReservationStatus")
		LIMIT 1`, authorID, sellerID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

type ReviewAIModerationRow struct {
	ID             int32
	Text           *string
	Rating         float64
	CreatedAt      time.Time
	ReviewedByID   int32
	ReviewedByName string
	ReviewedByMail string
	ReviewedUserID int32
	ReviewedName   string
	ReviewedMail   string
}

type ReviewAppealRow struct {
	ID            int64
	ReviewID      int32
	UserID        int32
	Reason        string
	Status        string
	ModeratorID   *int32
	ModeratorNote *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (r *ReviewPG) ListPendingAIModeration(ctx context.Context, limit int) ([]ReviewAIModerationRow, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx, `
		SELECT r.id, r.text, r.rating, r."createdAt",
		       a.id, a."fullName", a.email,
		       s.id, s."fullName", s.email
		FROM "Review" r
		JOIN "User" a ON a.id = r."reviewedById"
		JOIN "User" s ON s.id = r."reviewedUserId"
		WHERE r."moderateState" = 'MODERATE'::"ReviewModerate"
		ORDER BY r."createdAt" ASC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ReviewAIModerationRow
	for rows.Next() {
		var x ReviewAIModerationRow
		if err := rows.Scan(&x.ID, &x.Text, &x.Rating, &x.CreatedAt, &x.ReviewedByID, &x.ReviewedByName, &x.ReviewedByMail, &x.ReviewedUserID, &x.ReviewedName, &x.ReviewedMail); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *ReviewPG) GetReviewParties(ctx context.Context, reviewID int32) (authorName, authorEmail, sellerName, sellerEmail string, err error) {
	err = r.pool.QueryRow(ctx, `
		SELECT a."fullName", a.email, s."fullName", s.email
		FROM "Review" r
		JOIN "User" a ON a.id = r."reviewedById"
		JOIN "User" s ON s.id = r."reviewedUserId"
		WHERE r.id = $1`, reviewID).Scan(&authorName, &authorEmail, &sellerName, &sellerEmail)
	return
}

func (r *ReviewPG) CreateReviewAppeal(ctx context.Context, reviewID, userID int32, reason string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO "ReviewAppeal" ("reviewId","userId",reason,status,"createdAt","updatedAt")
		VALUES ($1,$2,$3,'OPEN',NOW(),NOW())`, reviewID, userID, reason)
	return err
}

func (r *ReviewPG) ListReviewAppeals(ctx context.Context, userID *int32) ([]ReviewAppealRow, error) {
	base := `SELECT id,"reviewId","userId",reason,status,"moderatorId","moderatorNote","createdAt","updatedAt" FROM "ReviewAppeal"`
	var rows pgx.Rows
	var err error
	if userID != nil {
		rows, err = r.pool.Query(ctx, base+` WHERE "userId" = $1 ORDER BY "createdAt" DESC`, *userID)
	} else {
		rows, err = r.pool.Query(ctx, base+` ORDER BY "createdAt" DESC`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ReviewAppealRow
	for rows.Next() {
		var x ReviewAppealRow
		if err := rows.Scan(&x.ID, &x.ReviewID, &x.UserID, &x.Reason, &x.Status, &x.ModeratorID, &x.ModeratorNote, &x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *ReviewPG) ResolveReviewAppeal(ctx context.Context, appealID int64, moderatorID int32, status string, note *string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE "ReviewAppeal"
		SET status = $2, "moderatorId" = $3, "moderatorNote" = $4, "updatedAt" = NOW()
		WHERE id = $1`, appealID, status, moderatorID, note)
	return err
}

func (r *ReviewPG) UserEmailAndName(ctx context.Context, userID int32) (name, email string, err error) {
	err = r.pool.QueryRow(ctx, `SELECT "fullName", email FROM "User" WHERE id = $1`, userID).Scan(&name, &email)
	return
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
