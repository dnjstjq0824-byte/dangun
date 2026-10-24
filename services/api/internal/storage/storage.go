package storage

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/booking-network/bn-near/services/api/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound             = errors.New("resource not found")
	ErrAlreadyExists        = errors.New("resource already exists")
	ErrUnauthorized         = errors.New("unauthorized action")
	ErrConflict             = errors.New("conflict: resource state does not permit operation")
	ErrBlocked              = errors.New("communication blocked by user")
	ErrInvalidInput         = errors.New("invalid input data")
	ErrVerificationRequired = errors.New("neighborhood verification required")
)

type MemoryStore struct {
	mu sync.RWMutex

	Users                 map[string]*models.User
	Profiles              map[string]*models.UserProfile
	Sessions              map[string]*models.UserSession
	Blocks                map[string]*models.UserBlock // blocker_blocked -> Block
	Countries             map[string]*models.Country
	Cities                map[string]*models.City
	Neighborhoods         map[string]*models.Neighborhood
	UserNeighborhoods     map[string]*models.UserNeighborhood // user_neighborhood -> UserNeighborhood
	VerificationAttempts  []*models.LocationVerificationAttempt
	Categories            map[string]*models.Category
	Listings              map[string]*models.Listing
	Favorites             map[string]*models.ListingFavorite // user_listing -> Favorite
	Conversations         map[string]*models.Conversation
	Messages              map[string][]*models.ChatMessage // conv_id -> messages
	MessageClientIndex    map[string]*models.ChatMessage   // conv_sender_clientMsgId -> Message
	Offers                map[string]*models.PriceOffer
	Deals                 map[string]*models.Deal
	DealReviews           map[string]*models.DealReview // deal_reviewer -> Review
	CommunityPosts        map[string]*models.CommunityPost
	CommunityComments     map[string][]*models.CommunityComment // post_id -> comments
	CommunityLikes        map[string]bool                       // user_post -> bool
	Businesses            map[string]*models.Business
	Reports               map[string]*models.Report
	ModerationCases       map[string]*models.ModerationCase
	Enforcements          map[string]*models.EnforcementAction
	Appeals               map[string]*models.Appeal
	Notifications         map[string][]*models.Notification // user_id -> notifications
	NotificationPrefs     map[string]*models.NotificationPreference
	AuditLogs             []*models.AuditLog
	LiveWSSessionsCount   int64
}

