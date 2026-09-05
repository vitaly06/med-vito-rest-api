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
	logs    *repository.LogPG
	chat    *repository.ChatPG
	payment *PaymentService
	reserve *repository.ReservationPG
	cdek    *CDEKService
	users   *repository.UserPG
}

func dealTrace(dealID, userID int32, stage string, format string, args ...any) {
	msg := strings.TrimSpace(fmt.Sprintf(format, args...))
	if msg == "" {
		log.Printf("deal_trace deal=%d user=%d stage=%s", dealID, userID, stage)
		return
	}
	log.Printf("deal_trace deal=%d user=%d stage=%s %s", dealID, userID, stage, msg)
}

func trimOrDash(v *string) string {
	if v == nil {
		return "-"
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return "-"
	}
	return s
}

func NewDealService(cfg config.Config, repo *repository.DealPG, logs *repository.LogPG, chat *repository.ChatPG, payment *PaymentService, reserve *repository.ReservationPG, cdek *CDEKService, users *repository.UserPG) *DealService {
	return &DealService{cfg: cfg, repo: repo, logs: logs, chat: chat, payment: payment, reserve: reserve, cdek: cdek, users: users}
}

type CreateDealRequest struct {
	ProductID         int32   `json:"productId"`
	DeliveryCost      int32   `json:"deliveryCost"`
	CDEKTariffCode    *int32  `json:"cdekTariffCode"`
	CDEKTariffName    *string `json:"cdekTariffName"`
	CDEKFromCity      *int32  `json:"cdekFromCityCode"`
	CDEKToCity        *int32  `json:"cdekToCityCode"`
	CDEKFromPVZ       *string `json:"cdekFromPvzCode"`
	CDEKToPVZ         *string `json:"cdekToPvzCode"`
	CDEKToAddress     *string `json:"cdekToAddress"`
	CDEKPackageWeight *int32  `json:"cdekPackageWeight"`
	CDEKPackageLength *int32  `json:"cdekPackageLength"`
	CDEKPackageWidth  *int32  `json:"cdekPackageWidth"`
	CDEKPackageHeight *int32  `json:"cdekPackageHeight"`
	CDEKRecipientMode *string `json:"cdekRecipientMode"`
}

type SetCdekHandoffRequest struct {
	Mode            string  `json:"mode"`
	CDEKFromPVZ     *string `json:"cdekFromPvzCode"`
	CDEKFromAddress *string `json:"cdekFromAddress"`
}

type MarkShippedRequest struct {
	CDEKOrderUUID   *string `json:"cdekOrderUuid"`
	CDEKTrackNumber *string `json:"cdekTrackNumber"`
}

// dealHasCdekAutoRoute — в сделке сохранён маршрут калькулятора CDEK (тариф + города).
func dealHasCdekAutoRoute(deal *repository.DealRow) bool {
	if deal == nil || deal.CDEKTariffCode == nil || *deal.CDEKTariffCode <= 0 {
		return false
	}
	if deal.CDEKFromCity == nil || *deal.CDEKFromCity <= 0 || deal.CDEKToCity == nil || *deal.CDEKToCity <= 0 {
		return false
	}
	return true
}

// sellerHandoffReady — продавец указал способ передачи (ПВЗ или курьер).
func sellerHandoffReady(deal *repository.DealRow) bool {
	if deal == nil || deal.CDEKSellerHandoff == nil {
		return false
	}
	mode := strings.TrimSpace(*deal.CDEKSellerHandoff)
	switch mode {
	case "pvz":
		return deal.CDEKFromPVZ != nil && strings.TrimSpace(*deal.CDEKFromPVZ) != ""
	case "courier":
		return deal.CDEKFromAddress != nil && strings.TrimSpace(*deal.CDEKFromAddress) != ""
	default:
		return false
	}
}

func dealPackageDims(deal *repository.DealRow) (weight, length, width, height int) {
	weight, length, width, height = 1000, 20, 20, 20
	if deal == nil {
		return
	}
	if deal.CDEKPackageWeight != nil && *deal.CDEKPackageWeight > 0 {
		weight = int(*deal.CDEKPackageWeight)
	}
	if deal.CDEKPackageLength != nil && *deal.CDEKPackageLength > 0 {
		length = int(*deal.CDEKPackageLength)
	}
	if deal.CDEKPackageWidth != nil && *deal.CDEKPackageWidth > 0 {
		width = int(*deal.CDEKPackageWidth)
	}
	if deal.CDEKPackageHeight != nil && *deal.CDEKPackageHeight > 0 {
		height = int(*deal.CDEKPackageHeight)
	}
	return
}

