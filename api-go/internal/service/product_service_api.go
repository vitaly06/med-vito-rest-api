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
		return nil, &AppError{400, "РўРѕРІР°СЂ РґР»СЏ СѓРґР°Р»РµРЅРёСЏ РЅРµ РЅР°Р№РґРµРЅ"}
	}
	if err != nil {
		return nil, err
	}
	if uid != userID {
		return nil, &AppError{403, "Р’С‹ РЅРµ РјРѕР¶РµС‚Рµ СѓРґР°Р»РёС‚СЊ С‡СѓР¶РѕР№ С‚РѕРІР°СЂ"}
	}
	if s.s3 != nil {
		for _, u := range imgs {
			_ = s.s3.DeleteByURL(ctx, u)
		}
	}
	if err := s.prod.DeleteProductByID(ctx, productID); err != nil {
		return nil, err
	}
	return map[string]any{"message": "РўРѕРІР°СЂ СѓСЃРїРµС€РЅРѕ СѓРґР°Р»С‘РЅ"}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, productID, userID int32, name, priceStr, quantityStr, state, description, address, videoStr, fieldJSON string, files []UploadedFile) (map[string]any, error) {
	if s.s3 == nil && len(files) > 0 {
		return nil, &AppError{500, "S3 РЅРµ РЅР°СЃС‚СЂРѕРµРЅ"}
	}
	uid, typeID, existingImages, modState, err := s.prod.ProductWithTypeForUpdate(ctx, productID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{400, "РўРѕРІР°СЂ РЅРµ РЅР°Р№РґРµРЅ"}
	}
	if err != nil {
		return nil, err
	}
	if uid != userID {
		return nil, &AppError{403, "Р’С‹ РЅРµ РјРѕР¶РµС‚Рµ СЂРµРґР°РєС‚РёСЂРѕРІР°С‚СЊ С‡СѓР¶РѕР№ С‚РѕРІР°СЂ"}
	}
	isDraft := modState == "DRAFT"
	fvMap, err := parseFieldValuesMap(fieldJSON)
	if err != nil {
		if !isDraft {
			return nil, &AppError{400, err.Error()}
		}
		fvMap = map[string]string{}
	}
	var fieldIDs []int32
	for k := range fvMap {
		if id64, err := strconv.ParseInt(k, 10, 32); err == nil {
			fieldIDs = append(fieldIDs, int32(id64))
		}
	}
	if len(fieldIDs) > 0 && !isDraft {
		if err := s.prod.ValidateTypeFieldIDs(ctx, fieldIDs); err != nil {
			return nil, &AppError{400, err.Error()}
		}
		if typeID != nil {
			badNames, err := s.prod.TypeFieldsBelongToType(ctx, *typeID, fieldIDs)
			if err != nil {
				return nil, err
			}
			if len(badNames) > 0 {
				return nil, &AppError{400, fmt.Sprintf("РџРѕР»СЏ %s РЅРµ РїСЂРёРЅР°РґР»РµР¶Р°С‚ С‚РёРїСѓ СЌС‚РѕРіРѕ С‚РѕРІР°СЂР°", strings.Join(badNames, ", "))}
			}
		}
	}

	var namePtr, statePtr, descPtr, addrPtr, vidPtr *string
	var quantityPtr *int32
	var pricePtr *int32
	if strings.TrimSpace(name) != "" {
		v := strings.TrimSpace(name)
		namePtr = &v
	}
	if strings.TrimSpace(priceStr) != "" {
		p, err := strconv.Atoi(strings.TrimSpace(priceStr))
		if err != nil || p < 1 {
			if !isDraft {
				return nil, &AppError{400, "РќРµРєРѕСЂСЂРµРєС‚РЅР°СЏ С†РµРЅР°"}
			}
		} else {
			pp := int32(p)
			pricePtr = &pp
		}
	}
	if strings.TrimSpace(state) != "" {
		if isDraft {
			v := strings.TrimSpace(strings.ToUpper(state))
			if v == "NEW" || v == "USED" {
				statePtr = &v
			}
		} else {
			v := strings.TrimSpace(state)
			statePtr = &v
		}
	}
	if strings.TrimSpace(quantityStr) != "" {
		q, err := strconv.Atoi(strings.TrimSpace(quantityStr))
		if err != nil || q < 1 {
			if !isDraft {
				return nil, &AppError{400, "Количество должно быть целым числом больше 0"}
			}
		} else {
			qq := int32(q)
			quantityPtr = &qq
		}
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
			return nil, &AppError{400, "РћС€РёР±РєР° Р·Р°РіСЂСѓР·РєРё РІ S3: " + err.Error()}
		}
		newImages = append(newImages, u)
	}
	var imgsArg []string
	if len(newImages) > 0 {
		imgsArg = append(append([]string{}, existingImages...), newImages...)
	}

	if err := s.prod.UpdateProductPartial(ctx, productID, namePtr, pricePtr, quantityPtr, statePtr, descPtr, addrPtr, vidPtr, imgsArg); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{400, "РўРѕРІР°СЂ РЅРµ РЅР°Р№РґРµРЅ"}
		}
		return nil, &AppError{400, "РћС€РёР±РєР° РїСЂРё РѕР±РЅРѕРІР»РµРЅРёРё: " + err.Error()}
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
	return map[string]any{"message": "РўРѕРІР°СЂ СѓСЃРїРµС€РЅРѕ РѕР±РЅРѕРІР»С‘РЅ", "product": prod}, nil
}

