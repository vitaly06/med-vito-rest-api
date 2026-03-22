package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"unicode/utf8"

	"med-vito/api-go/internal/repository"
)

var (
	validTicketThemes = map[string]bool{
		"TECHNICAL_ISSUE": true, "ACCOUNT_PROBLEM": true, "PAYMENT_ISSUE": true,
		"PRODUCT_QUESTION": true, "COMPLAINT": true, "SUGGESTION": true, "OTHER": true,
	}
	validTicketPriority = map[string]bool{
		"LOW": true, "MEDIUM": true, "HIGH": true, "URGENT": true,
	}
	validTicketStatus = map[string]bool{
		"OPEN": true, "IN_PROGRESS": true, "RESOLVED": true, "CLOSED": true,
	}
)

type SupportService struct {
	repo *repository.SupportPG
}

func NewSupportService(repo *repository.SupportPG) *SupportService {
	return &SupportService{repo: repo}
}

func (s *SupportService) CreateTicket(ctx context.Context, userID int32, theme, subject, message, priority string) (*repository.SupportTicketFull, error) {
	theme = strings.TrimSpace(theme)
	if !validTicketThemes[theme] {
		return nil, &AppError{400, "Некорректная тема обращения"}
	}
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return nil, &AppError{400, "Заголовок обязателен"}
	}
	if utf8.RuneCountInString(subject) > 200 {
		return nil, &AppError{400, "Заголовок не должен превышать 200 символов"}
	}
	message = strings.TrimSpace(message)
	if message == "" {
		return nil, &AppError{400, "Сообщение обязательно"}
	}
	if utf8.RuneCountInString(message) > 2000 {
		return nil, &AppError{400, "Сообщение не должно превышать 2000 символов"}
	}
	priority = strings.TrimSpace(priority)
	if priority == "" {
		priority = "MEDIUM"
	}
	if !validTicketPriority[priority] {
		return nil, &AppError{400, "Некорректный приоритет"}
	}
	return s.repo.CreateTicket(ctx, userID, theme, subject, priority, message)
}

type TicketWithMessages struct {
	repository.SupportTicketFull
	Messages []repository.SupportMessageOut `json:"messages"`
}

func (s *SupportService) GetTicket(ctx context.Context, ticketID, userID int32, isModerator bool) (*TicketWithMessages, error) {
	t, err := s.repo.GetTicketFull(ctx, ticketID)
	if errors.Is(err, repository.ErrSupportTicketNotFound) {
		return nil, &AppError{404, "Тикет не найден"}
	}
	if err != nil {
		return nil, err
	}
	if !isModerator && t.UserID != userID {
		return nil, &AppError{403, "Нет доступа к этому тикету"}
	}
	msgs, err := s.repo.ListMessagesForTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	return &TicketWithMessages{SupportTicketFull: *t, Messages: msgs}, nil
}

func (s *SupportService) GetUserTickets(ctx context.Context, userID int32, page, limit int, status, priority, theme *string, search *string) (map[string]any, error) {
	if err := validateOptionalFilters(status, priority, theme); err != nil {
		return nil, err
	}
	list, total, err := s.repo.ListTickets(ctx, &userID, status, priority, theme, search, false, page, limit)
	if err != nil {
		return nil, err
	}
	return listTicketsResponse(list, total, page, limit), nil
}

func (s *SupportService) GetAllTickets(ctx context.Context, page, limit int, status, priority, theme *string, search *string) (map[string]any, error) {
	if err := validateOptionalFilters(status, priority, theme); err != nil {
		return nil, err
	}
	list, total, err := s.repo.ListTickets(ctx, nil, status, priority, theme, search, true, page, limit)
	if err != nil {
		return nil, err
	}
	return listTicketsResponse(list, total, page, limit), nil
}

func validateOptionalFilters(status, priority, theme *string) error {
	if status != nil && *status != "" && !validTicketStatus[*status] {
		return &AppError{400, "Некорректный статус тикета"}
	}
	if priority != nil && *priority != "" && !validTicketPriority[*priority] {
		return &AppError{400, "Некорректный приоритет"}
	}
	if theme != nil && *theme != "" && !validTicketThemes[*theme] {
		return &AppError{400, "Некорректная тема"}
	}
	return nil
}

func listTicketsResponse(list []repository.SupportTicketListItem, total, page, limit int) map[string]any {
	if limit < 1 {
		limit = 10
	}
	pages := int(math.Ceil(float64(total) / float64(limit)))
	return map[string]any{
		"tickets": list,
		"pagination": map[string]any{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": pages,
		},
	}
}

func (s *SupportService) SendMessage(ctx context.Context, ticketID, userID int32, text string, isModerator bool) (*repository.SupportMessageOut, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, &AppError{400, "Сообщение обязательно"}
	}
	if utf8.RuneCountInString(text) > 2000 {
		return nil, &AppError{400, "Сообщение не должно превышать 2000 символов"}
	}
	ownerID, st, err := s.repo.GetTicketOwnerID(ctx, ticketID)
	if errors.Is(err, repository.ErrSupportTicketNotFound) {
		return nil, &AppError{404, "Тикет не найден"}
	}
	if err != nil {
		return nil, err
	}
	if !isModerator && ownerID != userID {
		return nil, &AppError{403, "Нет доступа к этому тикету"}
	}
	if st == "CLOSED" {
		return nil, &AppError{400, "Нельзя отправлять сообщения в закрытый тикет"}
	}
	return s.repo.SendMessageTx(ctx, ticketID, userID, text, isModerator, st)
}

func (s *SupportService) UpdateTicket(ctx context.Context, ticketID, moderatorID int32, status, priority *string) (*repository.SupportTicketFull, error) {
	if status != nil && *status != "" && !validTicketStatus[*status] {
		return nil, &AppError{400, "Некорректный статус тикета"}
	}
	if priority != nil && *priority != "" && !validTicketPriority[*priority] {
		return nil, &AppError{400, "Некорректный приоритет"}
	}
	var stPtr, prPtr *string
	if status != nil && *status != "" {
		stPtr = status
	}
	if priority != nil && *priority != "" {
		prPtr = priority
	}
	t, err := s.repo.UpdateTicket(ctx, ticketID, moderatorID, stPtr, prPtr)
	if errors.Is(err, repository.ErrSupportTicketNotFound) {
		return nil, &AppError{404, "Тикет не найден"}
	}
	return t, err
}

func (s *SupportService) AssignTicket(ctx context.Context, ticketID, moderatorID int32) (*repository.SupportTicketFull, error) {
	t, err := s.repo.AssignTicket(ctx, ticketID, moderatorID)
	if errors.Is(err, repository.ErrSupportTicketNotFound) {
		return nil, &AppError{404, "Тикет не найден"}
	}
	return t, err
}

func (s *SupportService) TicketStats(ctx context.Context) (*repository.SupportStats, error) {
	return s.repo.TicketStats(ctx)
}