type Storage interface {
	// Users
	CreateUser(user *models.User, rawPassword string) (*models.User, error)
	AuthenticateUser(email, password string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUserProfile(userID, bio, displayName string) (*models.UserProfile, error)
	SoftDeleteUser(userID string) error
	CreateSession(session *models.UserSession) error
	GetSession(sessionID string) (*models.UserSession, error)
	RevokeSession(sessionID string) error
	BlockUser(blockerID, blockedID, reason string) error
	UnblockUser(blockerID, blockedID string) error
	IsBlocked(userA, userB string) bool

	// Geo & Verification
	ListCountries() ([]*models.Country, error)
	ListCities(countryCode string) ([]*models.City, error)
	GetNeighborhoodByID(id string) (*models.Neighborhood, error)
	SearchNeighborhoods(query string, cityID string, lat, lng *float64) ([]*models.Neighborhood, error)
	VerifyNeighborhood(userID, neighborhoodID string, lat, lng, accuracy float64) (*models.UserNeighborhood, error)
	GetUserNeighborhoods(userID string) ([]*models.UserNeighborhood, error)
	IsUserVerifiedInNeighborhood(userID, neighborhoodID string) bool

	// Listings
	CreateListing(listing *models.Listing) (*models.Listing, error)
	GetListingByID(id, viewerUserID string) (*models.Listing, error)
	UpdateListing(listing *models.Listing) (*models.Listing, error)
	DeleteListing(id, userID string) error
	ListListings(neighborhoodID, categoryID, condition string, minPrice, maxPrice *int64, isFree *bool, lat, lng *float64, radiusMeters float64, sortOrder string, viewerUserID string) ([]*models.Listing, error)
	ToggleFavorite(userID, listingID string) (bool, error)
	GetUserFavorites(userID string) ([]*models.Listing, error)
	GetUserListings(userID string) ([]*models.Listing, error)

	// Chat
	GetOrCreateConversation(listingID, buyerUserID string) (*models.Conversation, error)
	GetUserConversations(userID string) ([]*models.Conversation, error)
	GetConversationByID(id, userID string) (*models.Conversation, error)
	SaveMessage(msg *models.ChatMessage) (*models.ChatMessage, error)
	GetMessages(conversationID string, sinceSequence int64) ([]*models.ChatMessage, error)
	MarkMessagesRead(conversationID, userID string, sequenceNo int64) error

	// Deals & Offers
	CreateOffer(offer *models.PriceOffer) (*models.PriceOffer, error)
	GetOfferByID(id string) (*models.PriceOffer, error)
	UpdateOfferStatus(id, userID, newStatus string) (*models.PriceOffer, error)
	GetDealByID(id string) (*models.Deal, error)
	GetDealByListingID(listingID string) (*models.Deal, error)
	ReserveDeal(listingID, sellerUserID, buyerUserID string) (*models.Deal, error)
	ScheduleMeeting(dealID, userID, placeName string, meetingTime time.Time) (*models.Deal, error)
	SellerCompleteDeal(dealID, sellerUserID string) (*models.Deal, error)
	BuyerConfirmDeal(dealID, buyerUserID string) (*models.Deal, error)
	CancelDeal(dealID, userID string) (*models.Deal, error)
	GetUserDeals(userID string) ([]*models.Deal, error)

	// Reviews & Reputation
	CreateDealReview(review *models.DealReview) (*models.DealReview, error)
	GetUserReviews(userID string) ([]*models.DealReview, error)
	GetUserReputation(userID string) (*models.UserProfile, error)

	// Community
	CreateCommunityPost(post *models.CommunityPost) (*models.CommunityPost, error)
	GetCommunityPostByID(id, viewerUserID string) (*models.CommunityPost, error)
	ListCommunityPosts(neighborhoodID, category string, viewerUserID string) ([]*models.CommunityPost, error)
	CreateCommunityComment(comment *models.CommunityComment) (*models.CommunityComment, error)
	GetCommunityComments(postID string) ([]*models.CommunityComment, error)
	ToggleCommunityLike(userID, postID string) (bool, error)

	// Business
	CreateBusiness(biz *models.Business) (*models.Business, error)
	GetBusinessByID(id string) (*models.Business, error)
	ListBusinesses(neighborhoodID string) ([]*models.Business, error)
	VerifyBusiness(bizID, adminID string, approve bool) (*models.Business, error)

	// Moderation & Safety
	CreateReport(report *models.Report) (*models.Report, error)
	ListReports(status string) ([]*models.Report, error)
	CreateModerationCase(reportID, targetType, targetID string) (*models.ModerationCase, error)
	ResolveModerationCase(caseID, adminID, decision, reason string) (*models.ModerationCase, error)
	ApplyEnforcementAction(action *models.EnforcementAction) (*models.EnforcementAction, error)
	SubmitAppeal(appeal *models.Appeal) (*models.Appeal, error)
	ListAppeals(status string) ([]*models.Appeal, error)

	// Notifications
	GetNotifications(userID string) ([]*models.Notification, error)
	MarkNotificationRead(id, userID string) error
	SendNotification(notif *models.Notification)

	// Admin
	GetDashboardStats() (*models.DashboardStats, error)
	RecordAuditLog(log *models.AuditLog)
	ListAuditLogs(limit int) ([]*models.AuditLog, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUserStatus(userID, status, adminID string) error
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		Users:                make(map[string]*models.User),
		Profiles:             make(map[string]*models.UserProfile),
		Sessions:             make(map[string]*models.UserSession),
		Blocks:               make(map[string]*models.UserBlock),
		Countries:            make(map[string]*models.Country),
		Cities:               make(map[string]*models.City),
		Neighborhoods:        make(map[string]*models.Neighborhood),
		UserNeighborhoods:    make(map[string]*models.UserNeighborhood),
		VerificationAttempts: make([]*models.LocationVerificationAttempt, 0),
		Categories:           make(map[string]*models.Category),
		Listings:             make(map[string]*models.Listing),
		Favorites:            make(map[string]*models.ListingFavorite),
		Conversations:        make(map[string]*models.Conversation),
		Messages:             make(map[string][]*models.ChatMessage),
		MessageClientIndex:   make(map[string]*models.ChatMessage),
		Offers:               make(map[string]*models.PriceOffer),
		Deals:                make(map[string]*models.Deal),
		DealReviews:          make(map[string]*models.DealReview),
		CommunityPosts:       make(map[string]*models.CommunityPost),
		CommunityComments:    make(map[string][]*models.CommunityComment),
		CommunityLikes:       make(map[string]bool),
		Businesses:           make(map[string]*models.Business),
		Reports:              make(map[string]*models.Report),
		ModerationCases:      make(map[string]*models.ModerationCase),
		Enforcements:         make(map[string]*models.EnforcementAction),
		Appeals:              make(map[string]*models.Appeal),
		Notifications:        make(map[string][]*models.Notification),
		NotificationPrefs:    make(map[string]*models.NotificationPreference),
		AuditLogs:            make([]*models.AuditLog, 0),
	}
	s.seedData()
	return s
}

func (s *MemoryStore) seedData() {
	// Countries
	s.Countries["KR"] = &models.Country{Code: "KR", Name: "South Korea", DefaultCurrency: "KRW", DefaultLocale: "ko-KR", Status: "ACTIVE"}
	s.Countries["US"] = &models.Country{Code: "US", Name: "United States", DefaultCurrency: "USD", DefaultLocale: "en-US", Status: "BETA"}
	s.Countries["JP"] = &models.Country{Code: "JP", Name: "Japan", DefaultCurrency: "JPY", DefaultLocale: "ja-JP", Status: "BETA"}

	// Cities
	incheonID := "a0000000-0000-0000-0000-000000000001"
	seoulID := "a0000000-0000-0000-0000-000000000002"
	s.Cities[incheonID] = &models.City{ID: incheonID, CountryCode: "KR", Name: "Incheon", Timezone: "Asia/Seoul", Status: "ACTIVE"}
	s.Cities[seoulID] = &models.City{ID: seoulID, CountryCode: "KR", Name: "Seoul", Timezone: "Asia/Seoul", Status: "ACTIVE"}

	// Neighborhoods
	cheongnaID := "b0000000-0000-0000-0000-000000000001"
	songdoID := "b0000000-0000-0000-0000-000000000002"
	yeoksamID := "b0000000-0000-0000-0000-000000000003"
	yeonnamID := "b0000000-0000-0000-0000-000000000004"

	// Cheongna boundary: roughly (37.520, 126.630) to (37.550, 126.670)
	s.Neighborhoods[cheongnaID] = &models.Neighborhood{
		ID:     cheongnaID,
		CityID: incheonID,
		Slug:   "cheongna-dong",
		Name:   "청라동 (Cheongna-dong)",
		BoundaryVertices: []models.Coordinate{
			{Latitude: 37.5200, Longitude: 126.6300},
			{Latitude: 37.5200, Longitude: 126.6700},
			{Latitude: 37.5500, Longitude: 126.6700},
			{Latitude: 37.5500, Longitude: 126.6300},
		},
		PublicCenter: models.Coordinate{Latitude: 37.5350, Longitude: 126.6500},
		Status:       "ACTIVE",
	}

	// Songdo boundary: (37.360, 126.620) to (37.410, 126.670)
	s.Neighborhoods[songdoID] = &models.Neighborhood{
		ID:     songdoID,
		CityID: incheonID,
		Slug:   "songdo-dong",
		Name:   "송도동 (Songdo-dong)",
		BoundaryVertices: []models.Coordinate{
			{Latitude: 37.3600, Longitude: 126.6200},
			{Latitude: 37.3600, Longitude: 126.6700},
			{Latitude: 37.4100, Longitude: 126.6700},
			{Latitude: 37.4100, Longitude: 126.6200},
		},
		PublicCenter: models.Coordinate{Latitude: 37.3850, Longitude: 126.6450},
		Status:       "ACTIVE",
	}

	// Yeoksam boundary: (37.495, 127.025) to (37.510, 127.045)
	s.Neighborhoods[yeoksamID] = &models.Neighborhood{
		ID:     yeoksamID,
		CityID: seoulID,
		Slug:   "yeoksam-dong",
		Name:   "역삼동 (Yeoksam-dong)",
		BoundaryVertices: []models.Coordinate{
			{Latitude: 37.4950, Longitude: 127.0250},
			{Latitude: 37.4950, Longitude: 127.0450},
			{Latitude: 37.5100, Longitude: 127.0450},
			{Latitude: 37.5100, Longitude: 127.0250},
		},
		PublicCenter: models.Coordinate{Latitude: 37.5020, Longitude: 127.0350},
		Status:       "ACTIVE",
	}

	// Yeonnam boundary: (37.560, 126.915) to (37.572, 126.930)
	s.Neighborhoods[yeonnamID] = &models.Neighborhood{
		ID:     yeonnamID,
		CityID: seoulID,
		Slug:   "yeonnam-dong",
		Name:   "연남동 (Yeonnam-dong)",
		BoundaryVertices: []models.Coordinate{
			{Latitude: 37.5600, Longitude: 126.9150},
			{Latitude: 37.5600, Longitude: 126.9300},
			{Latitude: 37.5720, Longitude: 126.9300},
			{Latitude: 37.5720, Longitude: 126.9150},
		},
		PublicCenter: models.Coordinate{Latitude: 37.5660, Longitude: 126.9220},
		Status:       "ACTIVE",
	}

	// Categories
	catDigitalID := "c0000000-0000-0000-0000-000000000001"
	catApplianceID := "c0000000-0000-0000-0000-000000000002"
	catFurnitureID := "c0000000-0000-0000-0000-000000000003"
	catBabyID := "c0000000-0000-0000-0000-000000000004"
	catFashionID := "c0000000-0000-0000-0000-000000000005"
	catSportsID := "c0000000-0000-0000-0000-000000000006"
	catBooksID := "c0000000-0000-0000-0000-000000000007"
	catPetID := "c0000000-0000-0000-0000-000000000008"
	catPlantsID := "c0000000-0000-0000-0000-000000000009"
	catServicesID := "c0000000-0000-0000-0000-000000000010"

	s.Categories[catDigitalID] = &models.Category{ID: catDigitalID, Slug: "digital-tech", Name: "디지털기기", IconName: "smartphone", SortOrder: 1, IsActive: true}
	s.Categories[catApplianceID] = &models.Category{ID: catApplianceID, Slug: "home-appliances", Name: "생활가전", IconName: "tv", SortOrder: 2, IsActive: true}
	s.Categories[catFurnitureID] = &models.Category{ID: catFurnitureID, Slug: "furniture-interior", Name: "가구/인테리어", IconName: "home", SortOrder: 3, IsActive: true}
	s.Categories[catBabyID] = &models.Category{ID: catBabyID, Slug: "baby-kids", Name: "유아동", IconName: "heart", SortOrder: 4, IsActive: true}
	s.Categories[catFashionID] = &models.Category{ID: catFashionID, Slug: "fashion-clothing", Name: "의류/잡화", IconName: "shopping-bag", SortOrder: 5, IsActive: true}
	s.Categories[catSportsID] = &models.Category{ID: catSportsID, Slug: "sports-leisure", Name: "스포츠/레저", IconName: "activity", SortOrder: 6, IsActive: true}
	s.Categories[catBooksID] = &models.Category{ID: catBooksID, Slug: "books-tickets", Name: "도서/음반/티켓", IconName: "book", SortOrder: 7, IsActive: true}
	s.Categories[catPetID] = &models.Category{ID: catPetID, Slug: "pet-supplies", Name: "반려동물용품", IconName: "smile", SortOrder: 8, IsActive: true}
	s.Categories[catPlantsID] = &models.Category{ID: catPlantsID, Slug: "plants", Name: "식물", IconName: "feather", SortOrder: 9, IsActive: true}
	s.Categories[catServicesID] = &models.Category{ID: catServicesID, Slug: "local-services", Name: "동네서비스", IconName: "tool", SortOrder: 10, IsActive: true}

	// Password hash for "Password123!"
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	pwdStr := string(hashedPwd)

	adminID := "u0000000-0000-0000-0000-000000000001"
	sellerID := "u0000000-0000-0000-0000-000000000002"
	buyerID := "u0000000-0000-0000-0000-000000000003"
	merchantID := "u0000000-0000-0000-0000-000000000004"

	now := time.Now()
	nowVerified := now
	expiry := now.Add(30 * 24 * time.Hour)

	// Admin User
	s.Users[adminID] = &models.User{
		ID: adminID, Email: "admin@bn-near.com", PasswordHash: pwdStr, Nickname: "NearAdmin",
		Locale: "ko-KR", CountryCode: "KR", Status: "ACTIVE", EmailVerifiedAt: &nowVerified, CreatedAt: now, UpdatedAt: now,
	}
	s.Profiles[adminID] = &models.UserProfile{
		UserID: adminID, Bio: "BN NEAR Official Super Admin", PublicDisplayName: "Platform Guardian",
		PreferredCurrency: "KRW", Timezone: "Asia/Seoul", MannerTemperature: 50.0, ProfileVisibility: "PUBLIC",
		CreatedAt: now, UpdatedAt: now,
	}

	// Seller User (David)
	s.Users[sellerID] = &models.User{
		ID: sellerID, Email: "seller@bn-near.com", PasswordHash: pwdStr, Nickname: "David_Cheongna",
		Locale: "ko-KR", CountryCode: "KR", Status: "ACTIVE", EmailVerifiedAt: &nowVerified, CreatedAt: now, UpdatedAt: now,
	}
	s.Profiles[sellerID] = &models.UserProfile{
		UserID: sellerID, Bio: "Tech enthusiast in Cheongna. Fast responses!", PublicDisplayName: "David K.",
		PreferredCurrency: "KRW", Timezone: "Asia/Seoul", MannerTemperature: 38.2, ProfileVisibility: "PUBLIC",
		CreatedAt: now, UpdatedAt: now,
	}
	s.UserNeighborhoods[sellerID+"_"+cheongnaID] = &models.UserNeighborhood{
		ID: uuid.NewString(), UserID: sellerID, NeighborhoodID: cheongnaID, NeighborhoodName: "청라동 (Cheongna-dong)",
		IsPrimary: true, VerificationStatus: "VERIFIED", VerifiedAt: &nowVerified, ExpiresAt: &expiry, CreatedAt: now,
	}

	// Buyer User (Minji)
	s.Users[buyerID] = &models.User{
		ID: buyerID, Email: "buyer@bn-near.com", PasswordHash: pwdStr, Nickname: "Minji_Cheongna",
		Locale: "ko-KR", CountryCode: "KR", Status: "ACTIVE", EmailVerifiedAt: &nowVerified, CreatedAt: now, UpdatedAt: now,
	}
	s.Profiles[buyerID] = &models.UserProfile{
		UserID: buyerID, Bio: "Resident of Cheongna Lake Park neighborhood", PublicDisplayName: "Minji P.",
		PreferredCurrency: "KRW", Timezone: "Asia/Seoul", MannerTemperature: 37.0, ProfileVisibility: "PUBLIC",
		CreatedAt: now, UpdatedAt: now,
	}
	s.UserNeighborhoods[buyerID+"_"+cheongnaID] = &models.UserNeighborhood{
		ID: uuid.NewString(), UserID: buyerID, NeighborhoodID: cheongnaID, NeighborhoodName: "청라동 (Cheongna-dong)",
		IsPrimary: true, VerificationStatus: "VERIFIED", VerifiedAt: &nowVerified, ExpiresAt: &expiry, CreatedAt: now,
	}

	// Merchant User (Master Kim)
	s.Users[merchantID] = &models.User{
		ID: merchantID, Email: "merchant@bn-near.com", PasswordHash: pwdStr, Nickname: "MasterKim_Bakery",
		Locale: "ko-KR", CountryCode: "KR", Status: "ACTIVE", EmailVerifiedAt: &nowVerified, CreatedAt: now, UpdatedAt: now,
	}
	s.Profiles[merchantID] = &models.UserProfile{
		UserID: merchantID, Bio: "Organic Sourdough & Artisan Coffee", PublicDisplayName: "Chef Kim",
		PreferredCurrency: "KRW", Timezone: "Asia/Seoul", MannerTemperature: 42.5, ProfileVisibility: "PUBLIC",
		CreatedAt: now, UpdatedAt: now,
	}
	s.UserNeighborhoods[merchantID+"_"+cheongnaID] = &models.UserNeighborhood{
		ID: uuid.NewString(), UserID: merchantID, NeighborhoodID: cheongnaID, NeighborhoodName: "청라동 (Cheongna-dong)",
		IsPrimary: true, VerificationStatus: "VERIFIED", VerifiedAt: &nowVerified, ExpiresAt: &expiry, CreatedAt: now,
	}

	// Demo Listings
	list1ID := "l0000000-0000-0000-0000-000000000001"
	s.Listings[list1ID] = &models.Listing{
		ID: list1ID, SellerUserID: sellerID, SellerNickname: "David_Cheongna", SellerMannerTemp: 38.2,
		NeighborhoodID: cheongnaID, NeighborhoodName: "청라동 (Cheongna-dong)",
		CategoryID: catDigitalID, CategoryName: "디지털기기",
		Title:       "애플 M2 맥북 에어 16GB 미드나이트",
		Description: "실사용 6개월 미만 배터리 효율 98% 상태 극상입니다. 박스 및 정품 충전기 풀세트 포함. 청라호수공원 정문 직거래 선호합니다.",
		Condition:   "LIKE_NEW", PriceMinor: 1150000, CurrencyCode: "KRW", IsFree: false, IsNegotiable: true,
		ListingStatus: "ACTIVE",
		PublicLocation: &models.Coordinate{Latitude: 37.5360, Longitude: 126.6520},
		Images:        []string{"https://images.unsplash.com/photo-1517336714731-489689fd1ca8?auto=format&fit=crop&w=800&q=80"},
		ViewCount:     142, FavoriteCount: 12, ChatCount: 3, Version: 1, PublishedAt: now.Add(-3 * time.Hour), CreatedAt: now.Add(-3 * time.Hour), UpdatedAt: now,
	}

	list2ID := "l0000000-0000-0000-0000-000000000002"
	s.Listings[list2ID] = &models.Listing{
		ID: list2ID, SellerUserID: sellerID, SellerNickname: "David_Cheongna", SellerMannerTemp: 38.2,
		NeighborhoodID: cheongnaID, NeighborhoodName: "청라동 (Cheongna-dong)",
		CategoryID: catFurnitureID, CategoryName: "가구/인테리어",
		Title:       "이케아 원목 식탁 4인용 세트 (의자 4개 포함)",
		Description: "이사로 인해 내놓습니다. 원목 오일케어 잘 되어 있으며 생활기스 약간 있습니다. 분해 가능합니다.",
		Condition:   "GOOD", PriceMinor: 80000, CurrencyCode: "KRW", IsFree: false, IsNegotiable: false,
		ListingStatus: "ACTIVE",
		PublicLocation: &models.Coordinate{Latitude: 37.5340, Longitude: 126.6480},
		Images:        []string{"https://images.unsplash.com/photo-1615066390971-03e4e1c36ddf?auto=format&fit=crop&w=800&q=80"},
		ViewCount:     88, FavoriteCount: 7, ChatCount: 1, Version: 1, PublishedAt: now.Add(-12 * time.Hour), CreatedAt: now.Add(-12 * time.Hour), UpdatedAt: now,
	}

	list3ID := "l0000000-0000-0000-0000-000000000003"
	s.Listings[list3ID] = &models.Listing{
		ID: list3ID, SellerUserID: buyerID, SellerNickname: "Minji_Cheongna", SellerMannerTemp: 37.0,
		NeighborhoodID: cheongnaID, NeighborhoodName: "청라동 (Cheongna-dong)",
		CategoryID: catBooksID, CategoryName: "도서/음반/티켓",
		Title:       "동화책 전집 30권 무료 나눔 (아이들이 잘 봤어요)",
		Description: "필요하신 이웃분께 나눔합니다. 청라 2동 주민센터 앞에서 전달 가능합니다.",
		Condition:   "FAIR", PriceMinor: 0, CurrencyCode: "KRW", IsFree: true, IsNegotiable: false,
		ListingStatus: "ACTIVE",
		PublicLocation: &models.Coordinate{Latitude: 37.5380, Longitude: 126.6550},
		Images:        []string{"https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c?auto=format&fit=crop&w=800&q=80"},
		ViewCount:     210, FavoriteCount: 19, ChatCount: 6, Version: 1, PublishedAt: now.Add(-24 * time.Hour), CreatedAt: now.Add(-24 * time.Hour), UpdatedAt: now,
	}

	// Demo Local Business
	bizID := "d0000000-0000-0000-0000-000000000001"
	s.Businesses[bizID] = &models.Business{
		ID: bizID, OwnerUserID: merchantID, Name: "청라 아티장 베이커리 (Cheongna Artisan Bakery)",
		CategoryID: catServicesID, CategoryName: "동네서비스", CountryCode: "KR", CityID: incheonID,
		NeighborhoodID: cheongnaID, NeighborhoodName: "청라동 (Cheongna-dong)",
		Description:   "매일 아침 유기농 통밀로 굽는 천연발효 사워도우와 스페셜티 커피 전문점.",
		PublicAddress: "인천 서구 청라커낼로 260번길 15 1층", PhoneNumber: "032-567-8901",
		ContactVisibility: true, VerificationStatus: "VERIFIED", Status: "ACTIVE",
		CreatedAt: now.Add(-48 * time.Hour), UpdatedAt: now,
	}

	// Demo Community Post
	postID := "p0000000-0000-0000-0000-000000000001"
	s.CommunityPosts[postID] = &models.CommunityPost{
		ID: postID, AuthorUserID: buyerID, AuthorNickname: "Minji_Cheongna",
		NeighborhoodID: cheongnaID, NeighborhoodName: "청라동 (Cheongna-dong)",
		Category: "동네질문", Title: "청라 호수공원 근처 주말에 아이 데리고 가기 좋은 소아과 추천 부탁드려요!",
		Body:      "최근에 청라로 이사 왔는데 주말 진료 잘 봐주시는 소아청소년과가 어디가 좋은지 이웃분들 조언 구합니다 :)",
		ViewCount: 52, LikeCount: 5, CommentCount: 2, Status: "ACTIVE",
		CreatedAt: now.Add(-6 * time.Hour), UpdatedAt: now,
	}
	s.CommunityComments[postID] = []*models.CommunityComment{
		{
			ID: "m0000000-0000-0000-0000-000000000001", PostID: postID, AuthorUserID: sellerID, AuthorNickname: "David_Cheongna",
			Body: "커낼웨이 수변상가 쪽 연세소아과 선생님 친절하시고 일요일 오전에도 진료해요!", Status: "ACTIVE", CreatedAt: now.Add(-4 * time.Hour), UpdatedAt: now,
		},
		{
			ID: "m0000000-0000-0000-0000-000000000002", PostID: postID, AuthorUserID: merchantID, AuthorNickname: "MasterKim_Bakery",
			Body: "저희 아이도 연세소아과 다닙니다. 주말엔 대기 조금 있으니 일찍 가시는 걸 추천드려요~", Status: "ACTIVE", CreatedAt: now.Add(-2 * time.Hour), UpdatedAt: now,
		},
	}
}

// ----------------- USERS -----------------

func (s *MemoryStore) CreateUser(u *models.User, rawPassword string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.Users {
		if strings.EqualFold(existing.Email, u.Email) {
			return nil, ErrAlreadyExists
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	u.PasswordHash = string(hashed)
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	if u.Status == "" {
		u.Status = "ACTIVE"
	}
	if u.Locale == "" {
		u.Locale = "ko-KR"
	}
	if u.CountryCode == "" {
		u.CountryCode = "KR"
	}

	s.Users[u.ID] = u

	s.Profiles[u.ID] = &models.UserProfile{
		UserID:            u.ID,
		PublicDisplayName: u.Nickname,
		PreferredCurrency: "KRW",
		Timezone:          "Asia/Seoul",
		MannerTemperature: 36.5,
		ProfileVisibility: "PUBLIC",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	copyUser := *u
	copyUser.Profile = s.Profiles[u.ID]
	return &copyUser, nil
}

func (s *MemoryStore) AuthenticateUser(email, password string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matched *models.User
	for _, u := range s.Users {
		if strings.EqualFold(u.Email, email) && u.DeletedAt == nil {
			matched = u
			break
		}
	}

	if matched == nil {
		return nil, ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(matched.PasswordHash), []byte(password)); err != nil {
		return nil, ErrUnauthorized
	}

	copyUser := *matched
	copyUser.Profile = s.Profiles[matched.ID]
	return &copyUser, nil
}

func (s *MemoryStore) GetUserByID(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.Users[id]
	if !ok || u.DeletedAt != nil {
		return nil, ErrNotFound
	}
	copyUser := *u
	copyUser.Profile = s.Profiles[id]
	return &copyUser, nil
}

func (s *MemoryStore) GetUserByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.Users {
		if strings.EqualFold(u.Email, email) && u.DeletedAt == nil {
			copyUser := *u
			copyUser.Profile = s.Profiles[u.ID]
			return &copyUser, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) UpdateUserProfile(userID, bio, displayName string) (*models.UserProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	prof, ok := s.Profiles[userID]
	if !ok {
		return nil, ErrNotFound
	}
	if bio != "" {
		prof.Bio = bio
	}
	if displayName != "" {
		prof.PublicDisplayName = displayName
		if u, exists := s.Users[userID]; exists {
			u.Nickname = displayName
		}
	}
	prof.UpdatedAt = time.Now()
	copyProf := *prof
	return &copyProf, nil
}

func (s *MemoryStore) SoftDeleteUser(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.Users[userID]
	if !ok {
		return ErrNotFound
	}
	now := time.Now()
	u.DeletedAt = &now
	u.Status = "DELETED"
	return nil
}

func (s *MemoryStore) CreateSession(session *models.UserSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Sessions[session.ID] = session
	return nil
}

func (s *MemoryStore) GetSession(sessionID string) (*models.UserSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.Sessions[sessionID]
	if !ok || sess.RevokedAt != nil || sess.ExpiresAt.Before(time.Now()) {
		return nil, ErrUnauthorized
	}
	copySess := *sess
	return &copySess, nil
}

func (s *MemoryStore) RevokeSession(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.Sessions[sessionID]
	if !ok {
		return ErrNotFound
	}
	now := time.Now()
	sess.RevokedAt = &now
	return nil
}

func (s *MemoryStore) BlockUser(blockerID, blockedID, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if blockerID == blockedID {
		return ErrInvalidInput
	}
	key := blockerID + "_" + blockedID
	s.Blocks[key] = &models.UserBlock{
		ID:             uuid.NewString(),
		BlockerUserID:  blockerID,
		BlockedUserID:  blockedID,
		ReasonOptional: reason,
		CreatedAt:      time.Now(),
	}
	return nil
}

func (s *MemoryStore) UnblockUser(blockerID, blockedID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := blockerID + "_" + blockedID
	delete(s.Blocks, key)
	return nil
}

func (s *MemoryStore) IsBlocked(userA, userB string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, aBlockedB := s.Blocks[userA+"_"+userB]
	_, bBlockedA := s.Blocks[userB+"_"+userA]
	return aBlockedB || bBlockedA
}

// ----------------- GEOGRAPHY & VERIFICATION -----------------

func (s *MemoryStore) ListCountries() ([]*models.Country, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Country, 0, len(s.Countries))
	for _, c := range s.Countries {
		res = append(res, c)
	}
	return res, nil
}

func (s *MemoryStore) ListCities(countryCode string) ([]*models.City, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.City, 0)
	for _, c := range s.Cities {
		if countryCode == "" || strings.EqualFold(c.CountryCode, countryCode) {
			res = append(res, c)
		}
	}
	return res, nil
}

func (s *MemoryStore) GetNeighborhoodByID(id string) (*models.Neighborhood, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n, ok := s.Neighborhoods[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

func (s *MemoryStore) SearchNeighborhoods(query string, cityID string, lat, lng *float64) ([]*models.Neighborhood, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Neighborhood, 0)
	q := strings.ToLower(query)

	for _, n := range s.Neighborhoods {
		if cityID != "" && n.CityID != cityID {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(n.Name), q) && !strings.Contains(strings.ToLower(n.Slug), q) {
			continue
		}
		res = append(res, n)
	}

	// Sort by proximity if lat/lng provided
	if lat != nil && lng != nil {
		sort.Slice(res, func(i, j int) bool {
			d1 := HaversineMeters(*lat, *lng, res[i].PublicCenter.Latitude, res[i].PublicCenter.Longitude)
			d2 := HaversineMeters(*lat, *lng, res[j].PublicCenter.Latitude, res[j].PublicCenter.Longitude)
			return d1 < d2
		})
	}

	return res, nil
}

// Point-in-polygon ray casting algorithm
func isPointInPolygon(lat, lng float64, vertices []models.Coordinate) bool {
	if len(vertices) < 3 {
		return false
	}
	inside := false
	j := len(vertices) - 1
	for i := 0; i < len(vertices); i++ {
		xi, yi := vertices[i].Longitude, vertices[i].Latitude
		xj, yj := vertices[j].Longitude, vertices[j].Latitude

		intersect := ((yi > lat) != (yj > lat)) &&
			(lng < (xj-xi)*(lat-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
		j = i
	}
	return inside
}

func (s *MemoryStore) VerifyNeighborhood(userID, neighborhoodID string, lat, lng, accuracy float64) (*models.UserNeighborhood, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, ok := s.Neighborhoods[neighborhoodID]
	if !ok {
		return nil, ErrNotFound
	}

	attempt := &models.LocationVerificationAttempt{
		ID:             uuid.NewString(),
		UserID:         userID,
		NeighborhoodID: neighborhoodID,
		AccuracyMeters: accuracy,
		AttemptedAt:    time.Now(),
	}

	// Accuracy check threshold: 150m
	if accuracy > 150.0 {
		attempt.VerificationResult = "LOW_ACCURACY"
		s.VerificationAttempts = append(s.VerificationAttempts, attempt)
		return nil, fmt.Errorf("%w: reported GPS accuracy (%.1fm) is lower than 150m requirement", ErrInvalidInput, accuracy)
	}

	// PostGIS / Polygon boundary check
	inside := isPointInPolygon(lat, lng, n.BoundaryVertices)
	if !inside {
		// Calculate distance to public centroid as fallback tolerance (within 3km of center)
		distCenter := HaversineMeters(lat, lng, n.PublicCenter.Latitude, n.PublicCenter.Longitude)
		if distCenter > 3500.0 {
			attempt.VerificationResult = "OUT_OF_BOUNDS"
			s.VerificationAttempts = append(s.VerificationAttempts, attempt)
			return nil, fmt.Errorf("%w: location is outside %s boundary (distance to center: %.0fm)", ErrInvalidInput, n.Name, distCenter)
		}
	}

	attempt.VerificationResult = "VERIFIED"
	s.VerificationAttempts = append(s.VerificationAttempts, attempt)

	now := time.Now()
	expires := now.Add(30 * 24 * time.Hour)

	key := userID + "_" + neighborhoodID
	un, exists := s.UserNeighborhoods[key]
	if !exists {
		un = &models.UserNeighborhood{
			ID:                 uuid.NewString(),
			UserID:             userID,
			NeighborhoodID:     neighborhoodID,
			NeighborhoodName:   n.Name,
			IsPrimary:          true,
			VerificationStatus: "VERIFIED",
			VerifiedAt:         &now,
			ExpiresAt:          &expires,
			CreatedAt:          now,
		}
		s.UserNeighborhoods[key] = un
	} else {
		un.VerificationStatus = "VERIFIED"
		un.VerifiedAt = &now
		un.ExpiresAt = &expires
	}

	copyUN := *un
	return &copyUN, nil
}

func (s *MemoryStore) GetUserNeighborhoods(userID string) ([]*models.UserNeighborhood, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.UserNeighborhood, 0)
	for _, un := range s.UserNeighborhoods {
		if un.UserID == userID {
			res = append(res, un)
		}
	}
	return res, nil
}

func (s *MemoryStore) IsUserVerifiedInNeighborhood(userID, neighborhoodID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := userID + "_" + neighborhoodID
	un, ok := s.UserNeighborhoods[key]
	if !ok {
		return false
	}
	return un.VerificationStatus == "VERIFIED" && un.ExpiresAt != nil && un.ExpiresAt.After(time.Now())
}

// ----------------- MARKETPLACE LISTINGS -----------------

func (s *MemoryStore) CreateListing(l *models.Listing) (*models.Listing, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if l.PriceMinor < 0 {
		return nil, fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}
	if l.IsFree && l.PriceMinor > 0 {
		return nil, fmt.Errorf("%w: free item price must be 0", ErrInvalidInput)
	}

	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	now := time.Now()
	l.CreatedAt = now
	l.UpdatedAt = now
	l.PublishedAt = now
	l.ListingStatus = "ACTIVE"
	l.Version = 1

	if n, ok := s.Neighborhoods[l.NeighborhoodID]; ok {
		l.NeighborhoodName = n.Name
		if l.PublicLocation == nil {
			l.PublicLocation = &models.Coordinate{
				Latitude:  n.PublicCenter.Latitude,
				Longitude: n.PublicCenter.Longitude,
			}
		}
	}
	if cat, ok := s.Categories[l.CategoryID]; ok {
		l.CategoryName = cat.Name
	}
	if seller, ok := s.Users[l.SellerUserID]; ok {
		l.SellerNickname = seller.Nickname
		if prof, pOk := s.Profiles[l.SellerUserID]; pOk {
			l.SellerMannerTemp = prof.MannerTemperature
		}
	}

	s.Listings[l.ID] = l
	copyListing := *l
	return &copyListing, nil
}

func (s *MemoryStore) GetListingByID(id, viewerUserID string) (*models.Listing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	l, ok := s.Listings[id]
	if !ok {
		return nil, ErrNotFound
	}
	copyL := *l
	if viewerUserID != "" {
		_, copyL.IsFavorited = s.Favorites[viewerUserID+"_"+id]
	}
	return &copyL, nil
}

func (s *MemoryStore) UpdateListing(l *models.Listing) (*models.Listing, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.Listings[l.ID]
	if !ok {
		return nil, ErrNotFound
	}
	if existing.SellerUserID != l.SellerUserID {
		return nil, ErrUnauthorized
	}

	existing.Title = l.Title
	existing.Description = l.Description
	existing.Condition = l.Condition
	existing.PriceMinor = l.PriceMinor
	existing.IsFree = l.IsFree
	existing.IsNegotiable = l.IsNegotiable
	if len(l.Images) > 0 {
		existing.Images = l.Images
	}
	existing.UpdatedAt = time.Now()
	existing.Version++

	copyL := *existing
	return &copyL, nil
}

func (s *MemoryStore) DeleteListing(id, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.Listings[id]
	if !ok {
		return ErrNotFound
	}
	if existing.SellerUserID != userID {
		return ErrUnauthorized
	}
	existing.ListingStatus = "REMOVED"
	existing.UpdatedAt = time.Now()
	return nil
}

func (s *MemoryStore) ListListings(neighborhoodID, categoryID, condition string, minPrice, maxPrice *int64, isFree *bool, lat, lng *float64, radiusMeters float64, sortOrder string, viewerUserID string) ([]*models.Listing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Listing, 0)

	for _, l := range s.Listings {
		if l.ListingStatus != "ACTIVE" && l.ListingStatus != "RESERVED" {
			continue
		}
		if neighborhoodID != "" && l.NeighborhoodID != neighborhoodID {
			continue
		}
		if categoryID != "" && l.CategoryID != categoryID {
			continue
		}
		if condition != "" && l.Condition != condition {
			continue
		}
		if minPrice != nil && l.PriceMinor < *minPrice {
			continue
		}
		if maxPrice != nil && l.PriceMinor > *maxPrice {
			continue
		}
		if isFree != nil && l.IsFree != *isFree {
			continue
		}

		copyL := *l
		if viewerUserID != "" {
			_, copyL.IsFavorited = s.Favorites[viewerUserID+"_"+l.ID]
		}

		if lat != nil && lng != nil && l.PublicLocation != nil {
			dist := HaversineMeters(*lat, *lng, l.PublicLocation.Latitude, l.PublicLocation.Longitude)
			if radiusMeters > 0 && dist > radiusMeters {
				continue
			}
			copyL.DistanceMeters = dist
		}

		res = append(res, &copyL)
	}

	// Sorting
	switch sortOrder {
	case "distance":
		sort.Slice(res, func(i, j int) bool {
			return res[i].DistanceMeters < res[j].DistanceMeters
		})
	case "price_asc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].PriceMinor < res[j].PriceMinor
		})
	case "price_desc":
		sort.Slice(res, func(i, j int) bool {
			return res[i].PriceMinor > res[j].PriceMinor
		})
	default: // "latest"
		sort.Slice(res, func(i, j int) bool {
			return res[i].PublishedAt.After(res[j].PublishedAt)
		})
	}

	return res, nil
}

func (s *MemoryStore) ToggleFavorite(userID, listingID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, ok := s.Listings[listingID]
	if !ok {
		return false, ErrNotFound
	}

	key := userID + "_" + listingID
	if _, exists := s.Favorites[key]; exists {
		delete(s.Favorites, key)
		if l.FavoriteCount > 0 {
			l.FavoriteCount--
		}
		return false, nil
	} else {
		s.Favorites[key] = &models.ListingFavorite{
			UserID:    userID,
			ListingID: listingID,
			CreatedAt: time.Now(),
		}
		l.FavoriteCount++
		return true, nil
	}
}

func (s *MemoryStore) GetUserFavorites(userID string) ([]*models.Listing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Listing, 0)
	for key := range s.Favorites {
		parts := strings.Split(key, "_")
		if len(parts) == 2 && parts[0] == userID {
			if l, ok := s.Listings[parts[1]]; ok {
				copyL := *l
				copyL.IsFavorited = true
				res = append(res, &copyL)
			}
		}
	}
	return res, nil
}

func (s *MemoryStore) GetUserListings(userID string) ([]*models.Listing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Listing, 0)
	for _, l := range s.Listings {
		if l.SellerUserID == userID {
			copyL := *l
			res = append(res, &copyL)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.After(res[j].CreatedAt)
	})
	return res, nil
}

