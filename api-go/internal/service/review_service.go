package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"med-vito/api-go/internal/config"
	mailpkg "med-vito/api-go/internal/pkg/mail"
	"med-vito/api-go/internal/repository"
)

var reviewAllowedRatings = []float64{1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5}

const reviewAIPrompt = `Ты модератор отзывов. Верни JSON:
{
  "decision": "APPROVED|MANUAL|DENIED",
  "reason": "краткая причина",
  "flags": {
    "profanity": "OK|SUSPICIOUS|VIOLATION",
    "spam": "OK|SUSPICIOUS|VIOLATION",
    "fake_template": "OK|SUSPICIOUS|VIOLATION"
  }
}
Правила:
- Явная нецензурщина, оскорбления, угрозы => DENIED.
- Спам/мусорный текст/реклама контактов => DENIED.
- Шаблонный фейковый отзыв без фактов => MANUAL.
- Нормальный отзыв => APPROVED.
Ответ строго JSON без markdown.`

func isAllowedReviewRating(r float64) bool {
	for _, a := range reviewAllowedRatings {
		if math.Abs(r-a) < 1e-6 {
			return true
		}
	}
	return false
}

type ReviewService struct {
	repo       *repository.ReviewPG
	chat       *repository.ChatPG
	cfg        config.Config
	httpClient *http.Client
}

func NewReviewService(repo *repository.ReviewPG, chat *repository.ChatPG, cfg config.Config) *ReviewService {
	return &ReviewService{
		repo:       repo,
		chat:       chat,
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

func (s *ReviewService) SendReview(ctx context.Context, authorID int32, reviewedUserID int32, rating float64, text *string) (map[string]any, error) {
	if reviewedUserID <= 0 {
		return nil, &AppError{400, "Id оцениваемого продавца должен быть положительным числом"}
	}
	if !isAllowedReviewRating(rating) {
		return nil, &AppError{400, "Рейтинг должен быть одним из: 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5"}
	}
	ok, err := s.repo.UserExists(ctx, reviewedUserID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &AppError{400, "Продавец с таким id не найден"}
	}
	if authorID == reviewedUserID {
		return nil, &AppError{400, "Нельзя оставить отзыв самому себе"}
	}
	eligible, err := s.repo.HasDealOrReservation(ctx, authorID, reviewedUserID)
	if err != nil {
		return nil, err
	}
	if !eligible {
		return nil, &AppError{403, "Оставить отзыв можно только после покупки или резерва"}
	}
	reviewID, err := s.repo.InsertReview(ctx, authorID, reviewedUserID, rating, text)
	if err != nil {
		if repository.IsPgUniqueViolation(err) {
			return nil, &AppError{400, "Вы уже оставили отзыв этому пользователю"}
		}
		return nil, err
	}

	s.notifyReviewStatusAsync(ctx, reviewID, "MODERATE")
	s.notifyReviewInChatAsync(ctx, reviewID)
	return map[string]any{"message": "Отзыв успешно оставлен и отправлен на модерацию"}, nil
}

func (s *ReviewService) GetUserReviews(ctx context.Context, userID int32) (map[string]any, error) {
	rows, err := s.repo.ListApprovedForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var sum float64
	for _, r := range rows {
		sum += r.Rating
	}
	var total float64
	if len(rows) > 0 {
		total = math.Round((sum/float64(len(rows)))*100) / 100
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		lt := r.CreatedAt.In(time.Local)
		list = append(list, map[string]any{
			"rating":   r.Rating,
			"text":     r.Text,
			"date":     fmt.Sprintf("%02d.%02d.%02d", lt.Day(), int(lt.Month()), lt.Year()%100),
			"fullName": r.FullName,
		})
	}
	return map[string]any{
		"totalRating":  total,
		"reviewsCount": len(rows),
		"reviews":      list,
	}, nil
}

func (s *ReviewService) ModerateReview(ctx context.Context, reviewID int32, status string) (map[string]any, error) {
	st := strings.TrimSpace(strings.ToUpper(status))
	if st != "APPROVED" && st != "DENIDED" {
		return nil, &AppError{400, "Неверный статус модерации. Доступные статусы: APPROVED, DENIDED"}
	}
	if err := s.repo.SetReviewModeration(ctx, reviewID, st); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Отзыв для модерации не найден"}
		}
		return nil, err
	}
	s.notifyReviewStatusAsync(ctx, reviewID, st)
	s.notifyReviewStatusInChatAsync(ctx, reviewID, st)
	if st == "APPROVED" {
		return map[string]any{"message": "Отзыв успешно опубликован"}, nil
	}
	return map[string]any{"message": "Отзыв успешно отклонен"}, nil
}

