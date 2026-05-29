package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"med-vito/api-go/internal/config"
)

type CDEKService struct {
	cfg          config.Config
	client       *http.Client
	mu           sync.Mutex
	accessToken  string
	tokenExpires time.Time
}

func NewCDEKService(cfg config.Config) *CDEKService {
	return &CDEKService{cfg: cfg, client: &http.Client{Timeout: 30 * time.Second}}
}

func (s *CDEKService) configured() bool {
	return s.cfg.CDEKClientID != "" && s.cfg.CDEKClientSecret != ""
}

func (s *CDEKService) Cities(ctx context.Context, city string, limit int) ([]map[string]any, error) {
	city = strings.TrimSpace(city)
	if city == "" {
		return nil, &AppError{400, "РќСѓР¶РЅРѕ СѓРєР°Р·Р°С‚СЊ city"}
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	out, err := s.fetchCities(ctx, city, limit)
	if err == nil {
		return out, nil
	}

	// CDEK sandbox can return 410 for Cyrillic city queries; retry with transliteration.
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Status == http.StatusGone && containsCyrillic(city) {
		latinCity := transliterateRU(city)
		if latinCity != "" && latinCity != city {
			return s.fetchCities(ctx, latinCity, limit)
		}
	}

	return nil, err
}

func (s *CDEKService) fetchCities(ctx context.Context, city string, limit int) ([]map[string]any, error) {
	q := url.Values{}
	q.Set("city", city)
	q.Set("size", strconv.Itoa(limit))
	var out []map[string]any
	if err := s.getJSON(ctx, "/location/cities?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *CDEKService) DeliveryPoints(ctx context.Context, cityCode int) ([]map[string]any, error) {
	if cityCode <= 0 {
		return nil, &AppError{400, "РќСѓР¶РµРЅ cityCode"}
	}
	q := url.Values{}
	q.Set("city_code", strconv.Itoa(cityCode))
	q.Set("type", "PVZ")
	var out []map[string]any
	if err := s.getJSON(ctx, "/deliverypoints?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

type CDEKCalculateRequest struct {
	TariffCode   int `json:"tariffCode"`
	FromCityCode int `json:"fromCityCode"`
	ToCityCode   int `json:"toCityCode"`
	Weight       int `json:"weight"`
	Length       int `json:"length"`
	Width        int `json:"width"`
	Height       int `json:"height"`
}

type CDEKTariffsRequest struct {
	FromCityCode int `json:"fromCityCode"`
	ToCityCode   int `json:"toCityCode"`
	Weight       int `json:"weight"`
	Length       int `json:"length"`
	Width        int `json:"width"`
	Height       int `json:"height"`
}

const cdekAllowedTariffCode = 136

func (s *CDEKService) Calculate(ctx context.Context, req CDEKCalculateRequest) (map[string]any, error) {
	if req.TariffCode <= 0 || req.FromCityCode <= 0 || req.ToCityCode <= 0 {
		return nil, &AppError{400, "РќСѓР¶РЅС‹ tariffCode, fromCityCode Рё toCityCode"}
	}
	if req.TariffCode != cdekAllowedTariffCode {
		return nil, &AppError{400, "Доступен только тариф 136 (склад-склад)"}
	}
	if req.Weight <= 0 {
		req.Weight = 1000
	}
	if req.Length <= 0 {
		req.Length = 20
	}
	if req.Width <= 0 {
		req.Width = 20
	}
	if req.Height <= 0 {
		req.Height = 20
	}
	payload := map[string]any{
		"tariff_code": req.TariffCode,
		"from_location": map[string]any{
			"code": req.FromCityCode,
		},
		"to_location": map[string]any{
			"code": req.ToCityCode,
		},
		"packages": []map[string]any{
			{
				"weight": req.Weight,
				"length": req.Length,
				"width":  req.Width,
				"height": req.Height,
			},
		},
	}
	var out map[string]any
	if err := s.postJSON(ctx, "/calculator/tariff", payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *CDEKService) Tariffs(ctx context.Context, req CDEKTariffsRequest) ([]map[string]any, error) {
	if req.FromCityCode <= 0 || req.ToCityCode <= 0 {
		return nil, &AppError{400, "Нужны fromCityCode и toCityCode"}
	}
	if req.Weight <= 0 {
		req.Weight = 1000
	}
	if req.Length <= 0 {
		req.Length = 20
	}
	if req.Width <= 0 {
		req.Width = 20
	}
	if req.Height <= 0 {
		req.Height = 20
	}
	payload := map[string]any{
		"from_location": map[string]any{"code": req.FromCityCode},
		"to_location":   map[string]any{"code": req.ToCityCode},
		"packages": []map[string]any{
			{
				"weight": req.Weight,
				"length": req.Length,
				"width":  req.Width,
				"height": req.Height,
			},
		},
	}
	var raw any
	if err := s.postJSON(ctx, "/calculator/tarifflist", payload, &raw); err != nil {
		return nil, err
	}
	normalized := normalizeTariffList(raw)
	filtered := make([]map[string]any, 0, len(normalized))
	for _, item := range normalized {
		if code, ok := item["tariffCode"].(int); ok && code == cdekAllowedTariffCode {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

// CDEKCreateOrderInput — минимальный набор для POST /v2/orders (безопасная сделка).
type CDEKCreateOrderInput struct {
	TariffCode      int
	FromCityCode    int
	ToCityCode      int
	FromPVZ         *string
	ToPVZ           *string
	ToAddress       *string
	ClientNumber    string
	Comment         string
	SenderName      string
	SenderPhone     string
	RecipientName   string
	RecipientPhone  string
	PackageName     string
	WareKey         string
	DeclaredCostRub float64
	WeightGrams     int
	LengthCm        int
	WidthCm         int
	HeightCm        int
	FromAddress     *string
}

// CDEKOrderDetails — трек и последний статус из GET /orders/{uuid}.
type CDEKOrderDetails struct {
	Track      *string
	StatusCode string
	StatusName string
}

// CDEKCreateOrderResult — uuid заказа в CDEK и трек, если API уже вернул.
type CDEKCreateOrderResult struct {
	OrderUUID *string
	Track     *string
}

// CreateOrder — регистрация отправки в CDEK; number = ClientNumber для идемпотентности.
// Для сценария ПВЗ→ПВЗ отправляем только shipment_point/delivery_point без адресов.
func (s *CDEKService) CreateOrder(ctx context.Context, in CDEKCreateOrderInput) (*CDEKCreateOrderResult, error) {
	if !s.configured() {
		return nil, &AppError{400, "CDEK не настроен (CDEK_CLIENT_ID / CDEK_CLIENT_SECRET)"}
	}
	if in.TariffCode <= 0 || in.FromCityCode <= 0 || in.ToCityCode <= 0 {
		return nil, &AppError{400, "Нужны tariffCode, fromCityCode и toCityCode для заказа CDEK"}
	}
	clientNumber := strings.TrimSpace(in.ClientNumber)
	if clientNumber == "" {
		return nil, &AppError{400, "Пустой client number для CDEK"}
	}

	includePvz := cdekInputHasPvz(in)
	log.Printf(
		"cdek create-order start clientNumber=%s tariff=%d fromCity=%d toCity=%d hasToPvz=%t hasFromPvz=%t hasToAddress=%t",
		strings.TrimSpace(in.ClientNumber),
		in.TariffCode,
		in.FromCityCode,
		in.ToCityCode,
		in.ToPVZ != nil && strings.TrimSpace(*in.ToPVZ) != "",
		in.FromPVZ != nil && strings.TrimSpace(*in.FromPVZ) != "",
		in.ToAddress != nil && strings.TrimSpace(*in.ToAddress) != "",
	)
	res, err := s.postCdekOrder(ctx, in, includePvz)
	if err == nil && res != nil && res.OrderUUID != nil {
		log.Printf("cdek create-order success clientNumber=%s orderUUID=%s", strings.TrimSpace(in.ClientNumber), strings.TrimSpace(*res.OrderUUID))
	}
	return res, err
}

func cdekInputHasPvz(in CDEKCreateOrderInput) bool {
	if in.ToPVZ != nil && strings.TrimSpace(*in.ToPVZ) != "" {
		return true
	}
	if in.FromPVZ != nil && strings.TrimSpace(*in.FromPVZ) != "" {
		return true
	}
	return false
}

func buildCdekOrderPayload(in CDEKCreateOrderInput, includePvz bool) map[string]any {
	w := in.WeightGrams
	if w <= 0 {
		w = 1000
	}
	length := in.LengthCm
	if length <= 0 {
		length = 20
	}
	width := in.WidthCm
	if width <= 0 {
		width = 20
	}
	height := in.HeightCm
	if height <= 0 {
		height = 20
	}
	cost := in.DeclaredCostRub
	if cost < 1 {
		cost = 1
	}
	pkgName := strings.TrimSpace(in.PackageName)
	if pkgName == "" {
		pkgName = "Товар"
	}
	wareKey := strings.TrimSpace(in.WareKey)
	if wareKey == "" {
		wareKey = "item-1"
	}
	senderName := strings.TrimSpace(in.SenderName)
	if senderName == "" {
		senderName = "Отправитель"
	}
	recipientName := strings.TrimSpace(in.RecipientName)
	if recipientName == "" {
		recipientName = "Получатель"
	}

	payload := map[string]any{
		"type":          1,
		"number":        strings.TrimSpace(in.ClientNumber),
		"tariff_code":   in.TariffCode,
		"comment":       strings.TrimSpace(in.Comment),
		"recipient":     map[string]any{"name": recipientName, "phones": []map[string]any{{"number": in.RecipientPhone}}},
		"sender":        map[string]any{"name": senderName, "phones": []map[string]any{{"number": in.SenderPhone}}},
		"from_location": map[string]any{"code": in.FromCityCode},
		"to_location":   map[string]any{"code": in.ToCityCode},
		"packages": []map[string]any{{
			"number": "1",
			"weight": w,
			"length": length, "width": width, "height": height,
			"items": []map[string]any{{
				"name":     pkgName,
				"ware_key": wareKey,
				"payment": map[string]any{
					"value": 0,
				},
				"cost":   cost,
				"weight": w,
				"amount": 1,
			}},
		}},
	}
	// В режиме ПВЗ CDEK не принимает одновременно address и *_point.
	// Поэтому адреса добавляем только если ПВЗ не используется.
	if !includePvz {
		if in.ToAddress != nil {
			if addr := strings.TrimSpace(*in.ToAddress); addr != "" {
				payload["to_location"] = map[string]any{
					"code":    in.ToCityCode,
					"address": addr,
				}
			}
		}
		if in.FromAddress != nil {
			if addr := strings.TrimSpace(*in.FromAddress); addr != "" {
				payload["from_location"] = map[string]any{
					"code":    in.FromCityCode,
					"address": addr,
				}
			}
		}
	}
	if includePvz {
		if in.ToPVZ != nil {
			if code := strings.TrimSpace(*in.ToPVZ); code != "" {
				payload["delivery_point"] = code
			}
		}
		if in.FromPVZ != nil {
			if code := strings.TrimSpace(*in.FromPVZ); code != "" {
				payload["shipment_point"] = code
			}
		}
	}
	return payload
}

func (s *CDEKService) postCdekOrder(ctx context.Context, in CDEKCreateOrderInput, includePvz bool) (*CDEKCreateOrderResult, error) {
	payload := buildCdekOrderPayload(in, includePvz)
	log.Printf(
		"cdek post /orders clientNumber=%s includePvz=%t hasDeliveryPoint=%t hasShipmentPoint=%t hasToAddress=%t",
		strings.TrimSpace(in.ClientNumber),
		includePvz,
		payload["delivery_point"] != nil,
		payload["shipment_point"] != nil,
		func() bool {
			to, ok := payload["to_location"].(map[string]any)
			if !ok {
				return false
			}
			_, ok = to["address"]
			return ok
		}(),
	)
	var raw map[string]any
	if err := s.postJSON(ctx, "/orders", payload, &raw); err != nil {
		return nil, err
	}
	res, soft := parseCDEKOrderCreateResponse(raw)
	if res != nil && res.OrderUUID != nil && strings.TrimSpace(*res.OrderUUID) != "" {
		return res, nil
	}
	if soft != "" {
		return nil, &AppError{400, "CDEK: " + soft}
	}
	return nil, &AppError{502, "CDEK вернул пустой uuid при создании заказа"}
}

// LookupOrderByClientNumber — если заказ с таким number уже создан, подтягиваем uuid/трек.
func (s *CDEKService) LookupOrderByClientNumber(ctx context.Context, clientNumber string) (*CDEKCreateOrderResult, error) {
	clientNumber = strings.TrimSpace(clientNumber)
	if clientNumber == "" {
		return nil, nil
	}
	q := url.Values{}
	q.Set("number", clientNumber)
	q.Set("size", "10")
	var raw map[string]any
	if err := s.getJSON(ctx, "/orders?"+q.Encode(), &raw); err != nil {
		return nil, err
	}
	if ent, ok := raw["entity"].(map[string]any); ok {
		u := normalizeAnyString(ent["uuid"])
		if u != nil {
			t := normalizeAnyString(ent["cdek_number"])
			return &CDEKCreateOrderResult{OrderUUID: u, Track: t}, nil
		}
	}
	if arr, ok := raw["entity"].([]any); ok {
		for _, it := range arr {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			u := normalizeAnyString(m["uuid"])
			if u == nil {
				continue
			}
			t := normalizeAnyString(m["cdek_number"])
			return &CDEKCreateOrderResult{OrderUUID: u, Track: t}, nil
		}
	}
	return nil, nil
}

func parseCDEKOrderCreateResponse(body map[string]any) (*CDEKCreateOrderResult, string) {
	if body == nil {
		return nil, ""
	}
	if ent, ok := body["entity"].(map[string]any); ok {
		u := normalizeAnyString(ent["uuid"])
		t := normalizeAnyString(ent["cdek_number"])
		if u != nil {
			return &CDEKCreateOrderResult{OrderUUID: u, Track: t}, ""
		}
	}
	if msg := collectCDEKRequestMessages(body); msg != "" {
		return nil, msg
	}
	return nil, ""
}

func collectCDEKRequestMessages(body map[string]any) string {
	reqs, ok := body["requests"].([]any)
	if !ok {
		return ""
	}
	var parts []string
	for _, r := range reqs {
		rm, ok := r.(map[string]any)
		if !ok {
			continue
		}
		if errs, ok := rm["errors"].([]any); ok {
			for _, e := range errs {
				em, ok := e.(map[string]any)
				if !ok {
					continue
				}
				if m := normalizeAnyString(em["message"]); m != nil {
					parts = append(parts, *m)
				}
			}
		}
	}
	return strings.Join(parts, "; ")
}

func (s *CDEKService) TrackNumberByOrderUUID(ctx context.Context, orderUUID string) *string {
	details := s.OrderDetailsByOrderUUID(ctx, orderUUID)
	if details == nil {
		return nil
	}
	return details.Track
}

func (s *CDEKService) OrderDetailsByOrderUUID(ctx context.Context, orderUUID string) *CDEKOrderDetails {
	orderUUID = strings.TrimSpace(orderUUID)
	if orderUUID == "" {
		return nil
	}

	var out map[string]any
	if err := s.getJSON(ctx, "/orders/"+orderUUID, &out); err != nil {
		return nil
	}

	entity := out
	if ent, ok := out["entity"].(map[string]any); ok {
		entity = ent
	}

	details := &CDEKOrderDetails{}
	if t := pickTrackFromOrderEntity(entity); t != nil {
		details.Track = t
	}
	if t := pickTrackFromOrderEntity(out); t != nil && details.Track == nil {
		details.Track = t
	}
	code, name := pickLastStatusFromOrderEntity(entity)
	if code == "" && name == "" {
		code, name = pickLastStatusFromOrderEntity(out)
	}
	details.StatusCode = code
	details.StatusName = name
	return details
}

func pickTrackFromOrderEntity(entity map[string]any) *string {
	if entity == nil {
		return nil
	}
	for _, key := range []string{"cdek_number", "tracking_number", "track_number"} {
		if v, ok := entity[key]; ok {
			if s := normalizeAnyString(v); s != nil {
				return s
			}
		}
	}
	return nil
}

func pickLastStatusFromOrderEntity(entity map[string]any) (code, name string) {
	if entity == nil {
		return "", ""
	}
	if statuses, ok := entity["statuses"].([]any); ok && len(statuses) > 0 {
		last := statuses[len(statuses)-1]
		if m, ok := last.(map[string]any); ok {
			if c := normalizeAnyString(m["code"]); c != nil {
				code = strings.ToUpper(*c)
			}
			if n := normalizeAnyString(m["name"]); n != nil {
				name = *n
			}
		}
	}
	if code == "" {
		if c := normalizeAnyString(entity["status"]); c != nil {
			code = strings.ToUpper(*c)
		}
	}
	return code, name
}

func (s *CDEKService) QRByOrderUUID(ctx context.Context, orderUUID string) (*string, *string) {
	orderUUID = strings.TrimSpace(orderUUID)
	if orderUUID == "" {
		return nil, nil
	}

	var out map[string]any
	if err := s.getJSON(ctx, "/orders/"+orderUUID, &out); err != nil {
		return nil, nil
	}

	data := pickFirstString(out,
		"qr_code",
		"qrCode",
		"barcode",
		"barcode_data",
		"barcodeData",
	)
	url := pickFirstString(out,
		"qr_code_url",
		"qrCodeUrl",
		"barcode_url",
		"barcodeUrl",
		"print_url",
		"printUrl",
		"url",
	)
	return data, url
}

func pickFirstString(v any, keys ...string) *string {
	if len(keys) == 0 {
		return nil
	}
	switch t := v.(type) {
	case map[string]any:
		for _, key := range keys {
			if raw, ok := t[key]; ok {
				if s := normalizeAnyString(raw); s != nil {
					return s
				}
			}
		}
		for _, raw := range t {
			if s := pickFirstString(raw, keys...); s != nil {
				return s
			}
		}
	case []any:
		for _, raw := range t {
			if s := pickFirstString(raw, keys...); s != nil {
				return s
			}
		}
	}
	return nil
}

func normalizeAnyString(v any) *string {
	switch t := v.(type) {
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return nil
		}
		return &s
	case fmt.Stringer:
		s := strings.TrimSpace(t.String())
		if s == "" {
			return nil
		}
		return &s
	default:
		return nil
	}
}

func (s *CDEKService) getJSON(ctx context.Context, path string, target any) error {
	token, err := s.token(ctx)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.cfg.CDEKAPIBase+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeCDEKResponse(resp, target)
}

func (s *CDEKService) postJSON(ctx context.Context, path string, payload any, target any) error {
	token, err := s.token(ctx)
	if err != nil {
		return err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.CDEKAPIBase+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeCDEKResponse(resp, target)
}

func (s *CDEKService) token(ctx context.Context) (string, error) {
	if !s.configured() {
		return "", &AppError{400, "CDEK РЅРµ РЅР°СЃС‚СЂРѕРµРЅ (CDEK_CLIENT_ID / CDEK_CLIENT_SECRET)"}
	}
	s.mu.Lock()
	if s.accessToken != "" && time.Now().Before(s.tokenExpires.Add(-time.Minute)) {
		token := s.accessToken
		s.mu.Unlock()
		return token, nil
	}
	s.mu.Unlock()

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", s.cfg.CDEKClientID)
	form.Set("client_secret", s.cfg.CDEKClientSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.CDEKAPIBase+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Some CDEK environments validate client auth via Basic Authorization header.
	cred := s.cfg.CDEKClientID + ":" + s.cfg.CDEKClientSecret
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(cred)))
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := decodeCDEKResponse(resp, &data); err != nil {
		return "", err
	}
	if data.AccessToken == "" {
		return "", fmt.Errorf("CDEK token response is empty")
	}
	expiresIn := data.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	s.mu.Lock()
	s.accessToken = data.AccessToken
	s.tokenExpires = time.Now().Add(time.Duration(expiresIn) * time.Second)
	s.mu.Unlock()
	return data.AccessToken, nil
}

func decodeCDEKResponse(resp *http.Response, target any) error {
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &AppError{Status: resp.StatusCode, Message: "CDEK API error: " + string(raw)}
	}
	if len(raw) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	return dec.Decode(target)
}

func containsCyrillic(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Cyrillic) {
			return true
		}
	}
	return false
}

func transliterateRU(s string) string {
	replacer := strings.NewReplacer(
		"а", "a", "б", "b", "в", "v", "г", "g", "д", "d",
		"е", "e", "ё", "e", "ж", "zh", "з", "z", "и", "i", "й", "y",
		"к", "k", "л", "l", "м", "m", "н", "n", "о", "o", "п", "p",
		"р", "r", "с", "s", "т", "t", "у", "u", "ф", "f", "х", "kh",
		"ц", "ts", "ч", "ch", "ш", "sh", "щ", "shch", "ъ", "", "ы", "y",
		"ь", "", "э", "e", "ю", "yu", "я", "ya",
	)
	lower := strings.ToLower(s)
	latin := replacer.Replace(lower)
	return strings.Join(strings.Fields(latin), " ")
}

func normalizeTariffList(raw any) []map[string]any {
	switch t := raw.(type) {
	case []any:
		out := make([]map[string]any, 0, len(t))
		for _, item := range t {
			if m, ok := item.(map[string]any); ok {
				out = append(out, normalizeTariffItem(m))
			}
		}
		return out
	case map[string]any:
		for _, key := range []string{"tariff_codes", "tariffs", "list"} {
			if arr, ok := t[key].([]any); ok {
				out := make([]map[string]any, 0, len(arr))
				for _, item := range arr {
					if m, ok := item.(map[string]any); ok {
						out = append(out, normalizeTariffItem(m))
					}
				}
				return out
			}
		}
		return []map[string]any{normalizeTariffItem(t)}
	default:
		return []map[string]any{}
	}
}

func normalizeTariffItem(in map[string]any) map[string]any {
	out := map[string]any{}
	if v, ok := in["tariff_code"]; ok {
		out["tariffCode"] = v
	} else if v, ok := in["code"]; ok {
		out["tariffCode"] = v
	}
	if v, ok := in["tariff_name"]; ok {
		out["tariffName"] = v
	} else if v, ok := in["name"]; ok {
		out["tariffName"] = v
	}
	if v, ok := in["delivery_mode"]; ok {
		out["deliveryMode"] = v
	}
	if v, ok := in["from_door"]; ok {
		out["fromDoor"] = v
	}
	if v, ok := in["to_door"]; ok {
		out["toDoor"] = v
	}
	if v, ok := in["period_min"]; ok {
		out["periodMin"] = v
	}
	if v, ok := in["period_max"]; ok {
		out["periodMax"] = v
	}
	if v, ok := in["total_sum"]; ok {
		out["totalSum"] = v
	}
	out["raw"] = in
	return out
}
