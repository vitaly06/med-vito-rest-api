package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationPG struct {
	pool *pgxpool.Pool
}

func NewReservationPG(pool *pgxpool.Pool) *ReservationPG {
	return &ReservationPG{pool: pool}
}

type ReservationRow struct {
	ID           int64
	ProductID    int32
	ProductName  string
	BuyerID      int32
	BuyerName    string
	SellerID     int32
	SellerName   string
	Status       string
	Hours        int32
	Note         *string
	CancelReason *string
	ExtendedOnce bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
	CancelledAt  *time.Time
}

type ProductReserveSettings struct {
	ProductID         int32
	ProductName       string
	SellerID          int32
	SellerName        string
	SellerEmail       string
	AllowReservations bool
	ReservationHours  int32
	ModerateState     string
	IsHide            bool
}

func (r *ReservationPG) GetProductReserveSettings(ctx context.Context, productID int32) (*ProductReserveSettings, error) {
	var out ProductReserveSettings
	err := r.pool.QueryRow(ctx, `
		SELECT p.id, p.name, p."userId", u."fullName", u.email,
		       COALESCE(p."allowReservations", true), COALESCE(p."reservationHours", 24),
		       p."moderateState"::text, p."isHide"
		FROM "Product" p
		JOIN "User" u ON u.id = p."userId"
		WHERE p.id = $1`, productID).
		Scan(&out.ProductID, &out.ProductName, &out.SellerID, &out.SellerName, &out.SellerEmail, &out.AllowReservations, &out.ReservationHours, &out.ModerateState, &out.IsHide)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &out, err
}

func (r *ReservationPG) IsUserVerified(ctx context.Context, userID int32) (bool, string, error) {
	var verified bool
	var name string
	err := r.pool.QueryRow(ctx, `SELECT "isEmailVerified", "fullName" FROM "User" WHERE id = $1`, userID).Scan(&verified, &name)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, "", ErrNotFound
	}
	return verified, name, err
}

func (r *ReservationPG) ActiveReservationByProduct(ctx context.Context, productID int32) (*ReservationRow, error) {
	row := r.pool.QueryRow(ctx, reservationSelectSQL+` WHERE r."productId" = $1 AND r.status = 'ACTIVE' ORDER BY r."createdAt" DESC LIMIT 1`, productID)
	res, err := scanReservation(row)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *ReservationPG) PublicActiveReservationByProduct(ctx context.Context, productID int32) (*ReservationRow, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT r.id, r."productId", p.name, r."buyerId", bu."fullName", r."sellerId", su."fullName",
		       r.status, r."hours", r.note, r."cancelReason", r."extendedOnce", r."expiresAt", r."createdAt", r."cancelledAt"
		FROM "ProductReservation" r
		JOIN "Product" p ON p.id = r."productId"
		JOIN "User" bu ON bu.id = r."buyerId"
		JOIN "User" su ON su.id = r."sellerId"
		WHERE r."productId" = $1 AND r.status = 'ACTIVE'
		ORDER BY r."createdAt" DESC
		LIMIT 1`, productID)
	res, err := scanReservation(row)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	return res, err
}

func (r *ReservationPG) ActiveReservationsCount(ctx context.Context, buyerID int32) (int, error) {
	var c int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM "ProductReservation" WHERE "buyerId" = $1 AND status = 'ACTIVE'`, buyerID).Scan(&c)
	return c, err
}

func (r *ReservationPG) DailyReservationsCount(ctx context.Context, buyerID int32, since time.Time) (int, error) {
	var c int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM "ProductReservation" WHERE "buyerId" = $1 AND "createdAt" >= $2`, buyerID, since).Scan(&c)
	return c, err
}

func (r *ReservationPG) HasReserveBlock(ctx context.Context, buyerID int32, now time.Time) (bool, *time.Time, error) {
	var until *time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT "blockedUntil"
		FROM "ReservationUserPenalty"
		WHERE "userId" = $1`, buyerID).Scan(&until)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	if until != nil && until.After(now) {
		return true, until, nil
	}
	return false, until, nil
}

func (r *ReservationPG) CreateReservation(ctx context.Context, productID, buyerID, sellerID int32, hours int32, note *string, now time.Time) (*ReservationRow, error) {
	var id int64
	expires := now.Add(time.Duration(hours) * time.Hour)
	err := r.pool.QueryRow(ctx, `
		INSERT INTO "ProductReservation" ("productId", "buyerId", "sellerId", status, "hours", note, "expiresAt", "createdAt", "updatedAt")
		VALUES ($1,$2,$3,'ACTIVE',$4,$5,$6,NOW(),NOW())
		RETURNING id`, productID, buyerID, sellerID, hours, note, expires).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *ReservationPG) GetByID(ctx context.Context, id int64) (*ReservationRow, error) {
	return scanReservation(r.pool.QueryRow(ctx, reservationSelectSQL+` WHERE r.id = $1`, id))
}