// ensureCdekOrderRegistered — создаёт заказ в CDEK после выбора передачи продавцом (идемпотентно по number).
func (s *DealService) ensureCdekOrderRegistered(ctx context.Context, deal *repository.DealRow) (*string, *string, error) {
	if deal == nil {
		return nil, nil, nil
	}
	dealTrace(deal.ID, deal.SellerID, "cdek_register_enter", "status=%s hasOrder=%t hasRoute=%t handoffReady=%t cdekOn=%t",
		deal.Status,
		deal.CDEKOrderUUID != nil && strings.TrimSpace(*deal.CDEKOrderUUID) != "",
		dealHasCdekAutoRoute(deal),
		sellerHandoffReady(deal),
		s.cdek != nil && s.cdek.configured(),
	)
	if deal.CDEKOrderUUID != nil && strings.TrimSpace(*deal.CDEKOrderUUID) != "" {
		dealTrace(deal.ID, deal.SellerID, "cdek_register_skip", "reason=order_exists orderUuid=%s", trimOrDash(deal.CDEKOrderUUID))
		return nil, nil, nil
	}
	if !dealHasCdekAutoRoute(deal) {
		dealTrace(deal.ID, deal.SellerID, "cdek_register_skip", "reason=no_route")
		return nil, nil, nil
	}
	if !sellerHandoffReady(deal) {
		dealTrace(deal.ID, deal.SellerID, "cdek_register_skip", "reason=handoff_not_ready mode=%s fromPvz=%s fromAddr=%s",
			trimOrDash(deal.CDEKSellerHandoff), trimOrDash(deal.CDEKFromPVZ), trimOrDash(deal.CDEKFromAddress))
		return nil, nil, nil
	}
	if s.cdek == nil || !s.cdek.configured() {
		return nil, nil, nil
	}
	if s.users == nil {
		return nil, nil, &AppError{400, "Автозаказ CDEK: нет доступа к данным пользователей"}
	}
	buyer, err := s.users.FindUserByID(ctx, deal.BuyerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, &AppError{404, "Покупатель не найден"}
		}
		return nil, nil, err
	}
	seller, err := s.users.FindUserByID(ctx, deal.SellerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, &AppError{404, "Продавец не найден"}
		}
		return nil, nil, err
	}
	bPhone, err := PhoneForCdekAPI(buyer.PhoneNumber)
	if err != nil {
		return nil, nil, &AppError{400, "Телефон покупателя для CDEK: " + err.Error()}
	}
	sPhone, err := PhoneForCdekAPI(seller.PhoneNumber)
	if err != nil {
		return nil, nil, &AppError{400, "Телефон продавца для CDEK: " + err.Error()}
	}
	clientNumber := fmt.Sprintf("med-vito-deal-%d", deal.ID)
	in := CDEKCreateOrderInput{
		TariffCode:      int(*deal.CDEKTariffCode),
		FromCityCode:    int(*deal.CDEKFromCity),
		ToCityCode:      int(*deal.CDEKToCity),
		FromPVZ:         deal.CDEKFromPVZ,
		ToPVZ:           deal.CDEKToPVZ,
		ToAddress:       deal.CDEKToAddress,
		ClientNumber:    clientNumber,
		Comment:         fmt.Sprintf("Безопасная сделка #%d — %s", deal.ID, deal.ProductName),
		SenderName:      seller.FullName,
		SenderPhone:     sPhone,
		RecipientName:   buyer.FullName,
		RecipientPhone:  bPhone,
		PackageName:     deal.ProductName,
		WareKey:         fmt.Sprintf("deal-%d", deal.ID),
		DeclaredCostRub: float64(deal.ProductAmount),
	}
	in.WeightGrams, in.LengthCm, in.WidthCm, in.HeightCm = dealPackageDims(deal)
	if deal.CDEKSellerHandoff != nil && strings.TrimSpace(*deal.CDEKSellerHandoff) == "courier" {
		in.FromAddress = deal.CDEKFromAddress
	}
	if in.TariffCode == cdekAllowedTariffCode {
		if in.FromPVZ == nil || strings.TrimSpace(*in.FromPVZ) == "" {
			return nil, nil, &AppError{400, "Для тарифа 136 укажи ПВЗ отправителя (shipment_point)"}
		}
		if in.ToPVZ == nil || strings.TrimSpace(*in.ToPVZ) == "" {
			return nil, nil, &AppError{400, "Для тарифа 136 укажи ПВЗ получателя (delivery_point)"}
		}
		in.FromAddress = nil
		in.ToAddress = nil
	}
	log.Printf(
		"cdek auto-register start deal=%d status=%s tariff=%d fromCity=%d toCity=%d hasToPvz=%t hasFromPvz=%t hasToAddress=%t",
		deal.ID,
		deal.Status,
		in.TariffCode,
		in.FromCityCode,
		in.ToCityCode,
		in.ToPVZ != nil && strings.TrimSpace(*in.ToPVZ) != "",
		in.FromPVZ != nil && strings.TrimSpace(*in.FromPVZ) != "",
		in.ToAddress != nil && strings.TrimSpace(*in.ToAddress) != "",
	)

	res, err := s.cdek.CreateOrder(ctx, in)
	if err != nil {
		dealTrace(deal.ID, deal.SellerID, "cdek_create_order_error", "clientNumber=%s err=%v", clientNumber, err)
		log.Printf("cdek auto-register create failed deal=%d clientNumber=%s err=%v", deal.ID, clientNumber, err)
		res2, lookErr := s.cdek.LookupOrderByClientNumber(ctx, clientNumber)
		if lookErr == nil && res2 != nil && res2.OrderUUID != nil && strings.TrimSpace(*res2.OrderUUID) != "" {
			log.Printf("cdek auto-register found existing order by number deal=%d clientNumber=%s orderUUID=%s", deal.ID, clientNumber, strings.TrimSpace(*res2.OrderUUID))
			res, err = res2, nil
		}
	}
	if err != nil {
		if cdekNeedsRecipientAddress(err) {
			dealTrace(deal.ID, deal.SellerID, "cdek_register_fail", "reason=recipient_address_required")
			return nil, nil, &AppError{400, "CDEK требует адрес получателя для выбранного режима доставки"}
		}
		dealTrace(deal.ID, deal.SellerID, "cdek_register_fail", "err=%v", err)
		return nil, nil, err
	}
	if res == nil || res.OrderUUID == nil {
		return nil, nil, &AppError{502, "CDEK не вернул uuid заказа"}
	}
	if byNumber, lookErr := s.cdek.LookupOrderByClientNumber(ctx, clientNumber); lookErr == nil && byNumber != nil {
		// Prefer canonical lookup by client number: CDEK can return a request UUID in create response,
		// while list lookup contains the final order UUID and sometimes track.
		if byNumber.OrderUUID != nil && strings.TrimSpace(*byNumber.OrderUUID) != "" {
			res.OrderUUID = byNumber.OrderUUID
		}
		if byNumber.Track != nil && strings.TrimSpace(*byNumber.Track) != "" {
			res.Track = byNumber.Track
		}
	}
	u := strings.TrimSpace(*res.OrderUUID)
	uu := u
	var tt *string
	if res.Track != nil {
		if t := strings.TrimSpace(*res.Track); t != "" {
			tt = &t
		}
	}
	if tt == nil {
		if live := s.cdek.TrackNumberByOrderUUID(ctx, u); live != nil {
			tt = live
		}
	}
	if err := s.repo.SetCDEKShipment(ctx, deal.ID, &uu, tt); err != nil {
		dealTrace(deal.ID, deal.SellerID, "cdek_register_fail", "save_shipment_failed err=%v", err)
		return nil, nil, err
	}
	dealTrace(deal.ID, deal.SellerID, "cdek_register_success", "orderUuid=%s track=%s", uu, trimOrDash(tt))
	log.Printf("cdek auto-register success deal=%d orderUUID=%s hasTrack=%t", deal.ID, uu, tt != nil && strings.TrimSpace(*tt) != "")
	return &uu, tt, nil
}

