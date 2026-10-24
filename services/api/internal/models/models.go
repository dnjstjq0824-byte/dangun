package models

import (
	"time"
)

// User Account
type User struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	PasswordHash    string     `json:"-"`
	Nickname        string     `json:"nickname"`
	AvatarAssetID   *string    `json:"avatar_asset_id,omitempty"`
	AvatarURL       string     `json:"avatar_url,omitempty"`
	Locale          string     `json:"locale"`
	CountryCode     string     `json:"country_code"`
	Status          string     `json:"status"` // PENDING_VERIFICATION, ACTIVE, RESTRICTED, SUSPENDED, DELETED
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	Profile         *UserProfile `json:"profile,omitempty"`
}

type UserProfile struct {
	UserID            string    `json:"user_id"`
	Bio               string    `json:"bio"`
	PublicDisplayName string    `json:"public_display_name"`
	PreferredCurrency string    `json:"preferred_currency"`
	Timezone          string    `json:"timezone"`
	MannerTemperature float64   `json:"manner_temperature"` // Base: 36.5
	ProfileVisibility string    `json:"profile_visibility"` // PUBLIC, NEIGHBORS_ONLY, PRIVATE
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type UserSession struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	RefreshTokenHash string     `json:"-"`
	DeviceID         string     `json:"device_id"`
	DeviceName       string     `json:"device_name"`
	IPAddress        string     `json:"ip_address"`
	ExpiresAt        time.Time  `json:"expires_at"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type UserBlock struct {
	ID             string    `json:"id"`
	BlockerUserID  string    `json:"blocker_user_id"`
	BlockedUserID  string    `json:"blocked_user_id"`
	ReasonOptional string    `json:"reason_optional,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// Geography
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Country struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	DefaultCurrency string `json:"default_currency"`
	DefaultLocale   string `json:"default_locale"`
	Status          string `json:"status"` // NOT_LAUNCHED, READ_ONLY, BETA, ACTIVE, SUSPENDED
}

type City struct {
	ID          string `json:"id"`
	CountryCode string `json:"country_code"`
	RegionID    string `json:"region_id,omitempty"`
	Name        string `json:"name"`
	Timezone    string `json:"timezone"`
	Status      string `json:"status"`
}

type Neighborhood struct {
	ID               string       `json:"id"`
	CityID           string       `json:"city_id"`
	Slug             string       `json:"slug"`
	Name             string       `json:"name"`
	BoundaryVertices []Coordinate `json:"boundary_vertices"`
	PublicCenter     Coordinate   `json:"public_center"`
	Status           string       `json:"status"`
}

type UserNeighborhood struct {
	ID                 string     `json:"id"`
	UserID             string     `json:"user_id"`
	NeighborhoodID     string     `json:"neighborhood_id"`
	NeighborhoodName   string     `json:"neighborhood_name,omitempty"`
	IsPrimary          bool       `json:"is_primary"`
	VerificationStatus string     `json:"verification_status"` // UNVERIFIED, VERIFIED, EXPIRED
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

type LocationVerificationAttempt struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	NeighborhoodID     string    `json:"neighborhood_id"`
	AccuracyMeters     float64   `json:"accuracy_meters"`
	VerificationResult string    `json:"verification_result"` // VERIFIED, OUT_OF_BOUNDS, LOW_ACCURACY
	RiskReasonCodes    string    `json:"risk_reason_codes,omitempty"`
	AttemptedAt        time.Time `json:"attempted_at"`
}