// ----------------- REAL-TIME CHAT -----------------

func (s *MemoryStore) GetOrCreateConversation(listingID, buyerUserID string) (*models.Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, ok := s.Listings[listingID]
	if !ok {
		return nil, ErrNotFound
	}
	if l.SellerUserID == buyerUserID {
		return nil, fmt.Errorf("%w: seller cannot chat with themselves", ErrInvalidInput)
	}

	// Check existing conversation
	for _, conv := range s.Conversations {
		if conv.ListingID == listingID && conv.BuyerUserID == buyerUserID {
			copyConv := *conv
			return &copyConv, nil
		}
	}

	// Create new
	now := time.Now()
	conv := &models.Conversation{
		ID:            uuid.NewString(),
		ListingID:     listingID,
		ListingTitle:  l.Title,
		ListingPrice:  l.PriceMinor,
		ListingStatus: l.ListingStatus,
		BuyerUserID:   buyerUserID,
		SellerUserID:  l.SellerUserID,
		Status:        "ACTIVE",
		LastMessage:   "대화가 시작되었습니다.",
		LastMessageAt: now,
		UnreadCount:   0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if len(l.Images) > 0 {
		conv.ListingImage = l.Images[0]
	}
	if buyer, bOk := s.Users[buyerUserID]; bOk {
		conv.BuyerNickname = buyer.Nickname
	}
	if seller, sOk := s.Users[l.SellerUserID]; sOk {
		conv.SellerNickname = seller.Nickname
	}

	s.Conversations[conv.ID] = conv
	s.Messages[conv.ID] = make([]*models.ChatMessage, 0)

	l.ChatCount++

	copyConv := *conv
	return &copyConv, nil
}

func (s *MemoryStore) GetUserConversations(userID string) ([]*models.Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Conversation, 0)
	for _, conv := range s.Conversations {
		if conv.BuyerUserID == userID || conv.SellerUserID == userID {
			copyConv := *conv
			if l, ok := s.Listings[conv.ListingID]; ok {
				copyConv.ListingStatus = l.ListingStatus
			}
			res = append(res, &copyConv)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].LastMessageAt.After(res[j].LastMessageAt)
	})
	return res, nil
}

