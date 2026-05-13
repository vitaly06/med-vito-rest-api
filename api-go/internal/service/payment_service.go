package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/pkg/tinkoff"
	"med-vito/api-go/internal/repository"
)

const tinkoffAPIBase = "https://securepay.tinkoff.ru/v2"

type PaymentService struct {
	cfg    config.Config
	repo   *repository.PaymentPG
	deals  *repository.DealPG
	client *http.Client
}

func NewPaymentService(cfg config.Config, repo *repository.PaymentPG, deals *repository.DealPG) *PaymentService {
	return &PaymentService{
		cfg:    cfg,
		repo:   repo,
		deals:  deals,
		client: &http.Client{Timeout: 45 * time.Second},
	}
}

func (s *PaymentService) tinkoffConfigured() bool {
	return s.cfg.TinkoffTerminalKey != "" && s.cfg.TinkoffSecretKey != ""
}

func parseTinkoffPaymentID(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", errors.New("empty PaymentId")
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s, nil
	}
	var n json.Number
	if err := json.Unmarshal(raw, &n); err == nil {
		return n.String(), nil
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return strconv.FormatInt(int64(f), 10), nil
	}
	return "", errors.New("PaymentId")
}

func tinkoffBoolSuccess(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return tinkoff.ScalarString(v) == "true"
}

// tokenParams вЂ” С‚РѕР»СЊРєРѕ РґР»СЏ РїРѕРґРїРёСЃРё (РІСЃРµ Р·РЅР°С‡РµРЅРёСЏ СЃС‚СЂРѕРєРё); jsonBody вЂ” С„Р°РєС‚РёС‡РµСЃРєРѕРµ С‚РµР»Рѕ (Amount С‡РёСЃР»РѕРј, РєР°Рє Сѓ Nest).
func (s *PaymentService) tinkoffRequest(ctx context.Context, path string, tokenParams map[string]string, jsonBody map[string]any) (map[string]interface{}, error) {
	token := tinkoff.Token(s.cfg.TinkoffSecretKey, tokenParams)
	payload := make(map[string]any, len(jsonBody)+1)
	for k, v := range jsonBody {
		payload[k] = v
	}
	payload["Token"] = token

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tinkoffAPIBase+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Рў-Р‘Р°РЅРє HTTP %d: %s", resp.StatusCode, truncateTinkoffBody(raw))
	}
	var data map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// CreatePayment вЂ” Init + Р·Р°РїРёСЃСЊ РІ Р‘Р” (РєР°Рє Nest createPayment).
func (s *PaymentService) CreatePayment(ctx context.Context, userID int32, amount float64, description string) (map[string]any, error) {
	if !s.tinkoffConfigured() {
		return nil, &AppError{400, "РџР»Р°С‚С‘Р¶РЅР°СЏ СЃРёСЃС‚РµРјР° РЅРµ РЅР°СЃС‚СЂРѕРµРЅР° (TINKOFF_TERMINAL_KEY / TINKOFF_SECRET_KEY)"}
	}
	if amount < 1 {
		return nil, &AppError{400, "РњРёРЅРёРјР°Р»СЊРЅР°СЏ СЃСѓРјРјР° РїРѕРїРѕР»РЅРµРЅРёСЏ 1 СЂСѓР±Р»СЊ"}
	}
	if description == "" {
		description = "РџРѕРїРѕР»РЅРµРЅРёРµ Р±Р°Р»Р°РЅСЃР°"
	}
	kopeks := int64(math.Round(amount * 100))
	orderID := fmt.Sprintf("%d-%d", userID, time.Now().UnixMilli())

	params := map[string]string{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"Amount":      strconv.FormatInt(kopeks, 10),
		"OrderId":     orderID,
		"Description": description,
		"PayType":     "T",
	}

	jsonBody := map[string]any{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"Amount":      kopeks,
		"OrderId":     orderID,
		"Description": description,
		"PayType":     "T",
	}
	data, err := s.tinkoffRequest(ctx, "/Init", params, jsonBody)
	if err != nil {
		return nil, &AppError{400, fmt.Sprintf("РћС€РёР±РєР° РїСЂРё СЃРѕР·РґР°РЅРёРё РїР»Р°С‚РµР¶Р°: %v", err)}
	}
	if !tinkoffBoolSuccess(data["Success"]) {
		msg := tinkoff.ScalarString(data["Message"])
		if msg == "" {
			msg = tinkoff.ScalarString(data["Details"])
		}
		return nil, &AppError{400, fmt.Sprintf("РћС€РёР±РєР° СЃРѕР·РґР°РЅРёСЏ РїР»Р°С‚РµР¶Р°: %s", msg)}
	}
	pidRaw, err := json.Marshal(data["PaymentId"])
	if err != nil {
		return nil, &AppError{400, "РћС€РёР±РєР° СЃРѕР·РґР°РЅРёСЏ РїР»Р°С‚РµР¶Р°: PaymentId"}
	}
	paymentID, err := parseTinkoffPaymentID(pidRaw)
	if err != nil {
		return nil, &AppError{400, "РћС€РёР±РєР° СЃРѕР·РґР°РЅРёСЏ РїР»Р°С‚РµР¶Р°: PaymentId"}
	}
	paymentURL := tinkoff.ScalarString(data["PaymentURL"])
	var urlPtr *string
	if paymentURL != "" {
		urlPtr = &paymentURL
	}
	if err := s.repo.Insert(ctx, orderID, paymentID, userID, amount, "PENDING", urlPtr); err != nil {
		return nil, &AppError{400, fmt.Sprintf("РћС€РёР±РєР° РїСЂРё СЃРѕР·РґР°РЅРёРё РїР»Р°С‚РµР¶Р°: %v", err)}
	}
	return map[string]any{
		"paymentId":  paymentID,
		"paymentUrl": paymentURL,
		"orderId":    orderID,
		"amount":     amount,
	}, nil
}

