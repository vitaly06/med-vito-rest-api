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
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/repository"
)

const (
	aiApprovedReason             = "Одобрено ИИ автоматически"
	visionTechnicalErrorReason   = "Ошибка анализа фото, требуется ручная проверка"
	textTechnicalErrorReason     = "Ошибка ИИ-сервиса, требуется ручная проверка"
	defaultManualReviewReason    = "Требуется ручная проверка"
	yandexTextEndpoint           = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
	yandexVisionEndpoint         = "https://ai.api.cloud.yandex.net/v1/chat/completions"
	visionImageDownloadTimeout   = 30 * time.Second
	visionRequestTimeout         = 60 * time.Second
	textRequestTimeout           = 25 * time.Second
	visionModerationUserAgent    = "MedVito-ModerationWorker/1.0"
	moderationWorkerDefaultDelay = 30 * time.Second
)

const textSystemPrompt = `Ты — автоматический модератор объявлений на маркетплейсе.


Запрещено:
- Контакты в тексте: телефон, email, ник в мессенджере (даже написанные словами: "восемь девятьсот", "собака", "тг @...")
- Внешние ссылки (http://, t.me/, wa.me/ и т.д.)
- Нецензурная лексика и оскорбления
- Мошеннические признаки: "предоплата", "переведи деньги", "только безнал", "аванс"
- SEO-спам: многократные повторы одних и тех же слов или ключевых фраз, бессмысленный набор текста, перечисление несвязанных слов для выдачи в поиске
- НЕ является спамом: вежливые фразы продавца ("советуем заглянуть", "в нашем профиле есть другие товары", "отличное качество"), стандартные описания товара, упоминание ассортимента магазина
- Подозрительно низкая цена (менее 10% от рыночной для категории)

Для каждого из 4 критериев укажи статус: OK / SUSPICIOUS / VIOLATION.
- OK — нарушений нет
- SUSPICIOUS — есть признаки, но неоднозначно (требуется модератор)
- VIOLATION — явное нарушение (автоотказ)

Итоговое решение:
- APPROVED — все критерии OK
- MANUAL — хотя бы один SUSPICIOUS, нет VIOLATION
- DENIED — хотя бы один VIOLATION

Отвечай СТРОГО в JSON без markdown-обёртки:
{
  "category": "APPROVED" | "MANUAL" | "DENIED",
  "reason": "Объяснение на русском (пустая строка если APPROVED)",
  "details": {
    "categorization": "OK" | "SUSPICIOUS" | "VIOLATION",
    "spam": "OK" | "SUSPICIOUS" | "VIOLATION",
    "fraud": "OK" | "SUSPICIOUS" | "VIOLATION",
    "contacts": "OK" | "SUSPICIOUS" | "VIOLATION"
  }
}`

const visionPrompt = `Это фото для объявления на маркетплейсе.

Проверь наличие ЛЮБОГО из следующих нарушений:
1. Оружие, боеприпасы, взрывчатка, ножи как основной товар
2. NSFW / откровенный контент / части тела
3. Насилие, кровь, шокирующие материалы
4. Скриншот стороннего сайта/приложения с контактами (телефон, email, ник)
5. Наркотики, алкоголь, табак

Если хотя бы одно нарушение есть — DENIED.
Если фото нечёткое, подозрительное или невозможно определить товар — MANUAL.
Если фото обычного товара (в т.ч. немедицинского бытового) без нарушений — APPROVED.

Ответь СТРОГО в JSON без markdown-обёртки:
{
  "decision": "APPROVED" | "MANUAL" | "DENIED",
  "reason": "Объяснение на русском что именно нарушено (пустая строка если APPROVED)"
}

DENIED — явное нарушение из списка выше
MANUAL — сомнительно, нужна ручная проверка
APPROVED — фото подходит`

