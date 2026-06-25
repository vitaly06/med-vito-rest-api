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

	"med-vito/api-go/internal/repository"
	"med-vito/api-go/internal/service"
)

// supportWSHub — in-memory hub для support WebSocket.
type supportWSHub struct {
	mu          sync.RWMutex
	ticketRooms map[int32]map[*supportWSClient]struct{} // ticketID → set of clients
	userClients map[int32]map[*supportWSClient]struct{} // userID → set of clients
	clientRooms map[*supportWSClient]map[int32]struct{} // client → set of ticketIDs
}

type supportWSClient struct {
	conn   *websocket.Conn
	userID int32
}

func newSupportWSHub() *supportWSHub {
	return &supportWSHub{
		ticketRooms: make(map[int32]map[*supportWSClient]struct{}),
		userClients: make(map[int32]map[*supportWSClient]struct{}),
		clientRooms: make(map[*supportWSClient]map[int32]struct{}),
	}
}

func (h *supportWSHub) addClient(c *supportWSClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.userClients[c.userID]; !ok {
		h.userClients[c.userID] = make(map[*supportWSClient]struct{})
	}
	h.userClients[c.userID][c] = struct{}{}
	h.clientRooms[c] = make(map[int32]struct{})
}

func (h *supportWSHub) removeClient(c *supportWSClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.userClients[c.userID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.userClients, c.userID)
		}
	}
	if rooms, ok := h.clientRooms[c]; ok {
		for tid := range rooms {
			if roomSet, ok2 := h.ticketRooms[tid]; ok2 {
				delete(roomSet, c)
				if len(roomSet) == 0 {
					delete(h.ticketRooms, tid)
				}
			}
		}
		delete(h.clientRooms, c)
	}
}

func (h *supportWSHub) joinTicket(c *supportWSClient, ticketID int32) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.ticketRooms[ticketID]; !ok {
		h.ticketRooms[ticketID] = make(map[*supportWSClient]struct{})
	}
	h.ticketRooms[ticketID][c] = struct{}{}
	if _, ok := h.clientRooms[c]; !ok {
		h.clientRooms[c] = make(map[int32]struct{})
	}
	h.clientRooms[c][ticketID] = struct{}{}
}

func (h *supportWSHub) leaveTicket(c *supportWSClient, ticketID int32) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if roomSet, ok := h.ticketRooms[ticketID]; ok {
		delete(roomSet, c)
		if len(roomSet) == 0 {
			delete(h.ticketRooms, ticketID)
		}
	}
	if rooms, ok := h.clientRooms[c]; ok {
		delete(rooms, ticketID)
	}
}

func (h *supportWSHub) broadcastToTicket(ticketID int32, event string, data any) {
	h.mu.RLock()
	clients := make([]*supportWSClient, 0)
	for c := range h.ticketRooms[ticketID] {
		clients = append(clients, c)
	}
	h.mu.RUnlock()
	for _, c := range clients {
		_ = c.conn.WriteJSON(wsOutgoing{Event: event, Data: data})
	}
}

// BroadcastToUser — рассылает событие всем WS-соединениям пользователя.
// Используется для push-уведомлений о балансе/лимите.
func (h *supportWSHub) broadcastToUser(userID int32, event string, data any) {
	h.mu.RLock()
	clients := make([]*supportWSClient, 0)
	for c := range h.userClients[userID] {
		clients = append(clients, c)
	}
	h.mu.RUnlock()
	for _, c := range clients {
		_ = c.conn.WriteJSON(wsOutgoing{Event: event, Data: data})
	}
}

var supportNativeHub = newSupportWSHub()

// BroadcastSupportNotification — вызывается из SupportService для push в WS.
func BroadcastSupportNotification(userID int32, ticketID int32, data any) {
	supportNativeHub.broadcastToTicket(ticketID, "newSupportMessage", data)
	// Кроме этого — уведомление пользователю, если он не в тикете
	supportNativeHub.broadcastToUser(userID, "supportNotification", data)
}

