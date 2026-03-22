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
	ID         int32
	Days       int32
	TotalPrice int32
	StartDate  time.Time
	EndDate    time.Time
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
	err = tx.QueryRow(ctx, `SELECT "userId" FROM "Product" WHERE id = $1 FOR UPDATE`, productID).Scan(&ownerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPromoProductNotFound
	}
	if err != nil {
		return nil, err
	}
	if ownerID != userID {
		return nil, ErrPromoNotOwner
	}

	var pricePerDay int32
	err = tx.QueryRow(ctx, `SELECT "pricePerDay" FROM "Promotion" WHERE id = $1 FOR UPDATE`, promotionID).Scan(&pricePerDay)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPromoTariffNotFound
	}
	if err != nil {
		return nil, err
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
	err = tx.QueryRow(ctx, `
		INSERT INTO "ProductPromotion" ("productId", "promotionId", "userId", days, "totalPrice", "startDate", "endDate", "isActive", "isPaid", "createdAt", "updatedAt")
		VALUES ($1,$2,$3,$4,$5,$6,$7,true,false,NOW(),NOW())
		RETURNING id`,
		productID, promotionID, userID, days, totalPrice, start, end,
	).Scan(&ppID)
	if err != nil {
		return nil, err
	}

	newBal := bal - balanceUsed
	newBonus := bonus - bonusUsed
	if _, err := tx.Exec(ctx, `UPDATE "User" SET balance = $2, "bonusBalance" = $3, "updatedAt" = NOW() WHERE id = $1`, userID, newBal, newBonus); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE "ProductPromotion" SET "isPaid" = true, "updatedAt" = NOW() WHERE id = $1`, ppID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &AddPromotionResult{
		ID:         ppID,
		Days:       days,
		TotalPrice: totalPrice,
		StartDate:  start,
		EndDate:    end,
	}, nil
}