type TextModerationResult struct {
	Category string `json:"category"`
	Reason   string `json:"reason"`
	Details  struct {
		Categorization string `json:"categorization"`
		Spam           string `json:"spam"`
		Fraud          string `json:"fraud"`
		Contacts       string `json:"contacts"`
	} `json:"details"`
}

type VisionModerationResult struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

type ModerationService struct {
	cfg          config.Config
	repo         *repository.ProductPG
	httpClient   *http.Client
	pollInterval time.Duration
	batchSize    int

	mu           sync.Mutex
	isProcessing bool
}

func NewModerationService(cfg config.Config, repo *repository.ProductPG) *ModerationService {
	poll := time.Duration(cfg.AIModerationPollInterval) * time.Second
	if poll <= 0 {
		poll = moderationWorkerDefaultDelay
	}
	batchSize := cfg.AIModerationBatchSize
	if batchSize <= 0 {
		batchSize = 5
	}
	return &ModerationService{
		cfg:          cfg,
		repo:         repo,
		httpClient:   &http.Client{},
		pollInterval: poll,
		batchSize:    batchSize,
	}
}

func (s *ModerationService) Enabled() bool {
	return s.cfg.AIModerationEnabled &&
		strings.TrimSpace(s.cfg.YandexFolderID) != "" &&
		strings.TrimSpace(s.cfg.YandexAPIKey) != ""
}

func (s *ModerationService) Start(ctx context.Context) {
	if !s.Enabled() {
		log.Printf("AI moderation worker disabled")
		return
	}
	log.Printf("AI moderation worker started. Poll interval: %ds", int(math.Round(s.pollInterval.Seconds())))
	go s.run(ctx)
}

func (s *ModerationService) run(ctx context.Context) {
	s.runPoll(ctx)
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("AI moderation worker stopped")
			return
		case <-ticker.C:
			s.runPoll(ctx)
		}
	}
}

func (s *ModerationService) runPoll(ctx context.Context) {
	s.mu.Lock()
	if s.isProcessing {
		s.mu.Unlock()
		log.Printf("AI moderation worker: previous poll still running, skipping tick")
		return
	}
	s.isProcessing = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.isProcessing = false
		s.mu.Unlock()
	}()

	products, err := s.repo.ListPendingAIModeration(ctx, s.batchSize)
	if err != nil {
		log.Printf("AI moderation worker poll error: %v", err)
		return
	}
	if len(products) == 0 {
		return
	}

	log.Printf("AI moderation worker: found %d product(s) pending moderation", len(products))
	for _, product := range products {
		s.processProduct(ctx, product)
	}
}

