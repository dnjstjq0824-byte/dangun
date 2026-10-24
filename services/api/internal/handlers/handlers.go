package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/booking-network/bn-near/services/api/internal/middleware"
	"github.com/booking-network/bn-near/services/api/internal/models"
	"github.com/booking-network/bn-near/services/api/internal/storage"
	"github.com/google/uuid"
)

type API struct {
	Store storage.Storage
}

func NewAPI(store storage.Storage) *API {
	return &API{Store: store}
}

func (api *API) RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (api *API) RespondError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	})
}

// ----------------- HEALTH & SYSTEM -----------------

func (api *API) Health(w http.ResponseWriter, r *http.Request) {
	api.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "bn-near-api",
		"version": "1.0.0",
		"uptime":  "running",
	})
}

func (api *API) Ready(w http.ResponseWriter, r *http.Request) {
	api.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"ready":    true,
		"database": "connected",
		"cache":    "connected",
	})
}

// ----------------- AUTH HANDLERS -----------------

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
	Locale      string `json:"locale"`
	CountryCode string `json:"country_code"`
}

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request payload")
		return
	}

	if req.Email == "" || len(req.Password) < 8 || len(req.Nickname) < 2 {
		api.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email, password (min 8 chars), and nickname (min 2 chars) are required")
		return
	}

	u := &models.User{
		Email:       req.Email,
		Nickname:    req.Nickname,
		Locale:      req.Locale,
		CountryCode: req.CountryCode,
	}

	created, err := api.Store.CreateUser(u, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			api.RespondError(w, http.StatusConflict, "USER_ALREADY_EXISTS", "A user with this email already exists")
			return
		}
		api.RespondError(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}

	token, refreshToken, expiresAt, _ := middleware.GenerateTokens(created)

	api.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"user":          created,
		"access_token":  token,
		"refresh_token": refreshToken,
		"expires_at":    expiresAt.Format(time.RFC3339),
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request payload")
		return
	}

	user, err := api.Store.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		api.RespondError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		return
	}

	token, refreshToken, expiresAt, _ := middleware.GenerateTokens(user)

	sess := &models.UserSession{
		ID:               uuid.NewString(),
		UserID:           user.ID,
		RefreshTokenHash: refreshToken,
		DeviceID:         req.DeviceID,
		ExpiresAt:        expiresAt,
		CreatedAt:        time.Now(),
	}
	_ = api.Store.CreateSession(sess)

	api.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"user":          user,
		"access_token":  token,
		"refresh_token": refreshToken,
		"expires_at":    expiresAt.Format(time.RFC3339),
	})
}

func (api *API) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		api.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Active session required")
		return
	}

	token, newRefresh, expiresAt, _ := middleware.GenerateTokens(user)
	api.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  token,
		"refresh_token": newRefresh,
		"expires_at":    expiresAt.Format(time.RFC3339),
	})
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	api.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

func (api *API) GetMe(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		api.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}
	u, _ := api.Store.GetUserByID(user.ID)
	api.RespondJSON(w, http.StatusOK, u)
}

func (api *API) UpdateMe(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		api.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	var req struct {
		Bio         string `json:"bio"`
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid payload")
		return
	}

	prof, err := api.Store.UpdateUserProfile(user.ID, req.Bio, req.DisplayName)
	if err != nil {
		api.RespondError(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, prof)
}

func (api *API) DeleteMe(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		api.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}
	_ = api.Store.SoftDeleteUser(user.ID)
	api.RespondJSON(w, http.StatusOK, map[string]string{"message": "Account marked for deletion"})
}

func (api *API) BlockUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	var req struct {
		BlockedUserID string `json:"blocked_user_id"`
		Reason        string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.BlockedUserID == "" {
		api.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "blocked_user_id required")
		return
	}
	_ = api.Store.BlockUser(user.ID, req.BlockedUserID, req.Reason)
	api.RespondJSON(w, http.StatusOK, map[string]string{"message": "User blocked"})
}

