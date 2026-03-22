package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"med-vito/api-go/internal/repository"
)

func (s *ProductService) DeleteProduct(ctx context.Context, productID, userID int32) (map[string]any, error) {
	uid, imgs, err := s.prod.GetProductOwnerAndImages(ctx, productID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{400, "Товар для удаления не найден"}
	}
	if err != nil {
		return nil, err
	}
	if uid != userID {
		return nil, &AppError{403, "Вы не можете удалить чужой товар"}
	}
	if s.s3 != nil {
		for _, u := range imgs {
			_ = s.s3.DeleteByURL(ctx, u)
		}
	}
	if err := s.prod.DeleteProductByID(ctx, productID); err != nil {
		return nil, err
	}
	return map[string]any{"message": "Товар успешно удалён"}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, productID, userID int32, name, priceStr, state, description, address, videoStr, fieldJSON string, files []UploadedFile) (map[string]any, error) {
	if s.s3 == nil && len(files) > 0 {
		return nil, &AppError{500, "S3 не настроен"}
	}
	uid, typeID, existingImages, err := s.prod.ProductWithTypeForUpdate(ctx, productID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{400, "Товар не найден"}
	}
	if err != nil {
		return nil, err
	}
	if uid != userID {
		return nil, &AppError{403, "Вы не можете редактировать чужой товар"}
	}
	fvMap, err := parseFieldValuesMap(fieldJSON)
	if err != nil {
		return nil, &AppError{400, err.Error()}
	}
	var fieldIDs []int32
	for k := range fvMap {
		if id64, err := strconv.ParseInt(k, 10, 32); err == nil {
			fieldIDs = append(fieldIDs, int32(id64))
		}
	}
	if len(fieldIDs) > 0 {
		if err := s.prod.ValidateTypeFieldIDs(ctx, fieldIDs); err != nil {
			return nil, &AppError{400, err.Error()}
		}
		if typeID != nil {
			badNames, err := s.prod.TypeFieldsBelongToType(ctx, *typeID, fieldIDs)
			if err != nil {
				return nil, err
			}
			if len(badNames) > 0 {
				return nil, &AppError{400, fmt.Sprintf("Поля %s не принадлежат типу этого товара", strings.Join(badNames, ", "))}
			}
		}
	}

	var namePtr, statePtr, descPtr, addrPtr, vidPtr *string
	var pricePtr *int32
	if strings.TrimSpace(name) != "" {
		v := strings.TrimSpace(name)
		namePtr = &v
	}
	if strings.TrimSpace(priceStr) != "" {
		p, err := strconv.Atoi(strings.TrimSpace(priceStr))
		if err != nil || p < 1 {
			return nil, &AppError{400, "Некорректная цена"}
		}
		pp := int32(p)
		pricePtr = &pp
	}
	if strings.TrimSpace(state) != "" {
		v := strings.TrimSpace(state)
		statePtr = &v
	}
	if description != "" {
		descPtr = &description
	}
	if strings.TrimSpace(address) != "" {
		v := strings.TrimSpace(address)
		addrPtr = &v
	}
	if videoStr != "" {
		v := strings.TrimSpace(videoStr)
		vidPtr = &v
	}

	var newImages []string
	for _, f := range files {
		u, err := s.s3.Upload(ctx, "products", f.Name, f.ContentType, f.Body)
		if err != nil {
			return nil, &AppError{400, "Ошибка загрузки в S3: " + err.Error()}
		}
		newImages = append(newImages, u)
	}
	var imgsArg []string
	if len(newImages) > 0 {
		imgsArg = append(append([]string{}, existingImages...), newImages...)
	}

	if err := s.prod.UpdateProductPartial(ctx, productID, namePtr, pricePtr, statePtr, descPtr, addrPtr, vidPtr, imgsArg); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{400, "Товар не найден"}
		}
		return nil, &AppError{400, "Ошибка при обновлении: " + err.Error()}
	}
	for k, v := range fvMap {
		fid64, err := strconv.ParseInt(k, 10, 32)
		if err != nil {
			continue
		}
		if err := s.prod.UpsertProductFieldValue(ctx, productID, int32(fid64), v); err != nil {
			return nil, err
		}
	}
	prod, err := s.prod.LoadProductWithRelations(ctx, productID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"message": "Товар успешно обновлён", "product": prod}, nil
}

func (s *ProductService) AvailableFilters(ctx context.Context, catSlug, subSlug, typeSlug *string) (map[string]any, error) {
	return s.prod.BuildAvailableFilters(ctx, catSlug, subSlug, typeSlug)
}

// ProductSearchQuery — query-параметры поиска.
type ProductSearchQuery struct {
	Search, CategorySlug, SubCategorySlug, TypeSlug *string
	MinPrice, MaxPrice                              *int32
	State, Region, ProfileType                      *string
	FieldValuesJSON                                 *string
	SortBy                                          string
	Page, Limit                                     int
}

