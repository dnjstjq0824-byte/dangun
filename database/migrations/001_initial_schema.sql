-- BN NEAR Master Database Schema
-- Target: PostgreSQL 16+ with PostGIS 3.4+
-- Encoding: UTF-8

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "citext";

-- ============================================================================
-- 1. USERS & IDENTITY
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email CITEXT UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    avatar_asset_id UUID NULL,
    locale VARCHAR(10) NOT NULL DEFAULT 'ko-KR',
    country_code VARCHAR(2) NOT NULL DEFAULT 'KR',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('PENDING_VERIFICATION', 'ACTIVE', 'RESTRICTED', 'SUSPENDED', 'DELETED')),
    email_verified_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT NULL,
    public_display_name VARCHAR(100) NULL,
    preferred_currency VARCHAR(3) NOT NULL DEFAULT 'KRW',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Seoul',
    manner_temperature NUMERIC(4, 1) NOT NULL DEFAULT 36.5 CHECK (manner_temperature >= 0.0 AND manner_temperature <= 99.9),
    profile_visibility VARCHAR(20) NOT NULL DEFAULT 'PUBLIC' CHECK (profile_visibility IN ('PUBLIC', 'NEIGHBORS_ONLY', 'PRIVATE')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL,
    device_id VARCHAR(100) NOT NULL,
    device_name VARCHAR(100) NULL,
    ip_address VARCHAR(45) NULL,
    user_agent TEXT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_blocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    blocker_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason_optional VARCHAR(255) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_user_blocks UNIQUE (blocker_user_id, blocked_user_id),
    CONSTRAINT chk_no_self_block CHECK (blocker_user_id <> blocked_user_id)
);

CREATE TABLE IF NOT EXISTS consent_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    consent_type VARCHAR(50) NOT NULL,
    policy_version VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'GRANTED' CHECK (status IN ('GRANTED', 'WITHDRAWN')),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    withdrawn_at TIMESTAMPTZ NULL
);

-- ============================================================================
-- 2. GEOGRAPHY & NEIGHBORHOODS
-- ============================================================================

CREATE TABLE IF NOT EXISTS countries (
    code VARCHAR(2) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    default_currency VARCHAR(3) NOT NULL,
    default_locale VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('NOT_LAUNCHED', 'READ_ONLY', 'BETA', 'ACTIVE', 'SUSPENDED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS regions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_code VARCHAR(2) NOT NULL REFERENCES countries(code) ON DELETE RESTRICT,
    parent_region_id UUID NULL REFERENCES regions(id),
    name VARCHAR(100) NOT NULL,
    region_type VARCHAR(50) NOT NULL DEFAULT 'PROVINCE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_code VARCHAR(2) NOT NULL REFERENCES countries(code) ON DELETE RESTRICT,
    region_id UUID NULL REFERENCES regions(id),
    name VARCHAR(100) NOT NULL,
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'BETA', 'INACTIVE')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS neighborhoods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    city_id UUID NOT NULL REFERENCES cities(id) ON DELETE RESTRICT,
    parent_neighborhood_id UUID NULL REFERENCES neighborhoods(id),
    slug VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    boundary_geom GEOMETRY(Polygon, 4326) NOT NULL,
    public_center_geom GEOMETRY(Point, 4326) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'BETA', 'SUSPENDED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_neighborhoods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    neighborhood_id UUID NOT NULL REFERENCES neighborhoods(id) ON DELETE RESTRICT,
    is_primary BOOLEAN NOT NULL DEFAULT TRUE,
    verification_status VARCHAR(20) NOT NULL DEFAULT 'UNVERIFIED' CHECK (verification_status IN ('UNVERIFIED', 'VERIFIED', 'EXPIRED')),
    verified_at TIMESTAMPTZ NULL,
    expires_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_user_neighborhood UNIQUE (user_id, neighborhood_id)
);

CREATE TABLE IF NOT EXISTS location_verification_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    neighborhood_id UUID NOT NULL REFERENCES neighborhoods(id) ON DELETE RESTRICT,
    accuracy_meters NUMERIC(6, 2) NOT NULL,
    verification_result VARCHAR(30) NOT NULL CHECK (verification_result IN ('VERIFIED', 'OUT_OF_BOUNDS', 'LOW_ACCURACY', 'RATE_LIMITED', 'SUSPICIOUS_MOCK')),
    risk_reason_codes VARCHAR(100) NULL,
    attempted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS country_feature_policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_code VARCHAR(2) NOT NULL REFERENCES countries(code) ON DELETE CASCADE,
    feature_code VARCHAR(50) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    policy_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    reviewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_country_feature UNIQUE (country_code, feature_code)
);

