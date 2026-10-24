# BN NEAR — 12. Database Migrations Guide

**Order ID:** BN-NEAR-MASTER-001  
**Engine:** PostgreSQL 16 + PostGIS 3.4  

---

## 1. Migration Strategy
All schema changes are tracked as linear, forward-only SQL migration scripts stored in `database/migrations/`:
- `001_initial_postgis_and_extensions.sql`: Enables `uuid-ossp`, `postgis`, `pg_trgm`, `citext`.
- `002_core_identity_and_geo.sql`: Creates `users`, `user_profiles`, `user_sessions`, `countries`, `cities`, `neighborhoods`, `user_neighborhoods`.
- `003_marketplace_and_search.sql`: Creates `categories`, `listings`, `listing_images`, `listing_favorites`, spatial and trgm indexes.
- `004_realtime_chat.sql`: Creates `conversations`, `conversation_participants`, `chat_messages`, unique constraints and sequence index.
- `005_deals_and_reputation.sql`: Creates `price_offers`, `deals`, `deal_meetings`, `deal_reviews`, `reputation_snapshots`.
- `006_community_and_business.sql`: Creates `community_posts`, `community_comments`, `community_likes`, `businesses`, `business_hours`.
- `007_safety_admin_and_audit.sql`: Creates `reports`, `moderation_cases`, `enforcement_actions`, `audit_logs`, `admin_roles`.

---

## 2. Spatial Indexing Rules
To guarantee sub-10ms response times for neighborhood radius queries:
```sql
CREATE INDEX idx_listings_geo ON listings USING GIST (public_location_geom);
CREATE INDEX idx_neighborhoods_boundary ON neighborhoods USING GIST (boundary_geom);
CREATE INDEX idx_listings_neighborhood_status ON listings (neighborhood_id, listing_status, published_at DESC);
CREATE INDEX idx_chat_messages_conversation ON chat_messages (conversation_id, sequence_no);
CREATE INDEX idx_deals_listing_status ON deals (listing_id, status);
CREATE INDEX idx_reports_status_priority ON reports (status, priority, created_at);
```

---

## 3. Seed Fixtures
Production-like seed fixtures are loaded from `database/fixtures/seed_data.sql` and include:
- Countries: KR (South Korea), US (United States), JP (Japan).
- Real cities & neighborhoods with valid polygon boundaries:
  - Incheon: Cheongna-dong (청라동), Songdo-dong (송도동).
  - Seoul: Yeoksam-dong (역삼동), Yeonnam-dong (연남동).
  - Tokyo: Shibuya (渋谷), Shinjuku (新宿).
  - New York: Brooklyn - Williamsburg, Manhattan - SoHo.
- 12 top-level product categories with localized names and SVG icons.
- Pre-seeded test accounts for buyers, sellers, merchants, and platform administrators.

