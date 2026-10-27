package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Sender отправляет SMS-уведомления через MTS или Notisend с fallback/логированием.
type Sender struct {
	MTSBearer       string
	NotisendAPIKey  string
	NotisendProject string
	client          *http.Client
}

// NewSender создает новый экземпляр Sender.
func NewSender(mtsBearer, notisendAPIKey, notisendProject string) *Sender {
	return &Sender{
		MTSBearer:       strings.TrimSpace(mtsBearer),
		NotisendAPIKey:  strings.TrimSpace(notisendAPIKey),
		NotisendProject: strings.TrimSpace(notisendProject),
		client:          &http.Client{Timeout: 15 * time.Second},
	}
}

// NormalizePhone очищает номер телефона для отправки SMS.
func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}
	// Игнорируем технические префиксы соцсетей (например VK_..., YANDEX_...)
	if strings.HasPrefix(strings.ToUpper(phone), "VK_") || strings.HasPrefix(strings.ToUpper(phone), "YA_") {
		return ""
	}
	var b strings.Builder
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if r == '+' && b.Len() == 0 {
			b.WriteRune(r)
		}
	}
	cleaned := b.String()
	// Если номер 11 цифр и начинается с 8, заменяем на +7
	if len(cleaned) == 11 && cleaned[0] == '8' {
		cleaned = "+7" + cleaned[1:]
	} else if len(cleaned) == 10 {
		cleaned = "+7" + cleaned
	} else if len(cleaned) == 11 && cleaned[0] == '7' {
		cleaned = "+" + cleaned
	}
	return cleaned
}

// Send отправляет SMS синхронно.
func (s *Sender) Send(ctx context.Context, phone, message string) error {
	phone = NormalizePhone(phone)
	message = strings.TrimSpace(message)
	if phone == "" || message == "" {
		return nil
	}

	if s.MTSBearer != "" {
		body := map[string]any{
			"submits": []any{map[string]any{"msid": phone, "message": message}},
			"naming":  "Torguisamru",
		}
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost,
			"https://api.mts.ru/client-omni-adapter_production/1.0.2/mcom/messageManagement/messages",
			bytes.NewReader(raw))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.MTSBearer)

		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("sms(mts) error to %s: %v", phone, err)
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			log.Printf("sms(mts) non-ok status %s to %s", resp.Status, phone)
			return fmt.Errorf("mts status %s", resp.Status)
		}
		log.Printf("sms(mts) sent successfully to %s: %s", phone, message)
		return nil
	}

	if s.NotisendAPIKey != "" {
		u := "https://sms.notisend.ru/api/message/send?project=" + url.QueryEscape(s.NotisendProject) +
			"&message=" + url.QueryEscape(message) +
			"&recipients=" + url.QueryEscape(phone) +
			"&apikey=" + url.QueryEscape(s.NotisendAPIKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return err
		}
		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("sms(notisend) error to %s: %v", phone, err)
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			log.Printf("sms(notisend) non-ok status %s to %s", resp.Status, phone)
			return fmt.Errorf("notisend status %s", resp.Status)
		}
		log.Printf("sms(notisend) sent successfully to %s: %s", phone, message)
		return nil
	}

	// В режиме разработки или если шлюз не настроен — логируем
	log.Printf("[SMS GATEWAY NOT CONFIGURED] Message to %s: %s", phone, message)
	return nil
}

// SendAsync отправляет SMS в фоновой горутине.
func (s *Sender) SendAsync(phone, message string) {
	phone = NormalizePhone(phone)
	if phone == "" || strings.TrimSpace(message) == "" {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		_ = s.Send(ctx, phone, message)
	}()
}
