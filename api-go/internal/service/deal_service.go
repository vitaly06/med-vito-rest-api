package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/repository"
)

type DealService struct {
	cfg     config.Config
	repo    *repository.DealPG
	chat    *repository.ChatPG
	payment *PaymentService
	reserve *repository.ReservationPG
	cdek    *CDEKService
}

func NewDealService(cfg config.Config, repo *repository.DealPG, chat *repository.ChatPG, payment *PaymentService, reserve *repository.ReservationPG, cdek *CDEKService) *DealService {
	return &DealService{cfg: cfg, repo: repo, chat: chat, payment: payment, reserve: reserve, cdek: cdek}
}

type CreateDealRequest struct {
	ProductID      int32   `json:"productId"`
	DeliveryCost   int32   `json:"deliveryCost"`
	CDEKTariffCode *int32  `json:"cdekTariffCode"`
	CDEKTariffName *string `json:"cdekTariffName"`
	CDEKFromCity   *int32  `json:"cdekFromCityCode"`
	CDEKToCity     *int32  `json:"cdekToCityCode"`
	CDEKFromPVZ    *string `json:"cdekFromPvzCode"`
	CDEKToPVZ      *string `json:"cdekToPvzCode"`
}

type MarkShippedRequest struct {
	CDEKOrderUUID   *string `json:"cdekOrderUuid"`
	CDEKTrackNumber *string `json:"cdekTrackNumber"`
}

func (s *DealService) CreateDeal(ctx context.Context, buyerID int32, req CreateDealRequest) (map[string]any, error) {
	if req.ProductID <= 0 {
		return nil, &AppError{400, "РќСѓР¶РµРЅ productId"}
	}
	product, err := s.repo.ProductInfo(ctx, req.ProductID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "РўРѕРІР°СЂ РЅРµ РЅР°Р№РґРµРЅ"}
	}
	if err != nil {
		return nil, err
	}
	if product.UserID == buyerID {
		return nil, &AppError{400, "РќРµР»СЊР·СЏ СЃРѕР·РґР°С‚СЊ СЃРґРµР»РєСѓ РЅР° СЃРІРѕР№ С‚РѕРІР°СЂ"}
	}
	if !product.Approved || product.IsHide {
		return nil, &AppError{400, "РўРѕРІР°СЂ РЅРµРґРѕСЃС‚СѓРїРµРЅ РґР»СЏ Р±РµР·РѕРїР°СЃРЅРѕР№ СЃРґРµР»РєРё"}
	}
	if req.DeliveryCost < 0 {
		return nil, &AppError{400, "РЎС‚РѕРёРјРѕСЃС‚СЊ РґРѕСЃС‚Р°РІРєРё РЅРµ РјРѕР¶РµС‚ Р±С‹С‚СЊ РѕС‚СЂРёС†Р°С‚РµР»СЊРЅРѕР№"}
	}

	feePercent := s.cfg.DealPlatformFeePercent
	if feePercent < 0 {
		feePercent = 0
	}
	platformFee := int32(int(product.Price) * feePercent / 100)
	sellerAmount := product.Price - platformFee
	totalAmount := product.Price + req.DeliveryCost

	deal, err := s.repo.Create(ctx, repository.CreateDealParams{
		ProductID:      product.ID,
		BuyerID:        buyerID,
		SellerID:       product.UserID,
		ProductAmount:  product.Price,
		DeliveryCost:   req.DeliveryCost,
		PlatformFee:    platformFee,
		SellerAmount:   sellerAmount,
		TotalAmount:    totalAmount,
		CDEKTariffCode: req.CDEKTariffCode,
		CDEKTariffName: normalizeStringPtr(req.CDEKTariffName),
		CDEKFromCity:   req.CDEKFromCity,
		CDEKToCity:     req.CDEKToCity,
		CDEKFromPVZ:    normalizeStringPtr(req.CDEKFromPVZ),
		CDEKToPVZ:      normalizeStringPtr(req.CDEKToPVZ),
	})
	if err != nil {
		return nil, err
	}
	if s.reserve != nil {
		_ = s.reserve.MarkDealCreated(ctx, deal.ProductID, buyerID)
	}
	s.notifyOrderInChat(ctx, *deal)
	return s.formatDeal(*deal), nil
}