func (s *ReviewService) AllReviewsToModerate(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.repo.ListModerateQueue(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]any{
			"id":        r.ID,
			"rating":    r.Rating,
			"text":      r.Text,
			"createdAt": r.CreatedAt,
			"reviewedBy": map[string]any{
				"id":       r.ReviewedByID,
				"fullName": r.ReviewedByName,
			},
			"reviewedUser": map[string]any{
				"id":       r.ReviewedUserID,
				"fullName": r.ReviewedName,
			},
		})
	}
	return out, nil
}

func (s *ReviewService) StartAIModerationWorker(ctx context.Context) {
	if strings.TrimSpace(s.cfg.YandexAPIKey) == "" || strings.TrimSpace(s.cfg.YandexFolderID) == "" {
		return
	}
	t := time.NewTicker(30 * time.Second)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				_ = s.runReviewAIPoll(ctx)
			}
		}
	}()
}

func (s *ReviewService) runReviewAIPoll(ctx context.Context) error {
	items, err := s.repo.ListPendingAIModeration(ctx, 20)
	if err != nil {
		return err
	}
	for _, it := range items {
		decision, _, err := s.aiReviewDecision(ctx, it.Text)
		if err != nil {
			continue
		}
		switch decision {
		case "APPROVED":
			_ = s.repo.SetReviewModeration(ctx, it.ID, "APPROVED")
			s.notifyReviewStatusAsync(ctx, it.ID, "APPROVED")
			s.notifyReviewStatusInChatAsync(ctx, it.ID, "APPROVED")
		case "DENIED":
			_ = s.repo.SetReviewModeration(ctx, it.ID, "DENIDED")
			s.notifyReviewStatusAsync(ctx, it.ID, "DENIDED")
			s.notifyReviewStatusInChatAsync(ctx, it.ID, "DENIDED")
		default:
			// Остается в ручной очереди MODERATE.
		}
	}
	return nil
}

func (s *ReviewService) aiReviewDecision(ctx context.Context, text *string) (string, string, error) {
	payload := map[string]any{
		"modelUri": fmt.Sprintf("gpt://%s/yandexgpt/latest", s.cfg.YandexFolderID),
		"completionOptions": map[string]any{
			"stream":      false,
			"temperature": 0.05,
			"maxTokens":   200,
		},
		"messages": []map[string]string{
			{"role": "system", "text": reviewAIPrompt},
			{"role": "user", "text": "Отзыв: " + valueOrEmpty(text)},
		},
	}
	raw, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewReader(raw))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Api-Key "+s.cfg.YandexAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("ai status %d", resp.StatusCode)
	}
	var r struct {
		Result struct {
			Alternatives []struct {
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"alternatives"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", "", err
	}
	if len(r.Result.Alternatives) == 0 {
		return "", "", errors.New("no alternatives")
	}
	clean := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(r.Result.Alternatives[0].Message.Text, "```json", ""), "```", ""))
	var parsed struct {
		Decision string `json:"decision"`
		Reason   string `json:"reason"`
	}
	if err := json.Unmarshal([]byte(clean), &parsed); err != nil {
		return "", "", err
	}
	d := strings.ToUpper(strings.TrimSpace(parsed.Decision))
	if d != "APPROVED" && d != "DENIED" && d != "MANUAL" {
		return "", "", errors.New("invalid decision")
	}
	return d, parsed.Reason, nil
}