func (s *ProductService) FindAll(ctx context.Context, viewer *int32, q ProductSearchQuery, useSearch bool) ([]map[string]any, error) {
	if !useSearch {
		rows, err := s.prod.ListProductsPublic(ctx, `p."createdAt" DESC`, 0, 0)
		if err != nil {
			return nil, err
		}
		return s.mapListWithFavorites(ctx, rows, viewer)
	}

	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	sp := repository.ProductSearchParams{SortBy: q.SortBy, Page: q.Page, Limit: q.Limit}
		if q.Search != nil {
			sp.Search = q.Search
		}
		if q.CategorySlug != nil && *q.CategorySlug != "" {
			id, err := s.prod.CategoryIDBySlug(ctx, *q.CategorySlug)
			if err != nil {
				return nil, err
			}
			sp.CategoryID = id
		}
		if q.SubCategorySlug != nil && *q.SubCategorySlug != "" {
			id, err := s.prod.SubCategoryIDBySlug(ctx, *q.SubCategorySlug)
			if err != nil {
				return nil, err
			}
			sp.SubCategoryID = id
		}
		if q.TypeSlug != nil && *q.TypeSlug != "" {
			id, err := s.prod.TypeIDBySlug(ctx, *q.TypeSlug)
			if err != nil {
				return nil, err
			}
			sp.TypeID = id
		}
		sp.MinPrice, sp.MaxPrice = q.MinPrice, q.MaxPrice
		sp.State, sp.Region, sp.ProfileType = q.State, q.Region, q.ProfileType
		if q.FieldValuesJSON != nil && strings.TrimSpace(*q.FieldValuesJSON) != "" {
			m, err := parseFieldValuesMap(*q.FieldValuesJSON)
			if err != nil {
				return nil, &AppError{400, err.Error()}
			}
			sp.FieldValues = m
		}
	rows, err := s.prod.SearchProducts(ctx, sp)
	if err != nil {
		return nil, err
	}
	return s.mapListWithFavorites(ctx, rows, viewer)
}

func (s *ProductService) RandomProducts(ctx context.Context, viewer *int32) ([]map[string]any, error) {
	rows, err := s.prod.RandomSubcategoriesWithProducts(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, pr := range rows {
		var fav bool
		if viewer != nil {
			fav, err = s.prod.IsFavorite(ctx, *viewer, pr.ID)
			if err != nil {
				return nil, err
			}
		}
		out = append(out, s.formatListItem(pr, fav, false, 0))
	}
	return out, nil
}

func (s *ProductService) ProductsByUserID(ctx context.Context, userID int32) ([]map[string]any, error) {
	ok, err := s.prod.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &AppError{400, "Пользователь не найден"}
	}
	rows, err := s.prod.ListProductsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, pr := range rows {
		fav, err := s.prod.IsFavorite(ctx, userID, pr.ID)
		if err != nil {
			return nil, err
		}
		hasPromo := pr.PromotionLevel > 0
		out = append(out, s.formatListItem(pr, fav, hasPromo, pr.PromotionLevel))
	}
	return out, nil
}

func (s *ProductService) AddFavorite(ctx context.Context, userID, productID int32) (map[string]any, error) {
	ok, err := s.prod.ProductExists(ctx, productID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &AppError{400, "Товар не найден"}
	}
	exists, err := s.prod.IsFavorite(ctx, userID, productID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &AppError{400, "Товар уже добавлен в избранное"}
	}
	if err := s.prod.AddFavorite(ctx, userID, productID); err != nil {
		return nil, err
	}
	s.prod.InsertFavoriteAction(ctx, userID, productID)
	return map[string]any{"message": "Товар успешно добавлен в избранное"}, nil
}

func (s *ProductService) RemoveFavorite(ctx context.Context, userID, productID int32) (map[string]any, error) {
	ok, err := s.prod.ProductExists(ctx, productID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &AppError{400, "Товар не найден"}
	}
	if err := s.prod.RemoveFavorite(ctx, userID, productID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{400, "Товар не в избранном"}
		}
		return nil, err
	}
	return map[string]any{"message": "Товар успешно удалён из избранного"}, nil
}

