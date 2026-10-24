# BN NEAR — 05. Master User Flows

**Order ID:** BN-NEAR-MASTER-001  

---

## Flow 1: Registration, Neighborhood Selection & GPS Verification
```text
[Open App/Web] 
       │
       ▼
[Language / Country Selection]
       │
       ▼
[Register Email & Password] ──► [Email Verification Token Sent] ──► [Account Activated]
       │
       ▼
[Select Home Neighborhood] (e.g. Cheongna-dong, Seo-gu, Incheon)
       │
       ▼
[Tap 'Verify My Neighborhood'] ──► Request Native Location Permission
       │
       ▼
[Obtain GPS Coordinate + Accuracy] (e.g., Lat: 37.532, Lng: 126.654, Acc: 18m)
       │
       ▼
[Send to Go API: POST /api/v1/me/neighborhood-verifications]
       │
       ▼
[PostGIS ST_Contains / ST_DWithin Boundary Check]
       │
       ├──► [Within Boundary & Acc < 100m] ──► VERIFIED (30-day validity window)
       └──► [Outside Boundary / High Error]  ──► FAILED with actionable guidance
```

---

## Flow 2: Marketplace Sell Flow
```text
[Tap 'Sell' Floating Action Button]
       │
       ▼
[Upload Photos] ──► Client strips EXIF GPS metadata ──► Server validates MIME & scans
       │
       ▼
[Input Details]:
  • Title (e.g., "Apple M2 MacBook Air 16GB")
  • Category (Electronics)
  • Condition (LIKE_NEW)
  • Price Minor (KRW 1,150,000)
  • Negotiable: TRUE
  • Target Neighborhood (Auto-selected Primary Verified Neighborhood)
       │
       ▼
[Tap 'Publish Listing'] ──► Saved to PostgreSQL `listings` table
       │
       ▼
[Background Worker updates Geo & Text Search Indexes]
       │
       ▼
[Listing immediately visible in local feed & push alerts triggered for saved searches]
```

---

## Flow 3: Buyer Inquiry, Real-Time Chat & Price Negotiation
```text
[Buyer opens Listing Detail] ──► [Tap 'Chat Now']
       │
       ▼
[API GET/POST /api/v1/conversations] 
       ├── Existing conversation found OR new conversation record created
       │
       ▼
[WebSocket Handshake with JWT: /api/v1/ws]
       │
       ▼
[Buyer sends message] ──► Client sends {"type": "chat.message.send", ...}
       │
       ▼
[Go Server validates conversation membership & active blocks]
       │
       ▼
[Server persists message to PostgreSQL `chat_messages` with sequence_no]
       │
       ▼
[Server returns ACK to Buyer & broadcasts to Seller via WebSocket]
       │
       ▼
[Buyer taps 'Make Offer'] ──► Submits KRW 1,050,000
       │
       ▼
[Seller receives Offer Card in Chat] ──► Taps 'Accept Offer'
       │
       ▼
[Deal status shifts to NEGOTIATING / OFFER_ACCEPTED]
```

---

## Flow 4: Reservation & Two-Phase Deal Completion
```text
[Seller selects Buyer in Chat] ──► Taps 'Reserve Item'
       │
       ▼
[API POST /api/v1/deals/{id}/reserve with DB row lock]
       │
       ├── Success: Listing status -> `RESERVED`, other buyers see reserved badge
       └── Conflict: Prevented if already reserved by another party
       │
       ▼
[Parties meet in person at safe public location and exchange goods]
       │
       ▼
[Seller taps 'Mark as Sold'] ──► Deal status -> `COMPLETED_CLAIMED`
       │
       ▼
[Buyer receives prompt: 'Did you receive the item?']
       │
       ▼
[Buyer taps 'Confirm Purchase'] ──► Deal status -> `COMPLETED_CONFIRMED`
       │
       ▼
[Mutual Review Opportunity Unlocked for both parties]
```

---

## Flow 5: Mutual Reviews & Manner Score Calculation
```text
[Deal Reaches COMPLETED_CONFIRMED]
       │
       ▼
[Review Prompt displays in Chat & Activity Tab]
       │
       ▼
[Party selects Rating (1-5 stars) + Manner Tags]
  • "Punctual", "Kind and responsive", "Item exactly as described"
       │
       ▼
[Submit Review to POST /api/v1/deals/{id}/reviews]
       │
       ▼
[Server validates deal ownership and single-submission constraint]
       │
       ▼
[Reputation Engine computes updated weighted Manner Temperature]
       │
       ▼
[Updated score reflected on user's public profile badge]
```
