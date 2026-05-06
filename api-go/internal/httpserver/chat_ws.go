package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"med-vito/api-go/internal/service"
)

type wsClient struct {
	conn     *websocket.Conn
	fullName string
	userID   int32
}

type wsHub struct {
	mu          sync.RWMutex
	clientRooms map[*wsClient]map[int32]struct{}
	rooms       map[int32]map[*wsClient]struct{}
	users       map[int32]map[*wsClient]struct{}
}

func newWSHub() *wsHub {
	return &wsHub{
		clientRooms: make(map[*wsClient]map[int32]struct{}),
		rooms:       make(map[int32]map[*wsClient]struct{}),
		users:       make(map[int32]map[*wsClient]struct{}),
	}
}

func (h *wsHub) addClient(c *wsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.users[c.userID]; !ok {
		h.users[c.userID] = make(map[*wsClient]struct{})
	}
	h.users[c.userID][c] = struct{}{}
	h.clientRooms[c] = make(map[int32]struct{})
}

func (h *wsHub) removeClient(c *wsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.users[c.userID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.users, c.userID)
		}
	}
	if rooms, ok := h.clientRooms[c]; ok {
		for chatID := range rooms {
			if roomSet, ok2 := h.rooms[chatID]; ok2 {
				delete(roomSet, c)
				if len(roomSet) == 0 {
					delete(h.rooms, chatID)
				}
			}
		}
		delete(h.clientRooms, c)
	}
}

func (h *wsHub) joinChat(c *wsClient, chatID int32) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[chatID]; !ok {
		h.rooms[chatID] = make(map[*wsClient]struct{})
	}
	h.rooms[chatID][c] = struct{}{}
	if _, ok := h.clientRooms[c]; !ok {
		h.clientRooms[c] = make(map[int32]struct{})
	}
	h.clientRooms[c][chatID] = struct{}{}
}

func (h *wsHub) leaveChat(c *wsClient, chatID int32) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if roomSet, ok := h.rooms[chatID]; ok {
		delete(roomSet, c)
		if len(roomSet) == 0 {
			delete(h.rooms, chatID)
		}
	}
	if rooms, ok := h.clientRooms[c]; ok {
		delete(rooms, chatID)
	}
}

func (h *wsHub) roomClients(chatID int32) []*wsClient {
	h.mu.RLock()
	defer h.mu.RUnlock()
	set := h.rooms[chatID]
	out := make([]*wsClient, 0, len(set))
	for c := range set {
		out = append(out, c)
	}
	return out
}

func (h *wsHub) userClients(userID int32) []*wsClient {
	h.mu.RLock()
	defer h.mu.RUnlock()
	set := h.users[userID]
	out := make([]*wsClient, 0, len(set))
	for c := range set {
		out = append(out, c)
	}
	return out
}

var chatNativeHub = newWSHub()

type wsIncoming struct {
	Data  json.RawMessage `json:"data"`
	Event string          `json:"event"`
}

type wsOutgoing struct {
	Data  any    `json:"data"`
	Event string `json:"event"`
}

type joinChatData struct {
	ChatID int32 `json:"chatId"`
}

type markAsReadData struct {
	ChatID int32 `json:"chatId"`
}

type sendMessageData struct {
	ChatID  int32  `json:"chatId"`
	Content string `json:"content"`
}

type typingData struct {
	ChatID   int32 `json:"chatId"`
	IsTyping bool  `json:"isTyping"`
}

func RegisterChatWS(app *fiber.App, auth *service.AuthService, chat *service.ChatService) {
	app.Use("/ws/chat", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/chat", websocket.New(func(c *websocket.Conn) {
		ctx := context.Background()
		sid := strings.TrimSpace(c.Query("session_id"))
		if sid == "" {
			sid = strings.TrimSpace(c.Cookies("session_id"))
		}
		if sid == "" {
			sid = strings.TrimSpace(c.Cookies("ws_session_id"))
		}
		if sid == "" {
			_ = c.WriteJSON(wsOutgoing{Event: "error", Data: map[string]any{"message": "unauthorized"}})
			_ = c.Close()
			return
		}

		u, err := auth.UserFromSession(ctx, sid)
		if err != nil || u == nil {
			_ = c.WriteJSON(wsOutgoing{Event: "error", Data: map[string]any{"message": "unauthorized"}})
			_ = c.Close()
			return
		}

		client := &wsClient{conn: c, userID: u.ID, fullName: u.FullName}
		chatNativeHub.addClient(client)
		log.Printf("Native WS chat connected: user=%d", u.ID)
		defer func() {
			chatNativeHub.removeClient(client)
			_ = c.Close()
			log.Printf("Native WS chat disconnected: user=%d", u.ID)
		}()

		_ = c.SetReadDeadline(time.Now().Add(120 * time.Second))
		c.SetPongHandler(func(string) error {
			_ = c.SetReadDeadline(time.Now().Add(120 * time.Second))
			return nil
		})

		for {
			_, b, err := c.ReadMessage()
			if err != nil {
				return
			}
			var in wsIncoming
			if err := json.Unmarshal(b, &in); err != nil {
				_ = c.WriteJSON(wsOutgoing{Event: "error", Data: map[string]any{"message": "invalid json"}})
				continue
			}
			if err := handleWSChatEvent(ctx, client, chat, in); err != nil {
				_ = c.WriteJSON(wsOutgoing{Event: "error", Data: map[string]any{"message": err.Error()}})
			}
		}
	}))
}

