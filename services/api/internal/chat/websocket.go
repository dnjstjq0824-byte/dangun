package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/booking-network/bn-near/services/api/internal/middleware"
	"github.com/booking-network/bn-near/services/api/internal/models"
	"github.com/booking-network/bn-near/services/api/internal/storage"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for API clients
	},
}

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	User     *models.User
	SendChan chan []byte
}

type Hub struct {
	mu         sync.RWMutex
	Clients    map[string]map[*Client]bool // userID -> set of active clients
	Store      storage.Storage
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewHub(store storage.Storage) *Hub {
	return &Hub{
		Clients:    make(map[string]map[*Client]bool),
		Store:      store,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, ok := h.Clients[client.User.ID]; !ok {
				h.Clients[client.User.ID] = make(map[*Client]bool)
			}
			h.Clients[client.User.ID][client] = true
			if mem, ok := h.Store.(*storage.MemoryStore); ok {
				mem.LiveWSSessionsCount++
			}
			h.mu.Unlock()

			// Send connection ready
			readyEvent, _ := json.Marshal(map[string]interface{}{
				"type":       "connection.ready",
				"user_id":    client.User.ID,
				"server_time": time.Now().UTC().Format(time.RFC3339),
			})
			client.SendChan <- readyEvent

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Clients[client.User.ID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.SendChan)
					if len(clients) == 0 {
						delete(h.Clients, client.User.ID)
					}
					if mem, ok := h.Store.(*storage.MemoryStore); ok {
						if mem.LiveWSSessionsCount > 0 {
							mem.LiveWSSessionsCount--
						}
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) SendToUser(userID string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if userClients, ok := h.Clients[userID]; ok {
		for c := range userClients {
			select {
			case c.SendChan <- payload:
			default:
				// Buffer full
			}
		}
	}
}

type WsIncomingFrame struct {
	Type            string          `json:"type"`
	ConversationID  string          `json:"conversation_id"`
	ClientMessageID string          `json:"client_message_id"`
	SequenceNo      int64           `json:"sequence_no"`
	SinceSequenceNo int64           `json:"since_sequence_no"`
	Payload         json.RawMessage `json:"payload"`
}

type MessagePayload struct {
	MessageType string `json:"message_type"`
	Text        string `json:"text"`
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(65536)
	_ = c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var frame WsIncomingFrame
		if err := json.Unmarshal(message, &frame); err != nil {
			continue
		}

		switch frame.Type {
		case "ping":
			pong, _ := json.Marshal(map[string]string{"type": "pong"})
			c.SendChan <- pong

		case "chat.message.send":
			var pl MessagePayload
			_ = json.Unmarshal(frame.Payload, &pl)
			if pl.MessageType == "" {
				pl.MessageType = "TEXT"
			}

			// Save to Storage
			msg := &models.ChatMessage{
				ConversationID:  frame.ConversationID,
				SenderUserID:    c.User.ID,
				ClientMessageID: frame.ClientMessageID,
				MessageType:     pl.MessageType,
				Body:            pl.Text,
			}

			saved, err := c.Hub.Store.SaveMessage(msg)
			if err != nil {
				errResp, _ := json.Marshal(map[string]interface{}{
					"type":              "error",
					"client_message_id": frame.ClientMessageID,
					"message":           err.Error(),
				})
				c.SendChan <- errResp
				continue
			}

			// Send ACK back to sender
			ackResp, _ := json.Marshal(map[string]interface{}{
				"type":              "chat.message.ack",
				"client_message_id": saved.ClientMessageID,
				"message_id":        saved.ID,
				"conversation_id":   saved.ConversationID,
				"sequence_no":       saved.SequenceNo,
				"status":            "STORED",
				"created_at":        saved.CreatedAt.UTC().Format(time.RFC3339),
			})
			c.SendChan <- ackResp

			// Broadcast to other participant
			conv, _ := c.Hub.Store.GetConversationByID(saved.ConversationID, c.User.ID)
			if conv != nil {
				recipientID := conv.BuyerUserID
				if c.User.ID == conv.BuyerUserID {
					recipientID = conv.SellerUserID
				}

				receiveEvent, _ := json.Marshal(map[string]interface{}{
					"type":            "chat.message.receive",
					"message_id":      saved.ID,
					"conversation_id": saved.ConversationID,
					"sender_user_id":  saved.SenderUserID,
					"sender_nickname": c.User.Nickname,
					"sequence_no":     saved.SequenceNo,
					"payload": map[string]interface{}{
						"message_type": saved.MessageType,
						"text":         saved.Body,
					},
					"created_at": saved.CreatedAt.UTC().Format(time.RFC3339),
				})
				c.Hub.SendToUser(recipientID, receiveEvent)

				// Also create in-app notification for recipient
				c.Hub.Store.SendNotification(&models.Notification{
					UserID: recipientID,
					Type:   "CHAT_MESSAGE",
					Title:  c.User.Nickname + "님의 메시지",
					Body:   saved.Body,
				})
			}

		case "chat.read.ack":
			_ = c.Hub.Store.MarkMessagesRead(frame.ConversationID, c.User.ID, frame.SequenceNo)

			conv, _ := c.Hub.Store.GetConversationByID(frame.ConversationID, c.User.ID)
			if conv != nil {
				otherUserID := conv.BuyerUserID
				if c.User.ID == conv.BuyerUserID {
					otherUserID = conv.SellerUserID
				}
				readUpdate, _ := json.Marshal(map[string]interface{}{
					"type":            "chat.read.update",
					"conversation_id": frame.ConversationID,
					"reader_user_id":  c.User.ID,
					"sequence_no":     frame.SequenceNo,
					"read_at":         time.Now().UTC().Format(time.RFC3339),
				})
				c.Hub.SendToUser(otherUserID, readUpdate)
			}

		case "chat.sync":
			msgs, _ := c.Hub.Store.GetMessages(frame.ConversationID, frame.SinceSequenceNo)
			syncResp, _ := json.Marshal(map[string]interface{}{
				"type":            "chat.sync.result",
				"conversation_id": frame.ConversationID,
				"messages":        msgs,
			})
			c.SendChan <- syncResp
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendChan:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token required for websocket connection", http.StatusUnauthorized)
		return
	}

	claims, err := middleware.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user, err := hub.Store.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}

	client := &Client{
		Hub:      hub,
		Conn:     conn,
		User:     user,
		SendChan: make(chan []byte, 256),
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
