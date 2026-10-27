package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPromoProductNotFound = errors.New("promotion: product not found")
	ErrPromoTariffNotFound  = errors.New("promotion: tariff not found")
	ErrPromoNotOwner        = errors.New("promotion: not product owner")
	ErrPromoDowngrade       = errors.New("promotion: downgrade not allowed")
	ErrPromoNotApproved     = errors.New("promotion: product not approved")
)

type PromotionInsufficientError struct {
	Required, Available float64
}

func (e *PromotionInsufficientError) Error() string { return "promotion: insufficient funds" }

type PromotionPG struct {
	pool *pgxpool.Pool
}

func NewPromotionPG(pool *pgxpool.Pool) *PromotionPG {
	return &PromotionPG{pool: pool}
}

type PromotionTariff struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	PricePerDay int32  `json:"pricePerDay"`
}

func (r *PromotionPG) ListTariffs(ctx context.Context) ([]PromotionTariff, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, "pricePerDay" FROM "Promotion" ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PromotionTariff
	for rows.Next() {
		var t PromotionTariff
		if err := rows.Scan(&t.ID, &t.Name, &t.PricePerDay); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

type AddPromotionResult struct {
	ID          int32
	ProductName string
	TariffName  string
	Days        int32
	TotalPrice  int32
	StartDate   time.Time
	EndDate     time.Time
}

func (r *PromotionPG) InsertFunnelEvent(ctx context.Context, userID *int32, step string, promotionID, productID *int32) {
	_, _ = r.pool.Exec(ctx, `
		INSERT INTO "TariffFunnelEvent" ("userId", step, "promotionId", "productId", "createdAt")
		VALUES ($1, $2, $3, $4, NOW())`,
		userID, step, promotionID, productID)
}

func (r *PromotionPG) AddProductPromotion(ctx context.Context, userID, productID, promotionID, days int32) (*AddPromotionResult, error) {
	if days < 1 {
		return nil, fmt.Errorf("days")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var ownerID int32
	var prodName string
	var moderateState string
	err = tx.QueryRow(ctx, `SELECT "userId", name, COALESCE("moderateState"::text, 'MODERATE') FROM "Product" WHERE id = $1 FOR UPDATE`, productID).Scan(&ownerID, &prodName, &moderateState)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPromoProductNotFound
	}
	if err != nil {
		return nil, err
	}
	if ownerID != userID {
		return nil, ErrPromoNotOwner
	}
	if moderateState != "APPROVED" {
		return nil, ErrPromoNotApproved
	}

	var pricePerDay int32
	var tariffName string
	err = tx.QueryRow(ctx, `SELECT "pricePerDay", name FROM "Promotion" WHERE id = $1 FOR UPDATE`, promotionID).Scan(&pricePerDay, &tariffName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPromoTariffNotFound
	}
	if err != nil {
		return nil, err
	}
	var currentMax int32
	_ = tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(pr."pricePerDay"),0)::int
		FROM "ProductPromotion" pp
		JOIN "Promotion" pr ON pr.id = pp."promotionId"
		WHERE pp."productId" = $1
		  AND pp."isActive" = true
		  AND pp."isPaid" = true
		  AND pp."endDate" >= NOW()`, productID).Scan(&currentMax)
	if currentMax > 0 && pricePerDay < currentMax {
		return nil, ErrPromoDowngrade
	}

	totalPrice := int32(int64(days) * int64(pricePerDay))
	totalF := float64(totalPrice)

	var bal, bonus float64
	err = tx.QueryRow(ctx, `SELECT balance, "bonusBalance" FROM "User" WHERE id = $1 FOR UPDATE`, userID).Scan(&bal, &bonus)
	if err != nil {
		return nil, err
	}
	avail := bal + bonus
	if avail < totalF {
		return nil, &PromotionInsufficientError{Required: totalF, Available: avail}
	}

	balanceUsed := bal
	if balanceUsed < 0 {
		balanceUsed = 0
	}
	if balanceUsed > totalF {
		balanceUsed = totalF
	}
	remaining := totalF - balanceUsed
	bonusUsed := remaining

	start := time.Now().UTC()
	end := start.AddDate(0, 0, int(days))

	var ppID int32
	var existed bool
	err = tx.QueryRow(ctx, `
		SELECT id
		FROM "ProductPromotion"
		WHERE "productId" = $1 AND "promotionId" = $2 AND "isActive" = true AND "isPaid" = true AND "endDate" >= NOW()
		ORDER BY "endDate" DESC
		LIMIT 1`, productID, promotionID).Scan(&ppID)
	if err == nil {
		existed = true
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if existed {
		var curEnd time.Time
		if err := tx.QueryRow(ctx, `SELECT "endDate" FROM "ProductPromotion" WHERE id = $1 FOR UPDATE`, ppID).Scan(&curEnd); err != nil {
			return nil, err
		}
		newEnd := curEnd.AddDate(0, 0, int(days))
		if _, err := tx.Exec(ctx, `
			UPDATE "ProductPromotion"
			SET days = days + $2,
			    "totalPrice" = "totalPrice" + $3,
			    "endDate" = $4,
			    "updatedAt" = NOW()
			WHERE id = $1`, ppID, days, totalPrice, newEnd); err != nil {
			return nil, err
		}
		start = curEnd
		end = newEnd
	} else {
		err = tx.QueryRow(ctx, `
			INSERT INTO "ProductPromotion" ("productId", "promotionId", "userId", days, "totalPrice", "startDate", "endDate", "isActive", "isPaid", "createdAt", "updatedAt")
			VALUES ($1,$2,$3,$4,$5,$6,$7,true,false,NOW(),NOW())
			RETURNING id`,
			productID, promotionID, userID, days, totalPrice, start, end,
		).Scan(&ppID)
		if err != nil {
			return nil, err
		}
	}

	newBal := bal - balanceUsed
	newBonus := bonus - bonusUsed
	if _, err := tx.Exec(ctx, `UPDATE "User" SET balance = $2, "bonusBalance" = $3, "updatedAt" = NOW() WHERE id = $1`, userID, newBal, newBonus); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE "ProductPromotion" SET "isPaid" = true, "updatedAt" = NOW() WHERE id = $1`, ppID); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE "Product"
		SET "isHide" = false,
		    "expiresAt" = GREATEST(COALESCE("expiresAt", NOW()), $2),
		    "updatedAt" = NOW()
		WHERE id = $1`, productID, end); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	r.InsertFunnelEvent(ctx, &userID, "payment", &promotionID, &productID)
	r.InsertFunnelEvent(ctx, &userID, "publication", &promotionID, &productID)
	return &AddPromotionResult{
		ID:          ppID,
		ProductName: prodName,
		TariffName:  tariffName,
		Days:        days,
		TotalPrice:  totalPrice,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

type ExpiringPromotionRow struct {
	ID          int32
	ProductID   int32
	UserID      int32
	Days        int32
	EndDate     time.Time
	ProductName string
	Phone       string
	Email       string
	DaysLeft    int
}

func (r *PromotionPG) EnsureNotificationTable(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS "ProductPromotionNotification" (
			"promotionId" INT NOT NULL,
			"daysInterval" INT NOT NULL,
			"notifiedAt" TIMESTAMP NOT NULL DEFAULT NOW(),
			PRIMARY KEY ("promotionId", "daysInterval")
		)`)
	return err
}