func (api *API) UnblockUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	blockedID := r.URL.Query().Get("user_id")
	_ = api.Store.UnblockUser(user.ID, blockedID)
	api.RespondJSON(w, http.StatusOK, map[string]string{"message": "User unblocked"})
}

// ----------------- GEOGRAPHY & LOCATION -----------------

func (api *API) ListCountries(w http.ResponseWriter, r *http.Request) {
	countries, _ := api.Store.ListCountries()
	api.RespondJSON(w, http.StatusOK, countries)
}

func (api *API) ListCities(w http.ResponseWriter, r *http.Request) {
	countryCode := r.URL.Query().Get("country_code")
	cities, _ := api.Store.ListCities(countryCode)
	api.RespondJSON(w, http.StatusOK, cities)
}

func (api *API) SearchNeighborhoods(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	cityID := r.URL.Query().Get("city_id")

	var latPtr, lngPtr *float64
	if latStr := r.URL.Query().Get("lat"); latStr != "" {
		if latVal, err := strconv.ParseFloat(latStr, 64); err == nil {
			latPtr = &latVal
		}
	}
	if lngStr := r.URL.Query().Get("lng"); lngStr != "" {
		if lngVal, err := strconv.ParseFloat(lngStr, 64); err == nil {
			lngPtr = &lngVal
		}
	}

	res, _ := api.Store.SearchNeighborhoods(q, cityID, latPtr, lngPtr)
	api.RespondJSON(w, http.StatusOK, res)
}

func (api *API) GetUserNeighborhoods(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	un, _ := api.Store.GetUserNeighborhoods(user.ID)
	api.RespondJSON(w, http.StatusOK, un)
}