// CreateDealPayment — Init для безопасной сделки (без записи в Payment).
func (s *PaymentService) CreateDealPayment(ctx context.Context, userID, dealID int32, amount float64, description string) (string, string, string, error) {
	if !s.tinkoffConfigured() {
		return "", "", "", &AppError{400, "РџР»Р°С‚С‘Р¶РЅР°СЏ СЃРёСЃС‚РµРјР° РЅРµ РЅР°СЃС‚СЂРѕРµРЅР° (TINKOFF_TERMINAL_KEY / TINKOFF_SECRET_KEY)"}
	}
	if amount < 1 {
		return "", "", "", &AppError{400, "РњРёРЅРёРјР°Р»СЊРЅР°СЏ СЃСѓРјРјР° РѕРїР»Р°С‚С‹ 1 СЂСѓР±Р»СЊ"}
	}
	if description == "" {
		description = "Р‘РµР·РѕРїР°СЃРЅР°СЏ СЃРґРµР»РєР°"
	}
	kopeks := int64(math.Round(amount * 100))
	orderID := fmt.Sprintf("deal-%d-%d", dealID, time.Now().UnixMilli())

	params := map[string]string{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"Amount":      strconv.FormatInt(kopeks, 10),
		"OrderId":     orderID,
		"Description": description,
		"PayType":     "T",
	}

	jsonBody := map[string]any{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"Amount":      kopeks,
		"OrderId":     orderID,
		"Description": description,
		"PayType":     "T",
	}
	data, err := s.tinkoffRequest(ctx, "/Init", params, jsonBody)
	if err != nil {
		return "", "", "", &AppError{400, fmt.Sprintf("Payment create error: %v", err)}
	}
	if !tinkoffBoolSuccess(data["Success"]) {
		msg := tinkoff.ScalarString(data["Message"])
		if msg == "" {
			msg = tinkoff.ScalarString(data["Details"])
		}
		if strings.Contains(msg, "Неверные параметры") {
			msg = msg + ". Проверь TerminalKey/SecretKey и включена ли двухстадийная оплата (hold/confirm) в терминале Тинькофф."
		}
		return "", "", "", &AppError{400, fmt.Sprintf("Payment create failed: %s", msg)}
	}
	pidRaw, err := json.Marshal(data["PaymentId"])
	if err != nil {
		return "", "", "", &AppError{400, "РћС€РёР±РєР° СЃРѕР·РґР°РЅРёСЏ РїР»Р°С‚РµР¶Р°: PaymentId"}
	}
	paymentID, err := parseTinkoffPaymentID(pidRaw)
	if err != nil {
		return "", "", "", &AppError{400, "РћС€РёР±РєР° СЃРѕР·РґР°РЅРёСЏ РїР»Р°С‚РµР¶Р°: PaymentId"}
	}
	paymentURL := tinkoff.ScalarString(data["PaymentURL"])
	return paymentID, paymentURL, orderID, nil
}

