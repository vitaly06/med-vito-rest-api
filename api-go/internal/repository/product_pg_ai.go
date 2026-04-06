package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type AIModerationProduct struct {
	ID                        int32
	Name                      string
	Description               *string
	Price                     int32
	Images                    []string
	CategoryName              string
	SubCategoryName           string
	ModerateState             string
	ModerationRejectionReason *string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	CategoryID                int32
	SubCategoryID             int32
	UserID                    int32
	UserFullName              string
	UserEmail                 string
	UserPhoneNumber           string
}

type ModerationListRow struct {
	ID                        int32
	Name                      string
	Price                     int32
	Images                    []string
	ModerateState             string
	ModerationRejectionReason *string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	CategoryID                int32
	CategoryName              string
	SubCategoryID             int32
	SubCategoryName           string
	UserID                    int32
	UserFullName              string
	UserEmail                 string
	UserPhoneNumber           string
}

type ModerationDetail struct {
	ID                        int32
	Name                      string
	Price                     int32
	Description               *string
	Images                    []string
	VideoURL                  *string
	ModerateState             string
	ModerationRejectionReason *string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	CategoryID                int32
	CategoryName              string
	SubCategoryID             int32
	SubCategoryName           string
	TypeID                    *int32
	TypeName                  *string
	UserID                    int32
	UserFullName              string
	UserEmail                 string
	UserPhoneNumber           string
	UserProfileType           string
	FieldValues               []ModerationFieldValue
}

type ModerationFieldValue struct {
	Value   string
	FieldID int32
	Field   string
}

func (r *ProductPG) ListPendingAIModeration(ctx context.Context, limit int) ([]AIModerationProduct, error) {
	if limit <= 0 {
		limit = 5
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			p.id,
			p.name,
			p.description,
			p.price,
			p.images,
			c.name,
			sc.name
		FROM "Product" p
		JOIN "Category" c ON c.id = p."categoryId"
		JOIN "SubCategory" sc ON sc.id = p."subCategoryId"
		WHERE p."moderateState" = 'MODERATE'::"ProductModerate"
		  AND p."isHide" = false
		  AND p."moderationRejectionReason" IS NULL
		ORDER BY p."createdAt" ASC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AIModerationProduct
	for rows.Next() {
		var row AIModerationProduct
		if err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Description,
			&row.Price,
			&row.Images,
			&row.CategoryName,
			&row.SubCategoryName,
		); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *ProductPG) ListProductsForModeration(ctx context.Context, filter string, page, pageSize int) ([]ModerationListRow, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	whereSQL, args := buildModerationFilterWhere(filter)
	countQuery := `SELECT COUNT(*) FROM "Product" p WHERE ` + whereSQL

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, offset)
	rows, err := r.pool.Query(ctx, `
		SELECT
			p.id,
			p.name,
			p.price,
			p.images,
			p."moderateState"::text,
			p."moderationRejectionReason",
			p."createdAt",
			p."updatedAt",
			c.id,
			c.name,
			sc.id,
			sc.name,
			u.id,
			u."fullName",
			u.email,
			u."phoneNumber"
		FROM "Product" p
		JOIN "Category" c ON c.id = p."categoryId"
		JOIN "SubCategory" sc ON sc.id = p."subCategoryId"
		JOIN "User" u ON u.id = p."userId"
		WHERE `+whereSQL+`
		ORDER BY p."updatedAt" DESC
		LIMIT $`+fmt.Sprint(len(args)-1)+` OFFSET $`+fmt.Sprint(len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []ModerationListRow
	for rows.Next() {
		var row ModerationListRow
		if err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Price,
			&row.Images,
			&row.ModerateState,
			&row.ModerationRejectionReason,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.CategoryID,
			&row.CategoryName,
			&row.SubCategoryID,
			&row.SubCategoryName,
			&row.UserID,
			&row.UserFullName,
			&row.UserEmail,
			&row.UserPhoneNumber,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, row)
	}
	return out, total, rows.Err()
}

func (r *ProductPG) GetModerationProductByID(ctx context.Context, productID int32) (*ModerationDetail, error) {
	var out ModerationDetail
	err := r.pool.QueryRow(ctx, `
		SELECT
			p.id,
			p.name,
			p.price,
			p.description,
			p.images,
			p."videoUrl",
			p."moderateState"::text,
			p."moderationRejectionReason",
			p."createdAt",
			p."updatedAt",
			c.id,
			c.name,
			sc.id,
			sc.name,
			t.id,
			t.name,
			u.id,
			u."fullName",
			u.email,
			u."phoneNumber",
			u."profileType"::text
		FROM "Product" p
		JOIN "Category" c ON c.id = p."categoryId"
		JOIN "SubCategory" sc ON sc.id = p."subCategoryId"
		LEFT JOIN "SubcategotyType" t ON t.id = p."typeId"
		JOIN "User" u ON u.id = p."userId"
		WHERE p.id = $1`, productID).Scan(
		&out.ID,
		&out.Name,
		&out.Price,
		&out.Description,
		&out.Images,
		&out.VideoURL,
		&out.ModerateState,
		&out.ModerationRejectionReason,
		&out.CreatedAt,
		&out.UpdatedAt,
		&out.CategoryID,
		&out.CategoryName,
		&out.SubCategoryID,
		&out.SubCategoryName,
		&out.TypeID,
		&out.TypeName,
		&out.UserID,
		&out.UserFullName,
		&out.UserEmail,
		&out.UserPhoneNumber,
		&out.UserProfileType,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT pfv.value, tf.id, tf.name
		FROM "ProductFieldValue" pfv
		JOIN "TypeField" tf ON tf.id = pfv."fieldId"
		WHERE pfv."productId" = $1
		ORDER BY tf.id ASC`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item ModerationFieldValue
		if err := rows.Scan(&item.Value, &item.FieldID, &item.Field); err != nil {
			return nil, err
		}
		out.FieldValues = append(out.FieldValues, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &out, nil
}

func buildModerationFilterWhere(filter string) (string, []any) {
	const aiApprovedReason = "Одобрено ИИ автоматически"

	switch strings.ToUpper(strings.TrimSpace(filter)) {
	case "DENIED":
		return `p."moderateState" = 'DENIDED'::"ProductModerate"`, nil
	case "APPROVED_AI":
		return `p."moderateState" = 'APPROVED'::"ProductModerate" AND p."isHide" = false AND p."moderationRejectionReason" = '` + aiApprovedReason + `'`, nil
	case "MANUAL":
		return `(p."moderateState" = 'AI_REVIEWED'::"ProductModerate" OR p."moderateState" = 'MODERATE'::"ProductModerate")`, nil
	default:
		return `(
			p."moderateState" = 'DENIDED'::"ProductModerate"
			OR p."moderateState" = 'AI_REVIEWED'::"ProductModerate"
			OR p."moderateState" = 'MODERATE'::"ProductModerate"
			OR (
				p."moderateState" = 'APPROVED'::"ProductModerate"
				AND p."isHide" = false
				AND p."moderationRejectionReason" = '` + aiApprovedReason + `'
			)
		)`, nil
	}
}