func (s *ProductService) AvailableFilters(ctx context.Context, catSlug, subSlug, typeSlug *string) (map[string]any, error) {
	return s.prod.BuildAvailableFilters(ctx, catSlug, subSlug, typeSlug)
}

// ProductSearchQuery вЂ” query-РїР°СЂР°РјРµС‚СЂС‹ РїРѕРёСЃРєР°.
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
		rows, err := s.prod.ListProductsPublic(ctx, `CASE
			WHEN COALESCE((
				SELECT MAX(pr."pricePerDay") FROM "ProductPromotion" pp
				JOIN "Promotion" pr ON pr.id = pp."promotionId"
				WHERE pp."productId" = p.id AND pp."isActive" AND pp."isPaid" AND pp."endDate" >= NOW()
			), 0) >= 100 THEN 1
			WHEN COALESCE((
				SELECT MAX(pr."pricePerDay") FROM "ProductPromotion" pp
				JOIN "Promotion" pr ON pr.id = pp."promotionId"
				WHERE pp."productId" = p.id AND pp."isActive" AND pp."isPaid" AND pp."endDate" >= NOW()
			), 0) > 0 THEN 2
			ELSE 3
		END ASC, p."createdAt" DESC`, 0, 0)
		if err != nil {
			return nil, err
		}
		if s.stat != nil {
			var uid *int32
			if viewer != nil {
				uid = viewer
			}
			s.stat.TrackSearch(ctx, uid, q.Search, q.Region, q.CategorySlug, q.SubCategorySlug, q.TypeSlug, len(rows))
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
	rows = enforcePaidSlots(rows, q.Page, q.Limit)
	if s.stat != nil {
		var uid *int32
		if viewer != nil {
			uid = viewer
		}
		s.stat.TrackSearch(ctx, uid, q.Search, q.Region, q.CategorySlug, q.SubCategorySlug, q.TypeSlug, len(rows))
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

func (s *ProductService) RecommendedBySubcategory(ctx context.Context, viewer *int32, subcategoryID int32, limit int) ([]map[string]any, error) {
	rows, err := s.prod.RecommendedBySubcategory(ctx, subcategoryID, limit)
	if err != nil {
		return nil, err
	}
	paidInTop := 0
	for i := 0; i < len(rows) && i < 6; i++ {
		if rows[i].PromotionLevel > 0 {
			paidInTop++
			if paidInTop > 2 {
				for j := i + 1; j < len(rows); j++ {
					if rows[j].PromotionLevel == 0 {
						rows[i], rows[j] = rows[j], rows[i]
						break
					}
				}
			}
		}
	}
	return s.mapListWithFavorites(ctx, rows, viewer)
}

func (s *ProductService) ProductsByUserID(ctx context.Context, viewer *int32, userID int32) ([]map[string]any, error) {
	ok, err := s.prod.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &AppError{400, "РџРѕР»СЊР·РѕРІР°С‚РµР»СЊ РЅРµ РЅР°Р№РґРµРЅ"}
	}
	rows, err := s.prod.ListProductsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, pr := range rows {
		if pr.ModerateState != nil && *pr.ModerateState == "DRAFT" && (viewer == nil || *viewer != userID) {
			continue
		}
		fav, err := s.prod.IsFavorite(ctx, userID, pr.ID)
		if err != nil {
			return nil, err
		}
		hasPromo := pr.PromotionLevel > 0
		out = append(out, s.formatListItem(pr, fav, hasPromo, pr.PromotionLevel))
	}
	return out, nil
}

func (s *ProductService) MyDrafts(ctx context.Context, userID int32) ([]map[string]any, error) {
	rows, err := s.prod.ListDraftProductsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, pr := range rows {
		out = append(out, s.formatListItem(pr, false, pr.PromotionLevel > 0, pr.PromotionLevel))
	}
	return out, nil
}

func (s *ProductService) PublishDraft(ctx context.Context, productID, userID int32) (map[string]any, error) {
	card, err := s.prod.GetProductCard(ctx, productID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "Черновик не найден"}
	}
	if err != nil {
		return nil, err
	}
	if card.UserID != userID {
		return nil, &AppError{403, "Нет доступа к этому объявлению"}
	}
	if card.ModerateState != "DRAFT" {
		return nil, &AppError{400, "Это не черновик"}
	}
	name := strings.TrimSpace(card.Name)
	if name == "" || name == "Черновик" {
		return nil, &AppError{400, "Укажите название объявления"}
	}
	addr := strings.TrimSpace(card.Address)
	if addr == "" || addr == "—" {
		return nil, &AppError{400, "Укажите адрес"}
	}
	if card.Price < 1 {
		return nil, &AppError{400, "Укажите цену"}
	}
	st := strings.TrimSpace(strings.ToUpper(card.State))
	if st != "NEW" && st != "USED" {
		return nil, &AppError{400, "Укажите состояние товара (NEW или USED)"}
	}

	if err := s.prod.PublishDraftTx(ctx, productID, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Черновик не найден"}
		}
		if errors.Is(err, repository.ErrInsufficientFunds) {
			return nil, &AppError{400, fmt.Sprintf("Недостаточно средств для публикации. Требуется %d руб.", repository.AdListingCost)}
		}
		return nil, err
	}
	prod, err := s.prod.LoadProductWithRelations(ctx, productID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"message":       "Черновик опубликован и отправлен на модерацию",
		"product":       prod,
		"isDraft":       false,
		"moderateState": "MODERATE",
	}, nil
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
		return nil, &AppError{400, "РўРѕРІР°СЂ РЅРµ РЅР°Р№РґРµРЅ"}
	}
	if err := s.prod.RemoveFavorite(ctx, userID, productID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{400, "РўРѕРІР°СЂ РЅРµ РІ РёР·Р±СЂР°РЅРЅРѕРј"}
		}
		return nil, err
	}
	return map[string]any{"message": "РўРѕРІР°СЂ СѓСЃРїРµС€РЅРѕ СѓРґР°Р»С‘РЅ РёР· РёР·Р±СЂР°РЅРЅРѕРіРѕ"}, nil
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
		return nil, &AppError{400, "РўРѕРІР°СЂ РЅРµ РЅР°Р№РґРµРЅ"}
	}
	if err != nil {
		return nil, err
	}
	if card.ModerateState == "DRAFT" && (viewer == nil || *viewer != card.UserID) {
		return nil, &AppError{404, "РўРѕРІР°СЂ РЅРµ РЅР°Р№РґРµРЅ"}
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
		"id": card.ID, "name": card.Name, "description": card.Description, "price": card.Price, "quantity": card.Quantity,
		"isHide": card.IsHide, "images": card.Images, "address": card.Address, "userId": card.UserID, "videoUrl": card.VideoURL,
		"category":    map[string]any{"id": card.CategoryID, "name": card.CategoryName, "slug": card.CategorySlug},
		"subCategory": map[string]any{"id": card.SubCatID, "name": card.SubCatName, "slug": card.SubCatSlug},
		"type": func() map[string]any {
			if card.TypeID == nil {
				return nil
			}
			return map[string]any{"id": *card.TypeID, "name": card.TypeName, "slug": card.TypeSlug}
		}(),
		"fieldValues": fvArr,
		"isFavorited": fav,
		"seller":      seller,
	}, nil
}