func (s *MemoryStore) GetConversationByID(id, userID string) (*models.Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, ok := s.Conversations[id]
	if !ok {
		return nil, ErrNotFound
	}
	if conv.BuyerUserID != userID && conv.SellerUserID != userID {
		return nil, ErrUnauthorized
	}
	copyConv := *conv
	if l, lOk := s.Listings[conv.ListingID]; lOk {
		copyConv.ListingStatus = l.ListingStatus
	}
	return &copyConv, nil
}

func (s *MemoryStore) SaveMessage(msg *models.ChatMessage) (*models.ChatMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.Conversations[msg.ConversationID]
	if !ok {
		return nil, ErrNotFound
	}

	// Block check
	otherUserID := conv.BuyerUserID
	if msg.SenderUserID == conv.BuyerUserID {
		otherUserID = conv.SellerUserID
	}
	if _, blocked := s.Blocks[otherUserID+"_"+msg.SenderUserID]; blocked {
		return nil, ErrBlocked
	}
	if _, blocked := s.Blocks[msg.SenderUserID+"_"+otherUserID]; blocked {
		return nil, ErrBlocked
	}

	// Idempotency check: ClientMessageID deduplication
	dedupKey := msg.ConversationID + "_" + msg.SenderUserID + "_" + msg.ClientMessageID
	if existing, exists := s.MessageClientIndex[dedupKey]; exists {
		copyMsg := *existing
		return &copyMsg, nil
	}

	msgs := s.Messages[msg.ConversationID]
	seq := int64(len(msgs) + 1)

	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}
	msg.SequenceNo = seq
	msg.Status = "STORED"
	now := time.Now()
	msg.CreatedAt = now

	s.Messages[msg.ConversationID] = append(s.Messages[msg.ConversationID], msg)
	s.MessageClientIndex[dedupKey] = msg

	conv.LastMessage = msg.Body
	conv.LastMessageAt = now
	conv.UnreadCount++
	conv.UpdatedAt = now

	copyMsg := *msg
	return &copyMsg, nil
}

