package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DealPG struct {
	pool *pgxpool.Pool
}

func NewDealPG(pool *pgxpool.Pool) *DealPG {
	return &DealPG{pool: pool}
}

type DealRow struct {
	ID              int32
	ProductID       int32
	ProductName     string
	ProductImages   []string
	BuyerID         int32
	BuyerName       string
	SellerID        int32
	SellerName      string
	Status          string
	ProductAmount   int32
	DeliveryCost    int32
	PlatformFee     int32
	SellerAmount    int32
	TotalAmount     int32
	PaymentID       *string
	OrderID         *string
	PaymentURL      *string
	CDEKTariffCode  *int32
	CDEKTariffName  *string
	CDEKFromCity    *int32
	CDEKToCity      *int32
	CDEKFromPVZ     *string
	CDEKToPVZ       *string
	CDEKOrderUUID   *string
	CDEKTrackNumber *string
	DisputeReason   *string
	PaidAt          *time.Time
	ShippedAt       *time.Time
	DeliveredAt     *time.Time
	PayoutAt        *time.Time
	CompletedAt     *time.Time
	CancelledAt     *time.Time
	RefundedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CreateDealParams struct {
	ProductID      int32
	BuyerID        int32
	SellerID       int32
	ProductAmount  int32
	DeliveryCost   int32
	PlatformFee    int32
	SellerAmount   int32
	TotalAmount    int32
	CDEKTariffCode *int32
	CDEKTariffName *string
	CDEKFromCity   *int32
	CDEKToCity     *int32
	CDEKFromPVZ    *string
	CDEKToPVZ      *string
}

type DealProductInfo struct {
	ID       int32
	Name     string
	Price    int32
	UserID   int32
	Approved bool
	IsHide   bool
}

func (r *DealPG) ProductInfo(ctx context.Context, productID int32) (*DealProductInfo, error) {
	var out DealProductInfo
	var moderateState string
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, price, "userId", "moderateState"::text, "isHide"
		FROM "Product" WHERE id = $1`, productID).Scan(&out.ID, &out.Name, &out.Price, &out.UserID, &moderateState, &out.IsHide)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	out.Approved = moderateState == "APPROVED"
	return &out, nil
}

func (r *DealPG) Create(ctx context.Context, p CreateDealParams) (*DealRow, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `
		INSERT INTO "ProductDeal" (
			"productId", "buyerId", "sellerId",
			"productAmount", "deliveryCost", "platformFee", "sellerAmount", "totalAmount",
			"cdekTariffCode", "cdekTariffName", "cdekFromCityCode", "cdekToCityCode", "cdekFromPvzCode", "cdekToPvzCode",
			"createdAt", "updatedAt"
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,NOW(),NOW())
		RETURNING id`,
		p.ProductID, p.BuyerID, p.SellerID,
		p.ProductAmount, p.DeliveryCost, p.PlatformFee, p.SellerAmount, p.TotalAmount,
		p.CDEKTariffCode, p.CDEKTariffName, p.CDEKFromCity, p.CDEKToCity, p.CDEKFromPVZ, p.CDEKToPVZ,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *DealPG) SetPayment(ctx context.Context, dealID int32, paymentID, orderID, paymentURL string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE "ProductDeal"
		SET "paymentId" = $2, "orderId" = $3, "paymentUrl" = $4, "updatedAt" = NOW()
		WHERE id = $1`, dealID, paymentID, orderID, nullString(paymentURL))
	return err
}

func (r *DealPG) SetStatus(ctx context.Context, dealID int32, fromStatuses []string, toStatus string, timestampColumn string) error {
	args := []any{dealID, toStatus, fromStatuses}
	query := `UPDATE "ProductDeal" SET status = $2::"DealStatus", "updatedAt" = NOW()`
	if timestampColumn != "" {
		query += `, "` + timestampColumn + `" = NOW()`
	}
	query += ` WHERE id = $1 AND status = ANY($3::"DealStatus"[])`
	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *DealPG) MarkDelivered(ctx context.Context, dealID int32, payoutAt time.Time) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE "ProductDeal"
		SET status = 'DELIVERED'::"DealStatus", "deliveredAt" = NOW(), "payoutAvailableAt" = $2, "updatedAt" = NOW()
		WHERE id = $1 AND status = 'SHIPPED'::"DealStatus"`, dealID, payoutAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *DealPG) OpenDispute(ctx context.Context, dealID int32, reason string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE "ProductDeal"
		SET status = 'DISPUTE'::"DealStatus", "disputeReason" = $2, "updatedAt" = NOW()
		WHERE id = $1 AND status IN ('PAID'::"DealStatus", 'SHIPPED'::"DealStatus", 'DELIVERED'::"DealStatus")`,
		dealID, reason)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *DealPG) FindByID(ctx context.Context, dealID int32) (*DealRow, error) {
	row := r.pool.QueryRow(ctx, dealSelectSQL+` WHERE d.id = $1`, dealID)
	return scanDeal(row)
}