func (s *DealService) PayDeal(ctx context.Context, buyerID, dealID int32) (map[string]any, error) {
	deal, err := s.getUserDeal(ctx, dealID, buyerID, "buyer")
	if err != nil {
		return nil, err
	}
	if deal.Status != "CREATED" {
		return nil, &AppError{400, "РћРїР»Р°С‚РёС‚СЊ РјРѕР¶РЅРѕ С‚РѕР»СЊРєРѕ СЃРѕР·РґР°РЅРЅСѓСЋ СЃРґРµР»РєСѓ"}
	}
	if deal.PaymentID != nil && deal.PaymentURL != nil {
		return map[string]any{"deal": s.formatDeal(*deal), "paymentId": *deal.PaymentID, "paymentUrl": *deal.PaymentURL, "orderId": deal.OrderID}, nil
	}
	paymentID, paymentURL, orderID, err := s.payment.CreateDealPayment(ctx, buyerID, deal.ID, float64(deal.TotalAmount), "Р‘РµР·РѕРїР°СЃРЅР°СЏ СЃРґРµР»РєР°: "+deal.ProductName)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SetPayment(ctx, deal.ID, paymentID, orderID, paymentURL); err != nil {
		return nil, err
	}
	updated, err := s.repo.FindByID(ctx, deal.ID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"deal": s.formatDeal(*updated), "paymentId": paymentID, "paymentUrl": paymentURL, "orderId": orderID}, nil
}

func (s *DealService) MarkShipped(ctx context.Context, sellerID, dealID int32, req MarkShippedRequest) (map[string]any, error) {
	deal, err := s.getUserDeal(ctx, dealID, sellerID, "seller")
	if err != nil {
		return nil, err
	}
	if deal.Status != "PAID" {
		return nil, &AppError{400, "РћС‚РїСЂР°РІРєСѓ РјРѕР¶РЅРѕ РїРѕРґС‚РІРµСЂРґРёС‚СЊ С‚РѕР»СЊРєРѕ РїРѕСЃР»Рµ РѕРїР»Р°С‚С‹"}
	}
	orderUUID := normalizeStringPtr(req.CDEKOrderUUID)
	trackNumber := normalizeStringPtr(req.CDEKTrackNumber)
	if trackNumber == nil && orderUUID != nil && s.cdek != nil {
		if fetchedTrack := s.cdek.TrackNumberByOrderUUID(ctx, *orderUUID); fetchedTrack != nil {
			trackNumber = fetchedTrack
		}
	}
	if orderUUID != nil || trackNumber != nil {
		if err := s.repo.SetCDEKShipment(ctx, dealID, orderUUID, trackNumber); err != nil {
			return nil, err
		}
	}
	if err := s.repo.SetStatus(ctx, dealID, []string{"PAID"}, "SHIPPED", "shippedAt"); err != nil {
		return nil, &AppError{400, "РќРµ СѓРґР°Р»РѕСЃСЊ РїРѕРґС‚РІРµСЂРґРёС‚СЊ РѕС‚РїСЂР°РІРєСѓ"}
	}
	return s.GetDeal(ctx, sellerID, dealID)
}

func (s *DealService) ConfirmDelivery(ctx context.Context, buyerID, dealID int32) (map[string]any, error) {
	deal, err := s.getUserDeal(ctx, dealID, buyerID, "buyer")
	if err != nil {
		return nil, err
	}
	if deal.Status != "SHIPPED" {
		return nil, &AppError{400, "РџРѕР»СѓС‡РµРЅРёРµ РјРѕР¶РЅРѕ РїРѕРґС‚РІРµСЂРґРёС‚СЊ С‚РѕР»СЊРєРѕ РїРѕСЃР»Рµ РѕС‚РїСЂР°РІРєРё"}
	}
	delay := s.cfg.DealPayoutDelayDays
	if delay < 0 {
		delay = 0
	}
	payoutAt := time.Now().AddDate(0, 0, delay)
	if err := s.repo.MarkDelivered(ctx, dealID, payoutAt); err != nil {
		return nil, &AppError{400, "РќРµ СѓРґР°Р»РѕСЃСЊ РїРѕРґС‚РІРµСЂРґРёС‚СЊ РїРѕР»СѓС‡РµРЅРёРµ"}
	}
	return s.GetDeal(ctx, buyerID, dealID)
}

