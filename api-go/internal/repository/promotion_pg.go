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