func (s *DealService) resolvePvzAddress(ctx context.Context, toCityCode int, pvzCode string) *string {
	if s.cdek == nil || toCityCode <= 0 {
		return nil
	}
	code := strings.TrimSpace(pvzCode)
	if code == "" {
		return nil
	}
	points, err := s.cdek.DeliveryPoints(ctx, toCityCode)
	if err != nil {
		log.Printf("cdek pvz-address resolve failed city=%d pvz=%s: %v", toCityCode, code, err)
		return nil
	}
	for _, p := range points {
		rawCode := strings.TrimSpace(fmt.Sprint(p["code"]))
		if !strings.EqualFold(rawCode, code) {
			continue
		}
		loc, _ := p["location"].(map[string]any)
		if loc == nil {
			return nil
		}
		addr := strings.TrimSpace(fmt.Sprint(loc["address"]))
		if addr == "" || addr == "<nil>" {
			return nil
		}
		log.Printf("cdek pvz-address resolved city=%d pvz=%s addr=%s", toCityCode, code, addr)
		return &addr
	}
	log.Printf("cdek pvz-address not found city=%d pvz=%s", toCityCode, code)
	return nil
}

// ensureCdekOrdersInList — для PAID без uuid пробуем создать заказ CDEK (ограничение, чтобы не долбить API).
func (s *DealService) ensureCdekOrdersInList(ctx context.Context, deals []repository.DealRow, max int) {
	if max <= 0 || s.cdek == nil || !s.cdek.configured() {
		return
	}
	used := 0
	for i := range deals {
		if used >= max {
			break
		}
		if deals[i].Status != "PAID" || !dealHasCdekAutoRoute(&deals[i]) || !sellerHandoffReady(&deals[i]) {
			continue
		}
		if deals[i].CDEKOrderUUID != nil && strings.TrimSpace(*deals[i].CDEKOrderUUID) != "" {
			continue
		}
		if _, _, err := s.ensureCdekOrderRegistered(ctx, &deals[i]); err != nil {
			log.Printf("cdek auto-order list deal=%d: %v", deals[i].ID, err)
			continue
		}
		if fresh, e := s.repo.FindByID(ctx, deals[i].ID); e == nil {
			deals[i] = *fresh
		}
		used++
	}
}

func (s *DealService) CreateDeal(ctx context.Context, buyerID int32, req CreateDealRequest) (map[string]any, error) {
	if req.ProductID <= 0 {
		return nil, &AppError{400, "Нужен productId"}
	}
	product, err := s.repo.ProductInfo(ctx, req.ProductID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "Товар не найден"}
	}
	if err != nil {
		return nil, err
	}
	if product.UserID == buyerID {
		return nil, &AppError{400, "Нельзя создать сделку на свой товар"}
	}
	if !product.Approved || product.IsHide {
		return nil, &AppError{400, "Товар недоступен для безопасной сделки"}
	}
	if req.DeliveryCost < 0 {
		return nil, &AppError{400, "Стоимость доставки не может быть отрицательной"}
	}

	if s.users != nil {
		buyer, err := s.users.FindUserByID(ctx, buyerID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, &AppError{404, "Покупатель не найден"}
			}
			return nil, err
		}
		if isDealPhoneSynthetic(buyer.PhoneNumber) {
			return nil, &AppError{400, "Нельзя создать сделку с подставным телефоном в профиле. Укажи реальный номер +7 в настройках аккаунта."}
		}
		seller, err := s.users.FindUserByID(ctx, product.UserID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, &AppError{404, "Продавец не найден"}
			}
			return nil, err
		}
		if isDealPhoneSynthetic(seller.PhoneNumber) {
			return nil, &AppError{400, "Продавец не указал реальный телефон в профиле — безопасная сделка недоступна, пока он не сохранит +7."}
		}
	}

	feePercent := s.cfg.DealPlatformFeePercent
	if feePercent < 0 {
		feePercent = 0
	}
	toPvz := normalizeStringPtr(req.CDEKToPVZ)
	if toPvz == nil {
		return nil, &AppError{400, "Укажи ПВЗ получателя (cdekToPvzCode)"}
	}
	if req.CDEKTariffCode != nil && *req.CDEKTariffCode != 136 {
		return nil, &AppError{400, "Доступен только тариф 136 (склад-склад)"}
	}
	recipientMode := "pvz"
	if req.CDEKPackageWeight != nil && *req.CDEKPackageWeight <= 0 {
		return nil, &AppError{400, "Вес посылки должен быть больше 0"}
	}

	platformFee := int32(int(product.Price) * feePercent / 100)
	sellerAmount := product.Price - platformFee
	totalAmount := product.Price + req.DeliveryCost
	recipientModePtr := recipientMode

	deal, err := s.repo.Create(ctx, repository.CreateDealParams{
		ProductID:         product.ID,
		BuyerID:           buyerID,
		SellerID:          product.UserID,
		ProductAmount:     product.Price,
		DeliveryCost:      req.DeliveryCost,
		PlatformFee:       platformFee,
		SellerAmount:      sellerAmount,
		TotalAmount:       totalAmount,
		CDEKTariffCode:    req.CDEKTariffCode,
		CDEKTariffName:    normalizeStringPtr(req.CDEKTariffName),
		CDEKFromCity:      req.CDEKFromCity,
		CDEKToCity:        req.CDEKToCity,
		CDEKFromPVZ:       normalizeStringPtr(req.CDEKFromPVZ),
		CDEKToPVZ:         toPvz,
		CDEKToAddress:     normalizeStringPtr(req.CDEKToAddress),
		CDEKPackageWeight: req.CDEKPackageWeight,
		CDEKPackageLength: req.CDEKPackageLength,
		CDEKPackageWidth:  req.CDEKPackageWidth,
		CDEKPackageHeight: req.CDEKPackageHeight,
		CDEKRecipientMode: &recipientModePtr,
	})
	if err != nil {
		return nil, err
	}
	if s.reserve != nil {
		_ = s.reserve.MarkDealCreated(ctx, deal.ProductID, buyerID)
	}
	s.writeDealLog(ctx, buyerID, deal.ID, "create", "deal created")
	s.notifyOrderInChat(ctx, *deal)
	return s.formatDeal(*deal), nil
}