type VerifyLocationRequest struct {
	NeighborhoodID string  `json:"neighborhood_id"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	AccuracyMeters float64 `json:"accuracy_meters"`
}

func (api *API) VerifyNeighborhoodLocation(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	var req VerifyLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid payload")
		return
	}

	if req.NeighborhoodID == "" || req.Latitude == 0 || req.Longitude == 0 {
		api.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "neighborhood_id, latitude, longitude required")
		return
	}

	result, err := api.Store.VerifyNeighborhood(user.ID, req.NeighborhoodID, req.Latitude, req.Longitude, req.AccuracyMeters)
	if err != nil {
		api.RespondError(w, http.StatusUnprocessableEntity, "VERIFICATION_FAILED", err.Error())
		return
	}

	api.RespondJSON(w, http.StatusOK, result)
}

// ----------------- MARKETPLACE LISTINGS -----------------

func (api *API) ListCategories(w http.ResponseWriter, r *http.Request) {
	if mem, ok := api.Store.(*storage.MemoryStore); ok {
		cats := make([]*models.Category, 0, len(mem.Categories))
		for _, c := range mem.Categories {
			cats = append(cats, c)
		}
		api.RespondJSON(w, http.StatusOK, cats)
		return
	}
	api.RespondJSON(w, http.StatusOK, []interface{}{})
}

func (api *API) ListListings(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	viewerID := ""
	if user != nil {
		viewerID = user.ID
	}

	neighborhoodID := r.URL.Query().Get("neighborhood_id")
	categoryID := r.URL.Query().Get("category_id")
	condition := r.URL.Query().Get("condition")
	sortOrder := r.URL.Query().Get("sort")

	var minPrice, maxPrice *int64
	if minStr := r.URL.Query().Get("min_price"); minStr != "" {
		if val, err := strconv.ParseInt(minStr, 10, 64); err == nil {
			minPrice = &val
		}
	}
	if maxStr := r.URL.Query().Get("max_price"); maxStr != "" {
		if val, err := strconv.ParseInt(maxStr, 10, 64); err == nil {
			maxPrice = &val
		}
	}

	var isFree *bool
	if freeStr := r.URL.Query().Get("is_free"); freeStr != "" {
		b := freeStr == "true" || freeStr == "1"
		isFree = &b
	}

	radius := 10000.0 // 10km default
	if rStr := r.URL.Query().Get("radius_meters"); rStr != "" {
		if rVal, err := strconv.ParseFloat(rStr, 64); err == nil && rVal > 0 {
			radius = rVal
		}
	}

	var latPtr, lngPtr *float64
	if latStr := r.URL.Query().Get("lat"); latStr != "" {
		if latVal, err := strconv.ParseFloat(latStr, 64); err == nil {
			latPtr = &latVal
		}
	}
	if lngStr := r.URL.Query().Get("lng"); lngStr != "" {
		if lngVal, err := strconv.ParseFloat(lngStr, 64); err == nil {
			lngPtr = &lngVal
		}
	}

	listings, err := api.Store.ListListings(neighborhoodID, categoryID, condition, minPrice, maxPrice, isFree, latPtr, lngPtr, radius, sortOrder, viewerID)
	if err != nil {
		api.RespondError(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, listings)
}

func (api *API) CreateListing(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	var req models.Listing
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid payload")
		return
	}

	if req.Title == "" || req.CategoryID == "" || req.NeighborhoodID == "" {
		api.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Title, category_id, neighborhood_id required")
		return
	}

	req.SellerUserID = user.ID
	created, err := api.Store.CreateListing(&req)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	api.RespondJSON(w, http.StatusCreated, created)
}

func (api *API) GetListing(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/listings/")
	id = strings.Split(id, "/")[0]

	user := middleware.GetUserFromContext(r.Context())
	viewerID := ""
	if user != nil {
		viewerID = user.ID
	}

	l, err := api.Store.GetListingByID(id, viewerID)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Listing not found")
		return
	}
	api.RespondJSON(w, http.StatusOK, l)
}

func (api *API) UpdateListing(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/listings/")
	id = strings.Split(id, "/")[0]

	var req models.Listing
	_ = json.NewDecoder(r.Body).Decode(&req)
	req.ID = id
	req.SellerUserID = user.ID

	updated, err := api.Store.UpdateListing(&req)
	if err != nil {
		if errors.Is(err, storage.ErrUnauthorized) {
			api.RespondError(w, http.StatusForbidden, "FORBIDDEN", "You do not own this listing")
			return
		}
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, updated)
}

func (api *API) DeleteListing(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/listings/")
	id = strings.Split(id, "/")[0]

	if err := api.Store.DeleteListing(id, user.ID); err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, map[string]string{"message": "Listing removed"})
}

func (api *API) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	var req struct {
		ListingID string `json:"listing_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	isFav, err := api.Store.ToggleFavorite(user.ID, req.ListingID)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"listing_id":   req.ListingID,
		"is_favorited": isFav,
	})
}

func (api *API) GetMyFavorites(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	favs, _ := api.Store.GetUserFavorites(user.ID)
	api.RespondJSON(w, http.StatusOK, favs)
}

func (api *API) GetMyListings(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	myListings, _ := api.Store.GetUserListings(user.ID)
	api.RespondJSON(w, http.StatusOK, myListings)
}

// ----------------- CHAT HANDLERS -----------------

func (api *API) ListConversations(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	convs, _ := api.Store.GetUserConversations(user.ID)
	api.RespondJSON(w, http.StatusOK, convs)
}