func (s *MemoryStore) GetMessages(conversationID string, sinceSequence int64) ([]*models.ChatMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msgs, ok := s.Messages[conversationID]
	if !ok {
		return make([]*models.ChatMessage, 0), nil
	}

	res := make([]*models.ChatMessage, 0)
	for _, m := range msgs {
		if m.SequenceNo > sinceSequence {
			res = append(res, m)
		}
	}
	return res, nil
}

func (s *MemoryStore) MarkMessagesRead(conversationID, userID string, sequenceNo int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.Conversations[conversationID]
	if !ok {
		return ErrNotFound
	}
	conv.LastReadSequence = sequenceNo
	conv.UnreadCount = 0
	return nil
}

// ----------------- DEALS & OFFERS -----------------

func (s *MemoryStore) CreateOffer(offer *models.PriceOffer) (*models.PriceOffer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, ok := s.Listings[offer.ListingID]
	if !ok {
		return nil, ErrNotFound
	}
	if !l.IsNegotiable {
		return nil, fmt.Errorf("%w: listing does not accept offers", ErrConflict)
	}

	if offer.ID == "" {
		offer.ID = uuid.NewString()
	}
	now := time.Now()
	offer.SellerUserID = l.SellerUserID
	offer.Status = "PENDING"
	offer.CurrencyCode = l.CurrencyCode
	offer.ExpiresAt = now.Add(48 * time.Hour)
	offer.CreatedAt = now
	offer.UpdatedAt = now

	s.Offers[offer.ID] = offer
	copyOffer := *offer
	return &copyOffer, nil
}