-- ============================================================================
-- 3. MARKETPLACE & LISTINGS
-- ============================================================================

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id UUID NULL REFERENCES categories(id) ON DELETE SET NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    icon_name VARCHAR(50) NOT NULL DEFAULT 'tag',
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS listings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    seller_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    neighborhood_id UUID NOT NULL REFERENCES neighborhoods(id) ON DELETE RESTRICT,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    title VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    condition VARCHAR(20) NOT NULL DEFAULT 'GOOD' CHECK (condition IN ('NEW', 'LIKE_NEW', 'GOOD', 'FAIR', 'FOR_PARTS')),
    price_minor BIGINT NOT NULL DEFAULT 0 CHECK (price_minor >= 0),
    currency_code VARCHAR(3) NOT NULL DEFAULT 'KRW',
    is_free BOOLEAN NOT NULL DEFAULT FALSE,
    is_negotiable BOOLEAN NOT NULL DEFAULT TRUE,
    listing_status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (listing_status IN ('DRAFT', 'PENDING_REVIEW', 'ACTIVE', 'RESERVED', 'SOLD', 'ARCHIVED', 'REJECTED', 'REMOVED')),
    public_location_geom GEOMETRY(Point, 4326) NULL,
    view_count BIGINT NOT NULL DEFAULT 0,
    favorite_count BIGINT NOT NULL DEFAULT 0,
    version BIGINT NOT NULL DEFAULT 1,
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reserved_at TIMESTAMPTZ NULL,
    sold_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT chk_free_price CHECK (NOT (is_free = TRUE AND price_minor > 0))
);

CREATE TABLE IF NOT EXISTS media_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    storage_key VARCHAR(255) NOT NULL,
    public_url VARCHAR(500) NOT NULL,
    mime_type VARCHAR(50) NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    sha256 VARCHAR(64) NOT NULL,
    moderation_status VARCHAR(20) NOT NULL DEFAULT 'APPROVED' CHECK (moderation_status IN ('PENDING', 'APPROVED', 'REJECTED')),
    scan_status VARCHAR(20) NOT NULL DEFAULT 'CLEAN' CHECK (scan_status IN ('PENDING', 'CLEAN', 'MALICIOUS')),
    width INT NULL,
    height INT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS listing_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    media_asset_id UUID NOT NULL REFERENCES media_assets(id) ON DELETE CASCADE,
    sort_order INT NOT NULL DEFAULT 0,
    is_cover BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS listing_attributes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    attribute_key VARCHAR(50) NOT NULL,
    attribute_value VARCHAR(255) NOT NULL,
    CONSTRAINT uq_listing_attribute UNIQUE (listing_id, attribute_key)
);

CREATE TABLE IF NOT EXISTS listing_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    from_status VARCHAR(20) NOT NULL,
    to_status VARCHAR(20) NOT NULL,
    changed_by UUID NOT NULL REFERENCES users(id),
    reason VARCHAR(255) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS listing_favorites (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, listing_id)
);

CREATE TABLE IF NOT EXISTS saved_searches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    query_text VARCHAR(100) NOT NULL,
    filters_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    notification_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- 4. REAL-TIME CHAT
-- ============================================================================

CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    buyer_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'MUTED', 'ARCHIVED', 'BLOCKED')),
    last_message_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_conversation_listing_buyer UNIQUE (listing_id, buyer_user_id),
    CONSTRAINT chk_no_self_conversation CHECK (buyer_user_id <> seller_user_id)
);

