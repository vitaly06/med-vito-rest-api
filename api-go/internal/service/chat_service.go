package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"med-vito/api-go/internal/config"
	"med-vito/api-go/internal/pkg/mail"
	"med-vito/api-go/internal/repository"
)

type ChatService struct {
	repo  *repository.ChatPG
	users *repository.UserPG
	cfg   config.Config
}

func NewChatService(repo *repository.ChatPG, users *repository.UserPG, cfg config.Config) *ChatService {
	return &ChatService{repo: repo, users: users, cfg: cfg}
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

func (s *ChatService) StartDirectChat(ctx context.Context, sellerID, buyerID int32) (map[string]any, error) {
	if sellerID <= 0 {
		return nil, &AppError{400, "Нужен sellerId"}
	}
	if sellerID == buyerID {
		return nil, &AppError{400, "Нельзя писать самому себе"}
	}
	existing, err := s.repo.FindDirectChatBetweenUsers(ctx, buyerID, sellerID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return s.formatStartChat(existing), nil
	}
	newID, err := s.repo.InsertDirectChat(ctx, buyerID, sellerID)
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
		if chat.IsModerationChat {
			compName = "Модерация платформы"
			compPhone = ""
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
			product = map[string]any{
				"id":    compID,
				"name":  "Чат с " + compName,
				"price": 0,
				"image": nil,
			}
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
	if c.IsModerationChat {
		comp["fullName"] = "Модерация платформы"
		comp["phoneNumber"] = ""
		comp["role"] = "moderation"
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

	s.sendAsyncChatMessageNotification(chatID, senderID, senderName, content)

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

func (s *ChatService) sendAsyncChatMessageNotification(chatID, senderID int32, senderName, content string) {
	if s == nil || s.users == nil {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		full, err := s.repo.GetChatFull(ctx, chatID)
		if err != nil || full == nil {
			return
		}

		recipientID := full.SellerID
		if senderID == full.SellerID {
			recipientID = full.BuyerID
		}

		recipient, err := s.users.FindUserByID(ctx, recipientID)
		if err != nil || recipient == nil {
			return
		}

		email := strings.TrimSpace(strings.ToLower(recipient.Email))
		if email == "" || strings.HasSuffix(email, "@oauth.local") || !strings.Contains(email, "@") {
			return
		}

		productTitle := "объявлению"
		if full.ProductName != nil && strings.TrimSpace(*full.ProductName) != "" {
			productTitle = fmt.Sprintf("товару «%s»", strings.TrimSpace(*full.ProductName))
		}

		subject := fmt.Sprintf("Новое сообщение от %s", senderName)
		baseUrl := strings.TrimRight(s.cfg.BaseURL, "/")
		if baseUrl == "" {
			baseUrl = "http://localhost:3000"
		}
		body := fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; padding: 20px; color: #333;">
				<h2 style="color: #2563eb;">Новое сообщение в чате</h2>
				<p>Здравствуйте, <b>%s</b>!</p>
				<p>Пользователь <b>%s</b> прислал вам сообщение по %s:</p>
				<blockquote style="border-left: 4px solid #2563eb; margin: 16px 0; padding: 10px 16px; background-color: #f3f4f6; border-radius: 4px;">
					%s
				</blockquote>
				<p style="margin-top: 24px;">
					<a href="%s/profile/messages/%d" style="background-color: #2563eb; color: #ffffff; text-decoration: none; padding: 10px 18px; border-radius: 6px; display: inline-block; font-weight: bold;">
						Перейти к переписке
					</a>
				</p>
			</div>
		`, recipient.FullName, senderName, productTitle, content, baseUrl, chatID)

		_ = mail.SendHTMLSmart(
			s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword,
			s.cfg.SMTPFrom, email, subject, body,
			s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure,
		)
	}()
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
