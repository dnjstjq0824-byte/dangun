package main_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/booking-network/bn-near/services/api/internal/chat"
	"github.com/booking-network/bn-near/services/api/internal/middleware"
	"github.com/booking-network/bn-near/services/api/internal/router"
	"github.com/booking-network/bn-near/services/api/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func TestWebSocketRealtimeChat(t *testing.T) {
	store := storage.NewMemoryStore()
	hub := chat.NewHub(store)
	go hub.Run()
	routes := router.NewRouter(store, hub)

	server := httptest.NewServer(routes)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws"

	// 1. Authenticate Seller and Buyer
	seller, _ := store.GetUserByEmail("seller@bn-near.com")
	buyer, _ := store.GetUserByEmail("buyer@bn-near.com")

	sellerToken, _, _, _ := middleware.GenerateTokens(seller)
	buyerToken, _, _, _ := middleware.GenerateTokens(buyer)

	// 2. Connect Buyer to WebSocket
	buyerWS, _, err := websocket.DefaultDialer.Dial(wsURL+"?token="+buyerToken, nil)
	if err != nil {
		t.Fatalf("Failed to dial buyer websocket: %v", err)
	}
	defer buyerWS.Close()

	// Read connection.ready
	var readyEvent map[string]interface{}
	_ = buyerWS.ReadJSON(&readyEvent)
	if readyEvent["type"] != "connection.ready" {
		t.Fatalf("Expected connection.ready event, got %v", readyEvent)
	}

	// 3. Connect Seller to WebSocket
	sellerWS, _, err := websocket.DefaultDialer.Dial(wsURL+"?token="+sellerToken, nil)
	if err != nil {
		t.Fatalf("Failed to dial seller websocket: %v", err)
	}
	defer sellerWS.Close()

	_ = sellerWS.ReadJSON(&readyEvent)

	// 4. Create conversation
	listingID := "l0000000-0000-0000-0000-000000000001"
	conv, err := store.GetOrCreateConversation(listingID, buyer.ID)
	if err != nil {
		t.Fatalf("Failed to get or create conversation: %v", err)
	}

	// 5. Buyer sends message via WebSocket
	clientMsgID := uuid.NewString()
	sendEvent := map[string]interface{}{
		"type":              "chat.message.send",
		"conversation_id":   conv.ID,
		"client_message_id": clientMsgID,
		"payload": map[string]string{
			"message_type": "TEXT",
			"text":         "실시간 웹소켓 채팅 테스트 메시지입니다.",
		},
	}
	if err := buyerWS.WriteJSON(sendEvent); err != nil {
		t.Fatalf("Failed to write websocket message: %v", err)
	}

	// 6. Buyer must receive ACK
	var ackEvent map[string]interface{}
	_ = buyerWS.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err := buyerWS.ReadJSON(&ackEvent); err != nil {
		t.Fatalf("Buyer did not receive ACK: %v", err)
	}
	if ackEvent["type"] != "chat.message.ack" {
		t.Fatalf("Expected chat.message.ack, got %v", ackEvent)
	}
	if ackEvent["client_message_id"] != clientMsgID {
		t.Fatalf("ACK client_message_id mismatch")
	}

	// 7. Seller must receive chat.message.receive
	var receiveEvent map[string]interface{}
	_ = sellerWS.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err := sellerWS.ReadJSON(&receiveEvent); err != nil {
		t.Fatalf("Seller did not receive broadcast: %v", err)
	}
	if receiveEvent["type"] != "chat.message.receive" {
		t.Fatalf("Expected chat.message.receive, got %v", receiveEvent)
	}
	pl := receiveEvent["payload"].(map[string]interface{})
	if pl["text"] != "실시간 웹소켓 채팅 테스트 메시지입니다." {
		t.Fatalf("Received message body mismatch: %v", pl)
	}

	// 8. Ping / Pong test
	if err := buyerWS.WriteJSON(map[string]string{"type": "ping"}); err != nil {
		t.Fatalf("Failed to send ping: %v", err)
	}
	var pongEvent map[string]interface{}
	_ = buyerWS.ReadJSON(&pongEvent)
	if pongEvent["type"] != "pong" {
		t.Fatalf("Expected pong event, got %v", pongEvent)
	}
}
