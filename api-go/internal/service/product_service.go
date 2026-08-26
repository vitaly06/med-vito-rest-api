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
	prod    *repository.ProductPG
	s3      *s3client.Client
	user    *UserService
	stat    *StatisticsService
	support *SupportService
}

func NewProductService(prod *repository.ProductPG, s3 *s3client.Client, user *UserService, stat *StatisticsService, support *SupportService) *ProductService {
	return &ProductService{prod: prod, s3: s3, user: user, stat: stat, support: support}
}

func formatProductDate(t time.Time) string {
	lt := t.In(time.Local)
	d := lt.Day()
	m := int(lt.Month())
	y := lt.Year() % 100
	h, mi := lt.Hour(), lt.Minute()
	return fmt.Sprintf("%02d.%02d.%02d в %02d:%02d", d, m, y, h, mi)
}

func appendProductLifetime(out map[string]any, createdAt time.Time) {
	expiresAt := createdAt.Add(30 * 24 * time.Hour)
	remainingDuration := time.Until(expiresAt)
	remainingDays := int(remainingDuration.Hours() / 24)
	if remainingDuration > 0 && remainingDuration%(24*time.Hour) != 0 {
		remainingDays++
	}
	if remainingDays < 0 {
		remainingDays = 0
	}
	out["expiresAt"] = expiresAt.Format(time.RFC3339)
	out["daysUntilExpiration"] = remainingDays
	out["isExpired"] = remainingDuration <= 0
}
func (s *ProductService) formatListItem(pr repository.ProductListRow, isFav bool, hasPromo bool, promoLevel int32) map[string]any {
	badges := make([]string, 0, 3)
	if pr.PromotionLevel >= 100 {
		badges = append(badges, "Премиум")
	}
	if pr.SellerVerified {
		badges = append(badges, "Верифицирован")
	}
	if pr.ViewsCount > 100 {
		badges = append(badges, "Срочно")
	}
	out := map[string]any{
		"id":              pr.ID,
		"images":          pr.Images,
		"name":            pr.Name,
		"address":         pr.Address,
		"createdAt":       formatProductDate(pr.CreatedAt),
		"price":           pr.Price,
		"quantity":        pr.Quantity,
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
		"promotionName":   pr.PromotionName,
		"sellerRating":    pr.SellerRating,
		"sellerVerified":  pr.SellerVerified,
		"viewsCount":      pr.ViewsCount,
		"popularityScore": pr.PopularityScore,
		"badges":          badges,
		"isReserved":      pr.IsReserved,
		"isPaid":          hasPromo,
		"highlighted":     promoLevel > 0,
		"hasBanner":       promoLevel >= 100,
	}
	if pr.ModerateState != nil {
		out["moderateState"] = *pr.ModerateState
	}
	if pr.ModerationRejectionReason != nil {
		out["moderationRejectionReason"] = *pr.ModerationRejectionReason
	}
	out["isHide"] = pr.IsHide

	// Expiration info
	if !pr.ExpiresAt.IsZero() {
		expires := pr.ExpiresAt
		out["expiresAt"] = expires.Format(time.RFC3339)
		remaining := time.Until(expires)
		remainingDays := int(remaining.Hours() / 24)
		if remaining > 0 && remaining%(24*time.Hour) != 0 {
			remainingDays++
		}
		if remainingDays < 0 {
			remainingDays = 0
		}
		out["daysUntilExpiration"] = remainingDays
		out["isExpired"] = remaining <= 0
	} else {
		appendProductLifetime(out, pr.CreatedAt)
	}

	return out
}

