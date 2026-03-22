package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentPG struct {
	pool *pgxpool.Pool
}

func NewPaymentPG(pool *pgxpool.Pool) *PaymentPG {
	return &PaymentPG{pool: pool}
}

// PaymentRow — запись для API (как Prisma payment findMany).
type PaymentRow struct {
	ID         int32     `json:"id"`
	OrderID    string    `json:"orderId"`
	PaymentID  string    `json:"paymentId"`
	UserID     int32     `json:"userId"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	PaymentURL *string   `json:"paymentUrl"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (r *PaymentPG) Insert(ctx context.Context, orderID, paymentID string, userID int32, amount float64, status string, paymentURL *string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO "Payment" ("orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt")
		VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())`,
		orderID, paymentID, userID, amount, status, paymentURL)
	return err
}

func (r *PaymentPG) ListByUserID(ctx context.Context, userID int32, limit int) ([]PaymentRow, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, "orderId", "paymentId", "userId", amount, status, "paymentUrl", "createdAt", "updatedAt"
		FROM "Payment" WHERE "userId" = $1 ORDER BY "createdAt" DESC LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PaymentRow
	for rows.Next() {
		var p PaymentRow
		if err := rows.Scan(&p.ID, &p.OrderID, &p.PaymentID, &p.UserID, &p.Amount, &p.Status, &p.PaymentURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// TryCompleteTopUp — CONFIRMED webhook: только из PENDING начисляем баланс и пишем Log (идемпотентность как в Nest).
func (r *PaymentPG) TryCompleteTopUp(ctx context.Context, paymentID string) (alreadyDone bool, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var uid int32
	var amount float64
	var status string
	var email string
	var bal, bonus float64

	err = tx.QueryRow(ctx, `
		SELECT p."userId", p.amount, p.status, u.email, u.balance, u."bonusBalance"
		FROM "Payment" p
		JOIN "User" u ON u.id = p."userId"
		WHERE p."paymentId" = $1
		FOR UPDATE OF p, u`, paymentID).Scan(&uid, &amount, &status, &email, &bal, &bonus)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, ErrNotFound
	}
	if err != nil {
		return false, err
	}
	if status != "PENDING" {
		return true, nil
	}

	newBal := bal + amount
	action := fmt.Sprintf(
		"Пополнение баланса: id: %d; email: %s;\nсумма пополнения: %g; баланс: %g; бонусный баланс: %g",
		uid, email, amount, newBal, bonus,
	)

	if _, err := tx.Exec(ctx, `UPDATE "Payment" SET status = 'COMPLETED', "updatedAt" = NOW() WHERE "paymentId" = $1`, paymentID); err != nil {
		return false, err
	}
	if _, err := tx.Exec(ctx, `UPDATE "User" SET balance = balance + $2, "updatedAt" = NOW() WHERE id = $1`, uid, amount); err != nil {
		return false, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO "Log" ("userId", action) VALUES ($1, $2)`, uid, action); err != nil {
		return false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return false, nil
}

// TryFailIfPending — REJECTED / CANCELED: PENDING → FAILED.
func (r *PaymentPG) TryFailIfPending(ctx context.Context, paymentID string) (alreadyDone bool, err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var status string
	err = tx.QueryRow(ctx, `SELECT status FROM "Payment" WHERE "paymentId" = $1 FOR UPDATE`, paymentID).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, ErrNotFound
	}
	if err != nil {
		return false, err
	}
	if status != "PENDING" {
		return true, nil
	}
	if _, err := tx.Exec(ctx, `UPDATE "Payment" SET status = 'FAILED', "updatedAt" = NOW() WHERE "paymentId" = $1`, paymentID); err != nil {
		return false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return false, nil
}
