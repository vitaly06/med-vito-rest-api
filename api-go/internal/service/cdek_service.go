package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

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
	if strings.TrimSpace(city) == "" {
		return nil, &AppError{400, "Нужно указать city"}
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := url.Values{}
	q.Set("city", strings.TrimSpace(city))
	q.Set("size", strconv.Itoa(limit))
	var out []map[string]any
	if err := s.getJSON(ctx, "/location/cities?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *CDEKService) DeliveryPoints(ctx context.Context, cityCode int) ([]map[string]any, error) {
	if cityCode <= 0 {
		return nil, &AppError{400, "Нужен cityCode"}
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
		return nil, &AppError{400, "Нужны tariffCode, fromCityCode и toCityCode"}
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
		return "", &AppError{400, "CDEK не настроен (CDEK_CLIENT_ID / CDEK_CLIENT_SECRET)"}
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