func (s *MemoryStore) GetOfferByID(id string) (*models.PriceOffer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	offer, ok := s.Offers[id]
	if !ok {
		return nil, ErrNotFound
	}
	copyOffer := *offer
	return &copyOffer, nil
}

func (s *MemoryStore) UpdateOfferStatus(id, userID, newStatus string) (*models.PriceOffer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	offer, ok := s.Offers[id]
	if !ok {
		return nil, ErrNotFound
	}
	if offer.SellerUserID != userID && offer.BuyerUserID != userID {
		return nil, ErrUnauthorized
	}

	offer.Status = newStatus
	offer.UpdatedAt = time.Now()
	copyOffer := *offer
	return &copyOffer, nil
}

func (s *MemoryStore) GetDealByID(id string) (*models.Deal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, ok := s.Deals[id]
	if !ok {
		return nil, ErrNotFound
	}
	copyD := *d
	return &copyD, nil
}

func (s *MemoryStore) GetDealByListingID(listingID string) (*models.Deal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, d := range s.Deals {
		if d.ListingID == listingID && d.Status != "CANCELLED" {
			copyD := *d
			return &copyD, nil
		}
	}
	return nil, ErrNotFound
}

// Reservation Collision Lock: Guarantees ONLY ONE active reservation per listing
func (s *MemoryStore) ReserveDeal(listingID, sellerUserID, buyerUserID string) (*models.Deal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, ok := s.Listings[listingID]
	if !ok {
		return nil, ErrNotFound
	}
	if l.SellerUserID != sellerUserID {
		return nil, ErrUnauthorized
	}

	// Collision check: Cannot reserve if already RESERVED or SOLD
	if l.ListingStatus == "RESERVED" {
		return nil, fmt.Errorf("%w: listing is already reserved by another buyer", ErrConflict)
	}
	if l.ListingStatus == "SOLD" {
		return nil, fmt.Errorf("%w: listing is already sold", ErrConflict)
	}

	now := time.Now()
	l.ListingStatus = "RESERVED"
	l.ReservedAt = &now
	l.Version++
	l.UpdatedAt = now

	// Find conversation
	convID := ""
	for _, c := range s.Conversations {
		if c.ListingID == listingID && c.BuyerUserID == buyerUserID {
			convID = c.ID
			break
		}
	}

	deal := &models.Deal{
		ID:             uuid.NewString(),
		ListingID:      listingID,
		BuyerUserID:    buyerUserID,
		SellerUserID:   sellerUserID,
		ConversationID: convID,
		Status:         "RESERVED",
		CurrencyCode:   l.CurrencyCode,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	s.Deals[deal.ID] = deal

	copyDeal := *deal
	return &copyDeal, nil
}

func (s *MemoryStore) ScheduleMeeting(dealID, userID, placeName string, meetingTime time.Time) (*models.Deal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.Deals[dealID]
	if !ok {
		return nil, ErrNotFound
	}
	if d.SellerUserID != userID && d.BuyerUserID != userID {
		return nil, ErrUnauthorized
	}

	d.MeetingPlaceName = placeName
	d.MeetingTime = &meetingTime
	d.Status = "MEETING_PLANNED"
	d.UpdatedAt = time.Now()

	copyD := *d
	return &copyD, nil
}

// Two-Phase Deal Phase 1: Seller marks sold
func (s *MemoryStore) SellerCompleteDeal(dealID, sellerUserID string) (*models.Deal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.Deals[dealID]
	if !ok {
		return nil, ErrNotFound
	}
	if d.SellerUserID != sellerUserID {
		return nil, ErrUnauthorized
	}

	now := time.Now()
	d.Status = "COMPLETED_CLAIMED"
	d.SellerCompletedAt = &now
	d.UpdatedAt = now

	if l, ok := s.Listings[d.ListingID]; ok {
		l.ListingStatus = "SOLD"
		l.SoldAt = &now
		l.UpdatedAt = now
	}

	copyD := *d
	return &copyD, nil
}

// Two-Phase Deal Phase 2: Buyer confirms receipt, unlocking mutual review
func (s *MemoryStore) BuyerConfirmDeal(dealID, buyerUserID string) (*models.Deal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.Deals[dealID]
	if !ok {
		return nil, ErrNotFound
	}
	if d.BuyerUserID != buyerUserID {
		return nil, ErrUnauthorized
	}

	now := time.Now()
	d.Status = "COMPLETED_CONFIRMED"
	d.BuyerConfirmedAt = &now
	d.CanReview = true
	d.UpdatedAt = now

	copyD := *d
	return &copyD, nil
}

func (s *MemoryStore) CancelDeal(dealID, userID string) (*models.Deal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.Deals[dealID]
	if !ok {
		return nil, ErrNotFound
	}
	if d.SellerUserID != userID && d.BuyerUserID != userID {
		return nil, ErrUnauthorized
	}

	now := time.Now()
	d.Status = "CANCELLED"
	d.CancelledAt = &now
	d.UpdatedAt = now

	if l, ok := s.Listings[d.ListingID]; ok {
		l.ListingStatus = "ACTIVE"
		l.ReservedAt = nil
		l.UpdatedAt = now
	}

	copyD := *d
	return &copyD, nil
}

func (s *MemoryStore) GetUserDeals(userID string) ([]*models.Deal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Deal, 0)
	for _, d := range s.Deals {
		if d.SellerUserID == userID || d.BuyerUserID == userID {
			copyD := *d
			res = append(res, &copyD)
		}
	}
	return res, nil
}

