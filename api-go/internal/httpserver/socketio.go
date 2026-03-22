package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"

	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/service"
)

const (
	chatIONamespace    = "/chat"
	supportIONamespace = "/support"
)

// socketUserCtx — данные после проверки session_id (как Nest AuthenticatedSocket).
type socketUserCtx struct {
	UserID   int32
	FullName string
	RoleName string
}

var (
	engineSocketServer *socketio.Server
	engineSocketMu     sync.RWMutex
)

func setEngineSocketServer(s *socketio.Server) {
	engineSocketMu.Lock()
	engineSocketServer = s
	engineSocketMu.Unlock()
}

func socketCtxFromUser(u *domain.UserEntity) *socketUserCtx {
	rn := ""
	if u.RoleName != nil {
		rn = *u.RoleName
	}
	return &socketUserCtx{UserID: u.ID, FullName: u.FullName, RoleName: rn}
}

func socketCtxUser(c socketio.Conn) *socketUserCtx {
	raw := c.Context()
	if raw == nil {
		return nil
	}
	sc, ok := raw.(*socketUserCtx)
	if !ok || sc == nil {
		return nil
	}
	return sc
}

func socketUserID(c socketio.Conn) int32 {
	sc := socketCtxUser(c)
	if sc == nil {
		return 0
	}
	return sc.UserID
}

func supportModeratorFromCtx(c socketio.Conn) bool {
	sc := socketCtxUser(c)
	if sc == nil {
		return false
	}
	return sc.RoleName == "moderator" || sc.RoleName == "admin"
}

// RegisterSocketIO — один движок /socket.io: неймспейсы /chat и /support (как Nest).
func RegisterSocketIO(app *fiber.App, corsOrigins string, auth *service.AuthService, chat *service.ChatService, sup *service.SupportService) {
	opts := &engineio.Options{
		RequestChecker: socketIOCORSChecker(corsOrigins),
	}
	server := socketio.NewServer(opts)
	setEngineSocketServer(server)

	registerChatNamespace(server, auth, chat)
	registerSupportNamespace(server, auth, sup)

	h := adaptor.HTTPHandler(server)
	app.All("/socket.io", h)
	app.All("/socket.io/*", h)
	log.Print("Socket.IO: /socket.io — namespaces /chat, /support")
}

func socketIOCORSChecker(corsOrigins string) func(*http.Request) (http.Header, error) {
	return func(r *http.Request) (http.Header, error) {
		h := http.Header{}
		origin := r.Header.Get("Origin")
		if origin == "" {
			return h, nil
		}
		if corsOrigins == "" {
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Credentials", "true")
			return h, nil
		}
		for _, o := range strings.Split(corsOrigins, ",") {
			if strings.TrimSpace(o) == origin {
				h.Set("Access-Control-Allow-Origin", origin)
				h.Set("Access-Control-Allow-Credentials", "true")
				break
			}
		}
		return h, nil
	}
}

