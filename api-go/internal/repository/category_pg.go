package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"med-vito/api-go/internal/domain"
)

var ErrNotFound = errors.New("not found")

// CategoryPG вЂ” SQL Рє С‚Р°Р±Р»РёС†Р°Рј Prisma (РєР°РІС‹С‡РєРё РІ РёРґРµРЅС‚РёС„РёРєР°С‚РѕСЂР°С… РєР°Рє РІ РјРёРіСЂР°С†РёРё).
type CategoryPG struct {
	pool *pgxpool.Pool
}

func NewCategoryPG(pool *pgxpool.Pool) *CategoryPG {
	return &CategoryPG{pool: pool}
}

func joinInt32(ids []int32) string {
	if len(ids) == 0 {
		return ""
	}
	parts := make([]string, len(ids))
	for i, v := range ids {
		parts[i] = strconv.FormatInt(int64(v), 10)
	}
	return strings.Join(parts, ",")
}

// --- Category CRUD ---

func (r *CategoryPG) FindCategoryBySlug(ctx context.Context, slug string) (*domain.CategoryScalars, error) {
	const q = `SELECT "id", "name", "slug", "createdAt", "updatedAt" FROM "Category" WHERE "slug" = $1`
	row := r.pool.QueryRow(ctx, q, slug)
	var c domain.CategoryScalars
	if err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryPG) FindCategoryByID(ctx context.Context, id int32) (*domain.CategoryScalars, error) {
	const q = `SELECT "id", "name", "slug", "createdAt", "updatedAt" FROM "Category" WHERE "id" = $1`
	row := r.pool.QueryRow(ctx, q, id)
	var c domain.CategoryScalars
	if err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryPG) ExistsCategoryBySlug(ctx context.Context, slug string) (bool, error) {
	const q = `SELECT 1 FROM "Category" WHERE "slug" = $1 LIMIT 1`
	var one int
	err := r.pool.QueryRow(ctx, q, slug).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *CategoryPG) ExistsCategoryBySlugExceptID(ctx context.Context, slug string, exceptID int32) (bool, error) {
	const q = `SELECT 1 FROM "Category" WHERE "slug" = $1 AND "id" <> $2 LIMIT 1`
	var one int
	err := r.pool.QueryRow(ctx, q, slug, exceptID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *CategoryPG) ListAllCategorySlugs(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT "slug" FROM "Category"`)
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

func (r *CategoryPG) ListCategorySlugsExceptID(ctx context.Context, exceptID int32) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT "slug" FROM "Category" WHERE "id" <> $1`, exceptID)
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

func (r *CategoryPG) ListCategoriesOrderByName(ctx context.Context) ([]domain.CategoryTree, error) {
	rows, err := r.pool.Query(ctx, `SELECT "id", "name", "slug" FROM "Category" ORDER BY "name" ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.CategoryTree
	for rows.Next() {
		var c domain.CategoryTree
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return nil, err
		}
		c.SubCategories = []domain.SubCategoryTree{}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CategoryPG) InsertCategory(ctx context.Context, name, slug string) (*domain.CategoryScalars, error) {
	const q = `
		INSERT INTO "Category" ("name", "slug")
		VALUES ($1, $2)
		RETURNING "id", "name", "slug", "createdAt", "updatedAt"`
	row := r.pool.QueryRow(ctx, q, name, slug)
	var c domain.CategoryScalars
	if err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryPG) UpdateCategory(ctx context.Context, id int32, name string, slug *string) (*domain.CategoryScalars, error) {
	var q string
	var row pgx.Row
	if slug != nil {
		q = `UPDATE "Category" SET "name" = $2, "slug" = $3, "updatedAt" = NOW() WHERE "id" = $1
			RETURNING "id", "name", "slug", "createdAt", "updatedAt"`
		row = r.pool.QueryRow(ctx, q, id, name, *slug)
	} else {
		q = `UPDATE "Category" SET "name" = $2, "updatedAt" = NOW() WHERE "id" = $1
			RETURNING "id", "name", "slug", "createdAt", "updatedAt"`
		row = r.pool.QueryRow(ctx, q, id, name)
	}
	var c domain.CategoryScalars
	if err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryPG) DeleteCategory(ctx context.Context, id int32) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM "Category" WHERE "id" = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// --- SubCategory ---

type subCategoryRow struct {
	ID         int32
	Name       string
	Slug       string
	CategoryID int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (r *CategoryPG) ListSubCategoriesByCategoryIDs(ctx context.Context, ids []int32) ([]subCategoryRow, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := fmt.Sprintf(
		`SELECT "id", "name", "slug", "categoryId", "createdAt", "updatedAt" FROM "SubCategory" WHERE "categoryId" IN (%s)`,
		joinInt32(ids),
	)
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []subCategoryRow
	for rows.Next() {
		var sc subCategoryRow
		if err := rows.Scan(&sc.ID, &sc.Name, &sc.Slug, &sc.CategoryID, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

func (r *CategoryPG) FindSubCategoryByID(ctx context.Context, id int32) (*domain.SubCategoryWithCategory, error) {
	const q = `SELECT "id", "name", "slug", "categoryId", "createdAt", "updatedAt" FROM "SubCategory" WHERE "id" = $1`
	row := r.pool.QueryRow(ctx, q, id)
	var sc domain.SubCategoryWithCategory
	if err := row.Scan(&sc.ID, &sc.Name, &sc.Slug, &sc.CategoryID, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	cat, err := r.FindCategoryByID(ctx, sc.CategoryID)
	if err != nil {
		return nil, err
	}
	sc.Category = *cat
	types, err := r.listSubcategoryTypesBySubcategoryID(ctx, sc.ID)
	if err != nil {
		return nil, err
	}
	sc.SubcategoryTypes = types
	return &sc, nil
}

// --- SubcategotyType (РѕРїРµС‡Р°С‚РєР° РІ Prisma) ---

type typeRow struct {
	ID            int32
	Name          string
	Slug          string
	SubcategoryID int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (r *CategoryPG) listSubcategoryTypesBySubcategoryIDs(ctx context.Context, subIDs []int32) ([]typeRow, error) {
	if len(subIDs) == 0 {
		return nil, nil
	}
	q := fmt.Sprintf(
		`SELECT "id", "name", "slug", "subcategoryId", "createdAt", "updatedAt" FROM "SubcategotyType" WHERE "subcategoryId" IN (%s)`,
		joinInt32(subIDs),
	)
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []typeRow
	for rows.Next() {
		var t typeRow
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.SubcategoryID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *CategoryPG) listSubcategoryTypesBySubcategoryID(ctx context.Context, subID int32) ([]domain.SubcategoryTypeDetail, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT "id", "name", "slug", "subcategoryId", "createdAt", "updatedAt" FROM "SubcategotyType" WHERE "subcategoryId" = $1`,
		subID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int32
	var list []domain.SubcategoryTypeDetail
	for rows.Next() {
		var t domain.SubcategoryTypeDetail
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.SubcategoryID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		t.Fields = []domain.TypeField{}
		ids = append(ids, t.ID)
		list = append(list, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return list, nil
	}
	fieldsByType, err := r.listFieldsByTypeIDsFull(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range list {
		f := fieldsByType[list[i].ID]
		if f == nil {
			f = []domain.TypeField{}
		}
		list[i].Fields = f
	}
	return list, nil
}

func (r *CategoryPG) FindSubcategoryTypeByID(ctx context.Context, id int32) (*domain.SubcategoryTypePathDetail, error) {
	const q = `SELECT "id", "name", "slug", "subcategoryId", "createdAt", "updatedAt" FROM "SubcategotyType" WHERE "id" = $1`
	row := r.pool.QueryRow(ctx, q, id)
	var st domain.SubcategoryTypePathDetail
	if err := row.Scan(&st.ID, &st.Name, &st.Slug, &st.SubcategoryID, &st.CreatedAt, &st.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	subQ := `SELECT "id", "name", "slug", "categoryId", "createdAt", "updatedAt" FROM "SubCategory" WHERE "id" = $1`
	srow := r.pool.QueryRow(ctx, subQ, st.SubcategoryID)
	var sn domain.SubcategoryNested
	if err := srow.Scan(&sn.ID, &sn.Name, &sn.Slug, &sn.CategoryID, &sn.CreatedAt, &sn.UpdatedAt); err != nil {
		return nil, err
	}
	cat, err := r.FindCategoryByID(ctx, sn.CategoryID)
	if err != nil {
		return nil, err
	}
	sn.Category = *cat
	st.Subcategory = sn
	fields, err := r.listFieldsByTypeIDsFull(ctx, []int32{st.ID})
	if err != nil {
		return nil, err
	}
	st.Fields = fields[st.ID]
	if st.Fields == nil {
		st.Fields = []domain.TypeField{}
	}
	return &st, nil
}

// --- TypeField ---

func (r *CategoryPG) listFieldsByTypeIDsMini(ctx context.Context, typeIDs []int32) (map[int32][]domain.TypeFieldMini, error) {
	out := make(map[int32][]domain.TypeFieldMini)
	if len(typeIDs) == 0 {
		return out, nil
	}
	q := fmt.Sprintf(
		`SELECT "id", "name", "typeId" FROM "TypeField" WHERE "typeId" IN (%s) ORDER BY "id" ASC`,
		joinInt32(typeIDs),
	)
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f domain.TypeFieldMini
		var typeID int32
		if err := rows.Scan(&f.ID, &f.Name, &typeID); err != nil {
			return nil, err
		}
		out[typeID] = append(out[typeID], f)
	}
	return out, rows.Err()
}

func (r *CategoryPG) listFieldsByTypeIDsFull(ctx context.Context, typeIDs []int32) (map[int32][]domain.TypeField, error) {
	out := make(map[int32][]domain.TypeField)
	if len(typeIDs) == 0 {
		return out, nil
	}
	q := fmt.Sprintf(
		`SELECT "id", "name", "isRequired", "typeId" FROM "TypeField" WHERE "typeId" IN (%s) ORDER BY "id" ASC`,
		joinInt32(typeIDs),
	)
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f domain.TypeField
		if err := rows.Scan(&f.ID, &f.Name, &f.IsRequired, &f.TypeID); err != nil {
			return nil, err
		}
		out[f.TypeID] = append(out[f.TypeID], f)
	}
	return out, rows.Err()
}

// BuildCategoryTree вЂ” СЃРѕР±РёСЂР°РµС‚ РґРµСЂРµРІРѕ РґР»СЏ findAll (РєР°Рє Prisma select).
func (r *CategoryPG) BuildCategoryTree(ctx context.Context, roots []domain.CategoryTree) ([]domain.CategoryTree, error) {
	if len(roots) == 0 {
		return roots, nil
	}
	ids := make([]int32, len(roots))
	for i := range roots {
		ids[i] = roots[i].ID
	}
	subs, err := r.ListSubCategoriesByCategoryIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	subIDs := make([]int32, len(subs))
	subByCat := make(map[int32][]domain.SubCategoryTree)
	for i, s := range subs {
		subIDs[i] = s.ID
		node := domain.SubCategoryTree{
			ID:               s.ID,
			Name:             s.Name,
			Slug:             s.Slug,
			SubcategoryTypes: []domain.SubcategoryTypeTree{},
		}
		subByCat[s.CategoryID] = append(subByCat[s.CategoryID], node)
	}
	types, err := r.listSubcategoryTypesBySubcategoryIDs(ctx, subIDs)
	if err != nil {
		return nil, err
	}
	typeIDs := make([]int32, len(types))
	typeBySub := make(map[int32][]domain.SubcategoryTypeTree)
	for i, t := range types {
		typeIDs[i] = t.ID
		node := domain.SubcategoryTypeTree{
			ID:     t.ID,
			Name:   t.Name,
			Slug:   t.Slug,
			Fields: []domain.TypeFieldMini{},
		}
		typeBySub[t.SubcategoryID] = append(typeBySub[t.SubcategoryID], node)
	}
	fieldsByType, err := r.listFieldsByTypeIDsMini(ctx, typeIDs)
	if err != nil {
		return nil, err
	}
	// РїСЂРёРєСЂРµРїРёС‚СЊ fields Рє types РІРЅСѓС‚СЂРё map sub -> types
	for catIdx := range roots {
		cid := roots[catIdx].ID
		subNodes := subByCat[cid]
		for si := range subNodes {
			sid := subNodes[si].ID
			typeNodes := typeBySub[sid]
			for ti := range typeNodes {
				tid := typeNodes[ti].ID
				f := fieldsByType[tid]
				if f == nil {
					f = []domain.TypeFieldMini{}
				}
				typeNodes[ti].Fields = f
			}
			subNodes[si].SubcategoryTypes = typeNodes
		}
		roots[catIdx].SubCategories = subNodes
	}
	return roots, nil
}

// BuildCategoryDetail вЂ” РїРѕР»РЅРѕРµ РґРµСЂРµРІРѕ РґР»СЏ РѕРґРЅРѕР№ РєР°С‚РµРіРѕСЂРёРё (findById / findBySlug).
func (r *CategoryPG) BuildCategoryDetail(ctx context.Context, base domain.CategoryScalars) (*domain.CategoryDetail, error) {
	out := &domain.CategoryDetail{
		ID:            base.ID,
		Name:          base.Name,
		Slug:          base.Slug,
		CreatedAt:     base.CreatedAt,
		UpdatedAt:     base.UpdatedAt,
		SubCategories: []domain.SubCategoryDetail{},
	}
	subs, err := r.ListSubCategoriesByCategoryIDs(ctx, []int32{base.ID})
	if err != nil {
		return nil, err
	}
	for _, s := range subs {
		types, err := r.listSubcategoryTypesBySubcategoryID(ctx, s.ID)
		if err != nil {
			return nil, err
		}
		out.SubCategories = append(out.SubCategories, domain.SubCategoryDetail{
			ID:               s.ID,
			Name:             s.Name,
			Slug:             s.Slug,
			CategoryID:       s.CategoryID,
			CreatedAt:        s.CreatedAt,
			UpdatedAt:        s.UpdatedAt,
			SubcategoryTypes: types,
		})
	}
	return out, nil
}

// --- Админ: SubCategory (пути как Nest /subcategory/*) ---

func (r *CategoryPG) AdminListSubcategories(ctx context.Context) ([]domain.SubcategoryAdminListItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT "id", "name", "categoryId" FROM "SubCategory"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.SubcategoryAdminListItem
	for rows.Next() {
		var it domain.SubcategoryAdminListItem
		if err := rows.Scan(&it.ID, &it.Name, &it.CategoryID); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *CategoryPG) AdminFindSubcategoryFlatByID(ctx context.Context, id int32) (*domain.SubcategoryAdminDetail, error) {
	const q = `SELECT "id", "name", "slug", "categoryId", "createdAt", "updatedAt" FROM "SubCategory" WHERE "id" = $1`
	row := r.pool.QueryRow(ctx, q, id)
	var d domain.SubcategoryAdminDetail
	if err := row.Scan(&d.ID, &d.Name, &d.Slug, &d.CategoryID, &d.CreatedAt, &d.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *CategoryPG) AdminInsertSubcategory(ctx context.Context, name string, categoryID int32) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO "SubCategory" ("name", "slug", "categoryId") VALUES ($1, '', $2)`,
		name, categoryID,
	)
	return err
}

func (r *CategoryPG) AdminUpdateSubcategory(ctx context.Context, id int32, name *string, categoryID *int32) error {
	sets := []string{`"updatedAt" = NOW()`}
	args := []any{id}
	n := 2
	if name != nil {
		sets = append(sets, fmt.Sprintf(`"name" = $%d`, n))
		args = append(args, *name)
		n++
	}
	if categoryID != nil {
		sets = append(sets, fmt.Sprintf(`"categoryId" = $%d`, n))
		args = append(args, *categoryID)
		n++
	}
	if len(sets) == 1 {
		return nil
	}
	q := fmt.Sprintf(`UPDATE "SubCategory" SET %s WHERE "id" = $1`, strings.Join(sets, ", "))
	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CategoryPG) AdminDeleteSubcategory(ctx context.Context, id int32) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM "SubCategory" WHERE "id" = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Админ: SubcategotyType ---

func (r *CategoryPG) AdminListSubcategoryTypes(ctx context.Context) ([]domain.SubcategoryTypeAdminListItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT "id", "name", "subcategoryId" FROM "SubcategotyType"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.SubcategoryTypeAdminListItem
	for rows.Next() {
		var it domain.SubcategoryTypeAdminListItem
		if err := rows.Scan(&it.ID, &it.Name, &it.SubcategoryID); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *CategoryPG) AdminFindSubcategoryTypeFlatByID(ctx context.Context, id int32) (*domain.SubcategoryTypeAdminDetail, error) {
	const q = `SELECT "id", "name", "slug", "subcategoryId", "createdAt", "updatedAt" FROM "SubcategotyType" WHERE "id" = $1`
	row := r.pool.QueryRow(ctx, q, id)
	var d domain.SubcategoryTypeAdminDetail
	if err := row.Scan(&d.ID, &d.Name, &d.Slug, &d.SubcategoryID, &d.CreatedAt, &d.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *CategoryPG) AdminInsertSubcategoryType(ctx context.Context, name string, subcategoryID int32) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO "SubcategotyType" ("name", "slug", "subcategoryId") VALUES ($1, '', $2)`,
		name, subcategoryID,
	)
	return err
}

func (r *CategoryPG) AdminUpdateSubcategoryType(ctx context.Context, id int32, name *string, subcategoryID *int32) error {
	sets := []string{`"updatedAt" = NOW()`}
	args := []any{id}
	n := 2
	if name != nil {
		sets = append(sets, fmt.Sprintf(`"name" = $%d`, n))
		args = append(args, *name)
		n++
	}
	if subcategoryID != nil {
		sets = append(sets, fmt.Sprintf(`"subcategoryId" = $%d`, n))
		args = append(args, *subcategoryID)
		n++
	}
	if len(sets) == 1 {
		return nil
	}
	q := fmt.Sprintf(`UPDATE "SubcategotyType" SET %s WHERE "id" = $1`, strings.Join(sets, ", "))
	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CategoryPG) AdminDeleteSubcategoryType(ctx context.Context, id int32) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM "SubcategotyType" WHERE "id" = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Админ: TypeField (в Prisma нет updatedAt у TypeField) ---

func (r *CategoryPG) AdminListTypeFields(ctx context.Context) ([]domain.TypeFieldAdminItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT "id", "name", "typeId" FROM "TypeField"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.TypeFieldAdminItem
	for rows.Next() {
		var it domain.TypeFieldAdminItem
		if err := rows.Scan(&it.ID, &it.Name, &it.TypeID); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *CategoryPG) AdminFindTypeFieldItemByID(ctx context.Context, id int32) (*domain.TypeFieldAdminItem, error) {
	const q = `SELECT "id", "name", "typeId" FROM "TypeField" WHERE "id" = $1`
	row := r.pool.QueryRow(ctx, q, id)
	var it domain.TypeFieldAdminItem
	if err := row.Scan(&it.ID, &it.Name, &it.TypeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &it, nil
}

func (r *CategoryPG) AdminInsertTypeField(ctx context.Context, name string, typeID int32, isRequired bool) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO "TypeField" ("name", "isRequired", "typeId") VALUES ($1, $2, $3)`,
		name, isRequired, typeID,
	)
	return err
}

func (r *CategoryPG) AdminUpdateTypeField(ctx context.Context, id int32, name *string, typeID *int32) error {
	sets := []string{}
	args := []any{id}
	n := 2
	if name != nil {
		sets = append(sets, fmt.Sprintf(`"name" = $%d`, n))
		args = append(args, *name)
		n++
	}
	if typeID != nil {
		sets = append(sets, fmt.Sprintf(`"typeId" = $%d`, n))
		args = append(args, *typeID)
		n++
	}
	if len(sets) == 0 {
		return nil
	}
	q := fmt.Sprintf(`UPDATE "TypeField" SET %s WHERE "id" = $1`, strings.Join(sets, ", "))
	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CategoryPG) AdminDeleteTypeField(ctx context.Context, id int32) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM "TypeField" WHERE "id" = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