func enforcePaidSlots(rows []repository.ProductListRow, page, limit int) []repository.ProductListRow {
	if limit <= 0 || len(rows) == 0 {
		return rows
	}

	maxPaid := 5
	minPaid := 2
	if limit < 20 {
		maxPaid = limit / 4
		if maxPaid < 1 {
			maxPaid = 1
		}
		minPaid = 0
	}

	var paid []repository.ProductListRow
	var free []repository.ProductListRow
	for _, r := range rows {
		if r.PromotionLevel > 0 {
			paid = append(paid, r)
		} else {
			free = append(free, r)
		}
	}

	targetPaid := len(paid)
	if targetPaid > maxPaid {
		targetPaid = maxPaid
	}
	if targetPaid < minPaid && len(paid) >= minPaid {
		targetPaid = minPaid
	}
	targetFree := limit - targetPaid
	if targetFree < 0 {
		targetFree = 0
	}
	if targetFree > len(free) {
		targetFree = len(free)
	}

	// Добираем до лимита оставшимися платными/бесплатными, если одна из групп закончилась.
	for targetPaid+targetFree < limit {
		if targetFree < len(free) {
			targetFree++
			continue
		}
		if targetPaid < len(paid) {
			targetPaid++
			continue
		}
		break
	}

	usePaid := paid[:targetPaid]
	useFree := free[:targetFree]
	out := make([]repository.ProductListRow, 0, targetPaid+targetFree)

	// Равномерно распределяем платные внутри ленты.
	pi := 0
	fi := 0
	segment := 0
	if targetPaid > 0 {
		segment = targetFree / targetPaid
	}
	if segment < 1 {
		segment = 1
	}
	for fi < len(useFree) || pi < len(usePaid) {
		for k := 0; k < segment && fi < len(useFree); k++ {
			out = append(out, useFree[fi])
			fi++
		}
		if pi < len(usePaid) {
			out = append(out, usePaid[pi])
			pi++
		}
	}

	if len(out) > limit {
		out = out[:limit]
	}
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
func (s *ProductService) CreateProduct(ctx context.Context, userID int32, name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, fieldJSON, videoStr string, files []UploadedFile) (map[string]any, error) {
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
	if err := validateProductContentLimits(len(files), description, false); err != nil {
		return nil, err
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
	quantity := int32(1)
	if strings.TrimSpace(quantityStr) != "" {
		q, err := strconv.Atoi(strings.TrimSpace(quantityStr))
		if err != nil || q < 1 {
			return nil, &AppError{400, "Количество должно быть целым числом больше 0"}
		}
		quantity = int32(q)
	}
	beforeRemaining := int32(-1)
	if s.user != nil {
		if before, err := s.user.GetRemainingFreeAds(ctx, userID); err == nil {
			if v, ok := before["remaining"].(int32); ok {
				beforeRemaining = v
			} else if v, ok := before["remaining"].(int); ok {
				beforeRemaining = int32(v)
			} else if v, ok := before["remaining"].(float64); ok {
				beforeRemaining = int32(v)
			}
		}
	}
	err = s.prod.CreateProductTx(ctx, pid, userID, name, int32(price), quantity, state, description, address, vptr, urls, int32(catID), int32(subID), typePtr, intMap)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{400, "Пользователь не найден"}
	}
	if errors.Is(err, repository.ErrInsufficientFunds) {
		return nil, &AppError{400, fmt.Sprintf("Недостаточно средств для размещения объявления. Пополните баланс или бонусы — требуется %d ₽.", repository.AdListingCost)}
	}
	if err != nil {
		return nil, &AppError{400, "Ошибка при создании продукта: " + err.Error()}
	}
	s.notifyAdListingBilling(ctx, userID, beforeRemaining)
	prod, err := s.prod.LoadProductWithRelations(ctx, pid)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"message":       "Продукт успешно создан",
		"product":       prod,
		"isDraft":       false,
		"moderateState": "MODERATE",
	}, nil
}

