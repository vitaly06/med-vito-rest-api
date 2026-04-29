package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type ModerationAuditRow struct {
	ID          int64     `json:"id"`
	ActorUserID *int32    `json:"actorUserId"`
	ActorRole   *string   `json:"actorRole"`
	TargetType  string    `json:"targetType"`
	TargetID    int64     `json:"targetId"`
	Action      string    `json:"action"`
	Payload     []byte    `json:"payload"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ModerationAppealRow struct {
	ID               int64      `json:"id"`
	ProductID        int32      `json:"productId"`
	UserID           int32      `json:"userId"`
	Reason           string     `json:"reason"`
	Status           string     `json:"status"`
	ReviewedByUserID *int32     `json:"reviewedByUserId"`
	ReviewComment    *string    `json:"reviewComment"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

func (r *ProductPG) InsertModerationAudit(ctx context.Context, actorUserID *int32, actorRole *string, targetType string, targetID int64, action string, payload []byte) {
	_, _ = r.pool.Exec(ctx, `
		INSERT INTO "ModerationAuditLog" ("actorUserId","actorRole","targetType","targetId",action,payload,"createdAt")
		VALUES ($1,$2,$3,$4,$5,$6::jsonb,NOW())`,
		actorUserID, actorRole, targetType, targetID, action, string(payload))
}

func (r *ProductPG) ListModerationAudit(ctx context.Context, limit int) ([]ModerationAuditRow, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, "actorUserId", "actorRole", "targetType", "targetId", action, COALESCE(payload::text,'{}')::text, "createdAt"
		FROM "ModerationAuditLog"
		ORDER BY "createdAt" DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ModerationAuditRow
	for rows.Next() {
		var x ModerationAuditRow
		var payload string
		if err := rows.Scan(&x.ID, &x.ActorUserID, &x.ActorRole, &x.TargetType, &x.TargetID, &x.Action, &payload, &x.CreatedAt); err != nil {
			return nil, err
		}
		x.Payload = []byte(payload)
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *ProductPG) CreateAppeal(ctx context.Context, productID, userID int32, reason string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO "ModerationAppeal" ("productId","userId",reason,status,"createdAt","updatedAt")
		VALUES ($1,$2,$3,'OPEN',NOW(),NOW())`, productID, userID, reason)
	return err
}

func (r *ProductPG) ListAppeals(ctx context.Context, userID *int32) ([]ModerationAppealRow, error) {
	q := `SELECT id, "productId", "userId", reason, status, "reviewedByUserId", "reviewComment", "createdAt", "updatedAt" FROM "ModerationAppeal"`
	var rows pgx.Rows
	var err error
	if userID != nil {
		rows, err = r.pool.Query(ctx, q+` WHERE "userId" = $1 ORDER BY "createdAt" DESC`, *userID)
	} else {
		rows, err = r.pool.Query(ctx, q+` ORDER BY "createdAt" DESC`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ModerationAppealRow
	for rows.Next() {
		var x ModerationAppealRow
		if err := rows.Scan(&x.ID, &x.ProductID, &x.UserID, &x.Reason, &x.Status, &x.ReviewedByUserID, &x.ReviewComment, &x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *ProductPG) ReviewAppeal(ctx context.Context, appealID int64, reviewerID int32, status string, comment *string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE "ModerationAppeal"
		SET status = $2, "reviewedByUserId" = $3, "reviewComment" = $4, "updatedAt" = NOW()
		WHERE id = $1`, appealID, status, reviewerID, comment)
	return err
}

func (r *ProductPG) ModerationSummary(ctx context.Context, days int) (denied, approvedAI, appealsOpen int64, err error) {
	if days <= 0 {
		days = 30
	}
	err = r.pool.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE "moderateState" = 'DENIDED'::"ProductModerate")::bigint,
			COUNT(*) FILTER (WHERE "moderateState" = 'APPROVED'::"ProductModerate" AND COALESCE("moderationRejectionReason",'') = 'Одобрено ИИ автоматически')::bigint
		FROM "Product"
		WHERE "updatedAt" >= NOW() - ($1::text || ' days')::interval`, days).Scan(&denied, &approvedAI)
	if err != nil {
		return
	}
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*)::bigint FROM "ModerationAppeal" WHERE status = 'OPEN'`).Scan(&appealsOpen)
	return
}
