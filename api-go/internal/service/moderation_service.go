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
	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/repository"
)

const (
	aiApprovedReason             = "РћРґРѕР±СЂРµРЅРѕ РР Р°РІС‚РѕРјР°С‚РёС‡РµСЃРєРё"
	visionTechnicalErrorReason   = "РћС€РёР±РєР° Р°РЅР°Р»РёР·Р° С„РѕС‚Рѕ, С‚СЂРµР±СѓРµС‚СЃСЏ СЂСѓС‡РЅР°СЏ РїСЂРѕРІРµСЂРєР°"
	textTechnicalErrorReason     = "РћС€РёР±РєР° РР-СЃРµСЂРІРёСЃР°, С‚СЂРµР±СѓРµС‚СЃСЏ СЂСѓС‡РЅР°СЏ РїСЂРѕРІРµСЂРєР°"
	defaultManualReviewReason    = "РўСЂРµР±СѓРµС‚СЃСЏ СЂСѓС‡РЅР°СЏ РїСЂРѕРІРµСЂРєР°"
	yandexTextEndpoint           = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
	yandexVisionEndpoint         = "https://ai.api.cloud.yandex.net/v1/chat/completions"
	visionImageDownloadTimeout   = 30 * time.Second
	visionRequestTimeout         = 60 * time.Second
	textRequestTimeout           = 25 * time.Second
	visionModerationUserAgent    = "MedVito-ModerationWorker/1.0"
	moderationWorkerDefaultDelay = 30 * time.Second
)

