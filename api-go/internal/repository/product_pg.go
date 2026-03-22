package repository

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdListingCost — стоимость платного объявления (как в Nest).
const AdListingCost = 50

// ErrInsufficientFunds — не хватает денег на платное объявление.
var ErrInsufficientFunds = errors.New("INSUFFICIENT_FUNDS")

type ProductPG struct {
	pool *pgxpool.Pool
}

func NewProductPG(pool *pgxpool.Pool) *ProductPG {
	return &ProductPG{pool: pool}
}

// GenerateUniqueProductID — 7-значный id товара, как в Nest id-generator.
func (r *ProductPG) GenerateUniqueProductID(ctx context.Context) (int32, error) {
	const minID, maxID = 1_000_000, 9_999_999
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for attempt := 0; attempt < 100; attempt++ {
		id := int32(minID + rng.Intn(maxID-minID+1))
		var one int
		err := r.pool.QueryRow(ctx, `SELECT 1 FROM "Product" WHERE id = $1`, id).Scan(&one)
		if errors.Is(err, pgx.ErrNoRows) {
			return id, nil
		}
		if err != nil {
			return 0, err
		}
	}
	return 0, fmt.Errorf("не удалось сгенерировать id товара")
}

func (r *ProductPG) SubCategoryBelongsToCategory(ctx context.Context, categoryID, subCategoryID int32) (bool, error) {
	const q = `SELECT 1 FROM "SubCategory" WHERE id = $1 AND "categoryId" = $2`
	var x int
	err := r.pool.QueryRow(ctx, q, subCategoryID, categoryID).Scan(&x)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *ProductPG) ValidateTypeFieldIDs(ctx context.Context, fieldIDs []int32) error {
	if len(fieldIDs) == 0 {
		return nil
	}
	rows, err := r.pool.Query(ctx, `SELECT id FROM "TypeField" WHERE id = ANY($1::int[])`, fieldIDs)
	if err != nil {
		return err
	}
	defer rows.Close()
	found := make(map[int32]struct{})
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return err
		}
		found[id] = struct{}{}
	}
	var missing []string
	for _, id := range fieldIDs {
		if _, ok := found[id]; !ok {
			missing = append(missing, fmt.Sprintf("%d", id))
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("поля с ID %s не найдены", strings.Join(missing, ", "))
	}
	return nil
}

func (r *ProductPG) TypeFieldsBelongToType(ctx context.Context, typeID int32, fieldIDs []int32) ([]string, error) {
	if len(fieldIDs) == 0 {
		return nil, nil
	}
	q := `SELECT name FROM "TypeField" WHERE id = ANY($1::int[]) AND "typeId" <> $2`
	rows, err := r.pool.Query(ctx, q, fieldIDs, typeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		names = append(names, n)
	}
	return names, rows.Err()
}

