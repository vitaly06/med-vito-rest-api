package domain

import "time"

// --- Дерево как в categoryService.findAllCategories (select без дат у корня) ---

type CategoryTree struct {
	ID            int32             `json:"id"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	SubCategories []SubCategoryTree `json:"subCategories"`
}

type SubCategoryTree struct {
	ID               int32                 `json:"id"`
	Name             string                `json:"name"`
	Slug             string                `json:"slug"`
	SubcategoryTypes []SubcategoryTypeTree `json:"subcategoryTypes"`
}

type SubcategoryTypeTree struct {
	ID     int32           `json:"id"`
	Name   string          `json:"name"`
	Slug   string          `json:"slug"`
	Fields []TypeFieldMini `json:"fields"`
}

type TypeFieldMini struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// --- Полные сущности (findById / findBySlug) — как include в Prisma ---

type TypeField struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	IsRequired bool   `json:"isRequired"`
	TypeID     int32  `json:"typeId"`
}

type SubcategoryTypeDetail struct {
	ID            int32       `json:"id"`
	Name          string      `json:"name"`
	Slug          string      `json:"slug"`
	SubcategoryID int32       `json:"subcategoryId"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	Fields        []TypeField `json:"fields"`
}

type SubCategoryDetail struct {
	ID               int32                   `json:"id"`
	Name             string                  `json:"name"`
	Slug             string                  `json:"slug"`
	CategoryID       int32                   `json:"categoryId"`
	CreatedAt        time.Time               `json:"createdAt"`
	UpdatedAt        time.Time               `json:"updatedAt"`
	SubcategoryTypes []SubcategoryTypeDetail `json:"subcategoryTypes"`
}

type CategoryDetail struct {
	ID            int32               `json:"id"`
	Name          string              `json:"name"`
	Slug          string              `json:"slug"`
	CreatedAt     time.Time           `json:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt"`
	SubCategories []SubCategoryDetail `json:"subCategories"`
}

// CategoryScalars — вложенная Category у SubCategory (Prisma category: true).
type CategoryScalars struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SubCategoryWithCategory — ответ findBySlugPath type=subcategory.
type SubCategoryWithCategory struct {
	ID               int32                   `json:"id"`
	Name             string                  `json:"name"`
	Slug             string                  `json:"slug"`
	CategoryID       int32                   `json:"categoryId"`
	CreatedAt        time.Time               `json:"createdAt"`
	UpdatedAt        time.Time               `json:"updatedAt"`
	Category         CategoryScalars         `json:"category"`
	SubcategoryTypes []SubcategoryTypeDetail `json:"subcategoryTypes"`
}

// SubcategoryNested — subcategory внутри type=subcategory-type.
type SubcategoryNested struct {
	ID         int32           `json:"id"`
	Name       string          `json:"name"`
	Slug       string          `json:"slug"`
	CategoryID int32           `json:"categoryId"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
	Category   CategoryScalars `json:"category"`
}

// SubcategoryTypePathDetail — ответ findBySlugPath type=subcategory-type.
type SubcategoryTypePathDetail struct {
	ID            int32             `json:"id"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	SubcategoryID int32             `json:"subcategoryId"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	Subcategory   SubcategoryNested `json:"subcategory"`
	Fields        []TypeField       `json:"fields"`
}

// SlugPathResponse — обёртка type + data.
type SlugPathResponse struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// CreateCategoryResponse / Update — как Nest.
type CreateCategoryResponse struct {
	Message  string          `json:"message"`
	Category CategoryScalars `json:"category"`
}

type UpdateCategoryResponse struct {
	Message  string          `json:"message"`
	Category CategoryScalars `json:"category"`
}

type DeleteCategoryResponse struct {
	Message string `json:"message"`
}

// --- Админские сущности как у Nest (subcategory / subcategory-type / type-field) ---

type SubcategoryAdminListItem struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	CategoryID int32  `json:"categoryId"`
}

type SubcategoryAdminDetail struct {
	ID         int32     `json:"id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	CategoryID int32     `json:"categoryId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type SubcategoryTypeAdminListItem struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	SubcategoryID int32  `json:"subcategoryId"`
}

type SubcategoryTypeAdminDetail struct {
	ID            int32     `json:"id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	SubcategoryID int32     `json:"subcategoryId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// TypeFieldAdminItem — в Nest в списке и find-by-id только id, name, typeId (без isRequired).
type TypeFieldAdminItem struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	TypeID int32  `json:"typeId"`
}

type CatalogMessage struct {
	Message string `json:"message"`
}
