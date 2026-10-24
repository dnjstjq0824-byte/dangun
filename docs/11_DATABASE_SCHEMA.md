# BN NEAR — 11. Database Schema Specification

**Order ID:** BN-NEAR-MASTER-001  
**Database Engine:** PostgreSQL 16 + PostGIS 3.4  
**Primary Key Standard:** UUIDv4 across all primary entities  
**Time Standard:** TIMESTAMPTZ (UTC)  
**Currency & Price Standard:** ISO 4217 code + BIGINT minor currency units (e.g. 1000 KRW = 1000, 10.50 USD = 1050)  

---

## 1. Relational Entity Domain Groups

### A. Identity & Access
- `users`: Core account record (UUID, email, password_hash, nickname, avatar_asset_id, locale, country_code, status, timestamps).
- `user_profiles`: Extended user profile (bio, display name, preferred_currency, timezone, manner_temperature).
- `user_sessions`: Active JWT refresh tokens, device identifiers, IP metadata, revocation timestamps.
- `user_blocks`: Unique pair `(blocker_user_id, blocked_user_id)` preventing chat and deal interactions.
- `consent_records`: User privacy policy and location terms consent version tracking.

### B. Geography & Boundaries
- `countries`: Country code (ISO 3166-1 alpha-2), name, currency, launch status (`ACTIVE`, `BETA`, `READ_ONLY`).
- `regions`: Administrative province/state subdivision.
- `cities`: Municipal city record with timezone and bounding box.
- `neighborhoods`: Hyperlocal neighborhood entity containing:
  - `boundary_geom`: PostGIS `GEOMETRY(Polygon, 4326)` representing the official boundary.
  - `public_center_geom`: PostGIS `GEOMETRY(Point, 4326)` representing the neighborhood center.
- `user_neighborhoods`: Association between user and neighborhood, primary flag, verification status (`VERIFIED`, `EXPIRED`, `UNVERIFIED`), verification expiry timestamp.
- `location_verification_attempts`: Transient verification audit log recording accuracy, verification result, and risk codes (without permanently storing raw residential GPS coordinates).

### C. Marketplace & Listings
- `categories`: Hierarchical category tree (slug, name, icon, sort order).
- `listings`:
  - `seller_user_id`: Reference to user.
  - `neighborhood_id`: Associated verified neighborhood.
  - `category_id`: Associated category.
  - `title`, `description`, `condition` (`NEW`, `LIKE_NEW`, `GOOD`, `FAIR`, `FOR_PARTS`).
  - `price_minor`: Price in smallest currency unit (must be >= 0).
  - `currency_code`: 3-letter currency code.
  - `is_free`: Boolean indicating free giveaway.
  - `is_negotiable`: Boolean indicating whether price offers are accepted.
  - `listing_status`: `DRAFT`, `ACTIVE`, `RESERVED`, `SOLD`, `ARCHIVED`, `REMOVED`.
  - `public_location_geom`: Approximate public coordinate for spatial radius search.
  - `view_count`, `favorite_count`: Monotonically updated statistics.
  - `version`: Version counter for optimistic locking.
- `listing_images`: Reference to media assets, sort order, cover flag.
- `listing_favorites`: Unique pair `(user_id, listing_id)`.
- `listing_status_history`: Audit trail of all listing status transitions.

### D. Real-Time Chat & Messages
- `conversations`: 1:1 conversation between buyer and seller tied to a specific listing. Unique constraint on `(listing_id, buyer_user_id)`.
- `conversation_participants`: Participant state (last read sequence number, mute preferences, left timestamp).
- `chat_messages`:
  - `conversation_id`: Associated conversation.
  - `sender_user_id`: Sending user.
  - `client_message_id`: Client-generated UUID for idempotency and deduplication.
  - `sequence_no`: Monotonically increasing sequence number per conversation.
  - `message_type`: `TEXT`, `IMAGE`, `OFFER_PROPOSAL`, `MEETING_PROPOSAL`, `SYSTEM_NOTICE`.
  - `body`: Text payload or structured JSON.
  - `created_at`: Server persistence timestamp.
- `chat_attachments`: Links message to media assets.

### E. Deals, Offers & Mutual Reviews
- `price_offers`: Price negotiation record with proposed amount, status (`PENDING`, `ACCEPTED`, `DECLINED`, `EXPIRED`).
- `deals`:
  - State machine: `INITIATED`, `NEGOTIATING`, `MEETING_PLANNED`, `RESERVED`, `COMPLETED_CLAIMED`, `COMPLETED_CONFIRMED`, `CANCELLED`, `DISPUTED`.
  - `seller_completed_at`: Timestamp when seller marked item as sold.
  - `buyer_confirmed_at`: Timestamp when buyer confirmed transaction completion.
- `deal_meetings`: Safe meeting place name, meeting time, and public coordinates.
- `deal_reviews`: Post-deal mutual feedback:
  - Rating: 1 to 5 stars.
  - Manner tags: JSON array of structured feedback tags.
  - Unique constraint on `(deal_id, reviewer_user_id)`.
- `reputation_snapshots`: Cached calculated manner temperature (default 36.5°C) and review count.

### F. Community, Business, Safety & Governance
- `community_posts`: Neighborhood-scoped discussion posts with category tag.
- `community_comments`: Threaded comments and replies.
- `community_likes`: Unique post upvotes per user.
- `businesses`: Local merchant profile (name, category, address, hours, verification status).
- `reports`: Content & user reports (target type, reason code, description, priority, status).
- `moderation_cases`: Triage and investigation workbench case with admin decision records.
- `enforcement_actions`: Sanctions (warning, restriction, suspension, ban).
- `audit_logs`: Tamper-proof log of administrative actions.