func (s *DealService) OpenDispute(ctx context.Context, userID, dealID int32, reason string) (map[string]any, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, &AppError{400, "РќСѓР¶РЅРѕ СѓРєР°Р·Р°С‚СЊ РїСЂРёС‡РёРЅСѓ СЃРїРѕСЂР°"}
	}
	if _, err := s.getUserDeal(ctx, dealID, userID, "participant"); err != nil {
		return nil, err
	}
	if err := s.repo.OpenDispute(ctx, dealID, strings.TrimSpace(reason)); err != nil {
		return nil, &AppError{400, "РќРµ СѓРґР°Р»РѕСЃСЊ РѕС‚РєСЂС‹С‚СЊ СЃРїРѕСЂ"}
	}
	return s.GetDeal(ctx, userID, dealID)
}

func (s *DealService) CancelDeal(ctx context.Context, userID, dealID int32) (map[string]any, error) {
	deal, err := s.getUserDeal(ctx, dealID, userID, "participant")
	if err != nil {
		return nil, err
	}
	if deal.Status != "CREATED" && deal.Status != "PAID" {
		return nil, &AppError{400, "Отменить можно только до оформления доставки"}
	}
	if deal.CDEKOrderUUID != nil || deal.CDEKTrackNumber != nil {
		return nil, &AppError{400, "Доставка уже оформлена, отмена недоступна"}
	}
	if err := s.repo.SetStatus(ctx, dealID, []string{"CREATED", "PAID"}, "CANCELLED", "cancelledAt"); err != nil {
		return nil, &AppError{400, "Не удалось отменить сделку"}
	}
	return s.GetDeal(ctx, userID, dealID)
}

func (s *DealService) GetDeal(ctx context.Context, userID, dealID int32) (map[string]any, error) {
	deal, err := s.getUserDeal(ctx, dealID, userID, "participant")
	if err != nil {
		return nil, err
	}
	deal = s.refreshDealTrackFromCDEK(ctx, deal)
	return s.formatDeal(*deal), nil
}

func (s *DealService) GetDealCDEKQR(ctx context.Context, userID, dealID int32) (map[string]any, error) {
	deal, err := s.getUserDeal(ctx, dealID, userID, "participant")
	if err != nil {
		return nil, err
	}
	deal = s.refreshDealTrackFromCDEK(ctx, deal)
	if deal.CDEKOrderUUID == nil || strings.TrimSpace(*deal.CDEKOrderUUID) == "" {
		return nil, &AppError{404, "Для сделки еще не сохранен orderUuid CDEK"}
	}
	if s.cdek == nil {
		return nil, &AppError{400, "CDEK сервис не инициализирован"}
	}
	qrData, qrURL := s.cdek.QRByOrderUUID(ctx, *deal.CDEKOrderUUID)
	if qrData == nil && qrURL == nil {
		return nil, &AppError{404, "CDEK не вернул QR по этой отправке"}
	}
	trackPending := deal.CDEKTrackNumber == nil || strings.TrimSpace(*deal.CDEKTrackNumber) == ""
	return map[string]any{
		"qrCodeData":   qrData,
		"qrCodeUrl":    qrURL,
		"trackNumber":  deal.CDEKTrackNumber,
		"trackingUrl":  buildCDEKTrackingURL(deal.CDEKTrackNumber),
		"orderUuid":    deal.CDEKOrderUUID,
		"trackPending": trackPending,
	}, nil
}

