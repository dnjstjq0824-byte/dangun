# BN NEAR — WebSocket Protocol Specification

**Order ID:** BN-NEAR-MASTER-001  
**Endpoint:** `GET /api/v1/ws?token=<jwt_access_token>`  
**Subprotocol:** `bn-near.v1+json`  

---

## 1. Connection Lifecycle & Handshake
1. **Upgrade Request:** Client connects via HTTP Upgrade with bearer token parameter or header.
2. **Authentication:** Server validates the JWT session. If invalid, returns HTTP 401 and closes connection.
3. **Session Initialization:** Server registers connection in the active client hub and sends `connection.ready`.
4. **Heartbeat:** Ping/Pong every 30 seconds to maintain connection through NAT/firewalls.

---

## 2. Client-to-Server Events

### Event: `chat.message.send`
```json
{
  "type": "chat.message.send",
  "client_message_id": "a9e63c0a-01d6-444f-bc1e-f3b1451f1224",
  "conversation_id": "4e75d04c-3521-42ab-873b-e0fa90251781",
  "payload": {
    "message_type": "TEXT",
    "text": "안녕하세요! 청라 호수공원 근처에서 거래 가능할까요?",
    "attachment_id": null
  }
}
```

### Event: `chat.read.ack`
```json
{
  "type": "chat.read.ack",
  "conversation_id": "4e75d04c-3521-42ab-873b-e0fa90251781",
  "sequence_no": 42
}
```

### Event: `chat.sync`
```json
{
  "type": "chat.sync",
  "conversation_id": "4e75d04c-3521-42ab-873b-e0fa90251781",
  "since_sequence_no": 35
}
```

---

## 3. Server-to-Client Events

### Event: `chat.message.ack` (Sent back to sender)
```json
{
  "type": "chat.message.ack",
  "client_message_id": "a9e63c0a-01d6-444f-bc1e-f3b1451f1224",
  "message_id": "78201a4e-1284-4841-8608-f4cb5038c35b",
  "conversation_id": "4e75d04c-3521-42ab-873b-e0fa90251781",
  "sequence_no": 43,
  "status": "STORED",
  "created_at": "2026-10-24T22:30:00Z"
}
```

### Event: `chat.message.receive` (Broadcast to recipient)
```json
{
  "type": "chat.message.receive",
  "message_id": "78201a4e-1284-4841-8608-f4cb5038c35b",
  "conversation_id": "4e75d04c-3521-42ab-873b-e0fa90251781",
  "sender_user_id": "c13d7fb6-e137-4d43-85bb-b5bf67c3bb05",
  "sequence_no": 43,
  "payload": {
    "message_type": "TEXT",
    "text": "안녕하세요! 청라 호수공원 근처에서 거래 가능할까요?",
    "attachment_url": null
  },
  "created_at": "2026-10-24T22:30:00Z"
}
```

### Event: `chat.read.update` (Broadcast to sender when recipient reads)
```json
{
  "type": "chat.read.update",
  "conversation_id": "4e75d04c-3521-42ab-873b-e0fa90251781",
  "reader_user_id": "8fa160c8-2da7-4e78-8547-ee872b7a9f93",
  "sequence_no": 43,
  "read_at": "2026-10-24T22:30:05Z"
}
```

---

## 4. Concurrency, Reconnection & Deduplication Rules
1. **Deduplication:** The tuple `(conversation_id, sender_user_id, client_message_id)` is enforced unique. Duplicate client submissions return the existing message sequence without re-persisting.
2. **Monotonic Sequences:** Every persisted message is assigned an incremental integer `sequence_no` within the transaction.
3. **Reconnection Catch-Up:** When a client reconnects, it sends `chat.sync` with its highest known sequence number. The server streams all missing messages in order.
