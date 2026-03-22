package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"med-vito/api-go/internal/repository"
)

type ChatService struct {
	repo *repository.ChatPG
}

func NewChatService(repo *repository.ChatPG) *ChatService {
	return &ChatService{repo: repo}
}

func formatChatShortDate(t time.Time) string {
	lt := t.In(time.Local)
	return fmt.Sprintf("%02d.%02d.%02d", lt.Day(), int(lt.Month()), lt.Year()%100)
}

func messageTimeString(t time.Time) string {
	return t.In(time.Local).Format("15:04")
}

func (s *ChatService) formatStartChat(c *repository.ChatFullRow) map[string]any {
	var product any
	if c.JoinProductID != nil {
		product = map[string]any{
			"id":     *c.JoinProductID,
			"name":   c.ProductName,
			"images": append([]string{}, c.ProductImages...),
			"price":  c.ProductPrice,
		}
	} else {
		product = nil
	}
	return map[string]any{
		"id":               c.ID,
		"productId":        c.ProductID,
		"buyerId":          c.BuyerID,
		"sellerId":         c.SellerID,
		"isModerationChat": c.IsModerationChat,
		"lastMessageAt":    c.LastMessageAt,
		"createdAt":        c.CreatedAt,
		"product":          product,
		"buyer":            map[string]any{"id": c.BuyerID, "fullName": c.BuyerName, "phoneNumber": c.BuyerPhone},
		"seller":           map[string]any{"id": c.SellerID, "fullName": c.SellerName, "phoneNumber": c.SellerPhone},
	}
}

func (s *ChatService) StartChat(ctx context.Context, productID, buyerID int32) (map[string]any, error) {
	sellerID, ok, err := s.repo.ProductSellerID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &AppError{404, "Товар не найден"}
	}
	if sellerID == buyerID {
		return nil, &AppError{400, "Нельзя писать самому себе"}
	}
	existing, err := s.repo.FindChatByProductParticipants(ctx, productID, buyerID, sellerID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return s.formatStartChat(existing), nil
	}
	newID, err := s.repo.InsertProductChat(ctx, productID, buyerID, sellerID)
	if err != nil {
		return nil, err
	}
	full, err := s.repo.GetChatFull(ctx, newID)
	if err != nil {
		return nil, err
	}
	return s.formatStartChat(full), nil
}

func (s *ChatService) GetUserChats(ctx context.Context, userID int32) ([]map[string]any, error) {
	rows, err := s.repo.ListChatsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, chat := range rows {
		isBuyer := chat.BuyerID == userID
		var compName, compPhone string
		var compID int32
		if isBuyer {
			compID, compName, compPhone = chat.SellerID, chat.SellerName, chat.SellerPhone
		} else {
			compID, compName, compPhone = chat.BuyerID, chat.BuyerName, chat.BuyerPhone
		}
		unread := chat.UnreadSeller
		if isBuyer {
			unread = chat.UnreadBuyer
		}
		var product any
		if chat.JoinProductID != nil {
			var img any
			if len(chat.ProductImages) > 0 {
				img = chat.ProductImages[0]
			} else {
				img = nil
			}
			product = map[string]any{
				"id":    *chat.JoinProductID,
				"name":  chat.ProductName,
				"price": chat.ProductPrice,
				"image": img,
			}
		} else {
			product = nil
		}
		var lastMsg any
		if chat.LMID != nil && chat.LMContent != nil && chat.LMCreated != nil && chat.LMSender != nil && chat.LMRead != nil {
			lastMsg = map[string]any{
				"content":       *chat.LMContent,
				"createdAt":     *chat.LMCreated,
				"formattedDate": formatChatShortDate(*chat.LMCreated),
				"isFromMe":      *chat.LMSender == userID,
				"isRead":        *chat.LMRead,
			}
		} else {
			lastMsg = nil
		}
		la := chat.LastMessageAt
		if la.IsZero() {
			la = chat.CreatedAt
		}
		out = append(out, map[string]any{
			"id":               chat.ID,
			"isModerationChat": chat.IsModerationChat,
			"product":          product,
			"companion": map[string]any{
				"id": compID, "fullName": compName, "phoneNumber": compPhone, "avatar": nil,
			},
			"lastMessage":  lastMsg,
			"unreadCount":  unread,
			"lastActivity": formatChatShortDate(la),
			"createdAt":    chat.CreatedAt,
		})
	}
	return out, nil
}