func (s *ModerationService) processProduct(ctx context.Context, product repository.AIModerationProduct) {
	log.Printf(`AI moderation worker: [#%d] Processing "%s"`, product.ID, product.Name)

	textResult, err := s.moderateText(ctx, product.Name, product.Description, product.CategoryName, product.SubCategoryName, product.Price)
	if err != nil {
		log.Printf("AI moderation worker: [#%d] text moderation error: %v", product.ID, err)
		return
	}
	log.Printf("AI moderation worker: [#%d] Text -> %s%s", product.ID, textResult.Category, formatReasonSuffix(textResult.Reason))

	if textResult.Category == "DENIED" {
		if err := s.applyDecision(ctx, product.ID, "DENIDED", nullableReason(textResult.Reason)); err != nil {
			log.Printf("AI moderation worker: [#%d] apply denied decision error: %v", product.ID, err)
		}
		return
	}

	visionDecision := "APPROVED"
	visionReason := ""
	visionTechnicalFailure := false

	for _, imageURL := range product.Images {
		visionResult, err := s.moderateImage(ctx, imageURL)
		if err != nil {
			log.Printf("AI moderation worker: [#%d] vision moderation error: %v", product.ID, err)
			continue
		}
		log.Printf("AI moderation worker: [#%d] Vision [%s] -> %s", product.ID, imageURL, visionResult.Decision)

		if visionResult.Decision == "MANUAL" && visionResult.Reason == visionTechnicalErrorReason {
			visionTechnicalFailure = true
			continue
		}
		if visionResult.Decision == "DENIED" {
			visionDecision = "DENIED"
			visionReason = visionResult.Reason
			break
		}
		if visionResult.Decision == "MANUAL" {
			visionDecision = "MANUAL"
			visionReason = visionResult.Reason
		}
	}

	if visionDecision == "DENIED" {
		if err := s.applyDecision(ctx, product.ID, "DENIDED", nullableReason(visionReason)); err != nil {
			log.Printf("AI moderation worker: [#%d] apply denied vision decision error: %v", product.ID, err)
		}
		return
	}

	var reasons []string
	if textResult.Category == "MANUAL" && strings.TrimSpace(textResult.Reason) != "" {
		reasons = append(reasons, "Текст: "+strings.TrimSpace(textResult.Reason))
	}
	if visionDecision == "MANUAL" && strings.TrimSpace(visionReason) != "" {
		reasons = append(reasons, "Фото: "+strings.TrimSpace(visionReason))
	}
	if visionTechnicalFailure && visionDecision != "MANUAL" {
		reasons = append(reasons, visionTechnicalErrorReason)
	}

	if len(reasons) > 0 {
		manualReason := strings.Join(reasons, " / ")
		if err := s.applyDecision(ctx, product.ID, "AI_REVIEWED", &manualReason); err != nil {
			log.Printf("AI moderation worker: [#%d] apply manual decision error: %v", product.ID, err)
			return
		}
		log.Printf("AI moderation worker: [#%d] -> MANUAL (%s)", product.ID, manualReason)
		return
	}

	if textResult.Category == "APPROVED" && visionDecision == "APPROVED" {
		if err := s.applyDecision(ctx, product.ID, "APPROVED", stringPtr(aiApprovedReason)); err != nil {
			log.Printf("AI moderation worker: [#%d] apply approved decision error: %v", product.ID, err)
			return
		}
		log.Printf("AI moderation worker: [#%d] -> APPROVED (%s)", product.ID, aiApprovedReason)
		return
	}

	if err := s.applyDecision(ctx, product.ID, "AI_REVIEWED", stringPtr(defaultManualReviewReason)); err != nil {
		log.Printf("AI moderation worker: [#%d] apply fallback manual decision error: %v", product.ID, err)
		return
	}
	log.Printf("AI moderation worker: [#%d] -> MANUAL (fallback)", product.ID)
}

func (s *ModerationService) applyDecision(ctx context.Context, productID int32, state string, reason *string) error {
	return s.repo.SetModerationState(ctx, productID, state, reason)
}