const textSystemPrompt = `РўС‹ вЂ” Р°РІС‚РѕРјР°С‚РёС‡РµСЃРєРёР№ РјРѕРґРµСЂР°С‚РѕСЂ РѕР±СЉСЏРІР»РµРЅРёР№ РЅР° РјР°СЂРєРµС‚РїР»РµР№СЃРµ.


Р—Р°РїСЂРµС‰РµРЅРѕ:
- РљРѕРЅС‚Р°РєС‚С‹ РІ С‚РµРєСЃС‚Рµ: С‚РµР»РµС„РѕРЅ, email, РЅРёРє РІ РјРµСЃСЃРµРЅРґР¶РµСЂРµ (РґР°Р¶Рµ РЅР°РїРёСЃР°РЅРЅС‹Рµ СЃР»РѕРІР°РјРё: "РІРѕСЃРµРјСЊ РґРµРІСЏС‚СЊСЃРѕС‚", "СЃРѕР±Р°РєР°", "С‚Рі @...")
- Р’РЅРµС€РЅРёРµ СЃСЃС‹Р»РєРё (http://, t.me/, wa.me/ Рё С‚.Рґ.)
- РќРµС†РµРЅР·СѓСЂРЅР°СЏ Р»РµРєСЃРёРєР° Рё РѕСЃРєРѕСЂР±Р»РµРЅРёСЏ
- РњРѕС€РµРЅРЅРёС‡РµСЃРєРёРµ РїСЂРёР·РЅР°РєРё: "РїСЂРµРґРѕРїР»Р°С‚Р°", "РїРµСЂРµРІРµРґРё РґРµРЅСЊРіРё", "С‚РѕР»СЊРєРѕ Р±РµР·РЅР°Р»", "Р°РІР°РЅСЃ"
- SEO-СЃРїР°Рј: РјРЅРѕРіРѕРєСЂР°С‚РЅС‹Рµ РїРѕРІС‚РѕСЂС‹ РѕРґРЅРёС… Рё С‚РµС… Р¶Рµ СЃР»РѕРІ РёР»Рё РєР»СЋС‡РµРІС‹С… С„СЂР°Р·, Р±РµСЃСЃРјС‹СЃР»РµРЅРЅС‹Р№ РЅР°Р±РѕСЂ С‚РµРєСЃС‚Р°, РїРµСЂРµС‡РёСЃР»РµРЅРёРµ РЅРµСЃРІСЏР·Р°РЅРЅС‹С… СЃР»РѕРІ РґР»СЏ РІС‹РґР°С‡Рё РІ РїРѕРёСЃРєРµ
- РќР• СЏРІР»СЏРµС‚СЃСЏ СЃРїР°РјРѕРј: РІРµР¶Р»РёРІС‹Рµ С„СЂР°Р·С‹ РїСЂРѕРґР°РІС†Р° ("СЃРѕРІРµС‚СѓРµРј Р·Р°РіР»СЏРЅСѓС‚СЊ", "РІ РЅР°С€РµРј РїСЂРѕС„РёР»Рµ РµСЃС‚СЊ РґСЂСѓРіРёРµ С‚РѕРІР°СЂС‹", "РѕС‚Р»РёС‡РЅРѕРµ РєР°С‡РµСЃС‚РІРѕ"), СЃС‚Р°РЅРґР°СЂС‚РЅС‹Рµ РѕРїРёСЃР°РЅРёСЏ С‚РѕРІР°СЂР°, СѓРїРѕРјРёРЅР°РЅРёРµ Р°СЃСЃРѕСЂС‚РёРјРµРЅС‚Р° РјР°РіР°Р·РёРЅР°
- РџРѕРґРѕР·СЂРёС‚РµР»СЊРЅРѕ РЅРёР·РєР°СЏ С†РµРЅР° (РјРµРЅРµРµ 10% РѕС‚ СЂС‹РЅРѕС‡РЅРѕР№ РґР»СЏ РєР°С‚РµРіРѕСЂРёРё)

Р”Р»СЏ РєР°Р¶РґРѕРіРѕ РёР· 4 РєСЂРёС‚РµСЂРёРµРІ СѓРєР°Р¶Рё СЃС‚Р°С‚СѓСЃ: OK / SUSPICIOUS / VIOLATION.
- OK вЂ” РЅР°СЂСѓС€РµРЅРёР№ РЅРµС‚
- SUSPICIOUS вЂ” РµСЃС‚СЊ РїСЂРёР·РЅР°РєРё, РЅРѕ РЅРµРѕРґРЅРѕР·РЅР°С‡РЅРѕ (С‚СЂРµР±СѓРµС‚СЃСЏ РјРѕРґРµСЂР°С‚РѕСЂ)
- VIOLATION вЂ” СЏРІРЅРѕРµ РЅР°СЂСѓС€РµРЅРёРµ (Р°РІС‚РѕРѕС‚РєР°Р·)

РС‚РѕРіРѕРІРѕРµ СЂРµС€РµРЅРёРµ:
- APPROVED вЂ” РІСЃРµ РєСЂРёС‚РµСЂРёРё OK
- MANUAL вЂ” С…РѕС‚СЏ Р±С‹ РѕРґРёРЅ SUSPICIOUS, РЅРµС‚ VIOLATION
- DENIED вЂ” С…РѕС‚СЏ Р±С‹ РѕРґРёРЅ VIOLATION

РћС‚РІРµС‡Р°Р№ РЎРўР РћР“Рћ РІ JSON Р±РµР· markdown-РѕР±С‘СЂС‚РєРё:
{
  "category": "APPROVED" | "MANUAL" | "DENIED",
  "reason": "РћР±СЉСЏСЃРЅРµРЅРёРµ РЅР° СЂСѓСЃСЃРєРѕРј (РїСѓСЃС‚Р°СЏ СЃС‚СЂРѕРєР° РµСЃР»Рё APPROVED)",
  "details": {
    "categorization": "OK" | "SUSPICIOUS" | "VIOLATION",
    "spam": "OK" | "SUSPICIOUS" | "VIOLATION",
    "fraud": "OK" | "SUSPICIOUS" | "VIOLATION",
    "contacts": "OK" | "SUSPICIOUS" | "VIOLATION"
  }
}`