func (s *PaymentService) ConfirmDealPayment(ctx context.Context, paymentID string, amountRub int32) error {
	if !s.tinkoffConfigured() {
		return &AppError{400, "РџР»Р°С‚РµР¶РЅР°СЏ СЃРёСЃС‚РµРјР° РЅРµ РЅР°СЃС‚СЂРѕРµРЅР°"}
	}
	pid := strings.TrimSpace(paymentID)
	if pid == "" {
		return &AppError{400, "PaymentId РґР»СЏ confirm РїСѓСЃС‚РѕР№"}
	}
	kopeks := int64(amountRub) * 100
	params := map[string]string{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"PaymentId":   pid,
		"Amount":      strconv.FormatInt(kopeks, 10),
	}
	jsonBody := map[string]any{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"PaymentId":   pid,
		"Amount":      kopeks,
	}
	data, err := s.tinkoffRequest(ctx, "/Confirm", params, jsonBody)
	if err != nil {
		return &AppError{400, fmt.Sprintf("РћС€РёР±РєР° confirm РїР»Р°С‚РµР¶Р°: %v", err)}
	}
	if !tinkoffBoolSuccess(data["Success"]) {
		msg := tinkoff.ScalarString(data["Message"])
		if msg == "" {
			msg = tinkoff.ScalarString(data["Details"])
		}
		return &AppError{400, fmt.Sprintf("РћС€РёР±РєР° confirm РїР»Р°С‚РµР¶Р°: %s", msg)}
	}
	return nil
}

func (s *PaymentService) CancelDealPayment(ctx context.Context, paymentID string) error {
	if !s.tinkoffConfigured() {
		return &AppError{400, "РџР»Р°С‚РµР¶РЅР°СЏ СЃРёСЃС‚РµРјР° РЅРµ РЅР°СЃС‚СЂРѕРµРЅР°"}
	}
	pid := strings.TrimSpace(paymentID)
	if pid == "" {
		return &AppError{400, "PaymentId РґР»СЏ cancel РїСѓСЃС‚РѕР№"}
	}
	params := map[string]string{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"PaymentId":   pid,
	}
	jsonBody := map[string]any{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"PaymentId":   pid,
	}
	data, err := s.tinkoffRequest(ctx, "/Cancel", params, jsonBody)
	if err != nil {
		return &AppError{400, fmt.Sprintf("РћС€РёР±РєР° РѕС‚РјРµРЅС‹ РїР»Р°С‚РµР¶Р°: %v", err)}
	}
	if !tinkoffBoolSuccess(data["Success"]) {
		msg := tinkoff.ScalarString(data["Message"])
		if msg == "" {
			msg = tinkoff.ScalarString(data["Details"])
		}
		return &AppError{400, fmt.Sprintf("РћС€РёР±РєР° РѕС‚РјРµРЅС‹ РїР»Р°С‚РµР¶Р°: %s", msg)}
	}
	return nil
}

