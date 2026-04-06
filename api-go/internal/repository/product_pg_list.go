package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// ProductListRow — строка для списков (all-products, поиск, избранное).
type ProductListRow struct {
	ID              int32
	Images          []string
	Name            string
	Address         string
	CreatedAt       time.Time
	IsHide          bool
	Price           int32
	UserID          int32
	VideoURL        *string
	CategoryID      int32
	CategoryName    string
	CategorySlug    string
	SubCategoryID   int32
	SubCategoryName string
	SubCategorySlug string
	TypeID          *int32
	TypeName        *string
	TypeSlug        *string
	PromotionLevel  int32
	ModerateState   *string // только где нужно (модерация)
}

const productListSelect = `
	p.id, p.images, p.name, p.address, p."createdAt", p."isHide", p.price, p."userId", p."videoUrl",
	c.id, c.name, c.slug,
	sc.id, sc.name, sc.slug,
	t.id, t.name, t.slug,
	COALESCE((
		SELECT pr."pricePerDay" FROM "ProductPromotion" pp
		JOIN "Promotion" pr ON pr.id = pp."promotionId"
		WHERE pp."productId" = p.id AND pp."isActive" AND pp."isPaid" AND pp."endDate" >= NOW()
		ORDER BY pr."pricePerDay" DESC LIMIT 1
	), 0)::int,
	p."moderateState"::text`

const productListFrom = `
	FROM "Product" p
	JOIN "Category" c ON c.id = p."categoryId"
	JOIN "SubCategory" sc ON sc.id = p."subCategoryId"
	LEFT JOIN "SubcategotyType" t ON t.id = p."typeId"`