const visionPrompt = `Р­С‚Рѕ С„РѕС‚Рѕ РґР»СЏ РѕР±СЉСЏРІР»РµРЅРёСЏ РЅР° РјР°СЂРєРµС‚РїР»РµР№СЃРµ.

РџСЂРѕРІРµСЂСЊ РЅР°Р»РёС‡РёРµ Р›Р®Р‘РћР“Рћ РёР· СЃР»РµРґСѓСЋС‰РёС… РЅР°СЂСѓС€РµРЅРёР№:
1. РћСЂСѓР¶РёРµ, Р±РѕРµРїСЂРёРїР°СЃС‹, РІР·СЂС‹РІС‡Р°С‚РєР°, РЅРѕР¶Рё РєР°Рє РѕСЃРЅРѕРІРЅРѕР№ С‚РѕРІР°СЂ
2. NSFW / РѕС‚РєСЂРѕРІРµРЅРЅС‹Р№ РєРѕРЅС‚РµРЅС‚ / С‡Р°СЃС‚Рё С‚РµР»Р°
3. РќР°СЃРёР»РёРµ, РєСЂРѕРІСЊ, С€РѕРєРёСЂСѓСЋС‰РёРµ РјР°С‚РµСЂРёР°Р»С‹
4. РЎРєСЂРёРЅС€РѕС‚ СЃС‚РѕСЂРѕРЅРЅРµРіРѕ СЃР°Р№С‚Р°/РїСЂРёР»РѕР¶РµРЅРёСЏ СЃ РєРѕРЅС‚Р°РєС‚Р°РјРё (С‚РµР»РµС„РѕРЅ, email, РЅРёРє)
5. РќР°СЂРєРѕС‚РёРєРё, Р°Р»РєРѕРіРѕР»СЊ, С‚Р°Р±Р°Рє

Р•СЃР»Рё С…РѕС‚СЏ Р±С‹ РѕРґРЅРѕ РЅР°СЂСѓС€РµРЅРёРµ РµСЃС‚СЊ вЂ” DENIED.
Р•СЃР»Рё С„РѕС‚Рѕ РЅРµС‡С‘С‚РєРѕРµ, РїРѕРґРѕР·СЂРёС‚РµР»СЊРЅРѕРµ РёР»Рё РЅРµРІРѕР·РјРѕР¶РЅРѕ РѕРїСЂРµРґРµР»РёС‚СЊ С‚РѕРІР°СЂ вЂ” MANUAL.
Р•СЃР»Рё С„РѕС‚Рѕ РѕР±С‹С‡РЅРѕРіРѕ С‚РѕРІР°СЂР° (РІ С‚.С‡. РЅРµРјРµРґРёС†РёРЅСЃРєРѕРіРѕ Р±С‹С‚РѕРІРѕРіРѕ) Р±РµР· РЅР°СЂСѓС€РµРЅРёР№ вЂ” APPROVED.

РћС‚РІРµС‚СЊ РЎРўР РћР“Рћ РІ JSON Р±РµР· markdown-РѕР±С‘СЂС‚РєРё:
{
  "decision": "APPROVED" | "MANUAL" | "DENIED",
  "reason": "РћР±СЉСЏСЃРЅРµРЅРёРµ РЅР° СЂСѓСЃСЃРєРѕРј С‡С‚Рѕ РёРјРµРЅРЅРѕ РЅР°СЂСѓС€РµРЅРѕ (РїСѓСЃС‚Р°СЏ СЃС‚СЂРѕРєР° РµСЃР»Рё APPROVED)"
}

DENIED вЂ” СЏРІРЅРѕРµ РЅР°СЂСѓС€РµРЅРёРµ РёР· СЃРїРёСЃРєР° РІС‹С€Рµ
MANUAL вЂ” СЃРѕРјРЅРёС‚РµР»СЊРЅРѕ, РЅСѓР¶РЅР° СЂСѓС‡РЅР°СЏ РїСЂРѕРІРµСЂРєР°
APPROVED вЂ” С„РѕС‚Рѕ РїРѕРґС…РѕРґРёС‚`

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
		reasons = append(reasons, "РўРµРєСЃС‚: "+strings.TrimSpace(textResult.Reason))
	}
	if visionDecision == "MANUAL" && strings.TrimSpace(visionReason) != "" {
		reasons = append(reasons, "Р¤РѕС‚Рѕ: "+strings.TrimSpace(visionReason))
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
	if err := s.repo.SetModerationState(ctx, productID, state, reason); err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]any{"state": state, "reason": reason})
	actorRole := "AI"
	s.repo.InsertModerationAudit(ctx, nil, &actorRole, "product", int64(productID), "AI_MODERATION_DECISION", payload)
	return nil
}

func (s *ModerationService) moderateText(ctx context.Context, name string, description *string, categoryName, subcategoryName string, price int32) (*TextModerationResult, error) {
	userPrompt := fmt.Sprintf("РљР°С‚РµРіРѕСЂРёСЏ: %s\nРџРѕРґРєР°С‚РµРіРѕСЂРёСЏ: %s\nР¦РµРЅР°: %d СЂСѓР±.\nРќР°Р·РІР°РЅРёРµ: %s\nРћРїРёСЃР°РЅРёРµ: %s",
		categoryName,
		subcategoryName,
		price,
		name,
		descriptionOrDefault(description, "РЅРµ СѓРєР°Р·Р°РЅРѕ"),
	)
	userPrompt += "\n\nВажно: низкая цена сама по себе не является нарушением. Нельзя отклонять объявление только из-за цены без явных дополнительных признаков мошенничества, контактов, ссылок, предоплаты или спама."
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
	parsed = normalizePriceOnlyModerationDecision(parsed)
	return &parsed, nil
}