func (s *DealService) PayDeal(ctx context.Context, buyerID, dealID int32) (map[string]any, error) {
	dealTrace(dealID, buyerID, "pay_enter", "init")
	deal, err := s.getUserDeal(ctx, dealID, buyerID, "buyer")
	if err != nil {
		dealTrace(dealID, buyerID, "pay_fail", "get_user_deal_err=%v", err)
		return nil, err
	}
	if deal.Status != "CREATED" {
		dealTrace(dealID, buyerID, "pay_skip", "reason=bad_status status=%s", deal.Status)
		return nil, &AppError{400, "Оплатить можно только созданную сделку"}
	}
	// Демо: без Tinkoff сразу PAID — иначе CDEK и «отправить» никогда не откроются.
	if !s.payment.tinkoffConfigured() && s.cfg.DealAllowMockPayment {
		ok, err := s.repo.TryMarkPaidByDealID(ctx, deal.ID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, &AppError{400, "Не удалось отметить оплату (сделка не в статусе CREATED)"}
		}
		updated, err := s.repo.FindByID(ctx, deal.ID)
		if err != nil {
			return nil, err
		}
		s.writeDealLog(ctx, buyerID, deal.ID, "pay", "deal paid (mock)")
		mockPID := fmt.Sprintf("mock-%d", deal.ID)
		return map[string]any{
			"deal":        s.formatDeal(*updated),
			"paymentId":   mockPID,
			"paymentUrl":  "",
			"orderId":     "mock",
			"mockPayment": true,
		}, nil
	}
	if deal.PaymentID != nil && deal.PaymentURL != nil {
		dealTrace(dealID, buyerID, "pay_skip", "reason=payment_already_initialized paymentId=%s", trimOrDash(deal.PaymentID))
		return map[string]any{"deal": s.formatDeal(*deal), "paymentId": *deal.PaymentID, "paymentUrl": *deal.PaymentURL, "orderId": deal.OrderID}, nil
	}
	// Temporary test mode: always initialize deal payment for 1 RUB.
	paymentID, paymentURL, orderID, err := s.payment.CreateDealPayment(ctx, buyerID, deal.ID, 1, "Безопасная сделка: "+deal.ProductName)
	if err != nil {
		dealTrace(dealID, buyerID, "pay_fail", "create_deal_payment_err=%v", err)
		return nil, err
	}
	if err := s.repo.SetPayment(ctx, deal.ID, paymentID, orderID, paymentURL); err != nil {
		dealTrace(dealID, buyerID, "pay_fail", "set_payment_err=%v", err)
		return nil, err
	}
	updated, err := s.repo.FindByID(ctx, deal.ID)
	if err != nil {
		return nil, err
	}
	s.writeDealLog(ctx, buyerID, deal.ID, "payment_init", fmt.Sprintf("payment initialized paymentId=%s orderId=%s", paymentID, orderID))
	dealTrace(dealID, buyerID, "pay_success", "mode=tinkoff paymentId=%s orderId=%s", paymentID, orderID)
	return map[string]any{"deal": s.formatDeal(*updated), "paymentId": paymentID, "paymentUrl": paymentURL, "orderId": orderID}, nil
}

// SyncDealPayment — если webhook Тинькофф не дошёл (например localhost), покупатель дергает после оплаты.
func (s *DealService) SyncDealPayment(ctx context.Context, buyerID, dealID int32) (map[string]any, error) {
	dealTrace(dealID, buyerID, "sync_payment_enter", "init")
	deal, err := s.getUserDeal(ctx, dealID, buyerID, "buyer")
	if err != nil {
		dealTrace(dealID, buyerID, "sync_payment_fail", "get_user_deal_err=%v", err)
		return nil, err
	}
	if deal.Status != "CREATED" {
		return nil, &AppError{400, "Сделка уже не в статусе «создана» — синхронизация не нужна"}
	}
	if deal.PaymentID == nil || strings.TrimSpace(*deal.PaymentID) == "" {
		dealTrace(dealID, buyerID, "sync_payment_fail", "reason=no_payment_id")
		return nil, &AppError{400, "Сначала нажми оплату и получи paymentId"}
	}
	pid := strings.TrimSpace(*deal.PaymentID)
	dealTrace(dealID, buyerID, "sync_payment_check", "paymentId=%s", pid)
	if strings.HasPrefix(strings.ToLower(pid), "mock-") {
		ok, err := s.repo.TryMarkPaidByDealID(ctx, deal.ID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, &AppError{400, "Не удалось отметить демо-оплату"}
		}
		updated, err := s.repo.FindByID(ctx, deal.ID)
		if err != nil {
			return nil, err
		}
		s.writeDealLog(ctx, buyerID, deal.ID, "pay_sync", "deal paid via sync (mock)")
		return map[string]any{"deal": s.formatDeal(*updated)}, nil
	}
	st, err := s.payment.CheckPaymentStatus(ctx, pid)
	if err != nil {
		dealTrace(dealID, buyerID, "sync_payment_fail", "check_status_err=%v", err)
		return nil, err
	}
	status := strings.TrimSpace(fmt.Sprint(st["status"]))
	dealTrace(dealID, buyerID, "sync_payment_status", "paymentId=%s status=%s", pid, status)
	if status != "AUTHORIZED" && status != "CONFIRMED" {
		return nil, &AppError{400, fmt.Sprintf("Платёж ещё не подтверждён (статус: %s)", status)}
	}
	ok, err := s.repo.TryMarkPaidByPaymentID(ctx, pid)
	if err != nil {
		dealTrace(dealID, buyerID, "sync_payment_fail", "mark_paid_by_payment_err=%v", err)
		return nil, err
	}
	if !ok {
		dealTrace(dealID, buyerID, "sync_payment_fail", "mark_paid_by_payment_ok=false")
		return nil, &AppError{400, "Не удалось перевести сделку в оплаченную"}
	}
	updated, err := s.repo.FindByID(ctx, deal.ID)
	if err != nil {
		return nil, err
	}
	s.writeDealLog(ctx, buyerID, deal.ID, "pay_sync", "deal paid via sync")
	dealTrace(dealID, buyerID, "sync_payment_success", "mode=tinkoff status=%s", updated.Status)
	return map[string]any{"deal": s.formatDeal(*updated)}, nil
}