// ----------------- REVIEWS & REPUTATION -----------------

func (s *MemoryStore) CreateDealReview(review *models.DealReview) (*models.DealReview, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.Deals[review.DealID]
	if !ok {
		return nil, ErrNotFound
	}

	// Must be completed
	if d.Status != "COMPLETED_CONFIRMED" && d.Status != "COMPLETED_CLAIMED" {
		return nil, fmt.Errorf("%w: reviews permitted only upon completed transaction", ErrConflict)
	}

	// Reviewer eligibility
	if review.ReviewerUserID != d.BuyerUserID && review.ReviewerUserID != d.SellerUserID {
		return nil, ErrUnauthorized
	}

	targetUserID := d.SellerUserID
	if review.ReviewerUserID == d.SellerUserID {
		targetUserID = d.BuyerUserID
	}
	review.ReviewedUserID = targetUserID

	// Deduplication: One review per party per deal
	key := review.DealID + "_" + review.ReviewerUserID
	if _, exists := s.DealReviews[key]; exists {
		return nil, fmt.Errorf("%w: review already submitted for this deal", ErrAlreadyExists)
	}

	if review.Rating < 1 || review.Rating > 5 {
		return nil, fmt.Errorf("%w: rating must be between 1 and 5", ErrInvalidInput)
	}

	if review.ID == "" {
		review.ID = uuid.NewString()
	}
	review.Status = "PUBLISHED"
	review.CreatedAt = time.Now()

	s.DealReviews[key] = review

	// Recalculate Manner Temperature: base 36.5
	if prof, ok := s.Profiles[targetUserID]; ok {
		delta := 0.0
		switch review.Rating {
		case 5:
			delta = 0.5 + float64(len(review.Tags))*0.1
		case 4:
			delta = 0.2
		case 3:
			delta = 0.0
		case 2:
			delta = -0.5
		case 1:
			delta = -1.5
		}
		newTemp := prof.MannerTemperature + delta
		if newTemp < 0.0 {
			newTemp = 0.0
		}
		if newTemp > 99.9 {
			newTemp = 99.9
		}
		prof.MannerTemperature = math.Round(newTemp*10) / 10
		prof.UpdatedAt = time.Now()
	}

	copyR := *review
	return &copyR, nil
}

func (s *MemoryStore) GetUserReviews(userID string) ([]*models.DealReview, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.DealReview, 0)
	for _, r := range s.DealReviews {
		if r.ReviewedUserID == userID && r.Status == "PUBLISHED" {
			copyR := *r
			res = append(res, &copyR)
		}
	}
	return res, nil
}

func (s *MemoryStore) GetUserReputation(userID string) (*models.UserProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prof, ok := s.Profiles[userID]
	if !ok {
		return nil, ErrNotFound
	}
	copyProf := *prof
	return &copyProf, nil
}

// ----------------- COMMUNITY -----------------

func (s *MemoryStore) CreateCommunityPost(post *models.CommunityPost) (*models.CommunityPost, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if post.ID == "" {
		post.ID = uuid.NewString()
	}
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now
	post.Status = "ACTIVE"

	if n, ok := s.Neighborhoods[post.NeighborhoodID]; ok {
		post.NeighborhoodName = n.Name
	}
	if u, ok := s.Users[post.AuthorUserID]; ok {
		post.AuthorNickname = u.Nickname
	}

	s.CommunityPosts[post.ID] = post
	s.CommunityComments[post.ID] = make([]*models.CommunityComment, 0)

	copyPost := *post
	return &copyPost, nil
}

func (s *MemoryStore) GetCommunityPostByID(id, viewerUserID string) (*models.CommunityPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.CommunityPosts[id]
	if !ok {
		return nil, ErrNotFound
	}
	copyP := *p
	if viewerUserID != "" {
		copyP.IsLiked = s.CommunityLikes[viewerUserID+"_"+id]
	}
	return &copyP, nil
}

func (s *MemoryStore) ListCommunityPosts(neighborhoodID, category string, viewerUserID string) ([]*models.CommunityPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.CommunityPost, 0)
	for _, p := range s.CommunityPosts {
		if p.Status != "ACTIVE" {
			continue
		}
		if neighborhoodID != "" && p.NeighborhoodID != neighborhoodID {
			continue
		}
		if category != "" && p.Category != category {
			continue
		}
		copyP := *p
		if viewerUserID != "" {
			copyP.IsLiked = s.CommunityLikes[viewerUserID+"_"+p.ID]
		}
		res = append(res, &copyP)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.After(res[j].CreatedAt)
	})
	return res, nil
}

func (s *MemoryStore) CreateCommunityComment(comment *models.CommunityComment) (*models.CommunityComment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.CommunityPosts[comment.PostID]
	if !ok {
		return nil, ErrNotFound
	}

	if comment.ID == "" {
		comment.ID = uuid.NewString()
	}
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now
	comment.Status = "ACTIVE"

	if u, ok := s.Users[comment.AuthorUserID]; ok {
		comment.AuthorNickname = u.Nickname
	}

	s.CommunityComments[comment.PostID] = append(s.CommunityComments[comment.PostID], comment)
	post.CommentCount++

	copyC := *comment
	return &copyC, nil
}

func (s *MemoryStore) GetCommunityComments(postID string) ([]*models.CommunityComment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comments, ok := s.CommunityComments[postID]
	if !ok {
		return make([]*models.CommunityComment, 0), nil
	}
	res := make([]*models.CommunityComment, len(comments))
	copy(res, comments)
	return res, nil
}

func (s *MemoryStore) ToggleCommunityLike(userID, postID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.CommunityPosts[postID]
	if !ok {
		return false, ErrNotFound
	}

	key := userID + "_" + postID
	if s.CommunityLikes[key] {
		delete(s.CommunityLikes, key)
		if post.LikeCount > 0 {
			post.LikeCount--
		}
		return false, nil
	} else {
		s.CommunityLikes[key] = true
		post.LikeCount++
		return true, nil
	}
}

