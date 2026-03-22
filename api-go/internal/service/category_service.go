package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/pkg/slug"
	"med-vito/api-go/internal/repository"
)

// AppError — маппится в HTTP в хендлере (аналог Nest HttpException).
type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string { return e.Message }

type CategoryService struct {
	repo *repository.CategoryPG
}

func NewCategoryService(repo *repository.CategoryPG) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, name string, slugOpt *string) (*domain.CreateCategoryResponse, error) {
	if slug.IsBlankOrSpace(name) {
		return nil, &AppError{400, "Название категории обязательно для заполнения"}
	}
	name = strings.TrimSpace(name)

	var finalSlug string
	if slugOpt != nil && *slugOpt != "" {
		if !slug.IsValidSlug(*slugOpt) {
			return nil, &AppError{400, "Slug должен содержать только латинские буквы в нижнем регистре, цифры и дефисы"}
		}
		finalSlug = *slugOpt
	} else {
		finalSlug = slug.GenerateSlug(name)
	}

	exists, err := s.repo.ExistsCategoryBySlug(ctx, finalSlug)
	if err != nil {
		return nil, err
	}
	if exists {
		if slugOpt != nil && *slugOpt != "" {
			return nil, &AppError{400, fmt.Sprintf("Категория с slug '%s' уже существует", finalSlug)}
		}
		all, err := s.repo.ListAllCategorySlugs(ctx)
		if err != nil {
			return nil, err
		}
		finalSlug = slug.MakeUniqueSlug(finalSlug, all)
	}

	cat, err := s.repo.InsertCategory(ctx, name, finalSlug)
	if err != nil {
		return nil, err
	}
	return &domain.CreateCategoryResponse{
		Message:  "Категория успешно создана",
		Category: *cat,
	}, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id int32, name string, slugOpt *string) (*domain.UpdateCategoryResponse, error) {
	if slug.IsBlankOrSpace(name) {
		return nil, &AppError{400, "Название категории обязательно для заполнения"}
	}
	name = strings.TrimSpace(name)

	cur, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Категория с таким id не найдена"}
		}
		return nil, err
	}

	var newSlug *string
	if slugOpt != nil && *slugOpt != "" {
		if !slug.IsValidSlug(*slugOpt) {
			return nil, &AppError{400, "Slug должен содержать только латинские буквы в нижнем регистре, цифры и дефисы"}
		}
		exists, err := s.repo.ExistsCategoryBySlugExceptID(ctx, *slugOpt, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, &AppError{400, fmt.Sprintf("Категория с slug '%s' уже существует", *slugOpt)}
		}
		s := *slugOpt
		newSlug = &s
	} else if name != cur.Name {
		gen := slug.GenerateSlug(name)
		exists, err := s.repo.ExistsCategoryBySlugExceptID(ctx, gen, id)
		if err != nil {
			return nil, err
		}
		if exists {
			others, err := s.repo.ListCategorySlugsExceptID(ctx, id)
			if err != nil {
				return nil, err
			}
			u := slug.MakeUniqueSlug(gen, others)
			newSlug = &u
		} else {
			newSlug = &gen
		}
	}

	updated, err := s.repo.UpdateCategory(ctx, id, name, newSlug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Категория с таким id не найдена"}
		}
		return nil, err
	}
	return &domain.UpdateCategoryResponse{
		Message:  "Категория успешно обновлена",
		Category: *updated,
	}, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int32) (*domain.DeleteCategoryResponse, error) {
	err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Категория для удаления не найдена"}
		}
		return nil, err
	}
	return &domain.DeleteCategoryResponse{Message: "Категория успешно удалена"}, nil
}

func (s *CategoryService) FindAllCategories(ctx context.Context) ([]domain.CategoryTree, error) {
	roots, err := s.repo.ListCategoriesOrderByName(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.BuildCategoryTree(ctx, roots)
}

func (s *CategoryService) FindByID(ctx context.Context, id int32) (*domain.CategoryDetail, error) {
	base, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Категория не найдена"}
		}
		return nil, err
	}
	return s.repo.BuildCategoryDetail(ctx, *base)
}

func (s *CategoryService) FindBySlug(ctx context.Context, slugStr string) (*domain.CategoryDetail, error) {
	base, err := s.repo.FindCategoryBySlug(ctx, slugStr)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, fmt.Sprintf("Категория '%s' не найдена", slugStr)}
		}
		return nil, err
	}
	return s.repo.BuildCategoryDetail(ctx, *base)
}