func (s *DealService) MyPurchases(ctx context.Context, buyerID int32) ([]map[string]any, error) {
	deals, err := s.repo.ListByBuyer(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	deals, _ = s.refreshDealsTracksFromCDEK(ctx, deals, 5)
	return s.formatDeals(deals), nil
}

func (s *DealService) MySales(ctx context.Context, sellerID int32) ([]map[string]any, error) {
	deals, err := s.repo.ListBySeller(ctx, sellerID)
	if err != nil {
		return nil, err
	}
	deals, _ = s.refreshDealsTracksFromCDEK(ctx, deals, 5)
	return s.formatDeals(deals), nil
}

func (s *DealService) MyAllDeals(ctx context.Context, userID int32) ([]map[string]any, error) {
	purchases, err := s.repo.ListByBuyer(ctx, userID)
	if err != nil {
		return nil, err
	}
	sales, err := s.repo.ListBySeller(ctx, userID)
	if err != nil {
		return nil, err
	}
	purchases, usedBudget := s.refreshDealsTracksFromCDEK(ctx, purchases, 5)
	remainingBudget := 5 - usedBudget
	if remainingBudget < 0 {
		remainingBudget = 0
	}
	sales, _ = s.refreshDealsTracksFromCDEK(ctx, sales, remainingBudget)
	all := make([]map[string]any, 0, len(purchases)+len(sales))
	for _, deal := range purchases {
		item := s.formatDeal(deal)
		item["myRole"] = "buyer"
		all = append(all, item)
	}
	for _, deal := range sales {
		item := s.formatDeal(deal)
		item["myRole"] = "seller"
		all = append(all, item)
	}
	return all, nil
}

func (s *DealService) StartPayoutWorker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				autoDays := s.cfg.DealAutoCompleteDays
				if autoDays <= 0 {
					autoDays = 14
				}
				payoutDelay := s.cfg.DealPayoutDelayDays
				if payoutDelay < 0 {
					payoutDelay = 0
				}
				autoDelivered, err := s.repo.AutoDeliverExpiredShipped(ctx, time.Now().AddDate(0, 0, -autoDays), time.Now().AddDate(0, 0, payoutDelay), 20)
				if err != nil {
					log.Printf("deal auto-complete worker: %v", err)
				} else if autoDelivered > 0 {
					log.Printf("deal auto-complete worker: auto-delivered %d deal(s)", autoDelivered)
				}
				n, err := s.repo.CompleteDuePayouts(ctx, 20)
				if err != nil {
					log.Printf("deal payout worker: %v", err)
					continue
				}
				if n > 0 {
					log.Printf("deal payout worker: completed %d deal(s)", n)
				}
			}
		}
	}()
}

func (s *DealService) getUserDeal(ctx context.Context, dealID, userID int32, role string) (*repository.DealRow, error) {
	deal, err := s.repo.FindByID(ctx, dealID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "РЎРґРµР»РєР° РЅРµ РЅР°Р№РґРµРЅР°"}
	}
	if err != nil {
		return nil, err
	}
	switch role {
	case "buyer":
		if deal.BuyerID != userID {
			return nil, &AppError{403, "Р­С‚Рѕ РЅРµ РІР°С€Р° РїРѕРєСѓРїРєР°"}
		}
	case "seller":
		if deal.SellerID != userID {
			return nil, &AppError{403, "Р­С‚Рѕ РЅРµ РІР°С€Р° РїСЂРѕРґР°Р¶Р°"}
		}
	default:
		if deal.BuyerID != userID && deal.SellerID != userID {
			return nil, &AppError{403, "Р’С‹ РЅРµ СѓС‡Р°СЃС‚РЅРёРє СЃРґРµР»РєРё"}
		}
	}
	return deal, nil
}

func (s *DealService) formatDeals(deals []repository.DealRow) []map[string]any {
	out := make([]map[string]any, 0, len(deals))
	for _, deal := range deals {
		out = append(out, s.formatDeal(deal))
	}
	return out
}

func (s *DealService) formatDeal(deal repository.DealRow) map[string]any {
	return map[string]any{
		"id":         deal.ID,
		"status":     localizeDealStatus(deal.Status),
		"statusCode": deal.Status,
		"product": map[string]any{
			"id":     deal.ProductID,
			"name":   deal.ProductName,
			"images": deal.ProductImages,
		},
		"buyer": map[string]any{
			"id":       deal.BuyerID,
			"fullName": deal.BuyerName,
		},
		"seller": map[string]any{
			"id":       deal.SellerID,
			"fullName": deal.SellerName,
		},
		"amounts": map[string]any{
			"productAmount": deal.ProductAmount,
			"deliveryCost":  deal.DeliveryCost,
			"platformFee":   deal.PlatformFee,
			"sellerAmount":  deal.SellerAmount,
			"totalAmount":   deal.TotalAmount,
		},
		"paymentId":     deal.PaymentID,
		"orderId":       deal.OrderID,
		"paymentUrl":    deal.PaymentURL,
		"cdek":          formatDealCDEK(deal),
		"disputeReason": deal.DisputeReason,
		"paidAt":        deal.PaidAt,
		"shippedAt":     deal.ShippedAt,
		"deliveredAt":   deal.DeliveredAt,
		"payoutAt":      deal.PayoutAt,
		"completedAt":   deal.CompletedAt,
		"cancelledAt":   deal.CancelledAt,
		"refundedAt":    deal.RefundedAt,
		"createdAt":     deal.CreatedAt,
		"updatedAt":     deal.UpdatedAt,
	}
}