func (s *ProductService) MyFavorites(ctx context.Context, userID int32) ([]map[string]any, error) {
	rows, err := s.prod.ListFavoriteProducts(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, pr := range rows {
		out = append(out, s.formatListItem(pr, true, pr.PromotionLevel > 0, pr.PromotionLevel))
	}
	return out, nil
}

func (s *ProductService) GetProductCard(ctx context.Context, productID int32, viewer *int32) (map[string]any, error) {
	card, err := s.prod.GetProductCard(ctx, productID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{400, "Товар не найден"}
	}
	if err != nil {
		return nil, err
	}
	if viewer != nil && *viewer != card.UserID {
		s.prod.UpsertProductView(ctx, *viewer, productID)
	}
	fvArr := make([]map[string]any, 0, len(card.FieldPairs))
	for _, p := range card.FieldPairs {
		fvArr = append(fvArr, map[string]any{p.FieldName: p.Value})
	}
	var fav bool
	if viewer != nil {
		fav, _ = s.prod.IsFavorite(ctx, *viewer, productID)
	}
	seller, err := s.user.GetUserInfo(ctx, card.UserID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id": card.ID, "name": card.Name, "description": card.Description, "price": card.Price,
		"isHide": card.IsHide, "images": card.Images, "address": card.Address, "userId": card.UserID, "videoUrl": card.VideoURL,
		"category": map[string]any{"id": card.CategoryID, "name": card.CategoryName, "slug": card.CategorySlug},
		"subCategory": map[string]any{"id": card.SubCatID, "name": card.SubCatName, "slug": card.SubCatSlug},
		"type": func() map[string]any {
			if card.TypeID == nil {
				return nil
			}
			return map[string]any{"id": *card.TypeID, "name": card.TypeName, "slug": card.TypeSlug}
		}(),
		"fieldValues":  fvArr,
		"isFavorited":  fav,
		"seller":       seller,
	}, nil
}

func (s *ProductService) ToggleProduct(ctx context.Context, productID, userID int32) (map[string]any, error) {
	uid, _, err := s.prod.GetProductOwnerAndImages(ctx, productID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "Товар не надйен"}
	}
	if err != nil {
		return nil, err
	}
	if uid != userID {
		return nil, &AppError{403, "Вы не можете редактировать не свой товар"}
	}
	if err := s.prod.ToggleProductHide(ctx, productID); err != nil {
		return nil, err
	}
	return map[string]any{"message": "Статус активности товара сменен"}, nil
}

func (s *ProductService) ModerateProduct(ctx context.Context, productID int32, status, reason string) error {
	st := strings.TrimSpace(strings.ToUpper(status))
	if st != "APPROVED" && st != "DENIDED" {
		return &AppError{400, "Неверный статус модерации. Доступные статуты: APPROVED, DENIDED"}
	}
	name, sellerID, err := s.prod.GetProductNameAndSeller(ctx, productID)
	if errors.Is(err, repository.ErrNotFound) {
		return &AppError{404, "Товар для модерации не найден"}
	}
	if err != nil {
		return err
	}
	if st == "DENIDED" && strings.TrimSpace(reason) == "" {
		return &AppError{400, "Необходимо указать причину отказа в модерации"}
	}
	var rptr *string
	if st == "DENIDED" {
		r := strings.TrimSpace(reason)
		rptr = &r
	}
	if err := s.prod.SetModerationState(ctx, productID, st, rptr); err != nil {
		return err
	}
	if st == "DENIDED" && rptr != nil {
		adminID, err := s.prod.FirstAdminUserID(ctx)
		if err != nil || adminID == nil {
			return nil
		}
		chatID, err := s.prod.FindModerationChat(ctx, *adminID, sellerID)
		if err != nil {
			return nil
		}
		var cid int32
		if chatID == nil {
			cid, err = s.prod.CreateModerationChat(ctx, *adminID, sellerID)
			if err != nil {
				return nil
			}
		} else {
			cid = *chatID
		}
		buyerID, _, err := s.prod.GetChatBuyerSeller(ctx, cid)
		if err != nil {
			return nil
		}
		msg := fmt.Sprintf("❌ Ваш товар \"%s\" был отклонен модерацией.\n\nПричина отказа: %s", name, *rptr)
		mid, err := s.prod.InsertChatMessage(ctx, cid, *adminID, msg, productID)
		if err != nil {
			return nil
		}
		_ = s.prod.UpdateChatAfterMessage(ctx, cid, mid, *adminID, buyerID)
	}
	return nil
}

func (s *ProductService) AllProductsToModerate(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.prod.ListProductsModerate(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, pr := range rows {
		out = append(out, s.formatListItem(pr, false, pr.PromotionLevel > 0, pr.PromotionLevel))
	}
	return out, nil
}

func (s *ProductService) AllPromotedProducts(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.prod.ListPromotedProducts(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, pr := range rows {
		out = append(out, map[string]any{
			"id": pr.ProductID, "name": pr.Name, "price": pr.Price, "images": pr.Images,
			"category": pr.Category, "subCategory": pr.SubCategory, "type": pr.Type, "user": pr.User,
			"createdAt": pr.CreatedAt, "promotions": pr.Promotions,
		})
	}
	return out, nil
}

func (s *ProductService) TogglePromotion(ctx context.Context, promotionID int32) (map[string]any, error) {
	pid, pName, promoName, active, start, end, err := s.prod.TogglePromotionActive(ctx, promotionID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "Продвижение не найдено"}
	}
	if err != nil {
		return nil, err
	}
	msg := "отключено"
	if active {
		msg = "включено"
	}
	return map[string]any{
		"message": fmt.Sprintf("Продвижение %s", msg),
		"promotion": map[string]any{
			"id": promotionID, "productId": pid, "productName": pName, "promotionType": promoName,
			"isActive": active, "startDate": start, "endDate": end,
		},
	}, nil
}