func registerChatNamespace(server *socketio.Server, auth *service.AuthService, chat *service.ChatService) {
	server.OnConnect(chatIONamespace, func(c socketio.Conn) error {
		ctx := context.Background()
		u, err := auth.SocketUserFromCookie(ctx, c.RemoteHeader().Get("Cookie"))
		if err != nil || u == nil {
			_ = c.Close()
			return fmt.Errorf("unauthorized")
		}
		c.SetContext(socketCtxFromUser(u))
		c.Join(fmt.Sprintf("user_%d", u.ID))
		return nil
	})
	server.OnDisconnect(chatIONamespace, func(_ socketio.Conn, _ string) {})

	server.OnEvent(chatIONamespace, "joinChat", func(c socketio.Conn, payload map[string]interface{}) interface{} {
		ctx := context.Background()
		uid := socketUserID(c)
		if uid == 0 {
			return map[string]interface{}{"success": false, "message": "Нет авторизации"}
		}
		chatID := extractInt32FromPayload(payload, "chatId")
		if chatID <= 0 {
			return map[string]interface{}{"success": false, "message": "Некорректный chatId"}
		}
		if _, err := chat.GetChatInfo(ctx, chatID, uid); err != nil {
			return map[string]interface{}{"success": false, "message": errWS(err)}
		}
		c.Join(fmt.Sprintf("chat_%d", chatID))
		return map[string]interface{}{"success": true, "message": "Подключен к чату"}
	})

	server.OnEvent(chatIONamespace, "sendMessage", func(c socketio.Conn, payload map[string]interface{}) interface{} {
		ctx := context.Background()
		uid := socketUserID(c)
		if uid == 0 {
			return map[string]interface{}{"success": false, "message": "Нет авторизации"}
		}
		chatID := extractInt32FromPayload(payload, "chatId")
		if chatID <= 0 {
			return map[string]interface{}{"success": false, "message": "Некорректный chatId"}
		}
		content, _ := payload["content"].(string)
		msg, err := chat.SendMessage(ctx, chatID, uid, content)
		if err != nil {
			return map[string]interface{}{"success": false, "message": errWS(err)}
		}
		emit := map[string]interface{}{
			"id": msg["id"], "content": msg["content"], "senderId": msg["senderId"], "sender": msg["sender"],
			"createdAt": msg["createdAt"], "timeString": msg["timeString"], "chatId": chatID,
		}
		server.BroadcastToRoom(chatIONamespace, fmt.Sprintf("chat_%d", chatID), "newMessage", emit)
		info, err := chat.GetChatInfo(ctx, chatID, uid)
		if err == nil {
			if rid, ok := companionIDFromInfo(info); ok {
				server.BroadcastToRoom(chatIONamespace, fmt.Sprintf("user_%d", rid), "newChatMessage", map[string]interface{}{
					"chatId": chatID, "message": msg, "product": info["product"], "companion": info["companion"],
				})
			}
		}
		return map[string]interface{}{"success": true, "message": "Сообщение отправлено"}
	})

	server.OnEvent(chatIONamespace, "markAsRead", func(c socketio.Conn, payload map[string]interface{}) interface{} {
		ctx := context.Background()
		uid := socketUserID(c)
		if uid == 0 {
			return map[string]interface{}{"success": false, "message": "Нет авторизации"}
		}
		chatID := extractInt32FromPayload(payload, "chatId")
		if chatID <= 0 {
			return map[string]interface{}{"success": false, "message": "Некорректный chatId"}
		}
		if err := chat.MarkMessagesAsRead(ctx, chatID, uid); err != nil {
			return map[string]interface{}{"success": false, "message": errWS(err)}
		}
		server.BroadcastToRoom(chatIONamespace, fmt.Sprintf("chat_%d", chatID), "messagesRead", map[string]interface{}{
			"chatId": chatID, "readBy": uid,
		})
		return map[string]interface{}{"success": true}
	})
}

func registerSupportNamespace(server *socketio.Server, auth *service.AuthService, sup *service.SupportService) {
	server.OnConnect(supportIONamespace, func(c socketio.Conn) error {
		ctx := context.Background()
		u, err := auth.SocketUserFromCookie(ctx, c.RemoteHeader().Get("Cookie"))
		if err != nil || u == nil {
			_ = c.Close()
			return fmt.Errorf("unauthorized")
		}
		c.SetContext(socketCtxFromUser(u))
		return nil
	})
	server.OnDisconnect(supportIONamespace, func(_ socketio.Conn, _ string) {})

	server.OnEvent(supportIONamespace, "joinTicket", func(c socketio.Conn, payload map[string]interface{}) {
		ctx := context.Background()
		sc := socketCtxUser(c)
		if sc == nil {
			c.Emit("error", map[string]interface{}{"message": "Нет авторизации"})
			return
		}
		tid := extractInt32FromPayload(payload, "ticketId")
		if tid <= 0 {
			c.Emit("error", map[string]interface{}{"message": "Некорректный ticketId"})
			return
		}
		_, err := sup.GetTicket(ctx, tid, sc.UserID, supportModeratorFromCtx(c))
		if err != nil {
			c.Emit("error", map[string]interface{}{"message": errWS(err)})
			return
		}
		room := fmt.Sprintf("support_ticket_%d", tid)
		c.Join(room)
		c.Emit("joinedTicket", map[string]interface{}{"ticketId": tid, "message": "Присоединились к тикету"})
	})

	server.OnEvent(supportIONamespace, "leaveTicket", func(c socketio.Conn, payload map[string]interface{}) {
		tid := extractInt32FromPayload(payload, "ticketId")
		if tid <= 0 {
			return
		}
		room := fmt.Sprintf("support_ticket_%d", tid)
		c.Leave(room)
		c.Emit("leftTicket", map[string]interface{}{"ticketId": tid, "message": "Покинули тикет"})
	})

	server.OnEvent(supportIONamespace, "sendSupportMessage", func(c socketio.Conn, payload map[string]interface{}) {
		ctx := context.Background()
		sc := socketCtxUser(c)
		if sc == nil {
			c.Emit("error", map[string]interface{}{"message": "Нет авторизации"})
			return
		}
		tid := extractInt32FromPayload(payload, "ticketId")
		if tid <= 0 {
			c.Emit("error", map[string]interface{}{"message": "Некорректный ticketId"})
			return
		}
		text := supportExtractMessageText(payload)
		isMod := supportModeratorFromCtx(c)
		msg, err := sup.SendMessage(ctx, tid, sc.UserID, text, isMod)
		if err != nil {
			c.Emit("error", map[string]interface{}{"message": errWS(err)})
			return
		}
		room := fmt.Sprintf("support_ticket_%d", tid)
		server.BroadcastToRoom(supportIONamespace, room, "newSupportMessage", map[string]interface{}{
			"ticketId": tid, "message": msg,
		})
		if isMod {
			server.BroadcastToRoom(supportIONamespace, room, "moderatorResponse", map[string]interface{}{
				"ticketId": tid, "message": "Модератор ответил на ваш запрос",
			})
		}
	})

	server.OnEvent(supportIONamespace, "supportTyping", func(c socketio.Conn, payload map[string]interface{}) {
		sc := socketCtxUser(c)
		if sc == nil {
			return
		}
		tid := extractInt32FromPayload(payload, "ticketId")
		if tid <= 0 {
			return
		}
		isTyping := false
		if v, ok := payload["isTyping"].(bool); ok {
			isTyping = v
		}
		room := fmt.Sprintf("support_ticket_%d", tid)
		senderID := c.ID()
		payloadOut := map[string]interface{}{
			"ticketId": tid, "userId": sc.UserID, "userName": sc.FullName, "isTyping": isTyping,
		}
		_ = server.ForEach(supportIONamespace, room, func(other socketio.Conn) {
			if other.ID() == senderID {
				return
			}
			other.Emit("supportUserTyping", payloadOut)
		})
	})

	server.OnEvent(supportIONamespace, "updateTicketStatus", func(c socketio.Conn, payload map[string]interface{}) {
		sc := socketCtxUser(c)
		if sc == nil {
			c.Emit("error", map[string]interface{}{"message": "Нет авторизации"})
			return
		}
		if sc.RoleName != "moderator" && sc.RoleName != "admin" {
			c.Emit("error", map[string]interface{}{"message": "Нет прав для изменения статуса тикета"})
			return
		}
		tid := extractInt32FromPayload(payload, "ticketId")
		status, _ := payload["status"].(string)
		if tid <= 0 || strings.TrimSpace(status) == "" {
			c.Emit("error", map[string]interface{}{"message": "Некорректные данные"})
			return
		}
		room := fmt.Sprintf("support_ticket_%d", tid)
		server.BroadcastToRoom(supportIONamespace, room, "ticketStatusUpdated", map[string]interface{}{
			"ticketId": tid, "status": status, "updatedBy": sc.FullName,
		})
	})
}