func normalizePriceOnlyModerationDecision(in TextModerationResult) TextModerationResult {
	if !isPriceOnlyModerationDecision(in) {
		return in
	}

	in.Category = "APPROVED"
	in.Reason = ""
	in.Details.Categorization = "OK"
	in.Details.Spam = "OK"
	in.Details.Fraud = "OK"
	in.Details.Contacts = "OK"
	return in
}

func isPriceOnlyModerationDecision(in TextModerationResult) bool {
	if strings.EqualFold(strings.TrimSpace(in.Category), "APPROVED") {
		return false
	}

	reason := strings.ToLower(strings.TrimSpace(in.Reason))
	if reason == "" {
		return false
	}

	priceMarkers := []string{"низк", "дешев", "занижен", "цена", "стоим", "рын", "price", "cheap"}
	hasPriceMarker := false
	for _, marker := range priceMarkers {
		if strings.Contains(reason, marker) {
			hasPriceMarker = true
			break
		}
	}
	if !hasPriceMarker {
		return false
	}

	explicitViolationMarkers := []string{
		"предоплат", "аванс", "перевод", "безнал", "обман", "мошенн",
		"telegram", "whatsapp", "телефон", "email", "почт", "ссылка",
		"http", "t.me", "wa.me", "@", "контакт", "мат", "оскорб", "спам",
	}
	for _, marker := range explicitViolationMarkers {
		if strings.Contains(reason, marker) {
			return false
		}
	}

	if strings.EqualFold(strings.TrimSpace(in.Details.Contacts), "VIOLATION") {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(in.Details.Spam), "VIOLATION") {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(in.Details.Categorization), "VIOLATION") {
		return false
	}

	return true
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
			return nil, &AppError{404, "РўРѕРІР°СЂ РЅРµ РЅР°Р№РґРµРЅ"}
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

func (s *ModerationAdminService) AddAppeal(ctx context.Context, userID, productID int32, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return &AppError{400, "РџСЂРёС‡РёРЅР° Р°РїРµР»Р»СЏС†РёРё РѕР±СЏР·Р°С‚РµР»СЊРЅР°"}
	}
	if err := s.repo.CreateAppeal(ctx, productID, userID, reason); err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]any{"reason": reason})
	uid := userID
	role := "USER"
	s.repo.InsertModerationAudit(ctx, &uid, &role, "appeal", int64(productID), "APPEAL_CREATED", payload)
	return nil
}

func (s *ModerationAdminService) Appeals(ctx context.Context, me *domain.UserEntity, onlyMine bool) ([]repository.ModerationAppealRow, error) {
	if onlyMine {
		uid := me.ID
		return s.repo.ListAppeals(ctx, &uid)
	}
	return s.repo.ListAppeals(ctx, nil)
}

func (s *ModerationAdminService) ReviewAppeal(ctx context.Context, me *domain.UserEntity, appealID int64, status string, comment *string) error {
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != "APPROVED" && status != "REJECTED" {
		return &AppError{400, "status РґРѕР»Р¶РµРЅ Р±С‹С‚СЊ APPROVED РёР»Рё REJECTED"}
	}
	if err := s.repo.ReviewAppeal(ctx, appealID, me.ID, status, comment); err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]any{"status": status, "comment": comment})
	role := ""
	if me.RoleName != nil {
		role = *me.RoleName
	}
	uid := me.ID
	s.repo.InsertModerationAudit(ctx, &uid, &role, "appeal", appealID, "APPEAL_REVIEWED", payload)
	return nil
}

func (s *ModerationAdminService) AuditLogs(ctx context.Context, limit int) ([]map[string]any, error) {
	rows, err := s.repo.ListModerationAudit(ctx, limit)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]any{
			"id":          r.ID,
			"actorUserId": r.ActorUserID,
			"actorRole":   r.ActorRole,
			"targetType":  r.TargetType,
			"targetId":    r.TargetID,
			"action":      r.Action,
			"payload":     json.RawMessage(r.Payload),
			"createdAt":   r.CreatedAt,
		})
	}
	return out, nil
}

func (s *ModerationAdminService) Summary(ctx context.Context, days int) (map[string]any, error) {
	denied, approvedAI, appealsOpen, err := s.repo.ModerationSummary(ctx, days)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"days":        days,
		"denied":      denied,
		"approvedAI":  approvedAI,
		"appealsOpen": appealsOpen,
		"topComplaints": []map[string]any{
			{"type": "appeals_open", "count": appealsOpen},
			{"type": "denied_ads", "count": denied},
		},
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