func (s *ProductService) ToggleProduct(ctx context.Context, productID, userID int32) (map[string]any, error) {
	uid, _, err := s.prod.GetProductOwnerAndImages(ctx, productID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "РўРѕРІР°СЂ РЅРµ РЅР°РґР№РµРЅ"}
	}
	if err != nil {
		return nil, err
	}
	if uid != userID {
		return nil, &AppError{403, "Р’С‹ РЅРµ РјРѕР¶РµС‚Рµ СЂРµРґР°РєС‚РёСЂРѕРІР°С‚СЊ РЅРµ СЃРІРѕР№ С‚РѕРІР°СЂ"}
	}
	if err := s.prod.ToggleProductHide(ctx, productID); err != nil {
		return nil, err
	}
	return map[string]any{"message": "РЎС‚Р°С‚СѓСЃ Р°РєС‚РёРІРЅРѕСЃС‚Рё С‚РѕРІР°СЂР° СЃРјРµРЅРµРЅ"}, nil
}

func (s *ProductService) ModerateProduct(ctx context.Context, productID int32, status, reason string) error {
	st := strings.TrimSpace(strings.ToUpper(status))
	if st != "APPROVED" && st != "DENIDED" {
		return &AppError{400, "РќРµРІРµСЂРЅС‹Р№ СЃС‚Р°С‚СѓСЃ РјРѕРґРµСЂР°С†РёРё. Р”РѕСЃС‚СѓРїРЅС‹Рµ СЃС‚Р°С‚СѓС‚С‹: APPROVED, DENIDED"}
	}
	name, sellerID, err := s.prod.GetProductNameAndSeller(ctx, productID)
	if errors.Is(err, repository.ErrNotFound) {
		return &AppError{404, "РўРѕРІР°СЂ РґР»СЏ РјРѕРґРµСЂР°С†РёРё РЅРµ РЅР°Р№РґРµРЅ"}
	}
	if err != nil {
		return err
	}
	if st == "DENIDED" && strings.TrimSpace(reason) == "" {
		return &AppError{400, "РќРµРѕР±С…РѕРґРёРјРѕ СѓРєР°Р·Р°С‚СЊ РїСЂРёС‡РёРЅСѓ РѕС‚РєР°Р·Р° РІ РјРѕРґРµСЂР°С†РёРё"}
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
		msg := fmt.Sprintf("вќЊ Р’Р°С€ С‚РѕРІР°СЂ \"%s\" Р±С‹Р» РѕС‚РєР»РѕРЅРµРЅ РјРѕРґРµСЂР°С†РёРµР№.\n\nРџСЂРёС‡РёРЅР° РѕС‚РєР°Р·Р°: %s", name, *rptr)
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
		return nil, &AppError{404, "РџСЂРѕРґРІРёР¶РµРЅРёРµ РЅРµ РЅР°Р№РґРµРЅРѕ"}
	}
	if err != nil {
		return nil, err
	}
	msg := "РѕС‚РєР»СЋС‡РµРЅРѕ"
	if active {
		msg = "РІРєР»СЋС‡РµРЅРѕ"
	}
	return map[string]any{
		"message": fmt.Sprintf("РџСЂРѕРґРІРёР¶РµРЅРёРµ %s", msg),
		"promotion": map[string]any{
			"id": promotionID, "productId": pid, "productName": pName, "promotionType": promoName,
			"isActive": active, "startDate": start, "endDate": end,
		},
	}, nil
}