// ----------------- BUSINESS -----------------

func (s *MemoryStore) CreateBusiness(biz *models.Business) (*models.Business, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if biz.ID == "" {
		biz.ID = uuid.NewString()
	}
	now := time.Now()
	biz.CreatedAt = now
	biz.UpdatedAt = now
	biz.Status = "ACTIVE"
	biz.VerificationStatus = "PENDING"

	if n, ok := s.Neighborhoods[biz.NeighborhoodID]; ok {
		biz.NeighborhoodName = n.Name
	}
	if cat, ok := s.Categories[biz.CategoryID]; ok {
		biz.CategoryName = cat.Name
	}

	s.Businesses[biz.ID] = biz
	copyB := *biz
	return &copyB, nil
}

func (s *MemoryStore) GetBusinessByID(id string) (*models.Business, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.Businesses[id]
	if !ok {
		return nil, ErrNotFound
	}
	copyB := *b
	return &copyB, nil
}

func (s *MemoryStore) ListBusinesses(neighborhoodID string) ([]*models.Business, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Business, 0)
	for _, b := range s.Businesses {
		if b.Status != "ACTIVE" {
			continue
		}
		if neighborhoodID != "" && b.NeighborhoodID != neighborhoodID {
			continue
		}
		copyB := *b
		res = append(res, &copyB)
	}
	return res, nil
}

func (s *MemoryStore) VerifyBusiness(bizID, adminID string, approve bool) (*models.Business, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.Businesses[bizID]
	if !ok {
		return nil, ErrNotFound
	}
	if approve {
		b.VerificationStatus = "VERIFIED"
	} else {
		b.VerificationStatus = "REJECTED"
	}
	b.UpdatedAt = time.Now()
	copyB := *b
	return &copyB, nil
}

// ----------------- MODERATION -----------------

func (s *MemoryStore) CreateReport(r *models.Report) (*models.Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	r.Status = "PENDING"
	if r.Priority == "" {
		r.Priority = "NORMAL"
	}
	r.CreatedAt = time.Now()
	s.Reports[r.ID] = r

	// Auto-create triage case
	caseID := uuid.NewString()
	s.ModerationCases[caseID] = &models.ModerationCase{
		ID:         caseID,
		ReportID:   &r.ID,
		TargetType: r.TargetType,
		TargetID:   r.TargetID,
		Status:     "OPEN",
		CreatedAt:  time.Now(),
	}

	copyR := *r
	return &copyR, nil
}

func (s *MemoryStore) ListReports(status string) ([]*models.Report, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Report, 0)
	for _, r := range s.Reports {
		if status == "" || r.Status == status {
			copyR := *r
			res = append(res, &copyR)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.After(res[j].CreatedAt)
	})
	return res, nil
}

func (s *MemoryStore) CreateModerationCase(reportID, targetType, targetID string) (*models.ModerationCase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := &models.ModerationCase{
		ID:         uuid.NewString(),
		ReportID:   &reportID,
		TargetType: targetType,
		TargetID:   targetID,
		Status:     "OPEN",
		CreatedAt:  time.Now(),
	}
	s.ModerationCases[c.ID] = c
	copyC := *c
	return &copyC, nil
}

func (s *MemoryStore) ResolveModerationCase(caseID, adminID, decision, reason string) (*models.ModerationCase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.ModerationCases[caseID]
	if !ok {
		return nil, ErrNotFound
	}
	now := time.Now()
	c.AssignedAdminID = &adminID
	c.Status = "RESOLVED"
	c.Decision = &decision
	c.DecisionReason = &reason
	c.ResolvedAt = &now

	if c.ReportID != nil {
		if r, rOk := s.Reports[*c.ReportID]; rOk {
			r.Status = "ACTIONED"
			r.ResolvedAt = &now
		}
	}

	copyC := *c
	return &copyC, nil
}

func (s *MemoryStore) ApplyEnforcementAction(a *models.EnforcementAction) (*models.EnforcementAction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	a.CreatedAt = time.Now()
	a.AppealStatus = "NONE"
	s.Enforcements[a.ID] = a

	if u, ok := s.Users[a.UserID]; ok {
		if a.ActionType == "BAN" || a.ActionType == "SUSPENSION" {
			u.Status = "SUSPENDED"
		} else if a.ActionType == "RESTRICTION" {
			u.Status = "RESTRICTED"
		}
	}

	copyA := *a
	return &copyA, nil
}

func (s *MemoryStore) SubmitAppeal(app *models.Appeal) (*models.Appeal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	app.Status = "PENDING"
	app.CreatedAt = time.Now()
	s.Appeals[app.ID] = app

	if enf, ok := s.Enforcements[app.EnforcementActionID]; ok {
		enf.AppealStatus = "PENDING"
	}

	copyApp := *app
	return &copyApp, nil
}

func (s *MemoryStore) ListAppeals(status string) ([]*models.Appeal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.Appeal, 0)
	for _, a := range s.Appeals {
		if status == "" || a.Status == status {
			copyA := *a
			res = append(res, &copyA)
		}
	}
	return res, nil
}

// ----------------- NOTIFICATIONS -----------------

func (s *MemoryStore) GetNotifications(userID string) ([]*models.Notification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notifs, ok := s.Notifications[userID]
	if !ok {
		return make([]*models.Notification, 0), nil
	}
	res := make([]*models.Notification, len(notifs))
	copy(res, notifs)
	return res, nil
}

func (s *MemoryStore) MarkNotificationRead(id, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if notifs, ok := s.Notifications[userID]; ok {
		for _, n := range notifs {
			if n.ID == id {
				n.Status = "READ"
				return nil
			}
		}
	}
	return ErrNotFound
}

func (s *MemoryStore) SendNotification(n *models.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	n.Status = "UNREAD"
	n.CreatedAt = time.Now()
	s.Notifications[n.UserID] = append(s.Notifications[n.UserID], n)
}

// ----------------- ADMIN & AUDIT -----------------

func (s *MemoryStore) GetDashboardStats() (*models.DashboardStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var activeUsers, activeListings, completedDeals, pendingReports int64
	for _, u := range s.Users {
		if u.Status == "ACTIVE" {
			activeUsers++
		}
	}
	for _, l := range s.Listings {
		if l.ListingStatus == "ACTIVE" {
			activeListings++
		}
	}
	for _, d := range s.Deals {
		if d.Status == "COMPLETED_CONFIRMED" || d.Status == "COMPLETED_CLAIMED" {
			completedDeals++
		}
	}
	for _, r := range s.Reports {
		if r.Status == "PENDING" {
			pendingReports++
		}
	}

	return &models.DashboardStats{
		TotalUsers:            int64(len(s.Users)),
		ActiveUsers:           activeUsers,
		TotalListings:         int64(len(s.Listings)),
		ActiveListings:        activeListings,
		CompletedDeals:        completedDeals,
		ActiveConversations:   int64(len(s.Conversations)),
		PendingReports:        pendingReports,
		VerifiedNeighborhoods: int64(len(s.Neighborhoods)),
		LiveWSSessions:        s.LiveWSSessionsCount,
	}, nil
}

func (s *MemoryStore) RecordAuditLog(log *models.AuditLog) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if log.ID == "" {
		log.ID = uuid.NewString()
	}
	log.CreatedAt = time.Now()
	s.AuditLogs = append(s.AuditLogs, log)
}

func (s *MemoryStore) ListAuditLogs(limit int) ([]*models.AuditLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := len(s.AuditLogs)
	if limit <= 0 || limit > n {
		limit = n
	}
	res := make([]*models.AuditLog, limit)
	for i := 0; i < limit; i++ {
		res[i] = s.AuditLogs[n-1-i]
	}
	return res, nil
}

func (s *MemoryStore) GetAllUsers() ([]*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*models.User, 0, len(s.Users))
	for _, u := range s.Users {
		copyU := *u
		copyU.Profile = s.Profiles[u.ID]
		res = append(res, &copyU)
	}
	return res, nil
}

func (s *MemoryStore) UpdateUserStatus(userID, status, adminID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.Users[userID]
	if !ok {
		return ErrNotFound
	}
	u.Status = status
	u.UpdatedAt = time.Now()

	s.AuditLogs = append(s.AuditLogs, &models.AuditLog{
		ID:           uuid.NewString(),
		ActorType:    "ADMIN",
		ActorID:      adminID,
		Action:       "USER_STATUS_UPDATE",
		ResourceType: "USER",
		ResourceID:   userID,
		Result:       "SUCCESS",
		CreatedAt:    time.Now(),
	})
	return nil
}

// Haversine formula calculation in meters
func HaversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000.0 // Earth radius in meters
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
