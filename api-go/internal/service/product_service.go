package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"med-vito/api-go/internal/pkg/s3client"
	"med-vito/api-go/internal/repository"
)

// ProductService — логика Nest ProductService (без WebSocket).
type ProductService struct {
	prod *repository.ProductPG
	s3   *s3client.Client
	user *UserService
}

func NewProductService(prod *repository.ProductPG, s3 *s3client.Client, user *UserService) *ProductService {
	return &ProductService{prod: prod, s3: s3, user: user}
}

func formatProductDate(t time.Time) string {
	lt := t.In(time.Local)
	d := lt.Day()
	m := int(lt.Month())
	y := lt.Year() % 100
	h, mi := lt.Hour(), lt.Minute()
	return fmt.Sprintf("%02d.%02d.%02d в %02d:%02d", d, m, y, h, mi)
}

func (s *ProductService) formatListItem(pr repository.ProductListRow, isFav bool, hasPromo bool, promoLevel int32) map[string]any {
	out := map[string]any{
		"id":              pr.ID,
		"images":          pr.Images,
		"name":            pr.Name,
		"address":         pr.Address,
		"createdAt":       formatProductDate(pr.CreatedAt),
		"price":           pr.Price,
		"userId":          pr.UserID,
		"videoUrl":        pr.VideoURL,
		"isFavorited":     isFav,
		"hasPromotion":    hasPromo,
		"promotionLevel":  promoLevel,
		"categoryId":      pr.CategoryID,
		"categoryName":    pr.CategoryName,
		"categorySlug":    pr.CategorySlug,
		"subCategoryId":   pr.SubCategoryID,
		"subCategoryName": pr.SubCategoryName,
		"subCategorySlug": pr.SubCategorySlug,
		"typeId":          pr.TypeID,
		"typeName":        pr.TypeName,
		"typeSlug":        pr.TypeSlug,
	}
	if pr.ModerateState != nil {
		out["moderateState"] = *pr.ModerateState
	}
	out["isHide"] = pr.IsHide
	return out
}

func interleaveLuxRegular(rows []repository.ProductListRow) []repository.ProductListRow {
	var lux, reg []repository.ProductListRow
	for _, p := range rows {
		if p.PromotionLevel >= 100 {
			lux = append(lux, p)
		} else {
			reg = append(reg, p)
		}
	}
	var out []repository.ProductListRow
	li, ri := 0, 0
	for li < len(lux) || ri < len(reg) {
		if li < len(lux) {
			out = append(out, lux[li])
			li++
		}
		for i := 0; i < 5 && ri < len(reg); i++ {
			out = append(out, reg[ri])
			ri++
		}
	}
	return out
}

func (s *ProductService) mapListWithFavorites(ctx context.Context, rows []repository.ProductListRow, viewer *int32) ([]map[string]any, error) {
	ordered := interleaveLuxRegular(rows)
	out := make([]map[string]any, 0, len(ordered))
	for _, pr := range ordered {
		hasPromo := pr.PromotionLevel > 0
		var fav bool
		var err error
		if viewer != nil {
			fav, err = s.prod.IsFavorite(ctx, *viewer, pr.ID)
			if err != nil {
				return nil, err
			}
		}
		out = append(out, s.formatListItem(pr, fav, hasPromo, pr.PromotionLevel))
	}
	return out, nil
}

// CreateProduct — multipart + S3; fieldValues JSON map.
func (s *ProductService) CreateProduct(ctx context.Context, userID int32, name, priceStr, state, description, address, categoryStr, subStr, typeStr, fieldJSON, videoStr string, files []UploadedFile) (map[string]any, error) {
	if s.s3 == nil {
		return nil, &AppError{500, "S3 не настроен"}
	}
	if strings.TrimSpace(address) == "" {
		return nil, &AppError{400, "Адрес обязателен для заполнения"}
	}
	if strings.TrimSpace(name) == "" {
		return nil, &AppError{400, "Название обязательно"}
	}
	if strings.TrimSpace(state) == "" {
		return nil, &AppError{400, "Состояние обязательно (NEW или USED)"}
	}
	price, err := strconv.Atoi(strings.TrimSpace(priceStr))
	if err != nil || price < 1 {
		return nil, &AppError{400, "Цена должна быть числом больше 0"}
	}
	catID, err := strconv.ParseInt(strings.TrimSpace(categoryStr), 10, 32)
	if err != nil {
		return nil, &AppError{400, "Некорректный categoryId"}
	}
	subID, err := strconv.ParseInt(strings.TrimSpace(subStr), 10, 32)
	if err != nil {
		return nil, &AppError{400, "Некорректный subcategoryId"}
	}
	ok, err := s.prod.SubCategoryBelongsToCategory(ctx, int32(catID), int32(subID))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &AppError{400, "Подкатегория не найдена или не принадлежит указанной категории"}
	}
	var typePtr *int32
	if strings.TrimSpace(typeStr) != "" {
		tid64, err := strconv.ParseInt(strings.TrimSpace(typeStr), 10, 32)
		if err != nil {
			return nil, &AppError{400, "Некорректный typeId"}
		}
		t := int32(tid64)
		typePtr = &t
	}
	fvMap, err := parseFieldValuesMap(fieldJSON)
	if err != nil {
		return nil, &AppError{400, err.Error()}
	}
	var fieldIDs []int32
	for k := range fvMap {
		id64, err := strconv.ParseInt(k, 10, 32)
		if err != nil {
			continue
		}
		fieldIDs = append(fieldIDs, int32(id64))
	}
	if err := s.prod.ValidateTypeFieldIDs(ctx, fieldIDs); err != nil {
		return nil, &AppError{400, err.Error()}
	}

	var urls []string
	for _, f := range files {
		u, err := s.s3.Upload(ctx, "products", f.Name, f.ContentType, f.Body)
		if err != nil {
			return nil, &AppError{400, "Ошибка загрузки в S3: " + err.Error()}
		}
		urls = append(urls, u)
	}

	pid, err := s.prod.GenerateUniqueProductID(ctx)
	if err != nil {
		return nil, err
	}
	intMap := make(map[int32]string)
	for k, v := range fvMap {
		if id64, err := strconv.ParseInt(k, 10, 32); err == nil {
			intMap[int32(id64)] = v
		}
	}
	var vptr *string
	if strings.TrimSpace(videoStr) != "" {
		v := strings.TrimSpace(videoStr)
		vptr = &v
	}
	err = s.prod.CreateProductTx(ctx, pid, userID, name, int32(price), state, description, address, vptr, urls, int32(catID), int32(subID), typePtr, intMap)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{400, "Пользователь не найден"}
	}
	if errors.Is(err, repository.ErrInsufficientFunds) {
		return nil, &AppError{400, fmt.Sprintf("Недостаточно средств для размещения объявления. Требуется %d руб.", repository.AdListingCost)}
	}
	if err != nil {
		return nil, &AppError{400, "Ошибка при создании продукта: " + err.Error()}
	}
	prod, err := s.prod.LoadProductWithRelations(ctx, pid)
	if err != nil {
		return nil, err
	}
	return map[string]any{"message": "Продукт успешно создан", "product": prod}, nil
}

// UploadedFile — один файл из multipart (images).
type UploadedFile struct {
	Name        string
	ContentType string
	Body        []byte
}

func parseFieldValuesMap(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]string{}, nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, fmt.Errorf("некорректный fieldValues JSON")
	}
	return m, nil
}