func (s *ReviewService) notifyReviewStatusAsync(ctx context.Context, reviewID int32, status string) {
	authorName, authorEmail, authorPhone, sellerName, sellerEmail, sellerPhone, err := s.repo.GetReviewParties(ctx, reviewID)
	if err != nil {
		return
	}
	subjectAuthor := "Статус вашего отзыва"
	subjectSeller := "Новый/обновленный отзыв"
	bodyAuthor := fmt.Sprintf("<p>Здравствуйте, %s.</p><p>Статус отзыва: <b>%s</b>.</p>", authorName, status)
	bodySeller := fmt.Sprintf("<p>По вашему профилю обновился статус отзыва: <b>%s</b>.</p>", status)
	authorSMS := "Статус вашего отзыва: " + status
	sellerSMS := "По вашему профилю обновился статус отзыва: " + status
	if status == "MODERATE" {
		subjectAuthor = "Ваш отзыв отправлен на модерацию"
		bodyAuthor = fmt.Sprintf("<p>Здравствуйте, %s.</p><p>Ваш отзыв принят и отправлен на модерацию.</p>", authorName)
		subjectSeller = "Новый отзыв ожидает модерацию"
		bodySeller = fmt.Sprintf("<p>Здравствуйте, %s.</p><p>По вашему профилю оставили новый отзыв. После модерации он появится публично.</p>", sellerName)
		authorSMS = "Ваш отзыв отправлен на модерацию."
		sellerSMS = "Новый отзыв о вас ожидает модерацию."
	}
	go func() {
		if strings.TrimSpace(s.cfg.SMTPHost) != "" {
			_ = mailpkg.SendHTMLSmart(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPFrom, authorEmail, subjectAuthor, bodyAuthor, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure)
			_ = mailpkg.SendHTMLSmart(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPFrom, sellerEmail, subjectSeller, bodySeller, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure)
		}
		s.sendSMS(authorPhone, authorSMS)
		s.sendSMS(sellerPhone, sellerSMS)
	}()
}

func (s *ReviewService) notifyReviewInChatAsync(ctx context.Context, reviewID int32) {
	if s.chat == nil {
		return
	}
	meta, err := s.repo.GetReviewMeta(ctx, reviewID)
	if err != nil {
		return
	}
	go func() {
		chat, err := s.chat.FindDirectChatBetweenUsers(context.Background(), meta.ReviewedByID, meta.ReviewedUserID)
		if err != nil {
			log.Printf("review chat find failed review=%d: %v", reviewID, err)
			return
		}
		var chatID int32
		if chat == nil {
			chatID, err = s.chat.InsertDirectChat(context.Background(), meta.ReviewedByID, meta.ReviewedUserID)
			if err != nil {
				log.Printf("review chat create failed review=%d: %v", reviewID, err)
				return
			}
		} else {
			chatID = chat.ID
		}
		content := fmt.Sprintf("Оставлен новый отзыв (%0.1f★) пользователю \"%s\". Отзыв отправлен на модерацию.", meta.Rating, meta.ReviewedName)
		if _, _, _, _, _, err := s.chat.InsertChatMessage(context.Background(), chatID, meta.ReviewedByID, content); err != nil {
			log.Printf("review chat notify failed review=%d chat=%d: %v", reviewID, chatID, err)
		}
	}()
}