// CreateProductTx — списание/лимит + вставка товара и fieldValues в одной транзакции.
func (r *ProductPG) CreateProductTx(ctx context.Context, productID, userID int32, name string, price int32, state, description, address string, video *string,
	images []string, categoryID, subCategoryID int32, typeID *int32, fieldValues map[int32]string,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var freeLimit, usedAds int32
	var lastReset time.Time
	var balance, bonus float64
	err = tx.QueryRow(ctx, `
		SELECT "freeAdsLimit", "usedFreeAds", "lastAdLimitReset", "balance", "bonusBalance"
		FROM "User" WHERE id = $1 FOR UPDATE`, userID).Scan(&freeLimit, &usedAds, &lastReset, &balance, &bonus)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	lr := lastReset.UTC()
	if now.Month() != lr.Month() || now.Year() != lr.Year() {
		usedAds = 0
		_, err = tx.Exec(ctx, `UPDATE "User" SET "usedFreeAds" = 0, "lastAdLimitReset" = $2, "updatedAt" = NOW() WHERE id = $1`, userID, now)
		if err != nil {
			return err
		}
	}

	if usedAds >= freeLimit {
		total := balance + bonus
		if total < float64(AdListingCost) {
			return ErrInsufficientFunds
		}
		rem := float64(AdListingCost)
		bonusDeduct := bonus
		if bonusDeduct > rem {
			bonusDeduct = rem
		}
		rem -= bonusDeduct
		balDeduct := rem
		_, err = tx.Exec(ctx, `
			UPDATE "User" SET "bonusBalance" = "bonusBalance" - $2, "balance" = "balance" - $3, "updatedAt" = NOW() WHERE id = $1`,
			userID, bonusDeduct, balDeduct)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec(ctx, `UPDATE "User" SET "usedFreeAds" = "usedFreeAds" + 1, "updatedAt" = NOW() WHERE id = $1`, userID)
		if err != nil {
			return err
		}
	}

	var typeVal any
	if typeID != nil {
		typeVal = *typeID
	} else {
		typeVal = nil
	}
	descPtr := nullIfEmpty(description)
	_, err = tx.Exec(ctx, `
		INSERT INTO "Product" (id, name, price, state, description, address, "videoUrl", images, "isHide", "moderateState", "categoryId", "subCategoryId", "typeId", "userId", "createdAt", "updatedAt")
		VALUES ($1,$2,$3,$4::"ProductState",$5,$6,$7,$8::text[], false, 'MODERATE'::"ProductModerate",$9,$10,$11,$12,NOW(),NOW())`,
		productID, name, price, state, descPtr, address, video, images, categoryID, subCategoryID, typeVal, userID)
	if err != nil {
		return err
	}

	for fid, val := range fieldValues {
		_, err = tx.Exec(ctx, `INSERT INTO "ProductFieldValue" ("fieldId", "productId", value) VALUES ($1,$2,$3)`, fid, productID, val)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func nullIfEmpty(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

func (r *ProductPG) GetProductOwnerAndImages(ctx context.Context, productID int32) (userID int32, images []string, err error) {
	const q = `SELECT "userId", images FROM "Product" WHERE id = $1`
	err = r.pool.QueryRow(ctx, q, productID).Scan(&userID, &images)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil, ErrNotFound
	}
	return userID, images, err
}

func (r *ProductPG) DeleteProductByID(ctx context.Context, productID int32) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM "Product" WHERE id = $1`, productID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ProductWithTypeForUpdate — для PATCH: typeId, images, userId.
func (r *ProductPG) ProductWithTypeForUpdate(ctx context.Context, productID int32) (userID int32, typeID *int32, images []string, err error) {
	const q = `SELECT "userId", "typeId", images FROM "Product" WHERE id = $1`
	err = r.pool.QueryRow(ctx, q, productID).Scan(&userID, &typeID, &images)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil, nil, ErrNotFound
	}
	return userID, typeID, images, err
}

func (r *ProductPG) UpdateProductPartial(ctx context.Context, productID int32, name *string, price *int32, state *string, description *string, address *string, video *string, images []string) error {
	var sets []string
	var args []any
	n := 1
	if name != nil {
		sets = append(sets, fmt.Sprintf(`name = $%d`, n))
		args = append(args, *name)
		n++
	}
	if price != nil {
		sets = append(sets, fmt.Sprintf(`price = $%d`, n))
		args = append(args, *price)
		n++
	}
	if state != nil {
		sets = append(sets, fmt.Sprintf(`state = $%d::"ProductState"`, n))
		args = append(args, *state)
		n++
	}
	if description != nil {
		sets = append(sets, fmt.Sprintf(`description = $%d`, n))
		args = append(args, *description)
		n++
	}
	if address != nil {
		sets = append(sets, fmt.Sprintf(`address = $%d`, n))
		args = append(args, *address)
		n++
	}
	if video != nil {
		sets = append(sets, fmt.Sprintf(`"videoUrl" = $%d`, n))
		args = append(args, *video)
		n++
	}
	if images != nil {
		sets = append(sets, fmt.Sprintf(`images = $%d::text[]`, n))
		args = append(args, images)
		n++
	}
	if len(sets) == 0 {
		return nil
	}
	sets = append(sets, `"updatedAt" = NOW()`)
	args = append(args, productID)
	q := fmt.Sprintf(`UPDATE "Product" SET %s WHERE id = $%d`, strings.Join(sets, ", "), n)
	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProductPG) UpsertProductFieldValue(ctx context.Context, productID, fieldID int32, value string) error {
	const q = `
		INSERT INTO "ProductFieldValue" ("fieldId", "productId", value) VALUES ($1,$2,$3)
		ON CONFLICT ("fieldId", "productId") DO UPDATE SET value = EXCLUDED.value`
	_, err := r.pool.Exec(ctx, q, fieldID, productID, value)
	return err
}

func (r *ProductPG) CategoryIDBySlug(ctx context.Context, slug string) (*int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `SELECT id FROM "Category" WHERE slug = $1`, slug).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *ProductPG) SubCategoryIDBySlug(ctx context.Context, slug string) (*int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `SELECT id FROM "SubCategory" WHERE slug = $1 LIMIT 1`, slug).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *ProductPG) TypeIDBySlug(ctx context.Context, slug string) (*int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `SELECT id FROM "SubcategotyType" WHERE slug = $1 LIMIT 1`, slug).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *ProductPG) TypeFieldsForType(ctx context.Context, typeID int32) ([]struct {
	ID         int32
	Name       string
	IsRequired bool
}, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, "isRequired" FROM "TypeField" WHERE "typeId" = $1`, typeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct {
		ID         int32
		Name       string
		IsRequired bool
	}
	for rows.Next() {
		var row struct {
			ID         int32
			Name       string
			IsRequired bool
		}
		if err := rows.Scan(&row.ID, &row.Name, &row.IsRequired); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *ProductPG) DistinctFieldValuesForFilter(ctx context.Context, fieldID int32, whereSQL string, whereArgs []any) ([]string, error) {
	q := `SELECT DISTINCT pfv.value FROM "ProductFieldValue" pfv
		INNER JOIN "Product" p ON p.id = pfv."productId"
		WHERE pfv."fieldId" = $1 AND ` + whereSQL + ` AND trim(pfv.value) <> ''`
	args := append([]any{fieldID}, whereArgs...)
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vals []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, rows.Err()
}

func (r *ProductPG) AggregateProductPriceRange(ctx context.Context, whereSQL string, whereArgs []any) (min, max int32, err error) {
	q := `SELECT COALESCE(MIN(price),0)::int, COALESCE(MAX(price),0)::int FROM "Product" p WHERE ` + whereSQL
	err = r.pool.QueryRow(ctx, q, whereArgs...).Scan(&min, &max)
	return min, max, err
}

func (r *ProductPG) DistinctStates(ctx context.Context, whereSQL string, whereArgs []any) ([]string, error) {
	q := `SELECT DISTINCT state::text FROM "Product" p WHERE ` + whereSQL
	rows, err := r.pool.Query(ctx, q, whereArgs...)
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

func (r *ProductPG) DistinctSellerProfileTypes(ctx context.Context, whereSQL string, whereArgs []any) ([]string, error) {
	q := `SELECT DISTINCT u."profileType"::text FROM "Product" p JOIN "User" u ON u.id = p."userId" WHERE ` + whereSQL
	rows, err := r.pool.Query(ctx, q, whereArgs...)
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
		if s != "" {
			out = append(out, s)
		}
	}
	return out, rows.Err()
}

func approvedProductWhereFragments(categoryID, subCategoryID, typeID *int32) (sql string, args []any) {
	var parts []string
	n := 1
	parts = append(parts, `p."moderateState" = 'APPROVED'::"ProductModerate"`)
	parts = append(parts, `p."isHide" = false`)
	if categoryID != nil {
		parts = append(parts, fmt.Sprintf(`p."categoryId" = $%d`, n))
		args = append(args, *categoryID)
		n++
	}
	if subCategoryID != nil {
		parts = append(parts, fmt.Sprintf(`p."subCategoryId" = $%d`, n))
		args = append(args, *subCategoryID)
		n++
	}
	if typeID != nil {
		parts = append(parts, fmt.Sprintf(`p."typeId" = $%d`, n))
		args = append(args, *typeID)
		n++
	}
	return strings.Join(parts, " AND "), args
}

// BuildAvailableFilters — как getAvailableFilters в Nest.
func (r *ProductPG) BuildAvailableFilters(ctx context.Context, categorySlug, subCategorySlug, typeSlug *string) (map[string]any, error) {
	var catID, subID, typID *int32
	var err error
	if categorySlug != nil && *categorySlug != "" {
		catID, err = r.CategoryIDBySlug(ctx, *categorySlug)
		if err != nil {
			return nil, err
		}
	}
	if subCategorySlug != nil && *subCategorySlug != "" {
		subID, err = r.SubCategoryIDBySlug(ctx, *subCategorySlug)
		if err != nil {
			return nil, err
		}
	}
	if typeSlug != nil && *typeSlug != "" {
		typID, err = r.TypeIDBySlug(ctx, *typeSlug)
		if err != nil {
			return nil, err
		}
	}
	whereSQL, baseArgs := approvedProductWhereFragments(catID, subID, typID)

	fields := []map[string]any{}
	if typID != nil {
		tf, err := r.TypeFieldsForType(ctx, *typID)
		if err != nil {
			return nil, err
		}
		for _, f := range tf {
			vals, err := r.DistinctFieldValuesForFilter(ctx, f.ID, whereSQL, baseArgs)
			if err != nil {
				return nil, err
			}
			if len(vals) == 0 {
				continue
			}
			fields = append(fields, map[string]any{
				"fieldId": f.ID, "fieldName": f.Name, "isRequired": f.IsRequired, "values": vals,
			})
		}
	}
	minP, maxP, err := r.AggregateProductPriceRange(ctx, whereSQL, baseArgs)
	if err != nil {
		return nil, err
	}
	states, err := r.DistinctStates(ctx, whereSQL, baseArgs)
	if err != nil {
		return nil, err
	}
	profiles, err := r.DistinctSellerProfileTypes(ctx, whereSQL, baseArgs)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"fields":       fields,
		"priceRange":   map[string]any{"min": minP, "max": maxP},
		"states":       states,
		"profileTypes": profiles,
	}, nil
}

func (r *ProductPG) UserExists(ctx context.Context, id int32) (bool, error) {
	var one int
	err := r.pool.QueryRow(ctx, `SELECT 1 FROM "User" WHERE id = $1`, id).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *ProductPG) IsFavorite(ctx context.Context, userID, productID int32) (bool, error) {
	const q = `SELECT 1 FROM "_UserFavorites" WHERE "A" = $1 AND "B" = $2`
	var one int
	err := r.pool.QueryRow(ctx, q, productID, userID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *ProductPG) AddFavorite(ctx context.Context, userID, productID int32) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO "_UserFavorites" ("A","B") VALUES ($1,$2) ON CONFLICT ("A","B") DO NOTHING`, productID, userID)
	return err
}

func (r *ProductPG) RemoveFavorite(ctx context.Context, userID, productID int32) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM "_UserFavorites" WHERE "A" = $1 AND "B" = $2`, productID, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProductPG) InsertFavoriteAction(ctx context.Context, userID, productID int32) {
	_, _ = r.pool.Exec(ctx, `INSERT INTO "FavoriteAction" ("userId", "productId") VALUES ($1,$2) ON CONFLICT ("userId", "productId") DO NOTHING`, userID, productID)
}

func (r *ProductPG) ProductExists(ctx context.Context, id int32) (bool, error) {
	var one int
	err := r.pool.QueryRow(ctx, `SELECT 1 FROM "Product" WHERE id = $1`, id).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *ProductPG) UpsertProductView(ctx context.Context, viewerID, productID int32) {
	_, _ = r.pool.Exec(ctx, `
		INSERT INTO "ProductView" ("viewedById", "productId", "viewedAt") VALUES ($1,$2,NOW())
		ON CONFLICT ("viewedById", "productId") DO UPDATE SET "viewedAt" = EXCLUDED."viewedAt"`,
		viewerID, productID)
}

func (r *ProductPG) ToggleProductHide(ctx context.Context, productID int32) error {
	const q = `UPDATE "Product" SET "isHide" = NOT "isHide", "updatedAt" = NOW() WHERE id = $1`
	tag, err := r.pool.Exec(ctx, q, productID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProductPG) SetModerationState(ctx context.Context, productID int32, state string, reason *string) error {
	var q string
	var args []any
	if reason != nil {
		q = `UPDATE "Product" SET "moderateState" = $2::"ProductModerate", "moderationRejectionReason" = $3, "updatedAt" = NOW() WHERE id = $1`
		args = []any{productID, state, *reason}
	} else {
		q = `UPDATE "Product" SET "moderateState" = $2::"ProductModerate", "moderationRejectionReason" = NULL, "updatedAt" = NOW() WHERE id = $1`
		args = []any{productID, state}
	}
	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProductPG) GetProductNameAndSeller(ctx context.Context, productID int32) (name string, sellerID int32, err error) {
	const q = `SELECT name, "userId" FROM "Product" WHERE id = $1`
	err = r.pool.QueryRow(ctx, q, productID).Scan(&name, &sellerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", 0, ErrNotFound
	}
	return name, sellerID, err
}

func (r *ProductPG) FirstAdminUserID(ctx context.Context) (*int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `
		SELECT u.id FROM "User" u JOIN "Role" r ON r.id = u."roleId" WHERE r.name = 'admin' LIMIT 1`).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *ProductPG) FindModerationChat(ctx context.Context, adminID, userID int32) (*int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `
		SELECT id FROM "Chat"
		WHERE "isModerationChat" = true AND "productId" IS NULL
		AND (("buyerId" = $1 AND "sellerId" = $2) OR ("buyerId" = $2 AND "sellerId" = $1))
		LIMIT 1`, adminID, userID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *ProductPG) CreateModerationChat(ctx context.Context, adminID, userID int32) (int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `
		INSERT INTO "Chat" ("buyerId", "sellerId", "productId", "isModerationChat", "lastMessageAt", "createdAt", "updatedAt")
		VALUES ($1,$2,NULL,true,NOW(),NOW(),NOW()) RETURNING id`, adminID, userID).Scan(&id)
	return id, err
}