func (s *ChatService) GetChatInfo(ctx context.Context, chatID, userID int32) (map[string]any, error) {
	c, err := s.repo.GetChatFull(ctx, chatID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "Чат не найден"}
	}
	if err != nil {
		return nil, err
	}
	if c.BuyerID != userID && c.SellerID != userID {
		return nil, &AppError{403, "Нет доступа к этому чату"}
	}
	isBuyer := c.BuyerID == userID
	var comp map[string]any
	if isBuyer {
		comp = map[string]any{"id": c.SellerID, "fullName": c.SellerName, "phoneNumber": c.SellerPhone, "role": "seller"}
	} else {
		comp = map[string]any{"id": c.BuyerID, "fullName": c.BuyerName, "phoneNumber": c.BuyerPhone, "role": "buyer"}
	}
	var product any
	if c.JoinProductID != nil {
		var img any
		if len(c.ProductImages) > 0 {
			img = c.ProductImages[0]
		} else {
			img = nil
		}
		product = map[string]any{
			"id":          *c.JoinProductID,
			"name":        c.ProductName,
			"price":       c.ProductPrice,
			"image":       img,
			"description": c.ProductDesc,
		}
	} else {
		product = nil
	}
	ub, us, err := s.repo.ChatUnreadCounts(ctx, chatID)
	if err != nil {
		return nil, err
	}
	unreadCnt := us
	if isBuyer {
		unreadCnt = ub
	}
	return map[string]any{
		"id":               c.ID,
		"isModerationChat": c.IsModerationChat,
		"product":          product,
		"companion":        comp,
		"isUserBuyer":      isBuyer,
		"unreadCount":      unreadCnt,
		"createdAt":        c.CreatedAt,
	}, nil
}

// GetChatMessages — как Nest getChatMessages: DESC из БД, затем reverse; page/limit по умолчанию 1 и 50.
func (s *ChatService) GetChatMessages(ctx context.Context, chatID, userID int32, page, limit int) (map[string]any, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	buyerID, sellerID, err := s.repo.ChatParticipants(ctx, chatID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "Чат не найден"}
	}
	if err != nil {
		return nil, err
	}
	if buyerID != userID && sellerID != userID {
		return nil, &AppError{403, "Нет доступа к этому чату"}
	}
	offset := (page - 1) * limit
	rows, total, err := s.repo.ListMessagesPage(ctx, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	slices.Reverse(rows)
	msgs := make([]map[string]any, 0, len(rows))
	for _, m := range rows {
		var rel any
		if m.RelProductID != nil {
			rel = map[string]any{
				"id":            *m.RelProductID,
				"name":          m.RelName,
				"images":        m.RelImages,
				"price":         m.RelPrice,
				"moderateState": m.RelModerate,
			}
		} else {
			rel = nil
		}
		msgs = append(msgs, map[string]any{
			"id":             m.ID,
			"content":        m.Content,
			"senderId":       m.SenderID,
			"sender":         map[string]any{"id": m.SenderID, "fullName": m.SenderName},
			"relatedProduct": rel,
			"isFromMe":       m.SenderID == userID,
			"isRead":         m.IsRead,
			"readAt":         m.ReadAt,
			"createdAt":      m.CreatedAt,
			"timeString":     messageTimeString(m.CreatedAt),
		})
	}
	var pages int32
	if limit > 0 {
		pages = int32(math.Ceil(float64(total) / float64(limit)))
	}
	return map[string]any{
		"messages": msgs,
		"pagination": map[string]any{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": pages,
		},
	}, nil
}

// SendMessage — как Nest chat.sendMessage (для REST при необходимости и для Socket.IO).
func (s *ChatService) SendMessage(ctx context.Context, chatID, senderID int32, content string) (map[string]any, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, &AppError{400, "Пустое сообщение"}
	}
	msgID, senderName, createdAt, _, _, err := s.repo.InsertChatMessage(ctx, chatID, senderID, content)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, &AppError{404, "Чат не найден"}
	}
	if errors.Is(err, repository.ErrForbiddenChat) {
		return nil, &AppError{403, "Нет доступа к этому чату"}
	}
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id":         msgID,
		"content":    content,
		"senderId":   senderID,
		"sender":     map[string]any{"id": senderID, "fullName": senderName},
		"isRead":     false,
		"createdAt":  createdAt,
		"timeString": messageTimeString(createdAt),
	}, nil
}

// MarkMessagesAsRead — как Nest markMessagesAsRead.
func (s *ChatService) MarkMessagesAsRead(ctx context.Context, chatID, userID int32) error {
	err := s.repo.MarkChatMessagesRead(ctx, chatID, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return &AppError{404, "Чат не найден"}
	}
	if errors.Is(err, repository.ErrForbiddenChat) {
		return &AppError{403, "Нет доступа к этому чату"}
	}
	return err
}