func (s *DealService) SetCdekHandoff(ctx context.Context, sellerID, dealID int32, req SetCdekHandoffRequest) (map[string]any, error) {
	dealTrace(dealID, sellerID, "set_handoff_enter", "mode=%s fromPvz=%s hasFromAddress=%t", strings.TrimSpace(req.Mode), trimOrDash(req.CDEKFromPVZ), normalizeStringPtr(req.CDEKFromAddress) != nil)
	deal, err := s.getUserDeal(ctx, dealID, sellerID, "seller")
	if err != nil {
		dealTrace(dealID, sellerID, "set_handoff_fail", "get_user_deal_err=%v", err)
		return nil, err
	}
	if deal.Status != "PAID" {
		return nil, &AppError{400, "Способ передачи можно указать только после оплаты"}
	}
	mode := strings.ToLower(strings.TrimSpace(req.Mode))
	if mode != "pvz" && mode != "courier" {
		return nil, &AppError{400, "mode: pvz или courier"}
	}
	if deal.CDEKTariffCode != nil && int(*deal.CDEKTariffCode) == cdekAllowedTariffCode && mode != "pvz" {
		return nil, &AppError{400, "Для тарифа 136 доступна только передача через ПВЗ (shipment_point)"}
	}
	var fromPvz, fromAddr *string
	switch mode {
	case "pvz":
		fromPvz = normalizeStringPtr(req.CDEKFromPVZ)
		if fromPvz == nil {
			return nil, &AppError{400, "Укажи код ПВЗ СДЭК, куда сдашь посылку (cdekFromPvzCode)"}
		}
		if deal.CDEKFromCity == nil || *deal.CDEKFromCity <= 0 {
			return nil, &AppError{400, "Не удалось определить город отправителя для выбора ПВЗ"}
		}
		if s.cdek != nil {
			points, pErr := s.cdek.DeliveryPoints(ctx, int(*deal.CDEKFromCity))
			if pErr != nil {
				dealTrace(dealID, sellerID, "set_handoff_fail", "delivery_points_err=%v", pErr)
				return nil, &AppError{400, "Не удалось проверить ПВЗ отправителя"}
			}
			valid := false
			wanted := strings.TrimSpace(*fromPvz)
			for _, p := range points {
				if strings.EqualFold(strings.TrimSpace(fmt.Sprint(p["code"])), wanted) {
					valid = true
					break
				}
			}
			if !valid {
				return nil, &AppError{400, "Выбранный ПВЗ не найден в городе отправителя"}
			}
		}
	case "courier":
		fromAddr = normalizeStringPtr(req.CDEKFromAddress)
		if fromAddr == nil {
			return nil, &AppError{400, "Укажи адрес забора курьером (cdekFromAddress)"}
		}
	}
	if err := s.repo.SetCDEKSellerHandoff(ctx, dealID, mode, fromPvz, fromAddr); err != nil {
		dealTrace(dealID, sellerID, "set_handoff_fail", "save_handoff_err=%v", err)
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{400, "Не удалось сохранить способ передачи"}
		}
		return nil, err
	}
	fresh, err := s.repo.FindByID(ctx, dealID)
	if err != nil {
		return nil, err
	}
	if _, _, regErr := s.ensureCdekOrderRegistered(ctx, fresh); regErr != nil {
		dealTrace(dealID, sellerID, "set_handoff_fail", "register_after_handoff_err=%v", regErr)
		return nil, regErr
	}
	if updated, e := s.repo.FindByID(ctx, dealID); e == nil {
		fresh = updated
	}
	s.writeDealLog(ctx, sellerID, dealID, "cdek_handoff", "seller handoff: "+mode)
	dealTrace(dealID, sellerID, "set_handoff_success", "mode=%s orderUuid=%s track=%s", mode, trimOrDash(fresh.CDEKOrderUUID), trimOrDash(fresh.CDEKTrackNumber))
	return s.formatDeal(*s.refreshDealFromCDEK(ctx, fresh)), nil
}