// HandleNotification вЂ” webhook Рў-Р‘Р°РЅРє: РїСЂРѕРІРµСЂРєР° Token, РёРґРµРјРїРѕС‚РµРЅС‚РЅРѕРµ РЅР°С‡РёСЃР»РµРЅРёРµ.
func (s *PaymentService) HandleNotification(ctx context.Context, rawJSON []byte) (map[string]any, error) {
	if !s.tinkoffConfigured() {
		return nil, &AppError{400, "РџР»Р°С‚С‘Р¶РЅР°СЏ СЃРёСЃС‚РµРјР° РЅРµ РЅР°СЃС‚СЂРѕРµРЅР°"}
	}
	dec := json.NewDecoder(bytes.NewReader(rawJSON))
	dec.UseNumber()
	var root map[string]interface{}
	if err := dec.Decode(&root); err != nil {
		return nil, &AppError{400, "РќРµРєРѕСЂСЂРµРєС‚РЅРѕРµ С‚РµР»Рѕ СѓРІРµРґРѕРјР»РµРЅРёСЏ"}
	}
	got, _ := root["Token"].(string)
	if got == "" {
		return nil, &AppError{400, "РќРµРІРµСЂРЅР°СЏ РїРѕРґРїРёСЃСЊ СѓРІРµРґРѕРјР»РµРЅРёСЏ"}
	}
	want := tinkoff.Token(s.cfg.TinkoffSecretKey, tinkoff.ParamsFromNotificationMap(root))
	if got != want {
		return nil, &AppError{400, "РќРµРІРµСЂРЅР°СЏ РїРѕРґРїРёСЃСЊ СѓРІРµРґРѕРјР»РµРЅРёСЏ"}
	}
	paymentID := tinkoff.ScalarString(root["PaymentId"])
	if paymentID == "" {
		return nil, &AppError{400, "РџР»Р°С‚РµР¶ РЅРµ РЅР°Р№РґРµРЅ"}
	}
	status := tinkoff.ScalarString(root["Status"])

	switch status {
	case "AUTHORIZED":
		if s.deals != nil {
			marked, errDeal := s.deals.TryMarkPaidByPaymentID(ctx, paymentID)
			if errDeal != nil {
				return nil, errDeal
			}
			if marked {
				return map[string]any{"success": true, "message": "РЎРґРµР»РєР° РѕС‚РјРµС‡РµРЅР° РѕРїР»Р°С‡РµРЅРЅРѕР№ (hold)"}, nil
			}
		}
		return map[string]any{"success": true}, nil
	case "CONFIRMED":
		already, err := s.repo.TryCompleteTopUp(ctx, paymentID)
		if errors.Is(err, repository.ErrNotFound) {
			if s.deals != nil {
				marked, errDeal := s.deals.TryMarkPaidByPaymentID(ctx, paymentID)
				if errDeal != nil {
					return nil, errDeal
				}
				if marked {
					return map[string]any{"success": true, "message": "Сделка отмечена оплаченной"}, nil
				}
			}
			return nil, &AppError{400, "РџР»Р°С‚РµР¶ РЅРµ РЅР°Р№РґРµРЅ"}
		}
		if err != nil {
			return nil, err
		}
		if already {
			return map[string]any{"success": true, "message": "РџР»Р°С‚РµР¶ СѓР¶Рµ РѕР±СЂР°Р±РѕС‚Р°РЅ"}, nil
		}
		return map[string]any{"success": true, "message": "Р‘Р°Р»Р°РЅСЃ СѓСЃРїРµС€РЅРѕ РїРѕРїРѕР»РЅРµРЅ"}, nil
	case "REJECTED", "CANCELED":
		already, err := s.repo.TryFailIfPending(ctx, paymentID)
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{400, "РџР»Р°С‚РµР¶ РЅРµ РЅР°Р№РґРµРЅ"}
		}
		if err != nil {
			return nil, err
		}
		if already {
			return map[string]any{"success": true}, nil
		}
		return map[string]any{"success": true}, nil
	default:
		return map[string]any{"success": true}, nil
	}
}

func (s *PaymentService) GetUserPayments(ctx context.Context, userID int32) ([]repository.PaymentRow, error) {
	return s.repo.ListByUserID(ctx, userID, 50)
}

// CheckPaymentStatus вЂ” GetState РІ Рў-Р‘Р°РЅРє.
func (s *PaymentService) CheckPaymentStatus(ctx context.Context, paymentID string) (map[string]any, error) {
	if !s.tinkoffConfigured() {
		return nil, &AppError{400, "РџР»Р°С‚С‘Р¶РЅР°СЏ СЃРёСЃС‚РµРјР° РЅРµ РЅР°СЃС‚СЂРѕРµРЅР°"}
	}
	params := map[string]string{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"PaymentId":   paymentID,
	}
	jsonBody := map[string]any{
		"TerminalKey": s.cfg.TinkoffTerminalKey,
		"PaymentId":   paymentID,
	}
	data, err := s.tinkoffRequest(ctx, "/GetState", params, jsonBody)
	if err != nil {
		return nil, &AppError{400, fmt.Sprintf("РћС€РёР±РєР° РїСЂРё РїСЂРѕРІРµСЂРєРµ СЃС‚Р°С‚СѓСЃР°: %v", err)}
	}
	if !tinkoffBoolSuccess(data["Success"]) {
		return nil, &AppError{400, "РћС€РёР±РєР° РїСЂРѕРІРµСЂРєРё СЃС‚Р°С‚СѓСЃР° РїР»Р°С‚РµР¶Р°"}
	}
	var rub float64
	switch x := data["Amount"].(type) {
	case json.Number:
		k, _ := x.Float64()
		rub = k / 100
	case float64:
		rub = x / 100
	default:
		rub = 0
	}
	return map[string]any{
		"status":  tinkoff.ScalarString(data["Status"]),
		"amount":  rub,
		"orderId": tinkoff.ScalarString(data["OrderId"]),
	}, nil
}

const maxTinkoffErrBody = 400

func truncateTinkoffBody(b []byte) string {
	s := string(b)
	if len(s) <= maxTinkoffErrBody {
		return s
	}
	return s[:maxTinkoffErrBody] + "вЂ¦"
}