func (s *ModerationService) moderateText(ctx context.Context, name string, description *string, categoryName, subcategoryName string, price int32) (*TextModerationResult, error) {
	userPrompt := fmt.Sprintf("Категория: %s\nПодкатегория: %s\nЦена: %d руб.\nНазвание: %s\nОписание: %s",
		categoryName,
		subcategoryName,
		price,
		name,
		descriptionOrDefault(description, "не указано"),
	)
	payload := map[string]any{
		"modelUri": fmt.Sprintf("gpt://%s/yandexgpt/latest", s.cfg.YandexFolderID),
		"completionOptions": map[string]any{
			"stream":      false,
			"temperature": 0.05,
			"maxTokens":   250,
		},
		"messages": []map[string]string{
			{"role": "system", "text": textSystemPrompt},
			{"role": "user", "text": userPrompt},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, yandexTextEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Api-Key "+s.cfg.YandexAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := *s.httpClient
	client.Timeout = textRequestTimeout
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("AI moderation worker: text moderation transport error: %v", err)
		return fallbackTextResult(), nil
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("AI moderation worker: text moderation bad status %d body=%s", resp.StatusCode, string(rawBody))
		return fallbackTextResult(), nil
	}

	var response struct {
		Result struct {
			Alternatives []struct {
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"alternatives"`
		} `json:"result"`
	}
	if err := json.Unmarshal(rawBody, &response); err != nil {
		log.Printf("AI moderation worker: text moderation decode error: %v", err)
		return fallbackTextResult(), nil
	}
	if len(response.Result.Alternatives) == 0 {
		log.Printf("AI moderation worker: text moderation response missing alternatives")
		return fallbackTextResult(), nil
	}

	cleaned := cleanupModelJSON(response.Result.Alternatives[0].Message.Text)
	var parsed TextModerationResult
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		log.Printf("AI moderation worker: text moderation parse error: %v input=%s", err, cleaned)
		return fallbackTextResult(), nil
	}
	if !validTextCategory(parsed.Category) || strings.TrimSpace(parsed.Reason) == "" && parsed.Category != "APPROVED" {
		log.Printf("AI moderation worker: text moderation invalid payload: %s", cleaned)
		return fallbackTextResult(), nil
	}
	return &parsed, nil
}

func (s *ModerationService) moderateImage(ctx context.Context, imageURL string) (*VisionModerationResult, error) {
	base64Image, mimeType, err := s.downloadImageAsBase64(ctx, imageURL)
	if err != nil {
		log.Printf("AI moderation worker: [Vision] download error for %s: %v", imageURL, err)
		return fallbackVisionResult(), nil
	}

	payload := map[string]any{
		"model":       fmt.Sprintf("gpt://%s/gemma-3-27b-it/latest", s.cfg.YandexFolderID),
		"max_tokens":  150,
		"temperature": 0.05,
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image),
						},
					},
					{
						"type": "text",
						"text": visionPrompt,
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, yandexVisionEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Api-Key "+s.cfg.YandexAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-folder-id", s.cfg.YandexFolderID)

	client := *s.httpClient
	client.Timeout = visionRequestTimeout
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("AI moderation worker: [Vision] transport error: %v", err)
		return fallbackVisionResult(), nil
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("AI moderation worker: [Vision] bad status %d body=%s", resp.StatusCode, string(rawBody))
		return fallbackVisionResult(), nil
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(rawBody, &response); err != nil {
		log.Printf("AI moderation worker: [Vision] decode error: %v", err)
		return fallbackVisionResult(), nil
	}
	if len(response.Choices) == 0 || strings.TrimSpace(response.Choices[0].Message.Content) == "" {
		log.Printf("AI moderation worker: [Vision] response missing choices[0].message.content")
		return fallbackVisionResult(), nil
	}

	cleaned := cleanupModelJSON(response.Choices[0].Message.Content)
	var parsed VisionModerationResult
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		log.Printf("AI moderation worker: [Vision] parse error: %v input=%s", err, cleaned)
		return fallbackVisionResult(), nil
	}
	if !validVisionDecision(parsed.Decision) {
		log.Printf("AI moderation worker: [Vision] invalid payload: %s", cleaned)
		return fallbackVisionResult(), nil
	}
	return &parsed, nil
}

func (s *ModerationService) downloadImageAsBase64(ctx context.Context, imageURL string) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", visionModerationUserAgent)

	client := *s.httpClient
	client.Timeout = visionImageDownloadTimeout
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		rawBody, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("download status %d body=%s", resp.StatusCode, string(rawBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	mimeType := strings.TrimSpace(resp.Header.Get("Content-Type"))
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	return base64.StdEncoding.EncodeToString(body), mimeType, nil
}

type ModerationAdminService struct {
	repo *repository.ProductPG
}

func NewModerationAdminService(repo *repository.ProductPG) *ModerationAdminService {
	return &ModerationAdminService{repo: repo}
}

func (s *ModerationAdminService) ListProducts(ctx context.Context, filter string, page int) (map[string]any, error) {
	items, total, err := s.repo.ListProductsForModeration(ctx, filter, page, 20)
	if err != nil {
		return nil, err
	}
	pages := 0
	if total > 0 {
		pages = int(math.Ceil(float64(total) / 20))
	}

	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]any{
			"id":                        item.ID,
			"name":                      item.Name,
			"price":                     item.Price,
			"images":                    item.Images,
			"moderateState":             item.ModerateState,
			"moderationRejectionReason": item.ModerationRejectionReason,
			"createdAt":                 item.CreatedAt,
			"updatedAt":                 item.UpdatedAt,
			"category":                  map[string]any{"id": item.CategoryID, "name": item.CategoryName},
			"subCategory":               map[string]any{"id": item.SubCategoryID, "name": item.SubCategoryName},
			"user": map[string]any{
				"id":          item.UserID,
				"fullName":    item.UserFullName,
				"email":       item.UserEmail,
				"phoneNumber": item.UserPhoneNumber,
			},
		})
	}

	return map[string]any{
		"items": result,
		"total": total,
		"page":  page,
		"pages": pages,
	}, nil
}

func (s *ModerationAdminService) GetProduct(ctx context.Context, productID int32) (map[string]any, error) {
	item, err := s.repo.GetModerationProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Товар не найден"}
		}
		return nil, err
	}

	fieldValues := make([]map[string]any, 0, len(item.FieldValues))
	for _, fieldValue := range item.FieldValues {
		fieldValues = append(fieldValues, map[string]any{
			"value": fieldValue.Value,
			"field": map[string]any{
				"id":   fieldValue.FieldID,
				"name": fieldValue.Field,
			},
		})
	}

	var typeMap map[string]any
	if item.TypeID != nil {
		typeMap = map[string]any{"id": *item.TypeID, "name": item.TypeName}
	}

	return map[string]any{
		"id":                        item.ID,
		"name":                      item.Name,
		"price":                     item.Price,
		"description":               item.Description,
		"images":                    item.Images,
		"videoUrl":                  item.VideoURL,
		"moderateState":             item.ModerateState,
		"moderationRejectionReason": item.ModerationRejectionReason,
		"createdAt":                 item.CreatedAt,
		"updatedAt":                 item.UpdatedAt,
		"category":                  map[string]any{"id": item.CategoryID, "name": item.CategoryName},
		"subCategory":               map[string]any{"id": item.SubCategoryID, "name": item.SubCategoryName},
		"type":                      typeMap,
		"user": map[string]any{
			"id":          item.UserID,
			"fullName":    item.UserFullName,
			"email":       item.UserEmail,
			"phoneNumber": item.UserPhoneNumber,
			"profileType": item.UserProfileType,
		},
		"fieldValues": fieldValues,
	}, nil
}

func cleanupModelJSON(raw string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(raw, "```json", ""), "```", ""))
}

func fallbackTextResult() *TextModerationResult {
	result := &TextModerationResult{
		Category: "MANUAL",
		Reason:   textTechnicalErrorReason,
	}
	result.Details.Categorization = "OK"
	result.Details.Spam = "OK"
	result.Details.Fraud = "OK"
	result.Details.Contacts = "OK"
	return result
}

func fallbackVisionResult() *VisionModerationResult {
	return &VisionModerationResult{
		Decision: "MANUAL",
		Reason:   visionTechnicalErrorReason,
	}
}

func validTextCategory(v string) bool {
	switch v {
	case "APPROVED", "MANUAL", "DENIED":
		return true
	default:
		return false
	}
}

func validVisionDecision(v string) bool {
	switch v {
	case "APPROVED", "MANUAL", "DENIED":
		return true
	default:
		return false
	}
}

func stringPtr(v string) *string {
	return &v
}

func nullableReason(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}

func descriptionOrDefault(v *string, fallback string) string {
	if v == nil || strings.TrimSpace(*v) == "" {
		return fallback
	}
	return *v
}

func formatReasonSuffix(reason string) string {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return ""
	}
	return " (" + reason + ")"
}