func (s *DealService) MarkShipped(ctx context.Context, sellerID, dealID int32, req MarkShippedRequest) (map[string]any, error) {
	dealTrace(dealID, sellerID, "mark_shipped_enter", "reqOrderUuid=%s reqTrack=%s", trimOrDash(req.CDEKOrderUUID), trimOrDash(req.CDEKTrackNumber))
	deal, err := s.getUserDeal(ctx, dealID, sellerID, "seller")
	if err != nil {
		dealTrace(dealID, sellerID, "mark_shipped_fail", "get_user_deal_err=%v", err)
		return nil, err
	}
	if deal.Status != "PAID" {
		return nil, &AppError{400, "Отправку можно подтвердить только после оплаты"}
	}
	if dealHasCdekAutoRoute(deal) && !sellerHandoffReady(deal) {
		return nil, &AppError{400, "Сначала выбери способ передачи в СДЭК (ПВЗ или курьер)"}
	}
	orderUUID := normalizeStringPtr(req.CDEKOrderUUID)
	trackNumber := normalizeStringPtr(req.CDEKTrackNumber)
	if orderUUID == nil && deal.CDEKOrderUUID != nil {
		if u := strings.TrimSpace(*deal.CDEKOrderUUID); u != "" {
			orderUUID = &u
		}
	}
	if trackNumber == nil && deal.CDEKTrackNumber != nil {
		if t := strings.TrimSpace(*deal.CDEKTrackNumber); t != "" {
			trackNumber = &t
		}
	}
	hasRoute := dealHasCdekAutoRoute(deal)
	cdekOn := s.cdek != nil && s.cdek.configured()
	if hasRoute && cdekOn && orderUUID == nil && trackNumber == nil {
		dealTrace(dealID, sellerID, "mark_shipped_register", "reason=no_order_and_no_track")
		uu, tt, regErr := s.ensureCdekOrderRegistered(ctx, deal)
		if regErr != nil {
			dealTrace(dealID, sellerID, "mark_shipped_fail", "register_err=%v", regErr)
			return nil, regErr
		}
		if uu != nil {
			orderUUID = uu
		}
		if tt != nil {
			trackNumber = tt
		}
	}
	if hasRoute && cdekOn && orderUUID == nil {
		dealTrace(dealID, sellerID, "mark_shipped_fail", "reason=no_order_uuid_after_register")
		return nil, &AppError{400, "Сначала оформи передачу в СДЭК — должен появиться UUID заказа"}
	}
	if trackNumber == nil && orderUUID != nil && s.cdek != nil {
		if fetchedTrack := s.cdek.TrackNumberByOrderUUID(ctx, *orderUUID); fetchedTrack != nil {
			trackNumber = fetchedTrack
		}
	}
	if orderUUID != nil || trackNumber != nil {
		dealTrace(dealID, sellerID, "mark_shipped_save_shipment", "orderUuid=%s track=%s", trimOrDash(orderUUID), trimOrDash(trackNumber))
		if err := s.repo.SetCDEKShipment(ctx, dealID, orderUUID, trackNumber); err != nil {
			dealTrace(dealID, sellerID, "mark_shipped_fail", "set_shipment_err=%v", err)
			return nil, err
		}
	}
	if err := s.repo.SetStatus(ctx, dealID, []string{"PAID"}, "SHIPPED", "shippedAt"); err != nil {
		dealTrace(dealID, sellerID, "mark_shipped_fail", "set_status_err=%v", err)
		return nil, &AppError{400, "Не удалось подтвердить отправку"}
	}
	s.writeDealLog(ctx, sellerID, dealID, "ship", "deal marked as shipped")
	dealTrace(dealID, sellerID, "mark_shipped_success", "status=SHIPPED")
	return s.GetDeal(ctx, sellerID, dealID)
}

func (s *DealService) ConfirmDelivery(ctx context.Context, buyerID, dealID int32) (map[string]any, error) {
	deal, err := s.getUserDeal(ctx, dealID, buyerID, "buyer")
	if err != nil {
		return nil, err
	}
	if deal.Status != "SHIPPED" {
		return nil, &AppError{400, "Получение можно подтвердить только после отправки"}
	}
	if deal.PaymentID != nil {
		pid := strings.TrimSpace(*deal.PaymentID)
		if pid != "" && !strings.HasPrefix(strings.ToLower(pid), "mock-") {
			if err := s.payment.ConfirmDealPayment(ctx, pid, deal.TotalAmount); err != nil {
				return nil, err
			}
		}
	}
	delay := s.cfg.DealPayoutDelayDays
	if delay < 0 {
		delay = 0
	}
	payoutAt := time.Now().AddDate(0, 0, delay)
	if err := s.repo.MarkDelivered(ctx, dealID, payoutAt); err != nil {
		return nil, &AppError{400, "Не удалось подтвердить получение"}
	}
	s.writeDealLog(ctx, buyerID, dealID, "deliver", "deal marked as delivered")
	return s.GetDeal(ctx, buyerID, dealID)
}

func (s *DealService) OpenDispute(ctx context.Context, userID, dealID int32, reason string) (map[string]any, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, &AppError{400, "Нужно указать причину спора"}
	}
	if _, err := s.getUserDeal(ctx, dealID, userID, "participant"); err != nil {
		return nil, err
	}
	if err := s.repo.OpenDispute(ctx, dealID, strings.TrimSpace(reason)); err != nil {
		return nil, &AppError{400, "Не удалось открыть спор"}
	}
	s.writeDealLog(ctx, userID, dealID, "dispute", "dispute opened")
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
	if deal.Status == "PAID" {
		if deal.PaymentID != nil {
			pid := strings.TrimSpace(*deal.PaymentID)
			if pid != "" && !strings.HasPrefix(strings.ToLower(pid), "mock-") {
				if err := s.payment.CancelDealPayment(ctx, pid); err != nil {
					return nil, err
				}
			}
		}
		if err := s.repo.SetStatus(ctx, dealID, []string{"PAID"}, "REFUNDED", "refundedAt"); err != nil {
			return nil, &AppError{400, "Не удалось отменить сделку"}
		}
		s.writeDealLog(ctx, userID, dealID, "cancel", "deal cancelled with refund")
		return s.GetDeal(ctx, userID, dealID)
	}
	if err := s.repo.SetStatus(ctx, dealID, []string{"CREATED"}, "CANCELLED", "cancelledAt"); err != nil {
		return nil, &AppError{400, "Не удалось отменить сделку"}
	}
	s.writeDealLog(ctx, userID, dealID, "cancel", "deal cancelled")
	return s.GetDeal(ctx, userID, dealID)
}

func (s *DealService) GetDeal(ctx context.Context, userID, dealID int32) (map[string]any, error) {
	deal, err := s.getUserDeal(ctx, dealID, userID, "participant")
	if err != nil {
		return nil, err
	}
	deal = s.refreshDealFromCDEK(ctx, deal)
	return s.formatDeal(*deal), nil
}