func (r *ProductPG) InsertChatMessage(ctx context.Context, chatID, senderID int32, content string, relatedProductID int32) (int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `
		INSERT INTO "Message" (content, "senderId", "chatId", "relatedProductId", "createdAt", "updatedAt")
		VALUES ($1,$2,$3,$4,NOW(),NOW()) RETURNING id`, content, senderID, chatID, relatedProductID).Scan(&id)
	return id, err
}

func (r *ProductPG) UpdateChatAfterMessage(ctx context.Context, chatID, messageID, senderID, buyerID int32) error {
	var unreadBuyer, unreadSeller int
	if senderID == buyerID {
		unreadSeller = 1
	} else {
		unreadBuyer = 1
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE "Chat" SET "lastMessageId" = $2, "lastMessageAt" = NOW(),
		"unreadCountBuyer" = "unreadCountBuyer" + $3, "unreadCountSeller" = "unreadCountSeller" + $4, "updatedAt" = NOW()
		WHERE id = $1`, chatID, messageID, unreadBuyer, unreadSeller)
	return err
}

func (r *ProductPG) GetChatBuyerSeller(ctx context.Context, chatID int32) (buyerID, sellerID int32, err error) {
	err = r.pool.QueryRow(ctx, `SELECT "buyerId", "sellerId" FROM "Chat" WHERE id = $1`, chatID).Scan(&buyerID, &sellerID)
	return buyerID, sellerID, err
}

func (r *ProductPG) TogglePromotionActive(ctx context.Context, promotionRowID int32) (productID int32, productName string, promoName string, isActive bool, start, end time.Time, err error) {
	var curActive bool
	err = r.pool.QueryRow(ctx, `SELECT pp."isActive", pp."productId", pr.name, pp."startDate", pp."endDate", p.name
		FROM "ProductPromotion" pp
		JOIN "Promotion" pr ON pr.id = pp."promotionId"
		JOIN "Product" p ON p.id = pp."productId"
		WHERE pp.id = $1`, promotionRowID).Scan(&curActive, &productID, &promoName, &start, &end, &productName)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, "", "", false, time.Time{}, time.Time{}, ErrNotFound
	}
	if err != nil {
		return 0, "", "", false, time.Time{}, time.Time{}, err
	}
	newActive := !curActive
	_, err = r.pool.Exec(ctx, `UPDATE "ProductPromotion" SET "isActive" = $2, "updatedAt" = NOW() WHERE id = $1`, promotionRowID, newActive)
	if err != nil {
		return 0, "", "", false, time.Time{}, time.Time{}, err
	}
	return productID, productName, promoName, newActive, start, end, nil
}