func (s *ReviewService) notifyReviewStatusInChatAsync(ctx context.Context, reviewID int32, status string) {
	if s.chat == nil {
		return
	}
	meta, err := s.repo.GetReviewMeta(ctx, reviewID)
	if err != nil {
		return
	}
	go func() {
		chat, err := s.chat.FindDirectChatBetweenUsers(context.Background(), meta.ReviewedByID, meta.ReviewedUserID)
		if err != nil || chat == nil {
			return
		}
		content := fmt.Sprintf("Статус отзыва (%0.1f★) обновлен: %s.", meta.Rating, status)
		if _, _, _, _, _, err := s.chat.InsertChatMessage(context.Background(), chat.ID, meta.ReviewedUserID, content); err != nil {
			log.Printf("review chat status failed review=%d chat=%d: %v", reviewID, chat.ID, err)
		}
	}()
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *ReviewService) CreateAppeal(ctx context.Context, userID, reviewID int32, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return &AppError{400, "Причина апелляции обязательна"}
	}
	if err := s.repo.CreateReviewAppeal(ctx, reviewID, userID, reason); err != nil {
		return err
	}
	s.notifyAppealInChatAsync(ctx, userID, reviewID, reason)
	return nil
}

func (s *ReviewService) MyAppeals(ctx context.Context, userID int32) ([]repository.ReviewAppealRow, error) {
	uid := userID
	return s.repo.ListReviewAppeals(ctx, &uid)
}

func (s *ReviewService) AllAppeals(ctx context.Context) ([]repository.ReviewAppealRow, error) {
	return s.repo.ListReviewAppeals(ctx, nil)
}

func (s *ReviewService) ResolveAppeal(ctx context.Context, moderatorID int32, appealID int64, status string, note *string) error {
	st := strings.ToUpper(strings.TrimSpace(status))
	if st != "APPROVED" && st != "REJECTED" {
		return &AppError{400, "status должен быть APPROVED или REJECTED"}
	}
	return s.repo.ResolveReviewAppeal(ctx, appealID, moderatorID, st, note)
}

func (s *ReviewService) notifyAppealInChatAsync(ctx context.Context, userID, reviewID int32, reason string) {
	if s.chat == nil {
		return
	}
	meta, err := s.repo.GetReviewMeta(ctx, reviewID)
	if err != nil {
		return
	}
	go func() {
		chat, err := s.chat.FindDirectChatBetweenUsers(context.Background(), meta.ReviewedByID, meta.ReviewedUserID)
		if err != nil || chat == nil {
			return
		}
		senderID := meta.ReviewedByID
		if userID == meta.ReviewedUserID {
			senderID = meta.ReviewedUserID
		}
		content := "Подана апелляция по отзыву. Причина: " + strings.TrimSpace(reason)
		if _, _, _, _, _, err := s.chat.InsertChatMessage(context.Background(), chat.ID, senderID, content); err != nil {
			log.Printf("review appeal chat failed review=%d chat=%d: %v", reviewID, chat.ID, err)
		}
	}()
}

func (s *ReviewService) sendSMS(phone, message string) {
	phone = strings.TrimSpace(phone)
	message = strings.TrimSpace(message)
	if phone == "" || message == "" {
		return
	}
	if strings.TrimSpace(s.cfg.MTSBearer) != "" {
		body := map[string]any{
			"submits": []any{map[string]any{"msid": phone, "message": message}},
			"naming":  "Torguisamru",
		}
		if err := s.httpPostJSON(context.Background(), "https://api.mts.ru/client-omni-adapter_production/1.0.2/mcom/messageManagement/messages", body, s.cfg.MTSBearer); err != nil {
			log.Printf("review notify sms(mts) failed: phone=%s err=%v", phone, err)
		}
		return
	}
	if strings.TrimSpace(s.cfg.NotisendAPIKey) != "" {
		u := "https://sms.notisend.ru/api/message/send?project=" + url.QueryEscape(s.cfg.NotisendProject) +
			"&message=" + url.QueryEscape(message) +
			"&recipients=" + url.QueryEscape(phone) +
			"&apikey=" + url.QueryEscape(s.cfg.NotisendAPIKey)
		resp, err := http.Get(u)
		if err != nil {
			log.Printf("review notify sms(notisend) failed: phone=%s err=%v", phone, err)
			return
		}
		_ = resp.Body.Close()
	}
}

func (s *ReviewService) httpPostJSON(ctx context.Context, endpoint string, body any, bearer string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(bearer) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(bearer))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New("sms gateway status " + resp.Status)
	}
	return nil
}