func formatDealCDEK(deal repository.DealRow) map[string]any {
	trackPending := false
	if deal.CDEKOrderUUID != nil {
		orderUUID := strings.TrimSpace(*deal.CDEKOrderUUID)
		track := ""
		if deal.CDEKTrackNumber != nil {
			track = strings.TrimSpace(*deal.CDEKTrackNumber)
		}
		trackPending = orderUUID != "" && track == ""
	}
	return map[string]any{
		"tariffCode":   deal.CDEKTariffCode,
		"tariffName":   deal.CDEKTariffName,
		"fromCityCode": deal.CDEKFromCity,
		"toCityCode":   deal.CDEKToCity,
		"fromPvzCode":  deal.CDEKFromPVZ,
		"toPvzCode":    deal.CDEKToPVZ,
		"orderUuid":    deal.CDEKOrderUUID,
		"trackNumber":  deal.CDEKTrackNumber,
		"trackingUrl":  buildCDEKTrackingURL(deal.CDEKTrackNumber),
		"trackPending": trackPending,
	}
}

func buildCDEKTrackingURL(trackNumber *string) *string {
	if trackNumber == nil {
		return nil
	}
	track := strings.TrimSpace(*trackNumber)
	if track == "" {
		return nil
	}
	url := "https://www.cdek.ru/ru/tracking?order_id=" + track
	return &url
}

func (s *DealService) refreshDealTrackFromCDEK(ctx context.Context, deal *repository.DealRow) *repository.DealRow {
	if deal == nil || s.cdek == nil || deal.CDEKOrderUUID == nil {
		return deal
	}
	orderUUID := strings.TrimSpace(*deal.CDEKOrderUUID)
	if orderUUID == "" {
		return deal
	}

	liveTrack := s.cdek.TrackNumberByOrderUUID(ctx, orderUUID)
	if liveTrack == nil {
		return deal
	}
	currentTrack := ""
	if deal.CDEKTrackNumber != nil {
		currentTrack = strings.TrimSpace(*deal.CDEKTrackNumber)
	}
	if currentTrack == strings.TrimSpace(*liveTrack) {
		return deal
	}

	// Обновляем трек в БД, чтобы дальше фронт получал консистентные данные.
	if err := s.repo.SetCDEKShipment(ctx, deal.ID, nil, liveTrack); err != nil {
		return deal
	}
	updated, err := s.repo.FindByID(ctx, deal.ID)
	if err != nil {
		return deal
	}
	return updated
}

func (s *DealService) refreshDealsTracksFromCDEK(ctx context.Context, deals []repository.DealRow, maxLive int) ([]repository.DealRow, int) {
	if maxLive <= 0 || len(deals) == 0 {
		return deals, 0
	}
	used := 0
	for i := range deals {
		if used >= maxLive {
			break
		}
		if !dealNeedsLiveTrackSync(deals[i]) {
			continue
		}
		updated := s.refreshDealTrackFromCDEK(ctx, &deals[i])
		if updated != nil {
			deals[i] = *updated
		}
		used++
	}
	return deals, used
}

func dealNeedsLiveTrackSync(deal repository.DealRow) bool {
	if deal.CDEKOrderUUID == nil || strings.TrimSpace(*deal.CDEKOrderUUID) == "" {
		return false
	}
	switch deal.Status {
	case "PAID", "SHIPPED", "DELIVERED", "DISPUTE":
		return true
	default:
		return false
	}
}

func (s *DealService) notifyOrderInChat(ctx context.Context, deal repository.DealRow) {
	if s.chat == nil {
		return
	}
	chat, err := s.chat.FindChatByProductParticipants(ctx, deal.ProductID, deal.BuyerID, deal.SellerID)
	if err != nil {
		return
	}
	var chatID int32
	if chat == nil {
		chatID, err = s.chat.InsertProductChat(ctx, deal.ProductID, deal.BuyerID, deal.SellerID)
		if err != nil {
			return
		}
	} else {
		chatID = chat.ID
	}
	content := fmt.Sprintf(
		"Оформлен заказ по товару \"%s\". Сумма: %d ₽ (товар: %d ₽, доставка: %d ₽). Номер сделки: #%d.",
		deal.ProductName,
		deal.TotalAmount,
		deal.ProductAmount,
		deal.DeliveryCost,
		deal.ID,
	)
	_, _, _, _, _, _ = s.chat.InsertChatMessage(ctx, chatID, deal.BuyerID, content)
}

func normalizeStringPtr(v *string) *string {
	if v == nil {
		return nil
	}
	t := strings.TrimSpace(*v)
	if t == "" {
		return nil
	}
	return &t
}