// UploadedFile — один файл из multipart (images).
func (s *ProductService) CreateDraft(ctx context.Context, userID int32, name, priceStr, quantityStr, state, description, address, categoryStr, subStr, typeStr, fieldJSON, videoStr string, files []UploadedFile) (map[string]any, error) {
	// Черновик: без жёсткой валидации; в БД подставляем дефолты. S3 нужен только если реально грузим файлы.
	if len(files) > 0 && s.s3 == nil {
		return nil, &AppError{500, "S3 не настроен"}
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = "Черновик"
	}
	address = strings.TrimSpace(address)
	if address == "" {
		address = "—"
	}
	state = strings.TrimSpace(strings.ToUpper(state))
	if state != "NEW" && state != "USED" {
		state = "NEW"
	}

	price := 1
	if p := strings.TrimSpace(priceStr); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n >= 1 {
			price = n
		}
	}

	catID64, errCat := strconv.ParseInt(strings.TrimSpace(categoryStr), 10, 32)
	subID64, errSub := strconv.ParseInt(strings.TrimSpace(subStr), 10, 32)
	var catID, subID int32
	if errCat == nil && errSub == nil {
		ok, err := s.prod.SubCategoryBelongsToCategory(ctx, int32(catID64), int32(subID64))
		if err != nil {
			return nil, err
		}
		if ok {
			catID, subID = int32(catID64), int32(subID64)
		}
	}
	if catID == 0 || subID == 0 {
		dc, ds, err := s.prod.FirstCategoryWithSubcategory(ctx)
		if err != nil {
			return nil, err
		}
		catID, subID = dc, ds
	}

	var typePtr *int32
	if strings.TrimSpace(typeStr) != "" {
		if tid64, err := strconv.ParseInt(strings.TrimSpace(typeStr), 10, 32); err == nil {
			t := int32(tid64)
			typePtr = &t
		}
	}

	fvMap, err := parseFieldValuesMap(fieldJSON)
	if err != nil {
		fvMap = map[string]string{}
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
	var draftFieldIDs []int32
	for k, v := range fvMap {
		if id64, err := strconv.ParseInt(k, 10, 32); err == nil {
			fid := int32(id64)
			draftFieldIDs = append(draftFieldIDs, fid)
			intMap[fid] = v
		}
	}
	if len(draftFieldIDs) > 0 {
		if err := s.prod.ValidateTypeFieldIDs(ctx, draftFieldIDs); err != nil {
			intMap = make(map[int32]string)
		}
	}
	var vptr *string
	if strings.TrimSpace(videoStr) != "" {
		v := strings.TrimSpace(videoStr)
		vptr = &v
	}
	quantity := int32(1)
	if strings.TrimSpace(quantityStr) != "" {
		if q, err := strconv.Atoi(strings.TrimSpace(quantityStr)); err == nil && q >= 1 {
			quantity = int32(q)
		}
	}
	if err := s.prod.CreateDraftTx(ctx, pid, userID, name, int32(price), quantity, state, description, address, vptr, urls, catID, subID, typePtr, intMap); err != nil {
		return nil, err
	}
	prod, err := s.prod.LoadProductWithRelations(ctx, pid)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"message":       "Черновик сохранен",
		"product":       prod,
		"isDraft":       true,
		"moderateState": "DRAFT",
	}, nil
}

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

func (s *ProductService) notifyAdListingBilling(ctx context.Context, userID int32, beforeRemaining int32) {
	if s.support == nil || s.user == nil {
		return
	}
	after, err := s.user.GetRemainingFreeAds(ctx, userID)
	if err != nil {
		return
	}
	remaining := int32(0)
	total := int32(repository.FreeAdsMonthlyLimit)
	cost := int32(repository.AdListingCost)
	if v, ok := after["remaining"].(int32); ok {
		remaining = v
	} else if v, ok := after["remaining"].(int); ok {
		remaining = int32(v)
	}
	if v, ok := after["total"].(int32); ok {
		total = v
	} else if v, ok := after["total"].(int); ok {
		total = int32(v)
	}
	if v, ok := after["costPerAd"].(int32); ok {
		cost = v
	} else if v, ok := after["costPerAd"].(int); ok {
		cost = int32(v)
	}

	var text string
	if beforeRemaining <= 0 {
		text = fmt.Sprintf("Списано %d ₽ за размещение объявления (сначала бонусы, затем баланс). Бесплатный лимит (%d шт./мес.) исчерпан.", cost, total)
	} else if remaining == 0 && beforeRemaining == 1 {
		text = fmt.Sprintf("Использовано последнее бесплатное объявление в этом месяце (%d из %d). Следующие размещения — %d ₽.", total, total, cost)
	} else if beforeRemaining > remaining {
		text = fmt.Sprintf("Размещено бесплатное объявление. Осталось %d из %d бесплатных в этом месяце.", remaining, total)
	} else {
		return
	}
	s.support.NotifyUserBilling(ctx, userID, text)
}