func (r *ReservationPG) ListForUser(ctx context.Context, userID int32) ([]ReservationRow, error) {
	rows, err := r.pool.Query(ctx, reservationSelectSQL+` WHERE r."buyerId" = $1 OR r."sellerId" = $1 ORDER BY r."createdAt" DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanReservations(rows)
}

func (r *ReservationPG) FindParticipantReservation(ctx context.Context, id int64, userID int32) (*ReservationRow, error) {
	return scanReservation(r.pool.QueryRow(ctx, reservationSelectSQL+` WHERE r.id = $1 AND (r."buyerId" = $2 OR r."sellerId" = $2)`, id, userID))
}

func (r *ReservationPG) CancelByBuyer(ctx context.Context, id int64, buyerID int32) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE "ProductReservation"
		SET status = 'CANCELLED_BY_BUYER', "cancelledAt" = NOW(), "updatedAt" = NOW()
		WHERE id = $1 AND "buyerId" = $2 AND status = 'ACTIVE'`, id, buyerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ReservationPG) CancelBySeller(ctx context.Context, id int64, sellerID int32, reason string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE "ProductReservation"
		SET status = 'CANCELLED_BY_SELLER', "cancelReason" = $3, "cancelledAt" = NOW(), "updatedAt" = NOW()
		WHERE id = $1 AND "sellerId" = $2 AND status = 'ACTIVE'`, id, sellerID, reason)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ReservationPG) ExtendByBuyer(ctx context.Context, id int64, buyerID int32) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE "ProductReservation"
		SET "expiresAt" = "expiresAt" + (("hours"::text || ' hours')::interval),
		    "extendedOnce" = true,
		    "updatedAt" = NOW()
		WHERE id = $1 AND "buyerId" = $2 AND status = 'ACTIVE' AND "extendedOnce" = false`, id, buyerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ReservationPG) MarkDealCreated(ctx context.Context, productID, buyerID int32) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE "ProductReservation"
		SET status = 'DEAL_CREATED', "updatedAt" = NOW()
		WHERE "productId" = $1 AND "buyerId" = $2 AND status = 'ACTIVE'`, productID, buyerID)
	return err
}

func (r *ReservationPG) ExpireDue(ctx context.Context, now time.Time) (int, error) {
	tag, err := r.pool.Exec(ctx, `
		UPDATE "ProductReservation"
		SET status = 'EXPIRED', "updatedAt" = NOW()
		WHERE status = 'ACTIVE' AND "expiresAt" <= $1`, now)
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}

func (r *ReservationPG) ExpireDueAndReturn(ctx context.Context, now time.Time) ([]ReservationRow, error) {
	rows, err := r.pool.Query(ctx, `
		WITH upd AS (
			UPDATE "ProductReservation"
			SET status = 'EXPIRED', "updatedAt" = NOW()
			WHERE status = 'ACTIVE' AND "expiresAt" <= $1
			RETURNING id
		)
		SELECT
			r.id, r."productId", p.name, r."buyerId", bu."fullName", r."sellerId", su."fullName",
			r.status, r."hours", r.note, r."cancelReason", r."extendedOnce", r."expiresAt", r."createdAt", r."cancelledAt"
		FROM upd
		JOIN "ProductReservation" r ON r.id = upd.id
		JOIN "Product" p ON p.id = r."productId"
		JOIN "User" bu ON bu.id = r."buyerId"
		JOIN "User" su ON su.id = r."sellerId"
		ORDER BY r.id DESC`, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanReservations(rows)
}

func (r *ReservationPG) LastBuyerStatuses(ctx context.Context, buyerID int32, limit int) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT status FROM "ProductReservation"
		WHERE "buyerId" = $1
		ORDER BY "createdAt" DESC
		LIMIT $2`, buyerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *ReservationPG) SetReserveBlock(ctx context.Context, buyerID int32, until time.Time, reason string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO "ReservationUserPenalty" ("userId", "blockedUntil", reason, "updatedAt")
		VALUES ($1,$2,$3,NOW())
		ON CONFLICT ("userId")
		DO UPDATE SET "blockedUntil" = EXCLUDED."blockedUntil", reason = EXCLUDED.reason, "updatedAt" = NOW()`,
		buyerID, until, reason)
	return err
}

func (r *ReservationPG) UpdateProductReserveSettings(ctx context.Context, productID, sellerID int32, allow bool, hours int32) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE "Product"
		SET "allowReservations" = $3, "reservationHours" = $4, "updatedAt" = NOW()
		WHERE id = $1 AND "userId" = $2`, productID, sellerID, allow, hours)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ReservationPG) BuyerReliability(ctx context.Context, buyerID int32) (float64, error) {
	var total int
	var successful int
	if err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)::int,
		       COUNT(*) FILTER (WHERE status IN ('COMPLETED', 'DEAL_CREATED'))::int
		FROM "ProductReservation"
		WHERE "buyerId" = $1`, buyerID).Scan(&total, &successful); err != nil {
		return 0, err
	}
	if total == 0 {
		return 100, nil
	}
	return (float64(successful) / float64(total)) * 100.0, nil
}

func (r *ReservationPG) RecentExpiredBuyers(ctx context.Context, since time.Time) ([]int32, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT "buyerId" FROM "ProductReservation"
		WHERE status = 'EXPIRED' AND "updatedAt" >= $1`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

const reservationSelectSQL = `
	SELECT
		r.id, r."productId", p.name, r."buyerId", bu."fullName", r."sellerId", su."fullName",
		r.status, r."hours", r.note, r."cancelReason", r."extendedOnce", r."expiresAt", r."createdAt", r."cancelledAt"
	FROM "ProductReservation" r
	JOIN "Product" p ON p.id = r."productId"
	JOIN "User" bu ON bu.id = r."buyerId"
	JOIN "User" su ON su.id = r."sellerId"`

func scanReservation(row pgx.Row) (*ReservationRow, error) {
	var r ReservationRow
	err := row.Scan(
		&r.ID, &r.ProductID, &r.ProductName, &r.BuyerID, &r.BuyerName, &r.SellerID, &r.SellerName,
		&r.Status, &r.Hours, &r.Note, &r.CancelReason, &r.ExtendedOnce, &r.ExpiresAt, &r.CreatedAt, &r.CancelledAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func scanReservations(rows pgx.Rows) ([]ReservationRow, error) {
	var out []ReservationRow
	for rows.Next() {
		r, err := scanReservation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *r)
	}
	return out, rows.Err()
}