func (api *API) StartConversation(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	var req struct {
		ListingID string `json:"listing_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	conv, err := api.Store.GetOrCreateConversation(req.ListingID, user.ID)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "CHAT_INIT_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, conv)
}

func (api *API) GetConversation(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	conv, err := api.Store.GetConversationByID(id, user.ID)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Conversation not found or unauthorized")
		return
	}
	api.RespondJSON(w, http.StatusOK, conv)
}

func (api *API) GetMessages(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	convID := parts[len(parts)-2]

	conv, err := api.Store.GetConversationByID(convID, user.ID)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Conversation unauthorized")
		return
	}

	sinceSeq := int64(0)
	if sStr := r.URL.Query().Get("since_sequence"); sStr != "" {
		if sVal, err := strconv.ParseInt(sStr, 10, 64); err == nil {
			sinceSeq = sVal
		}
	}

	msgs, _ := api.Store.GetMessages(conv.ID, sinceSeq)
	api.RespondJSON(w, http.StatusOK, msgs)
}

func (api *API) SendMessageREST(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	convID := parts[len(parts)-2]

	var req struct {
		ClientMessageID string `json:"client_message_id"`
		MessageType     string `json:"message_type"`
		Body            string `json:"body"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.ClientMessageID == "" {
		req.ClientMessageID = uuid.NewString()
	}
	if req.MessageType == "" {
		req.MessageType = "TEXT"
	}

	msg := &models.ChatMessage{
		ConversationID:  convID,
		SenderUserID:    user.ID,
		ClientMessageID: req.ClientMessageID,
		MessageType:     req.MessageType,
		Body:            req.Body,
	}

	saved, err := api.Store.SaveMessage(msg)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "MESSAGE_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusCreated, saved)
}

func (api *API) MarkRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	convID := parts[len(parts)-2]

	var req struct {
		SequenceNo int64 `json:"sequence_no"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	_ = api.Store.MarkMessagesRead(convID, user.ID, req.SequenceNo)
	api.RespondJSON(w, http.StatusOK, map[string]string{"status": "READ"})
}

// ----------------- DEALS & OFFERS -----------------

func (api *API) CreatePriceOffer(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	listingID := parts[len(parts)-2]

	var req struct {
		AmountMinor int64 `json:"amount_minor"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	offer := &models.PriceOffer{
		ListingID:   listingID,
		BuyerUserID: user.ID,
		AmountMinor: req.AmountMinor,
	}
	created, err := api.Store.CreateOffer(offer)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "OFFER_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusCreated, created)
}

func (api *API) AcceptOffer(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	offerID := parts[len(parts)-2]

	updated, err := api.Store.UpdateOfferStatus(offerID, user.ID, "ACCEPTED")
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "OFFER_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, updated)
}

func (api *API) DeclineOffer(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	offerID := parts[len(parts)-2]

	updated, err := api.Store.UpdateOfferStatus(offerID, user.ID, "DECLINED")
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "OFFER_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, updated)
}

func (api *API) ReserveDeal(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	dealOrListingID := parts[len(parts)-2]

	var req struct {
		BuyerUserID string `json:"buyer_user_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	deal, err := api.Store.ReserveDeal(dealOrListingID, user.ID, req.BuyerUserID)
	if err != nil {
		if errors.Is(err, storage.ErrConflict) {
			api.RespondError(w, http.StatusConflict, "RESERVATION_CONFLICT", err.Error())
			return
		}
		api.RespondError(w, http.StatusBadRequest, "RESERVE_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, deal)
}

func (api *API) ScheduleMeeting(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	dealID := parts[len(parts)-2]

	var req struct {
		PlaceName   string `json:"place_name"`
		MeetingTime string `json:"meeting_time"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	t, _ := time.Parse(time.RFC3339, req.MeetingTime)
	if t.IsZero() {
		t = time.Now().Add(24 * time.Hour)
	}

	deal, err := api.Store.ScheduleMeeting(dealID, user.ID, req.PlaceName, t)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "MEETING_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, deal)
}

func (api *API) SellerCompleteDeal(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	dealID := parts[len(parts)-2]

	deal, err := api.Store.SellerCompleteDeal(dealID, user.ID)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "COMPLETE_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, deal)
}

func (api *API) BuyerConfirmDeal(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	dealID := parts[len(parts)-2]

	deal, err := api.Store.BuyerConfirmDeal(dealID, user.ID)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "CONFIRM_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, deal)
}

func (api *API) CancelDeal(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	dealID := parts[len(parts)-2]

	deal, err := api.Store.CancelDeal(dealID, user.ID)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "CANCEL_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, deal)
}

func (api *API) GetDeal(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	dealID := parts[len(parts)-1]

	deal, err := api.Store.GetDealByID(dealID)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Deal not found")
		return
	}
	api.RespondJSON(w, http.StatusOK, deal)
}

func (api *API) GetMyDeals(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	deals, _ := api.Store.GetUserDeals(user.ID)
	api.RespondJSON(w, http.StatusOK, deals)
}

// ----------------- REVIEWS & REPUTATION -----------------

func (api *API) CreateReview(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	dealID := parts[len(parts)-2]

	var req struct {
		Rating  int      `json:"rating"`
		Tags    []string `json:"tags"`
		Comment string   `json:"comment"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	review := &models.DealReview{
		DealID:         dealID,
		ReviewerUserID: user.ID,
		Rating:         req.Rating,
		Tags:           req.Tags,
		Comment:        req.Comment,
	}

	created, err := api.Store.CreateDealReview(review)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			api.RespondError(w, http.StatusConflict, "REVIEW_ALREADY_EXISTS", err.Error())
			return
		}
		if errors.Is(err, storage.ErrConflict) {
			api.RespondError(w, http.StatusPreconditionFailed, "DEAL_NOT_CONFIRMED", err.Error())
			return
		}
		api.RespondError(w, http.StatusBadRequest, "REVIEW_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusCreated, created)
}

func (api *API) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	userID := parts[len(parts)-2]

	reviews, _ := api.Store.GetUserReviews(userID)
	api.RespondJSON(w, http.StatusOK, reviews)
}

func (api *API) GetUserReputation(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	userID := parts[len(parts)-2]

	prof, err := api.Store.GetUserReputation(userID)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}
	api.RespondJSON(w, http.StatusOK, prof)
}

// ----------------- COMMUNITY -----------------

func (api *API) ListCommunityPosts(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	viewerID := ""
	if user != nil {
		viewerID = user.ID
	}
	neighborhoodID := r.URL.Query().Get("neighborhood_id")
	category := r.URL.Query().Get("category")

	posts, _ := api.Store.ListCommunityPosts(neighborhoodID, category, viewerID)
	api.RespondJSON(w, http.StatusOK, posts)
}

func (api *API) CreateCommunityPost(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	var req models.CommunityPost
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.NeighborhoodID == "" || req.Title == "" || req.Body == "" {
		api.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "neighborhood_id, title, and body required")
		return
	}

	req.AuthorUserID = user.ID
	created, err := api.Store.CreateCommunityPost(&req)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "POST_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusCreated, created)
}

func (api *API) GetCommunityPost(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	user := middleware.GetUserFromContext(r.Context())
	viewerID := ""
	if user != nil {
		viewerID = user.ID
	}

	p, err := api.Store.GetCommunityPostByID(id, viewerID)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	api.RespondJSON(w, http.StatusOK, p)
}

func (api *API) GetCommunityComments(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	postID := parts[len(parts)-2]

	comments, _ := api.Store.GetCommunityComments(postID)
	api.RespondJSON(w, http.StatusOK, comments)
}

func (api *API) CreateCommunityComment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	postID := parts[len(parts)-2]

	var req models.CommunityComment
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Body == "" {
		api.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Comment body required")
		return
	}

	req.PostID = postID
	req.AuthorUserID = user.ID
	created, err := api.Store.CreateCommunityComment(&req)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "COMMENT_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusCreated, created)
}

func (api *API) ToggleCommunityLike(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	postID := parts[len(parts)-2]

	liked, err := api.Store.ToggleCommunityLike(user.ID, postID)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"post_id": postID,
		"liked":   liked,
	})
}

// ----------------- BUSINESS -----------------

func (api *API) ListBusinesses(w http.ResponseWriter, r *http.Request) {
	neighborhoodID := r.URL.Query().Get("neighborhood_id")
	businesses, _ := api.Store.ListBusinesses(neighborhoodID)
	api.RespondJSON(w, http.StatusOK, businesses)
}

func (api *API) GetBusiness(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	b, err := api.Store.GetBusinessByID(id)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Business not found")
		return
	}
	api.RespondJSON(w, http.StatusOK, b)
}

func (api *API) CreateBusiness(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	var req models.Business
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.Name == "" || req.NeighborhoodID == "" || req.PublicAddress == "" {
		api.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Name, neighborhood_id, and public_address required")
		return
	}

	req.OwnerUserID = user.ID
	created, err := api.Store.CreateBusiness(&req)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, "BUSINESS_CREATE_FAILED", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusCreated, created)
}

// ----------------- MODERATION & REPORTS -----------------

func (api *API) CreateReport(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	var req models.Report
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.TargetType == "" || req.TargetID == "" || req.ReasonCode == "" {
		api.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Target type, ID, and reason required")
		return
	}

	req.ReporterUserID = user.ID
	created, _ := api.Store.CreateReport(&req)
	api.RespondJSON(w, http.StatusCreated, created)
}

func (api *API) SubmitAppeal(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	var req models.Appeal
	_ = json.NewDecoder(r.Body).Decode(&req)

	req.UserID = user.ID
	created, _ := api.Store.SubmitAppeal(&req)
	api.RespondJSON(w, http.StatusCreated, created)
}

// ----------------- NOTIFICATIONS -----------------

func (api *API) GetNotifications(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	notifs, _ := api.Store.GetNotifications(user.ID)
	api.RespondJSON(w, http.StatusOK, notifs)
}

func (api *API) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-2]

	_ = api.Store.MarkNotificationRead(id, user.ID)
	api.RespondJSON(w, http.StatusOK, map[string]string{"status": "READ"})
}

// ----------------- ADMIN HANDLERS -----------------

func (api *API) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := api.Store.GetDashboardStats()
	if err != nil {
		api.RespondError(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, stats)
}

func (api *API) AdminListUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := api.Store.GetAllUsers()
	api.RespondJSON(w, http.StatusOK, users)
}

func (api *API) AdminUpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	admin := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	userID := parts[len(parts)-2]

	var req struct {
		Status string `json:"status"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	_ = api.Store.UpdateUserStatus(userID, req.Status, admin.ID)
	api.RespondJSON(w, http.StatusOK, map[string]string{"message": "Status updated"})
}

func (api *API) AdminListReports(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	reports, _ := api.Store.ListReports(status)
	api.RespondJSON(w, http.StatusOK, reports)
}

func (api *API) AdminResolveCase(w http.ResponseWriter, r *http.Request) {
	admin := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	caseID := parts[len(parts)-2]

	var req struct {
		Decision string `json:"decision"`
		Reason   string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	c, err := api.Store.ResolveModerationCase(caseID, admin.ID, req.Decision, req.Reason)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, c)
}

func (api *API) AdminListAppeals(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	appeals, _ := api.Store.ListAppeals(status)
	api.RespondJSON(w, http.StatusOK, appeals)
}

func (api *API) AdminVerifyBusiness(w http.ResponseWriter, r *http.Request) {
	admin := middleware.GetUserFromContext(r.Context())
	parts := strings.Split(r.URL.Path, "/")
	bizID := parts[len(parts)-2]

	var req struct {
		Approve bool `json:"approve"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	biz, err := api.Store.VerifyBusiness(bizID, admin.ID, req.Approve)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	api.RespondJSON(w, http.StatusOK, biz)
}

func (api *API) AdminListAuditLogs(w http.ResponseWriter, r *http.Request) {
	logs, _ := api.Store.ListAuditLogs(50)
	api.RespondJSON(w, http.StatusOK, logs)
}