CREATE TABLE IF NOT EXISTS conversation_participants (
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_sequence BIGINT NOT NULL DEFAULT 0,
    muted_until TIMESTAMPTZ NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at TIMESTAMPTZ NULL,
    PRIMARY KEY (conversation_id, user_id)
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_message_id UUID NOT NULL,
    sequence_no BIGINT NOT NULL,
    message_type VARCHAR(30) NOT NULL DEFAULT 'TEXT' CHECK (message_type IN ('TEXT', 'IMAGE', 'OFFER_PROPOSAL', 'MEETING_PROPOSAL', 'SYSTEM_NOTICE')),
    body TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'STORED' CHECK (status IN ('STORED', 'DELIVERED', 'READ', 'DELETED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    edited_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_chat_message_client UNIQUE (conversation_id, sender_user_id, client_message_id),
    CONSTRAINT uq_chat_message_sequence UNIQUE (conversation_id, sequence_no)
);

CREATE TABLE IF NOT EXISTS chat_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id UUID NOT NULL REFERENCES chat_messages(id) ON DELETE CASCADE,
    media_asset_id UUID NOT NULL REFERENCES media_assets(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- 5. DEALS, OFFERS & REVIEWS
-- ============================================================================

CREATE TABLE IF NOT EXISTS price_offers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    buyer_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount_minor BIGINT NOT NULL CHECK (amount_minor >= 0),
    currency_code VARCHAR(3) NOT NULL DEFAULT 'KRW',
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'ACCEPTED', 'DECLINED', 'CANCELLED', 'EXPIRED')),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS deals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE RESTRICT,
    buyer_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    seller_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE RESTRICT,
    agreed_price_minor BIGINT NULL CHECK (agreed_price_minor IS NULL OR agreed_price_minor >= 0),
    currency_code VARCHAR(3) NOT NULL DEFAULT 'KRW',
    status VARCHAR(30) NOT NULL DEFAULT 'INITIATED' CHECK (status IN ('INITIATED', 'NEGOTIATING', 'MEETING_PLANNED', 'RESERVED', 'COMPLETED_CLAIMED', 'COMPLETED_CONFIRMED', 'CANCELLED', 'DISPUTED')),
    seller_completed_at TIMESTAMPTZ NULL,
    buyer_confirmed_at TIMESTAMPTZ NULL,
    cancelled_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS deal_meetings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    deal_id UUID NOT NULL REFERENCES deals(id) ON DELETE CASCADE,
    meeting_time TIMESTAMPTZ NOT NULL,
    public_meeting_place_name VARCHAR(150) NOT NULL,
    public_meeting_place_geom GEOMETRY(Point, 4326) NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'SCHEDULED' CHECK (status IN ('SCHEDULED', 'COMPLETED', 'CANCELLED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS deal_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    deal_id UUID NOT NULL REFERENCES deals(id) ON DELETE CASCADE,
    from_status VARCHAR(30) NOT NULL,
    to_status VARCHAR(30) NOT NULL,
    changed_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS deal_reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    deal_id UUID NOT NULL REFERENCES deals(id) ON DELETE CASCADE,
    reviewer_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reviewed_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    tags_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    comment TEXT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PUBLISHED' CHECK (status IN ('PUBLISHED', 'REPORTED', 'HIDDEN')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_deal_reviewer UNIQUE (deal_id, reviewer_user_id),
    CONSTRAINT chk_no_self_review CHECK (reviewer_user_id <> reviewed_user_id)
);

CREATE TABLE IF NOT EXISTS reputation_snapshots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    review_count INT NOT NULL DEFAULT 0,
    summary_score NUMERIC(4, 1) NOT NULL DEFAULT 36.5,
    calculation_version INT NOT NULL DEFAULT 1,
    computed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- 6. LOCAL COMMUNITY
-- ============================================================================

CREATE TABLE IF NOT EXISTS community_posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    neighborhood_id UUID NOT NULL REFERENCES neighborhoods(id) ON DELETE RESTRICT,
    category VARCHAR(50) NOT NULL DEFAULT 'GENERAL',
    title VARCHAR(150) NOT NULL,
    body TEXT NOT NULL,
    view_count BIGINT NOT NULL DEFAULT 0,
    like_count BIGINT NOT NULL DEFAULT 0,
    comment_count BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'FLAGGED', 'REMOVED', 'DELETED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS community_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES community_posts(id) ON DELETE CASCADE,
    author_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id UUID NULL REFERENCES community_comments(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'FLAGGED', 'REMOVED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS community_likes (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES community_posts(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

CREATE TABLE IF NOT EXISTS community_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    neighborhood_id UUID NOT NULL REFERENCES neighborhoods(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'SUSPENDED')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS group_memberships (
    group_id UUID NOT NULL REFERENCES community_groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'MEMBER' CHECK (role IN ('OWNER', 'MODERATOR', 'MEMBER')),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INVITED', 'BANNED')),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (group_id, user_id)
);

-- ============================================================================
-- 7. LOCAL BUSINESSES
-- ============================================================================

CREATE TABLE IF NOT EXISTS businesses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    name VARCHAR(150) NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    country_code VARCHAR(2) NOT NULL REFERENCES countries(code),
    city_id UUID NOT NULL REFERENCES cities(id),
    neighborhood_id UUID NOT NULL REFERENCES neighborhoods(id),
    description TEXT NULL,
    public_address VARCHAR(255) NOT NULL,
    public_location_geom GEOMETRY(Point, 4326) NULL,
    phone_number VARCHAR(30) NULL,
    contact_visibility BOOLEAN NOT NULL DEFAULT TRUE,
    verification_status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (verification_status IN ('PENDING', 'VERIFIED', 'REJECTED', 'SUSPENDED')),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'CLOSED', 'SUSPENDED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS business_hours (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    weekday INT NOT NULL CHECK (weekday >= 0 AND weekday <= 6),
    opens_at TIME NULL,
    closes_at TIME NULL,
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT uq_business_weekday UNIQUE (business_id, weekday)
);

CREATE TABLE IF NOT EXISTS business_posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    title VARCHAR(150) NOT NULL,
    body TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PUBLISHED' CHECK (status IN ('DRAFT', 'PUBLISHED', 'REMOVED')),
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- 8. SAFETY, MODERATION & APPEALS
-- ============================================================================

CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_type VARCHAR(30) NOT NULL CHECK (target_type IN ('LISTING', 'USER', 'CHAT_CONVERSATION', 'COMMUNITY_POST', 'COMMUNITY_COMMENT', 'BUSINESS', 'REVIEW')),
    target_id UUID NOT NULL,
    reason_code VARCHAR(50) NOT NULL,
    description TEXT NULL,
    priority VARCHAR(10) NOT NULL DEFAULT 'NORMAL' CHECK (priority IN ('LOW', 'NORMAL', 'HIGH', 'URGENT')),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'INVESTIGATING', 'ACTIONED', 'DISMISSED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS moderation_cases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NULL REFERENCES reports(id) ON DELETE SET NULL,
    target_type VARCHAR(30) NOT NULL,
    target_id UUID NOT NULL,
    assigned_admin_id UUID NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'IN_REVIEW', 'RESOLVED', 'CLOSED')),
    decision VARCHAR(30) NULL CHECK (decision IN ('NO_VIOLATION', 'WARNING_ISSUED', 'CONTENT_REMOVED', 'USER_RESTRICTED', 'USER_SUSPENDED')),
    decision_reason TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS enforcement_actions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action_type VARCHAR(30) NOT NULL CHECK (action_type IN ('WARNING', 'RESTRICTION', 'SUSPENSION', 'PERMANENT_BAN')),
    reason_code VARCHAR(50) NOT NULL,
    start_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_at TIMESTAMPTZ NULL,
    issued_by UUID NOT NULL REFERENCES users(id),
    appeal_status VARCHAR(20) NOT NULL DEFAULT 'NONE' CHECK (appeal_status IN ('NONE', 'PENDING', 'UPHELD', 'REVERSED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS appeals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    enforcement_action_id UUID NOT NULL REFERENCES enforcement_actions(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
    reviewed_by UUID NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ NULL
);

-- ============================================================================
-- 9. NOTIFICATIONS
-- ============================================================================

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(150) NOT NULL,
    body TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'UNREAD' CHECK (status IN ('UNREAD', 'READ', 'ARCHIVED')),
    metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS push_devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(20) NOT NULL CHECK (platform IN ('ANDROID', 'IOS', 'WEB')),
    push_token_reference VARCHAR(255) NOT NULL,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chat_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    deal_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    community_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    marketing_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- 10. ADMIN, RBAC & AUDIT LOGS
-- ============================================================================

CREATE TABLE IF NOT EXISTS admin_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL
);

