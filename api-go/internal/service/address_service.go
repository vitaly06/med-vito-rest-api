package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

const dadataSuggestURL = "https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address"

type AddressService struct {
	apiKey string
	cli    *http.Client
}

func NewAddressService(apiKey string) *AddressService {
	return &AddressService{
		apiKey: strings.TrimSpace(apiKey),
		cli:    &http.Client{Timeout: 15 * time.Second},
	}
}

// AddressSuggestionItem — как маппинг в Nest getSuggestions.
type AddressSuggestionItem struct {
	Value string `json:"value"`
	Lat   any    `json:"lat"`
	Lon   any    `json:"lon"`
}

type dadataSuggestResponse struct {
	Suggestions []struct {
		Value string `json:"value"`
		Data  struct {
			GeoLat any `json:"geo_lat"`
			GeoLon any `json:"geo_lon"`
		} `json:"data"`
	} `json:"suggestions"`
}

func (s *AddressService) GetSuggestions(ctx context.Context, query string, limit int) ([]AddressSuggestionItem, error) {
	if s.apiKey == "" {
		return nil, &AppError{500, "DADATA_API_KEY не найден в конфигурации"}
	}
	if limit < 1 {
		limit = 5
	}
	body, err := json.Marshal(map[string]any{"query": query, "count": limit})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, dadataSuggestURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token "+s.apiKey)

	resp, err := s.cli.Do(req)
	if err != nil {
		return nil, &AppError{502, "Ошибка получения подсказок адреса"}
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &AppError{502, "Ошибка получения подсказок адреса"}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errBody struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(raw, &errBody)
		msg := errBody.Message
		if msg == "" {
			msg = "Ошибка получения подсказок адреса"
		}
		st := resp.StatusCode
		if st < 400 || st > 599 {
			st = 502
		}
		return nil, &AppError{st, msg}
	}

	var parsed dadataSuggestResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, &AppError{502, "Ошибка разбора ответа DaData"}
	}
	out := make([]AddressSuggestionItem, 0, len(parsed.Suggestions))
	for _, su := range parsed.Suggestions {
		out = append(out, AddressSuggestionItem{
			Value: su.Value,
			Lat:   su.Data.GeoLat,
			Lon:   su.Data.GeoLon,
		})
	}
	return out, nil
}

// ValidateAddress — как Nest: при любой ошибке DaData — false; совпадение с первой подсказкой без учёта регистра.
func (s *AddressService) ValidateAddress(ctx context.Context, address string) bool {
	list, err := s.GetSuggestions(ctx, address, 1)
	if err != nil || len(list) == 0 {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(list[0].Value), strings.TrimSpace(address))
}