func (r *DealPG) ListByBuyer(ctx context.Context, buyerID int32) ([]DealRow, error) {
	rows, err := r.pool.Query(ctx, dealSelectSQL+` WHERE d."buyerId" = $1 ORDER BY d."createdAt" DESC`, buyerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDeals(rows)
}

func (r *DealPG) ListBySeller(ctx context.Context, sellerID int32) ([]DealRow, error) {
	rows, err := r.pool.Query(ctx, dealSelectSQL+` WHERE d."sellerId" = $1 ORDER BY d."createdAt" DESC`, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDeals(rows)
}

func (r *DealPG) CompleteDuePayouts(ctx context.Context, limit int) (int, error) {
	if limit <= 0 {
		limit = 20
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows, err := tx.Query(ctx, `
		SELECT id, "sellerId", "sellerAmount"
		FROM "ProductDeal"
		WHERE status = 'DELIVERED'::"DealStatus" AND "payoutAvailableAt" <= NOW()
		ORDER BY "payoutAvailableAt" ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED`, limit)
	if err != nil {
		return 0, err
	}
	type dueDeal struct {
		ID           int32
		SellerID     int32
		SellerAmount int32
	}
	var deals []dueDeal
	for rows.Next() {
		var d dueDeal
		if err := rows.Scan(&d.ID, &d.SellerID, &d.SellerAmount); err != nil {
			rows.Close()
			return 0, err
		}
		deals = append(deals, d)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	for _, d := range deals {
		if _, err := tx.Exec(ctx, `UPDATE "User" SET balance = balance + $2, "updatedAt" = NOW() WHERE id = $1`, d.SellerID, float64(d.SellerAmount)); err != nil {
			return 0, err
		}
		if _, err := tx.Exec(ctx, `
			UPDATE "ProductDeal"
			SET status = 'COMPLETED'::"DealStatus", "completedAt" = NOW(), "updatedAt" = NOW()
			WHERE id = $1`, d.ID); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return len(deals), nil
}

func (r *DealPG) AutoDeliverExpiredShipped(ctx context.Context, shippedBefore time.Time, payoutAt time.Time, limit int) (int, error) {
	if limit <= 0 {
		limit = 20
	}
	tag, err := r.pool.Exec(ctx, `
		UPDATE "ProductDeal"
		SET status = 'DELIVERED'::"DealStatus",
			"deliveredAt" = NOW(),
			"payoutAvailableAt" = $2,
			"updatedAt" = NOW()
		WHERE id IN (
			SELECT id
			FROM "ProductDeal"
			WHERE status = 'SHIPPED'::"DealStatus"
			  AND "shippedAt" IS NOT NULL
			  AND "shippedAt" <= $1
			ORDER BY "shippedAt" ASC
			LIMIT $3
		)`, shippedBefore, payoutAt, limit)
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}

const dealSelectSQL = `
	SELECT
		d.id, d."productId", p.name, p.images,
		d."buyerId", bu."fullName",
		d."sellerId", su."fullName",
		d.status::text,
		d."productAmount", d."deliveryCost", d."platformFee", d."sellerAmount", d."totalAmount",
		d."paymentId", d."orderId", d."paymentUrl",
		d."cdekTariffCode", d."cdekTariffName", d."cdekFromCityCode", d."cdekToCityCode", d."cdekFromPvzCode", d."cdekToPvzCode",
		d."cdekOrderUuid", d."cdekTrackNumber",
		d."disputeReason",
		d."paidAt", d."shippedAt", d."deliveredAt", d."payoutAvailableAt", d."completedAt", d."cancelledAt", d."refundedAt",
		d."createdAt", d."updatedAt"
	FROM "ProductDeal" d
	JOIN "Product" p ON p.id = d."productId"
	JOIN "User" bu ON bu.id = d."buyerId"
	JOIN "User" su ON su.id = d."sellerId"`

func scanDeal(row pgx.Row) (*DealRow, error) {
	var d DealRow
	err := row.Scan(
		&d.ID, &d.ProductID, &d.ProductName, &d.ProductImages,
		&d.BuyerID, &d.BuyerName,
		&d.SellerID, &d.SellerName,
		&d.Status,
		&d.ProductAmount, &d.DeliveryCost, &d.PlatformFee, &d.SellerAmount, &d.TotalAmount,
		&d.PaymentID, &d.OrderID, &d.PaymentURL,
		&d.CDEKTariffCode, &d.CDEKTariffName, &d.CDEKFromCity, &d.CDEKToCity, &d.CDEKFromPVZ, &d.CDEKToPVZ,
		&d.CDEKOrderUUID, &d.CDEKTrackNumber,
		&d.DisputeReason,
		&d.PaidAt, &d.ShippedAt, &d.DeliveredAt, &d.PayoutAt, &d.CompletedAt, &d.CancelledAt, &d.RefundedAt,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func scanDeals(rows pgx.Rows) ([]DealRow, error) {
	var out []DealRow
	for rows.Next() {
		d, err := scanDeal(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *d)
	}
	return out, rows.Err()
}

func nullString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
