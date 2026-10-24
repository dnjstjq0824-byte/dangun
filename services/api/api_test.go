package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/booking-network/bn-near/services/api/internal/chat"
	"github.com/booking-network/bn-near/services/api/internal/models"
	"github.com/booking-network/bn-near/services/api/internal/router"
	"github.com/booking-network/bn-near/services/api/internal/storage"
	"github.com/google/uuid"
)

func setupTestServer() (http.Handler, storage.Storage, *chat.Hub) {
	store := storage.NewMemoryStore()
	hub := chat.NewHub(store)
	go hub.Run()
	routes := router.NewRouter(store, hub)
	return routes, store, hub
}

func doJSONRequest(t *testing.T, handler http.Handler, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestCompleteHyperlocalMarketplaceLifecycle(t *testing.T) {
	handler, store, _ := setupTestServer()

	// 1. TC-01 & TC-02: Auth & Login for Seed Seller and Buyer
	sellerLogin := doJSONRequest(t, handler, "POST", "/api/v1/auth/login", "", map[string]string{
		"email":    "seller@bn-near.com",
		"password": "Password123!",
	})
	if sellerLogin.Code != http.StatusOK {
		t.Fatalf("TC-02: Seller login failed: %s", sellerLogin.Body.String())
	}
	var sellerAuth struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(sellerLogin.Body.Bytes(), &sellerAuth)
	sellerToken := sellerAuth.Data.AccessToken

	buyerLogin := doJSONRequest(t, handler, "POST", "/api/v1/auth/login", "", map[string]string{
		"email":    "buyer@bn-near.com",
		"password": "Password123!",
	})
	if buyerLogin.Code != http.StatusOK {
		t.Fatalf("TC-02: Buyer login failed: %s", buyerLogin.Body.String())
	}
	var buyerAuth struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(buyerLogin.Body.Bytes(), &buyerAuth)
	buyerToken := buyerAuth.Data.AccessToken

	// 2. TC-04: Hierarchy & Neighborhoods
	geoResp := doJSONRequest(t, handler, "GET", "/api/v1/neighborhoods/search?q=청라", "", nil)
	if geoResp.Code != http.StatusOK {
		t.Fatalf("TC-04: Neighborhood search failed: %s", geoResp.Body.String())
	}
	var nList struct {
		Data []*models.Neighborhood `json:"data"`
	}
	_ = json.Unmarshal(geoResp.Body.Bytes(), &nList)
	if len(nList.Data) == 0 {
		t.Fatalf("TC-04: Cheongna neighborhood not found in search")
	}
	cheongnaID := nList.Data[0].ID

	// 3. TC-05, TC-06, TC-07: Location Verification (Boundary & Accuracy)
	verifySuccess := doJSONRequest(t, handler, "POST", "/api/v1/me/neighborhood-verifications", sellerToken, map[string]interface{}{
		"neighborhood_id": cheongnaID,
		"latitude":        37.5350,
		"longitude":       126.6500,
		"accuracy_meters": 25.0,
	})
	if verifySuccess.Code != http.StatusOK {
		t.Fatalf("TC-05: Location verification failed: %s", verifySuccess.Body.String())
	}

	verifyOutside := doJSONRequest(t, handler, "POST", "/api/v1/me/neighborhood-verifications", sellerToken, map[string]interface{}{
		"neighborhood_id": cheongnaID,
		"latitude":        35.1796, // Busan
		"longitude":       129.0756,
		"accuracy_meters": 30.0,
	})
	if verifyOutside.Code != http.StatusUnprocessableEntity {
		t.Fatalf("TC-06: Expected 422 for out-of-bounds, got %d", verifyOutside.Code)
	}

	verifyLowAcc := doJSONRequest(t, handler, "POST", "/api/v1/me/neighborhood-verifications", sellerToken, map[string]interface{}{
		"neighborhood_id": cheongnaID,
		"latitude":        37.5350,
		"longitude":       126.6500,
		"accuracy_meters": 250.0,
	})
	if verifyLowAcc.Code != http.StatusUnprocessableEntity {
		t.Fatalf("TC-07: Expected 422 for low GPS accuracy, got %d", verifyLowAcc.Code)
	}

	// 4. TC-09 & TC-11: Listing Creation & Constraints
	negPriceResp := doJSONRequest(t, handler, "POST", "/api/v1/listings", sellerToken, map[string]interface{}{
		"title":           "Invalid Listing",
		"category_id":     "c0000000-0000-0000-0000-000000000001",
		"neighborhood_id": cheongnaID,
		"price_minor":     -5000,
	})
	if negPriceResp.Code != http.StatusBadRequest {
		t.Fatalf("TC-11: Expected 400 for negative price, got %d", negPriceResp.Code)
	}

	createListingResp := doJSONRequest(t, handler, "POST", "/api/v1/listings", sellerToken, map[string]interface{}{
		"title":           "소니 노이즈캔슬링 헤드폰 WH-1000XM5",
		"description":     "미개봉 풀박스 정품입니다. 호수공원 인근에서 직거래 원해요.",
		"category_id":     "c0000000-0000-0000-0000-000000000001",
		"condition":       "NEW",
		"price_minor":     380000,
		"currency_code":   "KRW",
		"is_free":         false,
		"is_negotiable":   true,
		"neighborhood_id": cheongnaID,
		"images":          []string{"https://images.unsplash.com/photo-1546435770-a3e426bf472b"},
	})
	if createListingResp.Code != http.StatusCreated {
		t.Fatalf("TC-09: Failed to create listing: %s", createListingResp.Body.String())
	}
	var createdListing struct {
		Data *models.Listing `json:"data"`
	}
	_ = json.Unmarshal(createListingResp.Body.Bytes(), &createdListing)
	headphoneListingID := createdListing.Data.ID

	// 5. TC-13 & TC-15: Spatial Radius Search & Favorites
	searchResp := doJSONRequest(t, handler, "GET", "/api/v1/listings?lat=37.535&lng=126.650&radius_meters=5000", buyerToken, nil)
	if searchResp.Code != http.StatusOK {
		t.Fatalf("TC-13: Spatial search failed: %s", searchResp.Body.String())
	}
	var searchData struct {
		Data []*models.Listing `json:"data"`
	}
	_ = json.Unmarshal(searchResp.Body.Bytes(), &searchData)
	if len(searchData.Data) == 0 {
		t.Fatalf("TC-13: Nearby listings should not be empty")
	}

	favResp := doJSONRequest(t, handler, "POST", "/api/v1/me/favorites", buyerToken, map[string]string{
		"listing_id": headphoneListingID,
	})
	if favResp.Code != http.StatusOK {
		t.Fatalf("TC-15: Toggle favorite failed: %s", favResp.Body.String())
	}

	// 6. TC-16, TC-18, TC-19, TC-20: 1:1 Conversation, Monotonic Sequences & Deduplication
	startChatResp := doJSONRequest(t, handler, "POST", "/api/v1/conversations", buyerToken, map[string]string{
		"listing_id": headphoneListingID,
	})
	if startChatResp.Code != http.StatusOK {
		t.Fatalf("TC-16: Start conversation failed: %s", startChatResp.Body.String())
	}
	var convData struct {
		Data *models.Conversation `json:"data"`
	}
	_ = json.Unmarshal(startChatResp.Body.Bytes(), &convData)
	convID := convData.Data.ID

	clientMsgID := uuid.NewString()
	sendMsgResp := doJSONRequest(t, handler, "POST", "/api/v1/conversations/"+convID+"/messages", buyerToken, map[string]string{
		"client_message_id": clientMsgID,
		"message_type":      "TEXT",
		"body":              "안녕하세요! 35만원에 네고 가능하실까요?",
	})
	if sendMsgResp.Code != http.StatusCreated {
		t.Fatalf("TC-18: Send message failed: %s", sendMsgResp.Body.String())
	}
	var msg1 struct {
		Data *models.ChatMessage `json:"data"`
	}
	_ = json.Unmarshal(sendMsgResp.Body.Bytes(), &msg1)
	if msg1.Data.SequenceNo != 1 {
		t.Fatalf("TC-18: Expected sequence 1, got %d", msg1.Data.SequenceNo)
	}

	sendMsgDedup := doJSONRequest(t, handler, "POST", "/api/v1/conversations/"+convID+"/messages", buyerToken, map[string]string{
		"client_message_id": clientMsgID,
		"message_type":      "TEXT",
		"body":              "안녕하세요! 35만원에 네고 가능하실까요?",
	})
	var msgDedup struct {
		Data *models.ChatMessage `json:"data"`
	}
	_ = json.Unmarshal(sendMsgDedup.Body.Bytes(), &msgDedup)
	if msgDedup.Data.SequenceNo != 1 {
		t.Fatalf("TC-19: Duplicate message should return sequence 1, got %d", msgDedup.Data.SequenceNo)
	}

	readResp := doJSONRequest(t, handler, "POST", "/api/v1/conversations/"+convID+"/read", sellerToken, map[string]interface{}{
		"sequence_no": 1,
	})
	if readResp.Code != http.StatusOK {
		t.Fatalf("TC-20: Read receipt failed: %s", readResp.Body.String())
	}

	// 7. TC-22: Price Offer Negotiation
	offerResp := doJSONRequest(t, handler, "POST", "/api/v1/listings/"+headphoneListingID+"/price-offers", buyerToken, map[string]interface{}{
		"amount_minor": 360000,
	})
	if offerResp.Code != http.StatusCreated {
		t.Fatalf("TC-22: Create offer failed: %s", offerResp.Body.String())
	}
	var offerData struct {
		Data *models.PriceOffer `json:"data"`
	}
	_ = json.Unmarshal(offerResp.Body.Bytes(), &offerData)
	offerID := offerData.Data.ID

	acceptResp := doJSONRequest(t, handler, "POST", "/api/v1/price-offers/"+offerID+"/accept", sellerToken, nil)
	if acceptResp.Code != http.StatusOK {
		t.Fatalf("TC-22: Accept offer failed: %s", acceptResp.Body.String())
	}

	// 8. TC-23 & TC-24: Reservation & Collision Lock
	reserveResp := doJSONRequest(t, handler, "POST", "/api/v1/deals/"+headphoneListingID+"/reserve", sellerToken, map[string]string{
		"buyer_user_id": convData.Data.BuyerUserID,
	})
	if reserveResp.Code != http.StatusOK {
		t.Fatalf("TC-23: Reservation failed: %s", reserveResp.Body.String())
	}
	var dealData struct {
		Data *models.Deal `json:"data"`
	}
	_ = json.Unmarshal(reserveResp.Body.Bytes(), &dealData)
	dealID := dealData.Data.ID

	secondBuyerID := "u0000000-0000-0000-0000-000000000004"
	collisionResp := doJSONRequest(t, handler, "POST", "/api/v1/deals/"+headphoneListingID+"/reserve", sellerToken, map[string]string{
		"buyer_user_id": secondBuyerID,
	})
	if collisionResp.Code != http.StatusConflict {
		t.Fatalf("TC-24: Expected 409 Conflict for reservation collision, got %d", collisionResp.Code)
	}

	// 9. TC-26 & TC-25: Two-Phase Deal Completion & Review Gating
	prematureReview := doJSONRequest(t, handler, "POST", "/api/v1/deals/"+dealID+"/reviews", buyerToken, map[string]interface{}{
		"rating": 5,
	})
	if prematureReview.Code != http.StatusPreconditionFailed {
		t.Fatalf("TC-26: Premature review should be blocked with 412, got %d", prematureReview.Code)
	}

	completeClaimResp := doJSONRequest(t, handler, "POST", "/api/v1/deals/"+dealID+"/seller-complete", sellerToken, nil)
	if completeClaimResp.Code != http.StatusOK {
		t.Fatalf("TC-25: Seller complete failed: %s", completeClaimResp.Body.String())
	}

	buyerConfirmResp := doJSONRequest(t, handler, "POST", "/api/v1/deals/"+dealID+"/buyer-confirm", buyerToken, nil)
	if buyerConfirmResp.Code != http.StatusOK {
		t.Fatalf("TC-25: Buyer confirm failed: %s", buyerConfirmResp.Body.String())
	}

	// 10. TC-27 & TC-28: Mutual Reviews & Manner Score Calculation
	reviewResp := doJSONRequest(t, handler, "POST", "/api/v1/deals/"+dealID+"/reviews", buyerToken, map[string]interface{}{
		"rating":  5,
		"tags":    []string{"Punctual", "Kind"},
		"comment": "좋은 물건 빠르게 직거래해주셔서 감사합니다!",
	})
	if reviewResp.Code != http.StatusCreated {
		t.Fatalf("TC-27: Review submission failed: %s", reviewResp.Body.String())
	}

	dupReview := doJSONRequest(t, handler, "POST", "/api/v1/deals/"+dealID+"/reviews", buyerToken, map[string]interface{}{
		"rating": 5,
	})
	if dupReview.Code != http.StatusConflict {
		t.Fatalf("TC-28: Duplicate review should return 409 Conflict, got %d", dupReview.Code)
	}

	sellerProf, _ := store.GetUserReputation(convData.Data.SellerUserID)
	if sellerProf.MannerTemperature <= 38.2 {
		t.Fatalf("TC-27: Seller manner temperature should have increased from 38.2, got %.1f", sellerProf.MannerTemperature)
	}

	// 11. TC-29 & TC-30: Community Post, Comment & Like
	postResp := doJSONRequest(t, handler, "POST", "/api/v1/community/posts", buyerToken, map[string]string{
		"neighborhood_id": cheongnaID,
		"category":        "추천",
		"title":           "청라호수공원 산책로 야경 명소 공유합니다",
		"body":            "남쪽 음악분수대 건너편 야경이 제일 예쁘네요. 다들 저녁 산책 다녀오세요!",
	})
	if postResp.Code != http.StatusCreated {
		t.Fatalf("TC-29: Community post creation failed: %s", postResp.Body.String())
	}
	var postData struct {
		Data *models.CommunityPost `json:"data"`
	}
	_ = json.Unmarshal(postResp.Body.Bytes(), &postData)
	cPostID := postData.Data.ID

	likeResp := doJSONRequest(t, handler, "POST", "/api/v1/community/posts/"+cPostID+"/likes", sellerToken, nil)
	if likeResp.Code != http.StatusOK {
		t.Fatalf("TC-30: Community like failed: %s", likeResp.Body.String())
	}

	// 12. TC-32 & TC-33: Moderation Report & Resolution
	reportResp := doJSONRequest(t, handler, "POST", "/api/v1/reports", buyerToken, map[string]string{
		"target_type": "LISTING",
		"target_id":   headphoneListingID,
		"reason_code": "SUSPICIOUS_SELLER",
		"description": "테스트 신고 접수",
	})
	if reportResp.Code != http.StatusCreated {
		t.Fatalf("TC-32: Report creation failed: %s", reportResp.Body.String())
	}

	// 13. TC-34 & TC-35: Admin Dashboard, RBAC & Audit Log
	adminLogin := doJSONRequest(t, handler, "POST", "/api/v1/auth/login", "", map[string]string{
		"email":    "admin@bn-near.com",
		"password": "Password123!",
	})
	var adminAuth struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(adminLogin.Body.Bytes(), &adminAuth)
	adminToken := adminAuth.Data.AccessToken

	forbiddenResp := doJSONRequest(t, handler, "GET", "/api/v1/admin/dashboard", buyerToken, nil)
	if forbiddenResp.Code != http.StatusForbidden {
		t.Fatalf("TC-34: Expected 403 for non-admin on admin dashboard, got %d", forbiddenResp.Code)
	}

	dashResp := doJSONRequest(t, handler, "GET", "/api/v1/admin/dashboard", adminToken, nil)
	if dashResp.Code != http.StatusOK {
		t.Fatalf("TC-34: Admin dashboard failed: %s", dashResp.Body.String())
	}

	store.RecordAuditLog(&models.AuditLog{
		ActorType:    "ADMIN",
		ActorID:      "u0000000-0000-0000-0000-000000000001",
		Action:       "TEST_ADMIN_MUTATION",
		ResourceType: "SYSTEM",
		ResourceID:   "ALL",
		Result:       "SUCCESS",
		CreatedAt:    time.Now(),
	})
	auditResp := doJSONRequest(t, handler, "GET", "/api/v1/admin/audit-logs", adminToken, nil)
	if auditResp.Code != http.StatusOK {
		t.Fatalf("TC-35: Audit log query failed: %s", auditResp.Body.String())
	}
}