func (s *CategoryService) FindBySlugPath(ctx context.Context, slugPath string) (*domain.SlugPathResponse, error) {
	parts := strings.Split(slugPath, "/")
	var slugs []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			slugs = append(slugs, p)
		}
	}
	if len(slugs) == 0 {
		return nil, &AppError{400, "Некорректный путь категории"}
	}

	catSlug := slugs[0]
	base, err := s.repo.FindCategoryBySlug(ctx, catSlug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, fmt.Sprintf("Категория '%s' не найдена", catSlug)}
		}
		return nil, err
	}
	detail, err := s.repo.BuildCategoryDetail(ctx, *base)
	if err != nil {
		return nil, err
	}

	if len(slugs) == 1 {
		return &domain.SlugPathResponse{Type: "category", Data: detail}, nil
	}

	subSlug := slugs[1]
	var subID int32
	var found bool
	for _, sc := range detail.SubCategories {
		if sc.Slug == subSlug {
			subID = sc.ID
			found = true
			break
		}
	}
	if !found {
		return nil, &AppError{404, fmt.Sprintf("Подкатегория '%s' не найдена в категории '%s'", subSlug, catSlug)}
	}

	if len(slugs) == 2 {
		withRel, err := s.repo.FindSubCategoryByID(ctx, subID)
		if err != nil {
			return nil, err
		}
		return &domain.SlugPathResponse{Type: "subcategory", Data: withRel}, nil
	}

	typeSlug := slugs[2]
	subFull, err := s.repo.FindSubCategoryByID(ctx, subID)
	if err != nil {
		return nil, err
	}
	var typeID int32
	found = false
	for _, t := range subFull.SubcategoryTypes {
		if t.Slug == typeSlug {
			typeID = t.ID
			found = true
			break
		}
	}
	if !found {
		return nil, &AppError{404, fmt.Sprintf("Тип '%s' не найден в подкатегории '%s'", typeSlug, subSlug)}
	}

	typeDetail, err := s.repo.FindSubcategoryTypeByID(ctx, typeID)
	if err != nil {
		return nil, err
	}
	return &domain.SlugPathResponse{Type: "subcategory-type", Data: typeDetail}, nil
}

// --- Админ: subcategory (как Nest SubcategoryController) ---

func (s *CategoryService) AdminListSubcategories(ctx context.Context) ([]domain.SubcategoryAdminListItem, error) {
	return s.repo.AdminListSubcategories(ctx)
}

func (s *CategoryService) AdminFindSubcategoryByID(ctx context.Context, id int32) (*domain.SubcategoryAdminDetail, error) {
	d, err := s.repo.AdminFindSubcategoryFlatByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Подкатегория не найдена"}
		}
		return nil, err
	}
	return d, nil
}

func (s *CategoryService) AdminCreateSubcategory(ctx context.Context, name string, categoryID int32) (*domain.CatalogMessage, error) {
	name = strings.TrimSpace(name)
	if slug.IsBlankOrSpace(name) {
		return nil, &AppError{400, "Название подкатегории обязательно для заполнения"}
	}
	if _, err := s.repo.FindCategoryByID(ctx, categoryID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{400, "Выбранная категория не существует"}
		}
		return nil, err
	}
	if err := s.repo.AdminInsertSubcategory(ctx, name, categoryID); err != nil {
		return nil, err
	}
	return &domain.CatalogMessage{Message: "Подкатегория успешно создана"}, nil
}

func (s *CategoryService) AdminUpdateSubcategory(ctx context.Context, id int32, name *string, categoryID *int32) (*domain.CatalogMessage, error) {
	if name == nil && categoryID == nil {
		return nil, &AppError{400, "Нет данных для обновления"}
	}
	if name != nil {
		t := strings.TrimSpace(*name)
		if slug.IsBlankOrSpace(t) {
			return nil, &AppError{400, "Название подкатегории не может быть пустым"}
		}
		name = &t
	}
	if categoryID != nil {
		if _, err := s.repo.FindCategoryByID(ctx, *categoryID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, &AppError{400, "Выбранная категория не существует"}
			}
			return nil, err
		}
	}
	if err := s.repo.AdminUpdateSubcategory(ctx, id, name, categoryID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Подкатегория с таким id не найдена"}
		}
		return nil, err
	}
	return &domain.CatalogMessage{Message: "Подкатегория успешно обновлена"}, nil
}

func (s *CategoryService) AdminDeleteSubcategory(ctx context.Context, id int32) (*domain.CatalogMessage, error) {
	if err := s.repo.AdminDeleteSubcategory(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{400, "Подкатегория для удаления не найдена"}
		}
		return nil, err
	}
	return &domain.CatalogMessage{Message: "Подкатегория успешно удалена"}, nil
}

// --- Админ: subcategory-type ---

func (s *CategoryService) AdminListSubcategoryTypes(ctx context.Context) ([]domain.SubcategoryTypeAdminListItem, error) {
	return s.repo.AdminListSubcategoryTypes(ctx)
}