CREATE TABLE IF NOT EXISTS admin_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS admin_user_roles (
    admin_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES admin_roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (admin_user_id, role_id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_type VARCHAR(20) NOT NULL DEFAULT 'ADMIN' CHECK (actor_type IN ('USER', 'ADMIN', 'SYSTEM')),
    actor_id UUID NOT NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id VARCHAR(100) NOT NULL,
    result VARCHAR(20) NOT NULL DEFAULT 'SUCCESS' CHECK (result IN ('SUCCESS', 'FAILURE', 'DENIED')),
    metadata_redacted JSONB NOT NULL DEFAULT '{}'::jsonb,
    request_id VARCHAR(100) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(100) UNIQUE NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    environment VARCHAR(20) NOT NULL DEFAULT 'production',
    country_code VARCHAR(2) NULL REFERENCES countries(code),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PROCESSED', 'FAILED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ NULL
);

-- ============================================================================
-- 11. CRITICAL PERFORMANCE & SPATIAL INDEXES
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_listings_geo ON listings USING GIST (public_location_geom);
CREATE INDEX IF NOT EXISTS idx_neighborhoods_boundary ON neighborhoods USING GIST (boundary_geom);
CREATE INDEX IF NOT EXISTS idx_listings_neighborhood_status ON listings (neighborhood_id, listing_status, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_listings_search_trgm ON listings USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_chat_messages_conversation ON chat_messages (conversation_id, sequence_no);
CREATE INDEX IF NOT EXISTS idx_deals_listing_status ON deals (listing_id, status);
CREATE INDEX IF NOT EXISTS idx_reports_status_priority ON reports (status, priority, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_actor ON audit_logs (actor_id, created_at DESC);