func (r *ProductPG) scanProductListRow(row pgx.Row) (*ProductListRow, error) {
	var pr ProductListRow
	var typeID *int32
	var typeName, typeSlug *string
	var mod *string
	err := row.Scan(
		&pr.ID, &pr.Images, &pr.Name, &pr.Address, &pr.CreatedAt, &pr.IsHide, &pr.Price, &pr.UserID, &pr.VideoURL,
		&pr.CategoryID, &pr.CategoryName, &pr.CategorySlug,
		&pr.SubCategoryID, &pr.SubCategoryName, &pr.SubCategorySlug,
		&typeID, &typeName, &typeSlug,
		&pr.PromotionLevel, &mod,
	)
	pr.TypeID, pr.TypeName, pr.TypeSlug = typeID, typeName, typeSlug
	pr.ModerateState = mod
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// ListProductsPublic — одобренные, не скрытые; limit<=0 — без LIMIT (как Nest findMany без take).
func (r *ProductPG) ListProductsPublic(ctx context.Context, orderClause string, limit, offset int) ([]ProductListRow, error) {
	base := fmt.Sprintf(`SELECT %s %s
		WHERE p."moderateState" = 'APPROVED'::"ProductModerate" AND p."isHide" = false
		ORDER BY %s`, productListSelect, productListFrom, orderClause)
	var q string
	var args []any
	if limit <= 0 {
		q = base
	} else {
		q = base + fmt.Sprintf(` LIMIT $1 OFFSET $2`)
		args = []any{limit, offset}
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProductListRows(rows)
}

func scanProductListRows(rows pgx.Rows) ([]ProductListRow, error) {
	var out []ProductListRow
	for rows.Next() {
		var pr ProductListRow
		var typeID *int32
		var typeName, typeSlug *string
		var mod *string
		err := rows.Scan(
			&pr.ID, &pr.Images, &pr.Name, &pr.Address, &pr.CreatedAt, &pr.IsHide, &pr.Price, &pr.UserID, &pr.VideoURL,
			&pr.CategoryID, &pr.CategoryName, &pr.CategorySlug,
			&pr.SubCategoryID, &pr.SubCategoryName, &pr.SubCategorySlug,
			&typeID, &typeName, &typeSlug,
			&pr.PromotionLevel, &mod,
		)
		if err != nil {
			return nil, err
		}
		pr.TypeID, pr.TypeName, pr.TypeSlug = typeID, typeName, typeSlug
		pr.ModerateState = mod
		out = append(out, pr)
	}
	return out, rows.Err()
}

// ListProductsModerate — на модерации.
func (r *ProductPG) ListProductsModerate(ctx context.Context) ([]ProductListRow, error) {
	q := fmt.Sprintf(`SELECT %s %s
		WHERE p."moderateState" IN ('MODERATE'::"ProductModerate", 'AI_REVIEWED'::"ProductModerate")
		ORDER BY p."createdAt" DESC`, productListSelect, productListFrom)
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProductListRows(rows)
}

func (r *ProductPG) ListProductsByUser(ctx context.Context, userID int32) ([]ProductListRow, error) {
	q := fmt.Sprintf(`SELECT %s %s WHERE p."userId" = $1 ORDER BY p."createdAt" DESC`, productListSelect, productListFrom)
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProductListRows(rows)
}

func (r *ProductPG) ListFavoriteProducts(ctx context.Context, userID int32) ([]ProductListRow, error) {
	q := fmt.Sprintf(`SELECT %s %s
		INNER JOIN "_UserFavorites" uf ON uf."A" = p.id AND uf."B" = $1
		ORDER BY p."createdAt" DESC`, productListSelect, productListFrom)
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProductListRows(rows)
}

// RandomSubcategoriesWithProducts — до 5 подкатегорий со случайным товаром (как Nest).
func (r *ProductPG) RandomSubcategoriesWithProducts(ctx context.Context) ([]ProductListRow, error) {
	const subQ = `
		SELECT sc.id FROM "SubCategory" sc
		WHERE EXISTS (
			SELECT 1 FROM "Product" p WHERE p."subCategoryId" = sc.id
		)
		ORDER BY RANDOM() LIMIT 5`
	rows, err := r.pool.Query(ctx, subQ)
	if err != nil {
		return nil, err
	}
	var subIDs []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		subIDs = append(subIDs, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var out []ProductListRow
	for _, sid := range subIDs {
		q := fmt.Sprintf(`SELECT %s %s WHERE p."subCategoryId" = $1 ORDER BY RANDOM() LIMIT 1`, productListSelect, productListFrom)
		row := r.pool.QueryRow(ctx, q, sid)
		pr, err := r.scanProductListRow(row)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			return nil, err
		}
		out = append(out, *pr)
	}
	return out, nil
}

// ProductSearchParams — фильтры для поиска.
type ProductSearchParams struct {
	Search        *string
	CategoryID    *int32
	SubCategoryID *int32
	TypeID        *int32
	MinPrice      *int32
	MaxPrice      *int32
	State         *string
	Region        *string
	ProfileType   *string
	FieldValues   map[string]string // fieldId -> contains
	SortBy        string
	Page, Limit   int
}

func (r *ProductPG) SearchProducts(ctx context.Context, p ProductSearchParams) ([]ProductListRow, error) {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	offset := (p.Page - 1) * p.Limit

	var where []string
	var args []any
	n := 1

	where = append(where, `p."moderateState" = 'APPROVED'::"ProductModerate"`)
	where = append(where, `p."isHide" = false`)

	if p.Search != nil && *p.Search != "" {
		if idVal, ok := parseInt32(*p.Search); ok {
			where = append(where, fmt.Sprintf(`p.id = $%d`, n))
			args = append(args, idVal)
			n++
		} else {
			pat := "%" + *p.Search + "%"
			where = append(where, fmt.Sprintf(`(p.name ILIKE $%d OR p.description ILIKE $%d)`, n, n+1))
			args = append(args, pat, pat)
			n += 2
		}
	}
	if p.CategoryID != nil {
		where = append(where, fmt.Sprintf(`p."categoryId" = $%d`, n))
		args = append(args, *p.CategoryID)
		n++
	}
	if p.SubCategoryID != nil {
		where = append(where, fmt.Sprintf(`p."subCategoryId" = $%d`, n))
		args = append(args, *p.SubCategoryID)
		n++
	}
	if p.TypeID != nil {
		where = append(where, fmt.Sprintf(`p."typeId" = $%d`, n))
		args = append(args, *p.TypeID)
		n++
	}
	if p.MinPrice != nil || p.MaxPrice != nil {
		if p.MinPrice != nil {
			where = append(where, fmt.Sprintf(`p.price >= $%d`, n))
			args = append(args, *p.MinPrice)
			n++
		}
		if p.MaxPrice != nil {
			where = append(where, fmt.Sprintf(`p.price <= $%d`, n))
			args = append(args, *p.MaxPrice)
			n++
		}
	}
	if p.State != nil && *p.State != "" {
		where = append(where, fmt.Sprintf(`p.state = $%d::"ProductState"`, n))
		args = append(args, *p.State)
		n++
	}
	if p.Region != nil && *p.Region != "" {
		where = append(where, fmt.Sprintf(`p.address ILIKE $%d`, n))
		args = append(args, "%"+*p.Region+"%")
		n++
	}
	if p.ProfileType != nil && *p.ProfileType != "" {
		where = append(where, fmt.Sprintf(`EXISTS (SELECT 1 FROM "User" u WHERE u.id = p."userId" AND u."profileType" = $%d::"ProfileType")`, n))
		args = append(args, *p.ProfileType)
		n++
	}
	if len(p.FieldValues) > 0 {
		var orParts []string
		for fid, val := range p.FieldValues {
			orParts = append(orParts, fmt.Sprintf(`(pfv."fieldId" = $%d AND pfv.value ILIKE $%d)`, n, n+1))
			id32, _ := parseInt32(fid)
			args = append(args, id32, "%"+val+"%")
			n += 2
		}
		where = append(where, `EXISTS (SELECT 1 FROM "ProductFieldValue" pfv WHERE pfv."productId" = p.id AND (`+strings.Join(orParts, " OR ")+`))`)
	}

	order := `p."createdAt" DESC`
	switch p.SortBy {
	case "price_asc":
		order = `p.price ASC`
	case "price_desc":
		order = `p.price DESC`
	case "date_asc":
		order = `p."createdAt" ASC`
	case "date_desc", "relevance", "":
		order = `p."createdAt" DESC`
	}

	whereSQL := strings.Join(where, " AND ")
	q := fmt.Sprintf(`SELECT %s %s WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d`,
		productListSelect, productListFrom, whereSQL, order, n, n+1)
	args = append(args, p.Limit, offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProductListRows(rows)
}

func parseInt32(s string) (int32, bool) {
	var v int64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &v)
	if err != nil || v <= 0 || v > 1<<31-1 {
		return 0, false
	}
	return int32(v), true
}

// ProductCardDB — сырые поля карточки товара.
type ProductCardDB struct {
	ID           int32
	Name         string
	Description  *string
	Price        int32
	IsHide       bool
	Images       []string
	Address      string
	UserID       int32
	VideoURL     *string
	CategoryID   int32
	CategoryName string
	CategorySlug string
	SubCatID     int32
	SubCatName   string
	SubCatSlug   string
	TypeID       *int32
	TypeName     *string
	TypeSlug     *string
	FieldPairs   []struct{ FieldName, Value string }
}

func (r *ProductPG) GetProductCard(ctx context.Context, productID int32) (*ProductCardDB, error) {
	const q = `
		SELECT p.id, p.name, p.description, p.price, p."isHide", p.images, p.address, p."userId", p."videoUrl",
			c.id, c.name, c.slug, sc.id, sc.name, sc.slug, t.id, t.name, t.slug
		FROM "Product" p
		JOIN "Category" c ON c.id = p."categoryId"
		JOIN "SubCategory" sc ON sc.id = p."subCategoryId"
		LEFT JOIN "SubcategotyType" t ON t.id = p."typeId"
		WHERE p.id = $1`
	var card ProductCardDB
	var typeID *int32
	var typeName, typeSlug *string
	err := r.pool.QueryRow(ctx, q, productID).Scan(
		&card.ID, &card.Name, &card.Description, &card.Price, &card.IsHide, &card.Images, &card.Address, &card.UserID, &card.VideoURL,
		&card.CategoryID, &card.CategoryName, &card.CategorySlug,
		&card.SubCatID, &card.SubCatName, &card.SubCatSlug,
		&typeID, &typeName, &typeSlug,
	)
	card.TypeID, card.TypeName, card.TypeSlug = typeID, typeName, typeSlug
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	frows, err := r.pool.Query(ctx, `
		SELECT tf.name, pfv.value FROM "ProductFieldValue" pfv
		JOIN "TypeField" tf ON tf.id = pfv."fieldId"
		WHERE pfv."productId" = $1`, productID)
	if err != nil {
		return nil, err
	}
	defer frows.Close()
	for frows.Next() {
		var fn, fv string
		if err := frows.Scan(&fn, &fv); err != nil {
			return nil, err
		}
		card.FieldPairs = append(card.FieldPairs, struct{ FieldName, Value string }{fn, fv})
	}
	return &card, frows.Err()
}

// LoadProductWithRelations — ответ после create/update (как Prisma include).
func (r *ProductPG) LoadProductWithRelations(ctx context.Context, productID int32) (map[string]any, error) {
	card, err := r.GetProductCard(ctx, productID)
	if err != nil {
		return nil, err
	}
	urow := r.pool.QueryRow(ctx, `SELECT id, "fullName", email, rating FROM "User" WHERE id = $1`, card.UserID)
	var uid int32
	var fullName, email string
	var rating *int32
	if err := urow.Scan(&uid, &fullName, &email, &rating); err != nil {
		return nil, err
	}
	fv := make([]map[string]any, 0, len(card.FieldPairs))
	for _, p := range card.FieldPairs {
		fv = append(fv, map[string]any{p.FieldName: p.Value})
	}
	out := map[string]any{
		"id": card.ID, "name": card.Name, "description": card.Description, "price": card.Price,
		"state":  nil, // заполним отдельно если нужно
		"images": card.Images, "address": card.Address, "videoUrl": card.VideoURL,
		"category":    map[string]any{"id": card.CategoryID, "name": card.CategoryName},
		"subCategory": map[string]any{"id": card.SubCatID, "name": card.SubCatName},
		"user":        map[string]any{"id": uid, "fullName": fullName, "email": email, "rating": rating},
		"fieldValues": fv,
	}
	// state из Product
	var st string
	_ = r.pool.QueryRow(ctx, `SELECT state::text FROM "Product" WHERE id = $1`, productID).Scan(&st)
	out["state"] = st
	return out, nil
}

// PromotedProductRow — админский список продвижений.
type PromotedProductRow struct {
	ProductID   int32
	Name        string
	Price       int32
	Images      []string
	Category    map[string]any
	SubCategory map[string]any
	Type        map[string]any
	User        map[string]any
	CreatedAt   time.Time
	Promotions  []map[string]any
}

func (r *ProductPG) ListPromotedProducts(ctx context.Context) ([]PromotedProductRow, error) {
	const q = `
		SELECT p.id FROM "Product" p
		WHERE EXISTS (
			SELECT 1 FROM "ProductPromotion" pp
			WHERE pp."productId" = p.id AND pp."isPaid" = true AND pp."endDate" >= NOW()
		)
		ORDER BY p."createdAt" DESC`
	ids, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer ids.Close()
	var pids []int32
	for ids.Next() {
		var id int32
		if err := ids.Scan(&id); err != nil {
			return nil, err
		}
		pids = append(pids, id)
	}
	if err := ids.Err(); err != nil {
		return nil, err
	}

	var out []PromotedProductRow
	for _, pid := range pids {
		row, err := r.loadOnePromoted(ctx, pid)
		if err != nil {
			return nil, err
		}
		out = append(out, *row)
	}
	return out, nil
}

func (r *ProductPG) loadOnePromoted(ctx context.Context, productID int32) (*PromotedProductRow, error) {
	const q = `
		SELECT p.id, p.name, p.price, p.images, p."createdAt",
			c.id, c.name, c.slug, sc.id, sc.name, sc.slug, t.id, t.name, t.slug,
			u.id, u."fullName", u.email
		FROM "Product" p
		JOIN "Category" c ON c.id = p."categoryId"
		JOIN "SubCategory" sc ON sc.id = p."subCategoryId"
		LEFT JOIN "SubcategotyType" t ON t.id = p."typeId"
		JOIN "User" u ON u.id = p."userId"
		WHERE p.id = $1`
	var pr PromotedProductRow
	var cid, scid, tid *int32
	var cname, cslug, scname, scslug, tname, tslug *string
	var uid int32
	var ufn, uem string
	err := r.pool.QueryRow(ctx, q, productID).Scan(
		&pr.ProductID, &pr.Name, &pr.Price, &pr.Images, &pr.CreatedAt,
		&cid, &cname, &cslug, &scid, &scname, &scslug, &tid, &tname, &tslug,
		&uid, &ufn, &uem,
	)
	if err != nil {
		return nil, err
	}
	pr.Category = map[string]any{"id": cid, "name": cname, "slug": cslug}
	pr.SubCategory = map[string]any{"id": scid, "name": scname, "slug": scslug}
	pr.Type = map[string]any{"id": tid, "name": tname, "slug": tslug}
	pr.User = map[string]any{"id": uid, "fullName": ufn, "email": uem}

	promRows, err := r.pool.Query(ctx, `
		SELECT pp.id, pr.name, pr."pricePerDay", pp."startDate", pp."endDate", pp."isActive", pp.days, pp."totalPrice"
		FROM "ProductPromotion" pp
		JOIN "Promotion" pr ON pr.id = pp."promotionId"
		WHERE pp."productId" = $1 AND pp."isPaid" = true AND pp."endDate" >= NOW()
		ORDER BY pr."pricePerDay" DESC`, productID)
	if err != nil {
		return nil, err
	}
	defer promRows.Close()
	for promRows.Next() {
		var id int32
		var pname string
		var ppd, days, total int32
		var start, end time.Time
		var active bool
		if err := promRows.Scan(&id, &pname, &ppd, &start, &end, &active, &days, &total); err != nil {
			return nil, err
		}
		pr.Promotions = append(pr.Promotions, map[string]any{
			"id": id, "promotionType": pname, "pricePerDay": ppd,
			"startDate": start, "endDate": end, "isActive": active, "days": days, "totalPrice": total,
		})
	}
	return &pr, promRows.Err()
}