// Marketplace
type Category struct {
	ID        string `json:"id"`
	ParentID  string `json:"parent_id,omitempty"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	IconName  string `json:"icon_name"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

type Listing struct {
	ID               string      `json:"id"`
	SellerUserID     string      `json:"seller_user_id"`
	SellerNickname   string      `json:"seller_nickname,omitempty"`
	SellerMannerTemp float64     `json:"seller_manner_temp,omitempty"`
	NeighborhoodID   string      `json:"neighborhood_id"`
	NeighborhoodName string      `json:"neighborhood_name,omitempty"`
	CategoryID       string      `json:"category_id"`
	CategoryName     string      `json:"category_name,omitempty"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	Condition        string      `json:"condition"` // NEW, LIKE_NEW, GOOD, FAIR, FOR_PARTS
	PriceMinor       int64       `json:"price_minor"`
	CurrencyCode     string      `json:"currency_code"`
	IsFree           bool        `json:"is_free"`
	IsNegotiable     bool        `json:"is_negotiable"`
	ListingStatus    string      `json:"listing_status"` // DRAFT, ACTIVE, RESERVED, SOLD, ARCHIVED, REMOVED
	PublicLocation   *Coordinate `json:"public_location,omitempty"`
	Images           []string    `json:"images"`
	ViewCount        int64       `json:"view_count"`
	FavoriteCount    int64       `json:"favorite_count"`
	ChatCount        int64       `json:"chat_count"`
	IsFavorited      bool        `json:"is_favorited,omitempty"`
	DistanceMeters   float64     `json:"distance_meters,omitempty"`
	Version          int64       `json:"version"`
	PublishedAt      time.Time   `json:"published_at"`
	ReservedAt       *time.Time  `json:"reserved_at,omitempty"`
	SoldAt           *time.Time  `json:"sold_at,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type ListingFavorite struct {
	UserID    string    `json:"user_id"`
	ListingID string    `json:"listing_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Real-Time Chat
type Conversation struct {
	ID               string    `json:"id"`
	ListingID        string    `json:"listing_id"`
	ListingTitle     string    `json:"listing_title,omitempty"`
	ListingPrice     int64     `json:"listing_price,omitempty"`
	ListingImage     string    `json:"listing_image,omitempty"`
	ListingStatus    string    `json:"listing_status,omitempty"`
	BuyerUserID      string    `json:"buyer_user_id"`
	BuyerNickname    string    `json:"buyer_nickname,omitempty"`
	SellerUserID     string    `json:"seller_user_id"`
	SellerNickname   string    `json:"seller_nickname,omitempty"`
	Status           string    `json:"status"` // ACTIVE, MUTED, ARCHIVED, BLOCKED
	LastMessage      string    `json:"last_message,omitempty"`
	LastMessageAt    time.Time `json:"last_message_at"`
	UnreadCount      int       `json:"unread_count"`
	LastReadSequence int64     `json:"last_read_sequence"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ChatMessage struct {
	ID              string    `json:"id"`
	ConversationID  string    `json:"conversation_id"`
	SenderUserID    string    `json:"sender_user_id"`
	ClientMessageID string    `json:"client_message_id"`
	SequenceNo      int64     `json:"sequence_no"`
	MessageType     string    `json:"message_type"` // TEXT, IMAGE, OFFER_PROPOSAL, MEETING_PROPOSAL, SYSTEM_NOTICE
	Body            string    `json:"body"`
	Status          string    `json:"status"` // STORED, DELIVERED, READ
	CreatedAt       time.Time `json:"created_at"`
}

// Deals & Offers
type PriceOffer struct {
	ID           string    `json:"id"`
	ListingID    string    `json:"listing_id"`
	BuyerUserID  string    `json:"buyer_user_id"`
	SellerUserID string    `json:"seller_user_id"`
	AmountMinor  int64     `json:"amount_minor"`
	CurrencyCode string    `json:"currency_code"`
	Status       string    `json:"status"` // PENDING, ACCEPTED, DECLINED, EXPIRED
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Deal struct {
	ID                     string     `json:"id"`
	ListingID              string     `json:"listing_id"`
	BuyerUserID            string     `json:"buyer_user_id"`
	SellerUserID           string     `json:"seller_user_id"`
	ConversationID         string     `json:"conversation_id"`
	AgreedPriceMinor       *int64     `json:"agreed_price_minor,omitempty"`
	CurrencyCode           string     `json:"currency_code"`
	Status                 string     `json:"status"` // INITIATED, NEGOTIATING, MEETING_PLANNED, RESERVED, COMPLETED_CLAIMED, COMPLETED_CONFIRMED, CANCELLED, DISPUTED
	MeetingTime            *time.Time `json:"meeting_time,omitempty"`
	MeetingPlaceName       string     `json:"meeting_place_name,omitempty"`
	SellerCompletedAt      *time.Time `json:"seller_completed_at,omitempty"`
	BuyerConfirmedAt       *time.Time `json:"buyer_confirmed_at,omitempty"`
	CancelledAt            *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	CanReview              bool       `json:"can_review,omitempty"`
	HasBuyerReviewed       bool       `json:"has_buyer_reviewed,omitempty"`
	HasSellerReviewed      bool       `json:"has_seller_reviewed,omitempty"`
}

type DealReview struct {
	ID             string    `json:"id"`
	DealID         string    `json:"deal_id"`
	ReviewerUserID string    `json:"reviewer_user_id"`
	ReviewedUserID string    `json:"reviewed_user_id"`
	Rating         int       `json:"rating"` // 1-5
	Tags           []string  `json:"tags"`
	Comment        string    `json:"comment,omitempty"`
	Status         string    `json:"status"` // PUBLISHED, REPORTED, HIDDEN
	CreatedAt      time.Time `json:"created_at"`
}

// Community
type CommunityPost struct {
	ID               string    `json:"id"`
	AuthorUserID     string    `json:"author_user_id"`
	AuthorNickname   string    `json:"author_nickname,omitempty"`
	NeighborhoodID   string    `json:"neighborhood_id"`
	NeighborhoodName string    `json:"neighborhood_name,omitempty"`
	Category         string    `json:"category"` // 질문, 정보, 추천, 분실물, 일상
	Title            string    `json:"title"`
	Body             string    `json:"body"`
	ViewCount        int64     `json:"view_count"`
	LikeCount        int64     `json:"like_count"`
	CommentCount     int64     `json:"comment_count"`
	IsLiked          bool      `json:"is_liked,omitempty"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CommunityComment struct {
	ID              string    `json:"id"`
	PostID          string    `json:"post_id"`
	AuthorUserID    string    `json:"author_user_id"`
	AuthorNickname  string    `json:"author_nickname,omitempty"`
	ParentCommentID string    `json:"parent_comment_id,omitempty"`
	Body            string    `json:"body"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Business
type Business struct {
	ID                 string    `json:"id"`
	OwnerUserID        string    `json:"owner_user_id"`
	Name               string    `json:"name"`
	CategoryID         string    `json:"category_id"`
	CategoryName       string    `json:"category_name,omitempty"`
	CountryCode        string    `json:"country_code"`
	CityID             string    `json:"city_id"`
	NeighborhoodID     string    `json:"neighborhood_id"`
	NeighborhoodName   string    `json:"neighborhood_name,omitempty"`
	Description        string    `json:"description"`
	PublicAddress      string    `json:"public_address"`
	PhoneNumber        string    `json:"phone_number,omitempty"`
	ContactVisibility  bool      `json:"contact_visibility"`
	VerificationStatus string    `json:"verification_status"` // PENDING, VERIFIED, REJECTED
	Status             string    `json:"status"`              // ACTIVE, CLOSED
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Safety, Moderation & Audit
type Report struct {
	ID             string     `json:"id"`
	ReporterUserID string     `json:"reporter_user_id"`
	TargetType     string     `json:"target_type"` // LISTING, USER, CHAT, COMMUNITY_POST, BUSINESS
	TargetID       string     `json:"target_id"`
	ReasonCode     string     `json:"reason_code"`
	Description    string     `json:"description,omitempty"`
	Priority       string     `json:"priority"` // LOW, NORMAL, HIGH, URGENT
	Status         string     `json:"status"`   // PENDING, INVESTIGATING, ACTIONED, DISMISSED
	CreatedAt      time.Time  `json:"created_at"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
}

type ModerationCase struct {
	ID              string     `json:"id"`
	ReportID        *string    `json:"report_id,omitempty"`
	TargetType      string     `json:"target_type"`
	TargetID        string     `json:"target_id"`
	AssignedAdminID *string    `json:"assigned_admin_id,omitempty"`
	Status          string     `json:"status"` // OPEN, IN_REVIEW, RESOLVED, CLOSED
	Decision        *string    `json:"decision,omitempty"`
	DecisionReason  *string    `json:"decision_reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
}

type EnforcementAction struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	ActionType   string     `json:"action_type"` // WARNING, RESTRICTION, SUSPENSION, BAN
	ReasonCode   string     `json:"reason_code"`
	StartAt      time.Time  `json:"start_at"`
	EndAt        *time.Time `json:"end_at,omitempty"`
	IssuedBy     string     `json:"issued_by"`
	AppealStatus string     `json:"appeal_status"` // NONE, PENDING, UPHELD, REVERSED
	CreatedAt    time.Time  `json:"created_at"`
}

type Appeal struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	EnforcementActionID string     `json:"enforcement_action_id"`
	Description         string     `json:"description"`
	Status              string     `json:"status"` // PENDING, APPROVED, REJECTED
	ReviewedBy          *string    `json:"reviewed_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	ResolvedAt          *time.Time `json:"resolved_at,omitempty"`
}

type AuditLog struct {
	ID               string    `json:"id"`
	ActorType        string    `json:"actor_type"` // USER, ADMIN, SYSTEM
	ActorID          string    `json:"actor_id"`
	Action           string    `json:"action"`
	ResourceType     string    `json:"resource_type"`
	ResourceID       string    `json:"resource_id"`
	Result           string    `json:"result"` // SUCCESS, FAILURE, DENIED
	MetadataRedacted string    `json:"metadata_redacted,omitempty"`
	RequestID        string    `json:"request_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// Notifications
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    string    `json:"status"` // UNREAD, READ
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationPreference struct {
	UserID           string    `json:"user_id"`
	ChatEnabled      bool      `json:"chat_enabled"`
	DealEnabled      bool      `json:"deal_enabled"`
	CommunityEnabled bool      `json:"community_enabled"`
	MarketingEnabled bool      `json:"marketing_enabled"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Admin Dashboard Summary
type DashboardStats struct {
	TotalUsers            int64 `json:"total_users"`
	ActiveUsers           int64 `json:"active_users"`
	TotalListings         int64 `json:"total_listings"`
	ActiveListings        int64 `json:"active_listings"`
	CompletedDeals        int64 `json:"completed_deals"`
	ActiveConversations   int64 `json:"active_conversations"`
	PendingReports        int64 `json:"pending_reports"`
	VerifiedNeighborhoods int64 `json:"verified_neighborhoods"`
	LiveWSSessions        int64 `json:"live_ws_sessions"`
}