func (s *DealService) GetDealCDEKQR(ctx context.Context, userID, dealID int32) (map[string]any, error) {
	deal, err := s.getUserDeal(ctx, dealID, userID, "participant")
	if err != nil {
		return nil, err
	}
	deal = s.refreshDealFromCDEK(ctx, deal)
	if deal.CDEKOrderUUID == nil || strings.TrimSpace(*deal.CDEKOrderUUID) == "" {
		return nil, &AppError{404, "Для сделки еще не сохранен orderUuid CDEK"}
	}
	if s.cdek == nil {
		return nil, &AppError{400, "CDEK сервис не инициализирован"}
	}
	qrData, qrURL := s.cdek.QRByOrderUUID(ctx, *deal.CDEKOrderUUID)
	trackPending := deal.CDEKTrackNumber == nil || strings.TrimSpace(*deal.CDEKTrackNumber) == ""

	barcodePdfUrl := fmt.Sprintf("/deals/%d/cdek-barcode-pdf", deal.ID)
	waybillPdfUrl := fmt.Sprintf("/deals/%d/cdek-waybill-pdf", deal.ID)

	return map[string]any{
		"qrCodeData":    qrData,
		"qrCodeUrl":     qrURL,
		"trackNumber":   deal.CDEKTrackNumber,
		"trackingUrl":   buildCDEKTrackingURL(deal.CDEKTrackNumber),
		"orderUuid":     deal.CDEKOrderUUID,
		"trackPending":  trackPending,
		"barcodePdfUrl": barcodePdfUrl,
		"waybillPdfUrl": waybillPdfUrl,
	}, nil
}

func (s *DealService) GetDealCDEKPrintPDF(ctx context.Context, userID, dealID int32, kind string) ([]byte, error) {
	deal, err := s.getUserDeal(ctx, dealID, userID, "participant")
	if err != nil {
		return nil, err
	}
	deal = s.refreshDealFromCDEK(ctx, deal)
	if deal.CDEKOrderUUID == nil || strings.TrimSpace(*deal.CDEKOrderUUID) == "" {
		return nil, &AppError{404, "Для сделки еще не сохранен orderUuid CDEK"}
	}
	if s.cdek == nil {
		return nil, &AppError{400, "CDEK сервис не инициализирован"}
	}
	return s.cdek.GetPDFForm(ctx, kind, *deal.CDEKOrderUUID)
}

