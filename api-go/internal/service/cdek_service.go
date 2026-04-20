package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (s *CDEKService) Calculate(ctx context.Context, req CDEKCalculateRequest) (map[string]any, error) {
	if req.TariffCode <= 0 || req.FromCityCode <= 0 || req.ToCityCode <= 0 {
		return nil, &AppError{400, "РќСѓР¶РЅС‹ tariffCode, fromCityCode Рё toCityCode"}
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