func (s *CategoryService) AdminFindSubcategoryTypeByID(ctx context.Context, id int32) (*domain.SubcategoryTypeAdminDetail, error) {
	d, err := s.repo.AdminFindSubcategoryTypeFlatByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Тип подкатегории не найден"}
		}
		return nil, err
	}
	return d, nil
}

func (s *CategoryService) AdminCreateSubcategoryType(ctx context.Context, name string, subcategoryID int32) (*domain.CatalogMessage, error) {
	name = strings.TrimSpace(name)
	if slug.IsBlankOrSpace(name) {
		return nil, &AppError{400, "Название типа подкатегории обязательно для заполнения"}
	}
	if _, err := s.repo.AdminFindSubcategoryFlatByID(ctx, subcategoryID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Подкатегория не найдена"}
		}
		return nil, err
	}
	if err := s.repo.AdminInsertSubcategoryType(ctx, name, subcategoryID); err != nil {
		return nil, err
	}
	return &domain.CatalogMessage{Message: "Тип подкатегории успешно создан"}, nil
}

func (s *CategoryService) AdminUpdateSubcategoryType(ctx context.Context, id int32, name *string, subcategoryID *int32) (*domain.CatalogMessage, error) {
	if name == nil && subcategoryID == nil {
		return nil, &AppError{400, "Нет данных для обновления"}
	}
	if name != nil {
		t := strings.TrimSpace(*name)
		if slug.IsBlankOrSpace(t) {
			return nil, &AppError{400, "Название типа подкатегории не может быть пустым"}
		}
		name = &t
	}
	if subcategoryID != nil {
		if _, err := s.repo.AdminFindSubcategoryFlatByID(ctx, *subcategoryID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, &AppError{404, "Подкатегория не найдена"}
			}
			return nil, err
		}
	}
	if err := s.repo.AdminUpdateSubcategoryType(ctx, id, name, subcategoryID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Тип подкатегории не найден"}
		}
		return nil, err
	}
	return &domain.CatalogMessage{Message: "Тип подкатегории успешно обновлён"}, nil
}

func (s *CategoryService) AdminDeleteSubcategoryType(ctx context.Context, id int32) (*domain.CatalogMessage, error) {
	if err := s.repo.AdminDeleteSubcategoryType(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Тип подкатегории для удаления не надйен"}
		}
		return nil, err
	}
	return &domain.CatalogMessage{Message: "Тип подкатегории успешно удалён"}, nil
}

// --- Админ: type-field ---

func (s *CategoryService) AdminListTypeFields(ctx context.Context) ([]domain.TypeFieldAdminItem, error) {
	return s.repo.AdminListTypeFields(ctx)
}

func (s *CategoryService) AdminFindTypeFieldByID(ctx context.Context, id int32) (*domain.TypeFieldAdminItem, error) {
	d, err := s.repo.AdminFindTypeFieldItemByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Характеристика не найдена"}
		}
		return nil, err
	}
	return d, nil
}

func (s *CategoryService) AdminCreateTypeField(ctx context.Context, name string, typeID int32) (*domain.CatalogMessage, error) {
	name = strings.TrimSpace(name)
	if slug.IsBlankOrSpace(name) {
		return nil, &AppError{400, "Название характеристики обязательно для заполнения"}
	}
	if _, err := s.repo.AdminFindSubcategoryTypeFlatByID(ctx, typeID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Тип подкатегории не найден"}
		}
		return nil, err
	}
	if err := s.repo.AdminInsertTypeField(ctx, name, typeID, false); err != nil {
		return nil, err
	}
	return &domain.CatalogMessage{Message: "Характеристика успешно добавлена"}, nil
}

func (s *CategoryService) AdminUpdateTypeField(ctx context.Context, id int32, name *string, typeID *int32) (*domain.CatalogMessage, error) {
	if name == nil && typeID == nil {
		return nil, &AppError{400, "Нет данных для обновления"}
	}
	if name != nil {
		t := strings.TrimSpace(*name)
		if slug.IsBlankOrSpace(t) {
			return nil, &AppError{400, "Название характеристики не может быть пустым"}
		}
		name = &t
	}
	if typeID != nil {
		if _, err := s.repo.AdminFindSubcategoryTypeFlatByID(ctx, *typeID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, &AppError{404, "Тип подкатегории не найден"}
			}
			return nil, err
		}
	}
	if err := s.repo.AdminUpdateTypeField(ctx, id, name, typeID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Характеристика не найдена"}
		}
		return nil, err
	}
	return &domain.CatalogMessage{Message: "Характеристика успешно обновлена"}, nil
}

func (s *CategoryService) AdminDeleteTypeField(ctx context.Context, id int32) (*domain.CatalogMessage, error) {
	if err := s.repo.AdminDeleteTypeField(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Характеристика для удаления не найдена"}
		}
		return nil, err
	}
	return &domain.CatalogMessage{Message: "Характеристика успешно удалена"}, nil
}