type joinTicketData struct {
	TicketID int32 `json:"ticketId"`
}

type supportSendData struct {
	TicketID int32  `json:"ticketId"`
	Message  struct {
		Text string `json:"text"`
	} `json:"message"`
}

func RegisterSupportWS(app *fiber.App, auth *service.AuthService, sup *service.SupportService) {
	// Подключаем WS-бродкаст: при каждом новом сообщении рассылаем подписчикам.
	sup.SetWSNotifier(func(ownerID int32, ticketID int32, msg *repository.SupportMessageOut) {
		// Маскируем имя модератора перед отправкой пользователю
		masked := *msg
		if masked.Author.Role != nil {
			rn := masked.Author.Role.Name
			if isModeratorRole(rn) {
				masked.Author.FullName = "Служба поддержки"
				masked.Author.Email = ""
			}
		}
		payload := map[string]any{"ticketId": ticketID, "message": &masked}
		supportNativeHub.broadcastToTicket(ticketID, "newSupportMessage", payload)
	})

	app.Use("/ws/support", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/support", websocket.New(func(c *websocket.Conn) {
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

		client := &supportWSClient{conn: c, userID: u.ID}
		supportNativeHub.addClient(client)
		log.Printf("Native WS support connected: user=%d", u.ID)
		defer func() {
			supportNativeHub.removeClient(client)
			_ = c.Close()
			log.Printf("Native WS support disconnected: user=%d", u.ID)
		}()

		_ = c.SetReadDeadline(time.Now().Add(120 * time.Second))
		c.SetPongHandler(func(string) error {
			_ = c.SetReadDeadline(time.Now().Add(120 * time.Second))
			return nil
		})

		isModerator := u.RoleName != nil && isModeratorRole(*u.RoleName)

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
			if err := handleWSSupportEvent(ctx, client, sup, in, isModerator); err != nil {
				_ = c.WriteJSON(wsOutgoing{Event: "error", Data: map[string]any{"message": err.Error()}})
			}
		}
	}))
}

func isModeratorRole(roleName string) bool {
	switch strings.ToUpper(roleName) {
	case "ADMIN", "SUPERADMIN", "SENIOR_MODERATOR", "MODERATOR":
		return true
	}
	return false
}

func handleWSSupportEvent(ctx context.Context, c *supportWSClient, sup *service.SupportService, in wsIncoming, isModerator bool) error {
	switch in.Event {
	case "joinTicket":
		var d joinTicketData
		if err := json.Unmarshal(in.Data, &d); err != nil || d.TicketID <= 0 {
			return errors.New("invalid ticketId")
		}
		_, err := sup.GetTicket(ctx, d.TicketID, c.userID, isModerator)
		if err != nil {
			return err
		}
		supportNativeHub.joinTicket(c, d.TicketID)
		_ = c.conn.WriteJSON(wsOutgoing{Event: "joinedTicket", Data: map[string]any{"ticketId": d.TicketID}})
		return nil

	case "leaveTicket":
		var d joinTicketData
		if err := json.Unmarshal(in.Data, &d); err != nil || d.TicketID <= 0 {
			return errors.New("invalid ticketId")
		}
		supportNativeHub.leaveTicket(c, d.TicketID)
		return nil

	case "sendSupportMessage":
		var d supportSendData
		if err := json.Unmarshal(in.Data, &d); err != nil || d.TicketID <= 0 {
			return errors.New("invalid payload")
		}
		text := strings.TrimSpace(d.Message.Text)
		if text == "" {
			return errors.New("empty message")
		}
		// SendMessage вызывает onNewMsg-колбэк, который делает broadcastToTicket.
		if _, err := sup.SendMessage(ctx, d.TicketID, c.userID, text, isModerator); err != nil {
			return err
		}
		return nil

	default:
		return nil
	}
}