func handleWSChatEvent(ctx context.Context, c *wsClient, chat *service.ChatService, in wsIncoming) error {
	switch in.Event {
	case "joinChat":
		var d joinChatData
		if err := json.Unmarshal(in.Data, &d); err != nil || d.ChatID <= 0 {
			return errors.New("invalid chatId")
		}
		if _, err := chat.GetChatInfo(ctx, d.ChatID, c.userID); err != nil {
			return err
		}
		chatNativeHub.joinChat(c, d.ChatID)
		_ = c.conn.WriteJSON(wsOutgoing{Event: "joinedChat", Data: map[string]any{"chatId": d.ChatID}})
		return nil
	case "leaveChat":
		var d joinChatData
		if err := json.Unmarshal(in.Data, &d); err != nil || d.ChatID <= 0 {
			return errors.New("invalid chatId")
		}
		chatNativeHub.leaveChat(c, d.ChatID)
		return nil
	case "typing":
		var d typingData
		if err := json.Unmarshal(in.Data, &d); err != nil || d.ChatID <= 0 {
			return errors.New("invalid payload")
		}
		for _, other := range chatNativeHub.roomClients(d.ChatID) {
			if other == c {
				continue
			}
			_ = other.conn.WriteJSON(wsOutgoing{
				Event: "userTyping",
				Data:  map[string]any{"chatId": d.ChatID, "isTyping": d.IsTyping, "userId": c.userID},
			})
		}
		return nil
	case "markAsRead":
		var d markAsReadData
		if err := json.Unmarshal(in.Data, &d); err != nil || d.ChatID <= 0 {
			return errors.New("invalid chatId")
		}
		if err := chat.MarkMessagesAsRead(ctx, d.ChatID, c.userID); err != nil {
			return err
		}
		for _, other := range chatNativeHub.roomClients(d.ChatID) {
			_ = other.conn.WriteJSON(wsOutgoing{
				Event: "messagesRead",
				Data:  map[string]any{"chatId": d.ChatID, "readBy": c.userID},
			})
		}
		return nil
	case "sendMessage":
		var d sendMessageData
		if err := json.Unmarshal(in.Data, &d); err != nil || d.ChatID <= 0 || strings.TrimSpace(d.Content) == "" {
			return errors.New("invalid payload")
		}
		msg, err := chat.SendMessage(ctx, d.ChatID, c.userID, d.Content)
		if err != nil {
			return err
		}
		emit := map[string]any{
			"id": msg["id"], "content": msg["content"], "senderId": msg["senderId"], "sender": msg["sender"],
			"createdAt": msg["createdAt"], "timeString": msg["timeString"], "chatId": d.ChatID,
		}
		for _, other := range chatNativeHub.roomClients(d.ChatID) {
			_ = other.conn.WriteJSON(wsOutgoing{Event: "newMessage", Data: emit})
		}
		info, err := chat.GetChatInfo(ctx, d.ChatID, c.userID)
		if err == nil {
			if rid, ok := companionIDFromInfo(info); ok {
				for _, uc := range chatNativeHub.userClients(rid) {
					_ = uc.conn.WriteJSON(wsOutgoing{
						Event: "newChatMessage",
						Data: map[string]any{
							"chatId": d.ChatID, "message": msg, "product": info["product"], "companion": info["companion"],
						},
					})
				}
			}
		}
		return nil
	default:
		log.Printf("Native WS chat unknown event: %s", in.Event)
		return nil
	}
}