func (r *PromotionPG) FindPromotionsExpiringInDays(ctx context.Context, daysLeft int) ([]ExpiringPromotionRow, error) {
	var q string
	switch daysLeft {
	case 7:
		q = `
			SELECT pp.id, pp."productId", pp."userId", pp.days, pp."endDate",
			       COALESCE(p.name, ''), COALESCE(u."phoneNumber", ''), COALESCE(u.email, '')
			FROM "ProductPromotion" pp
			JOIN "Product" p ON p.id = pp."productId"
			JOIN "User" u ON u.id = pp."userId"
			WHERE pp."isActive" = true
			  AND pp."isPaid" = true
			  AND pp.days >= 7
			  AND pp."endDate" <= NOW() + INTERVAL '7 days'
			  AND pp."endDate" > NOW() + INTERVAL '3 days'
			  AND NOT EXISTS (
			      SELECT 1 FROM "ProductPromotionNotification" n
			      WHERE n."promotionId" = pp.id AND n."daysInterval" = 7
			  )`
	case 3:
		q = `
			SELECT pp.id, pp."productId", pp."userId", pp.days, pp."endDate",
			       COALESCE(p.name, ''), COALESCE(u."phoneNumber", ''), COALESCE(u.email, '')
			FROM "ProductPromotion" pp
			JOIN "Product" p ON p.id = pp."productId"
			JOIN "User" u ON u.id = pp."userId"
			WHERE pp."isActive" = true
			  AND pp."isPaid" = true
			  AND pp.days >= 3
			  AND pp."endDate" <= NOW() + INTERVAL '3 days'
			  AND pp."endDate" > NOW() + INTERVAL '1 day'
			  AND NOT EXISTS (
			      SELECT 1 FROM "ProductPromotionNotification" n
			      WHERE n."promotionId" = pp.id AND n."daysInterval" = 3
			  )`
	case 1:
		q = `
			SELECT pp.id, pp."productId", pp."userId", pp.days, pp."endDate",
			       COALESCE(p.name, ''), COALESCE(u."phoneNumber", ''), COALESCE(u.email, '')
			FROM "ProductPromotion" pp
			JOIN "Product" p ON p.id = pp."productId"
			JOIN "User" u ON u.id = pp."userId"
			WHERE pp."isActive" = true
			  AND pp."isPaid" = true
			  AND pp."endDate" <= NOW() + INTERVAL '1 day'
			  AND pp."endDate" > NOW()
			  AND NOT EXISTS (
			      SELECT 1 FROM "ProductPromotionNotification" n
			      WHERE n."promotionId" = pp.id AND n."daysInterval" = 1
			  )`
	default:
		return nil, fmt.Errorf("unsupported daysLeft %d", daysLeft)
	}

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ExpiringPromotionRow
	for rows.Next() {
		var x ExpiringPromotionRow
		x.DaysLeft = daysLeft
		if err := rows.Scan(&x.ID, &x.ProductID, &x.UserID, &x.Days, &x.EndDate, &x.ProductName, &x.Phone, &x.Email); err != nil {
			continue
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *PromotionPG) RecordPromotionNotified(ctx context.Context, promotionID, daysInterval int) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO "ProductPromotionNotification" ("promotionId", "daysInterval", "notifiedAt")
		VALUES ($1, $2, NOW())
		ON CONFLICT ("promotionId", "daysInterval") DO NOTHING`, promotionID, daysInterval)
	return err
}

func (r *PromotionPG) DeactivateExpiredPromotions(ctx context.Context) ([]ExpiringPromotionRow, error) {
	const q = `
		WITH expired AS (
		    UPDATE "ProductPromotion" pp
		    SET "isActive" = false, "updatedAt" = NOW()
		    WHERE pp."isActive" = true AND pp."isPaid" = true AND pp."endDate" < NOW()
		    RETURNING pp.id, pp."productId", pp."userId", pp.days, pp."endDate"
		)
		SELECT e.id, e."productId", e."userId", e.days, e."endDate",
		       COALESCE(p.name, ''), COALESCE(u."phoneNumber", ''), COALESCE(u.email, '')
		FROM expired e
		JOIN "Product" p ON p.id = e."productId"
		JOIN "User" u ON u.id = e."userId"`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ExpiringPromotionRow
	for rows.Next() {
		var x ExpiringPromotionRow
		if err := rows.Scan(&x.ID, &x.ProductID, &x.UserID, &x.Days, &x.EndDate, &x.ProductName, &x.Phone, &x.Email); err != nil {
			continue
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