func supportExtractMessageText(payload map[string]interface{}) string {
	if payload == nil {
		return ""
	}
	m, ok := payload["message"].(map[string]interface{})
	if !ok {
		return ""
	}
	t, _ := m["text"].(string)
	return t
}

// NotifySupportNewTicket — как SupportGateway.notifyModeratorsNewTicket (после REST создания тикета).
func NotifySupportNewTicket(ticket interface{}) {
	engineSocketMu.RLock()
	s := engineSocketServer
	engineSocketMu.RUnlock()
	if s == nil {
		return
	}
	s.BroadcastToNamespace(supportIONamespace, "newSupportTicket", map[string]interface{}{
		"ticket": ticket, "message": "Создан новый тикет поддержки",
	})
}

// NotifySupportTicketAssigned — как notifyTicketAssigned.
func NotifySupportTicketAssigned(ticketID, moderatorID int32) {
	engineSocketMu.RLock()
	s := engineSocketServer
	engineSocketMu.RUnlock()
	if s == nil {
		return
	}
	room := fmt.Sprintf("support_ticket_%d", ticketID)
	s.BroadcastToRoom(supportIONamespace, room, "ticketAssigned", map[string]interface{}{
		"ticketId": ticketID, "moderatorId": moderatorID, "message": "Тикет назначен модератору",
	})
}

func extractInt32FromPayload(m map[string]interface{}, key string) int32 {
	if m == nil {
		return 0
	}
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch x := v.(type) {
	case float64:
		return int32(x)
	case float32:
		return int32(x)
	case int:
		return int32(x)
	case int32:
		return x
	case int64:
		return int32(x)
	default:
		return 0
	}
}

func companionIDFromInfo(info map[string]any) (int32, bool) {
	c, ok := info["companion"].(map[string]any)
	if !ok {
		return 0, false
	}
	v, ok := c["id"]
	if !ok {
		return 0, false
	}
	switch x := v.(type) {
	case float64:
		return int32(x), true
	case int:
		return int32(x), true
	case int32:
		return x, true
	case int64:
		return int32(x), true
	default:
		return 0, false
	}
}

func errWS(err error) string {
	var ae *service.AppError
	if errors.As(err, &ae) {
		return ae.Message
	}
	return err.Error()
}