func (s *DealService) MyPurchases(ctx context.Context, buyerID int32) ([]map[string]any, error) {
	deals, err := s.repo.ListByBuyer(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	s.ensureCdekOrdersInList(ctx, deals, 12)
	deals, _ = s.refreshDealsFromCDEK(ctx, deals, 5)
	return s.formatDeals(deals), nil
}

func (s *DealService) MySales(ctx context.Context, sellerID int32) ([]map[string]any, error) {
	deals, err := s.repo.ListBySeller(ctx, sellerID)
	if err != nil {
		return nil, err
	}
	s.ensureCdekOrdersInList(ctx, deals, 12)
	deals, _ = s.refreshDealsFromCDEK(ctx, deals, 5)
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
	s.ensureCdekOrdersInList(ctx, purchases, 12)
	s.ensureCdekOrdersInList(ctx, sales, 12)
	purchases, usedBudget := s.refreshDealsFromCDEK(ctx, purchases, 5)
	remainingBudget := 5 - usedBudget
	if remainingBudget < 0 {
		remainingBudget = 0
	}
	sales, _ = s.refreshDealsFromCDEK(ctx, sales, remainingBudget)
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
		return nil, &AppError{404, "Сделка не найдена"}
	}
	if err != nil {
		return nil, err
	}
	switch role {
	case "buyer":
		if deal.BuyerID != userID {
			return nil, &AppError{403, "Это не ваша покупка"}
		}
	case "seller":
		if deal.SellerID != userID {
			return nil, &AppError{403, "Это не ваша продажа"}
		}
	default:
		if deal.BuyerID != userID && deal.SellerID != userID {
			return nil, &AppError{403, "Вы не участник сделки"}
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
	cdek := formatDealCDEK(deal)
	if deal.Status == "CREATED" && s.cfg.DealAllowMockPayment && !s.payment.tinkoffConfigured() {
		cdek["registrationHint"] = "Демо без Тинькофф: покупатель нажимает «Оплатить» — сделка сразу станет оплаченной, затем подтянется заказ CDEK."
		cdek["mockPaymentAvailable"] = true
	}
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
		"cdek":          cdek,
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
	regHint := buildCdekRegistrationHint(deal)
	sellerNote := buildCdekSellerHandoffNote(deal)
	pkg := map[string]any{}
	if deal.CDEKPackageWeight != nil {
		pkg["weight"] = *deal.CDEKPackageWeight
	}
	if deal.CDEKPackageLength != nil {
		pkg["length"] = *deal.CDEKPackageLength
	}
	if deal.CDEKPackageWidth != nil {
		pkg["width"] = *deal.CDEKPackageWidth
	}
	if deal.CDEKPackageHeight != nil {
		pkg["height"] = *deal.CDEKPackageHeight
	}
	out := map[string]any{
		"tariffCode":        deal.CDEKTariffCode,
		"tariffName":        deal.CDEKTariffName,
		"fromCityCode":      deal.CDEKFromCity,
		"toCityCode":        deal.CDEKToCity,
		"fromPvzCode":       deal.CDEKFromPVZ,
		"toPvzCode":         deal.CDEKToPVZ,
		"toAddress":         deal.CDEKToAddress,
		"fromAddress":       deal.CDEKFromAddress,
		"recipientMode":     deal.CDEKRecipientMode,
		"sellerHandoff":     deal.CDEKSellerHandoff,
		"cdekStatus":        deal.CDEKStatus,
		"orderUuid":         deal.CDEKOrderUUID,
		"trackNumber":       deal.CDEKTrackNumber,
		"trackingUrl":       buildCDEKTrackingURL(deal.CDEKTrackNumber),
		"trackPending":      trackPending,
		"registrationHint":  regHint,
		"sellerHandoffHint": sellerNote,
		"deliveryStages":    buildCdekDeliveryStages(deal),
	}
	if len(pkg) > 0 {
		out["package"] = pkg
	}
	if regHint == "" {
		delete(out, "registrationHint")
	}
	return out
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

func (s *DealService) refreshDealFromCDEK(ctx context.Context, deal *repository.DealRow) *repository.DealRow {
	if deal == nil || s.cdek == nil || deal.CDEKOrderUUID == nil {
		return deal
	}
	orderUUID := strings.TrimSpace(*deal.CDEKOrderUUID)
	if orderUUID == "" {
		return deal
	}
	dealTrace(deal.ID, deal.SellerID, "cdek_refresh_enter", "orderUuid=%s currentTrack=%s currentStatus=%s", orderUUID, trimOrDash(deal.CDEKTrackNumber), trimOrDash(deal.CDEKStatus))

	details := s.cdek.OrderDetailsByOrderUUID(ctx, orderUUID)
	if details == nil {
		details = &CDEKOrderDetails{}
	}
	if details.Track == nil || strings.TrimSpace(*details.Track) == "" {
		if byNumber, err := s.cdek.LookupOrderByClientNumber(ctx, fmt.Sprintf("med-vito-deal-%d", deal.ID)); err == nil && byNumber != nil {
			if byNumber.Track != nil && strings.TrimSpace(*byNumber.Track) != "" {
				details.Track = byNumber.Track
			}
			if byNumber.OrderUUID != nil {
				lookupUUID := strings.TrimSpace(*byNumber.OrderUUID)
				if lookupUUID != "" && lookupUUID != orderUUID {
					_ = s.repo.SetCDEKShipment(ctx, deal.ID, &lookupUUID, nil)
					orderUUID = lookupUUID
				}
			}
		}
	}

	currentTrack := ""
	if deal.CDEKTrackNumber != nil {
		currentTrack = strings.TrimSpace(*deal.CDEKTrackNumber)
	}
	newTrack := ""
	if details.Track != nil {
		newTrack = strings.TrimSpace(*details.Track)
	}
	currentStatus := ""
	if deal.CDEKStatus != nil {
		currentStatus = strings.TrimSpace(*deal.CDEKStatus)
	}
	statusChanged := details.StatusCode != "" && details.StatusCode != currentStatus
	trackChanged := newTrack != "" && newTrack != currentTrack

	if !trackChanged && !statusChanged {
		dealTrace(deal.ID, deal.SellerID, "cdek_refresh_no_change", "orderUuid=%s", orderUUID)
		return deal
	}

	if trackChanged {
		dealTrace(deal.ID, deal.SellerID, "cdek_refresh_track_update", "old=%s new=%s", currentTrack, newTrack)
		if err := s.repo.SetCDEKShipment(ctx, deal.ID, nil, &newTrack); err != nil {
			dealTrace(deal.ID, deal.SellerID, "cdek_refresh_fail", "set_track_err=%v", err)
			return deal
		}
		s.notifyCdekTrackToBuyer(ctx, *deal, newTrack)
	}
	if statusChanged {
		dealTrace(deal.ID, deal.SellerID, "cdek_refresh_status_update", "old=%s new=%s", currentStatus, details.StatusCode)
		_ = s.repo.SetCDEKStatus(ctx, deal.ID, details.StatusCode)
	}
	updated, err := s.repo.FindByID(ctx, deal.ID)
	if err != nil {
		return deal
	}
	return updated
}

func (s *DealService) notifyCdekTrackToBuyer(ctx context.Context, deal repository.DealRow, track string) {
	if s.chat == nil || track == "" {
		return
	}
	chat, err := s.chat.FindChatByProductParticipants(ctx, deal.ProductID, deal.BuyerID, deal.SellerID)
	if err != nil || chat == nil {
		return
	}
	url := "https://www.cdek.ru/ru/tracking?order_id=" + track
	content := fmt.Sprintf(
		"Посылка по сделке #%d принята СДЭК. Трек: %s. Отслеживание: %s. Когда груз прибудет, СДЭК пришлёт SMS с кодом для получения.",
		deal.ID, track, url,
	)
	_, _, _, _, _, _ = s.chat.InsertChatMessage(ctx, chat.ID, deal.SellerID, content)
}

func (s *DealService) refreshDealsFromCDEK(ctx context.Context, deals []repository.DealRow, maxLive int) ([]repository.DealRow, int) {
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
		updated := s.refreshDealFromCDEK(ctx, &deals[i])
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

func cdekNeedsRecipientAddress(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "to_location.address") ||
		strings.Contains(msg, "recipient address and recipient delivery point")
}

func (s *DealService) AdminListDeals(ctx context.Context) ([]map[string]any, error) {
	deals, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return s.formatDeals(deals), nil
}

func (s *DealService) AdminGetDeal(ctx context.Context, dealID int32) (map[string]any, error) {
	deal, err := s.repo.FindByID(ctx, dealID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "Сделка не найдена"}
	}
	if err != nil {
		return nil, err
	}
	return s.formatDeal(*deal), nil
}

func (s *DealService) AdminSetStatus(ctx context.Context, actorUserID, dealID int32, status string) (map[string]any, error) {
	status = strings.ToUpper(strings.TrimSpace(status))
	switch status {
	case "CREATED", "PAID", "SHIPPED", "DELIVERED", "COMPLETED", "CANCELLED", "REFUNDED", "DISPUTE":
	default:
		return nil, &AppError{400, "Недопустимый статус сделки"}
	}
	if err := s.repo.AdminSetStatus(ctx, dealID, status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Сделка не найдена"}
		}
		return nil, err
	}
	s.writeDealLog(ctx, actorUserID, dealID, "admin_status", "admin set status to "+status)
	return s.AdminGetDeal(ctx, dealID)
}

func (s *DealService) AdminDealLogs(ctx context.Context, dealID int32) ([]map[string]any, error) {
	if s.logs == nil {
		return []map[string]any{}, nil
	}
	rows, err := s.logs.FindByDealID(ctx, dealID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		out = append(out, map[string]any{
			"id":        row.ID,
			"userId":    row.UserID,
			"action":    row.Action,
			"userName":  row.UserName,
			"userEmail": row.UserEmail,
		})
	}
	return out, nil
}

func (s *DealService) writeDealLog(ctx context.Context, userID, dealID int32, event, details string) {
	if s.logs == nil || userID <= 0 || dealID <= 0 {
		return
	}
	action := fmt.Sprintf("deal_id=%d event=%s %s", dealID, strings.TrimSpace(event), strings.TrimSpace(details))
	if err := s.logs.Insert(ctx, userID, action); err != nil {
		log.Printf("deal log write failed deal=%d user=%d: %v", dealID, userID, err)
	}
}
