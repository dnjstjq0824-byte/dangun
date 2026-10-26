package router

import (
	"net/http"
	"strings"

	"github.com/booking-network/bn-near/services/api/internal/chat"
	"github.com/booking-network/bn-near/services/api/internal/handlers"
	"github.com/booking-network/bn-near/services/api/internal/middleware"
	"github.com/booking-network/bn-near/services/api/internal/storage"
)

func NewRouter(store storage.Storage, hub *chat.Hub) http.Handler {
	api := handlers.NewAPI(store)
	mux := http.NewServeMux()

	// 1. System & Health
	mux.HandleFunc("/health", api.Health)
	mux.HandleFunc("/healthz", api.Health)
	mux.HandleFunc("/ready", api.Ready)
	mux.HandleFunc("/readyz", api.Ready)

	// 2. WebSocket Gateway
	mux.HandleFunc("/api/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})

	// 3. Dispatcher for REST APIs
	authMW := middleware.RequireAuth(store)
	adminMW := middleware.RequireAdmin(store)

	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		// Public Routes
		switch {
		case path == "/api/v1/auth/register" && method == http.MethodPost:
			api.Register(w, r)
			return
		case path == "/api/v1/auth/login" && method == http.MethodPost:
			api.Login(w, r)
			return
		case path == "/api/v1/countries" && method == http.MethodGet:
			api.ListCountries(w, r)
			return
		case path == "/api/v1/cities" && method == http.MethodGet:
			api.ListCities(w, r)
			return
		case (path == "/api/v1/health" || path == "/healthz") && method == http.MethodGet:
			api.Health(w, r)
			return
		case (path == "/api/v1/ready" || path == "/readyz") && method == http.MethodGet:
			api.Ready(w, r)
			return
		case (path == "/api/v1/neighborhoods" || path == "/api/v1/neighborhoods/search") && method == http.MethodGet:
			api.SearchNeighborhoods(w, r)
			return
		case path == "/api/v1/categories" && method == http.MethodGet:
			api.ListCategories(w, r)
			return
		case path == "/api/v1/listings" && method == http.MethodGet:
			middleware.OptionalAuth(store)(http.HandlerFunc(api.ListListings)).ServeHTTP(w, r)
			return
		case strings.HasPrefix(path, "/api/v1/listings/") && method == http.MethodGet && !strings.Contains(path, "/price-offers"):
			middleware.OptionalAuth(store)(http.HandlerFunc(api.GetListing)).ServeHTTP(w, r)
			return
		case path == "/api/v1/community/posts" && method == http.MethodGet:
			middleware.OptionalAuth(store)(http.HandlerFunc(api.ListCommunityPosts)).ServeHTTP(w, r)
			return
		case strings.HasPrefix(path, "/api/v1/community/posts/") && strings.HasSuffix(path, "/comments") && method == http.MethodGet:
			api.GetCommunityComments(w, r)
			return
		case strings.HasPrefix(path, "/api/v1/community/posts/") && method == http.MethodGet:
			middleware.OptionalAuth(store)(http.HandlerFunc(api.GetCommunityPost)).ServeHTTP(w, r)
			return
		case path == "/api/v1/businesses" && method == http.MethodGet:
			api.ListBusinesses(w, r)
			return
		case strings.HasPrefix(path, "/api/v1/businesses/") && method == http.MethodGet && !strings.HasPrefix(path, "/api/v1/admin/"):
			api.GetBusiness(w, r)
			return
		case strings.HasPrefix(path, "/api/v1/users/") && strings.HasSuffix(path, "/reviews") && method == http.MethodGet:
			api.GetUserReviews(w, r)
			return
		case strings.HasPrefix(path, "/api/v1/users/") && strings.HasSuffix(path, "/reputation") && method == http.MethodGet:
			api.GetUserReputation(w, r)
			return
		}

		// Authenticated Routes
		authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case path == "/api/v1/auth/refresh" && method == http.MethodPost:
				api.Refresh(w, r)
			case path == "/api/v1/auth/logout" && method == http.MethodPost:
				api.Logout(w, r)
			case path == "/api/v1/me" && method == http.MethodGet:
				api.GetMe(w, r)
			case path == "/api/v1/me" && method == http.MethodPatch:
				api.UpdateMe(w, r)
			case path == "/api/v1/me" && method == http.MethodDelete:
				api.DeleteMe(w, r)
			case path == "/api/v1/me/blocks" && method == http.MethodPost:
				api.BlockUser(w, r)
			case path == "/api/v1/me/blocks" && method == http.MethodDelete:
				api.UnblockUser(w, r)
			case path == "/api/v1/me/neighborhoods" && method == http.MethodGet:
				api.GetUserNeighborhoods(w, r)
			case path == "/api/v1/me/neighborhood-verifications" && method == http.MethodPost:
				api.VerifyNeighborhoodLocation(w, r)
			case path == "/api/v1/listings" && method == http.MethodPost:
				api.CreateListing(w, r)
			case strings.HasPrefix(path, "/api/v1/listings/") && method == http.MethodPatch:
				api.UpdateListing(w, r)
			case strings.HasPrefix(path, "/api/v1/listings/") && method == http.MethodDelete:
				api.DeleteListing(w, r)
			case path == "/api/v1/me/favorites" && method == http.MethodPost:
				api.ToggleFavorite(w, r)
			case path == "/api/v1/me/favorites" && method == http.MethodGet:
				api.GetMyFavorites(w, r)
			case path == "/api/v1/me/listings" && method == http.MethodGet:
				api.GetMyListings(w, r)
			case path == "/api/v1/conversations" && method == http.MethodGet:
				api.ListConversations(w, r)
			case path == "/api/v1/conversations" && method == http.MethodPost:
				api.StartConversation(w, r)
			case strings.HasPrefix(path, "/api/v1/conversations/") && strings.HasSuffix(path, "/messages") && method == http.MethodGet:
				api.GetMessages(w, r)
			case strings.HasPrefix(path, "/api/v1/conversations/") && strings.HasSuffix(path, "/messages") && method == http.MethodPost:
				api.SendMessageREST(w, r)
			case strings.HasPrefix(path, "/api/v1/conversations/") && strings.HasSuffix(path, "/read") && method == http.MethodPost:
				api.MarkRead(w, r)
			case strings.HasPrefix(path, "/api/v1/conversations/") && method == http.MethodGet:
				api.GetConversation(w, r)
			case strings.HasPrefix(path, "/api/v1/listings/") && strings.HasSuffix(path, "/price-offers") && method == http.MethodPost:
				api.CreatePriceOffer(w, r)
			case strings.HasPrefix(path, "/api/v1/price-offers/") && strings.HasSuffix(path, "/accept") && method == http.MethodPost:
				api.AcceptOffer(w, r)
			case strings.HasPrefix(path, "/api/v1/price-offers/") && strings.HasSuffix(path, "/decline") && method == http.MethodPost:
				api.DeclineOffer(w, r)
			case strings.HasPrefix(path, "/api/v1/deals/") && strings.HasSuffix(path, "/reserve") && method == http.MethodPost:
				api.ReserveDeal(w, r)
			case strings.HasPrefix(path, "/api/v1/deals/") && strings.HasSuffix(path, "/meeting") && method == http.MethodPost:
				api.ScheduleMeeting(w, r)
			case strings.HasPrefix(path, "/api/v1/deals/") && strings.HasSuffix(path, "/seller-complete") && method == http.MethodPost:
				api.SellerCompleteDeal(w, r)
			case strings.HasPrefix(path, "/api/v1/deals/") && strings.HasSuffix(path, "/buyer-confirm") && method == http.MethodPost:
				api.BuyerConfirmDeal(w, r)
			case strings.HasPrefix(path, "/api/v1/deals/") && strings.HasSuffix(path, "/cancel") && method == http.MethodPost:
				api.CancelDeal(w, r)
			case strings.HasPrefix(path, "/api/v1/deals/") && strings.HasSuffix(path, "/reviews") && method == http.MethodPost:
				api.CreateReview(w, r)
			case path == "/api/v1/me/deals" && method == http.MethodGet:
				api.GetMyDeals(w, r)
			case strings.HasPrefix(path, "/api/v1/deals/") && method == http.MethodGet:
				api.GetDeal(w, r)
			case path == "/api/v1/community/posts" && method == http.MethodPost:
				api.CreateCommunityPost(w, r)
			case strings.HasPrefix(path, "/api/v1/community/posts/") && strings.HasSuffix(path, "/comments") && method == http.MethodPost:
				api.CreateCommunityComment(w, r)
			case strings.HasPrefix(path, "/api/v1/community/posts/") && strings.HasSuffix(path, "/likes") && method == http.MethodPost:
				api.ToggleCommunityLike(w, r)
			case path == "/api/v1/me/businesses" && method == http.MethodPost:
				api.CreateBusiness(w, r)
			case path == "/api/v1/reports" && method == http.MethodPost:
				api.CreateReport(w, r)
			case path == "/api/v1/appeals" && method == http.MethodPost:
				api.SubmitAppeal(w, r)
			case path == "/api/v1/me/notifications" && method == http.MethodGet:
				api.GetNotifications(w, r)
			case strings.HasPrefix(path, "/api/v1/me/notifications/") && strings.HasSuffix(path, "/read") && method == http.MethodPatch:
				api.MarkNotificationRead(w, r)

			// Admin operations
			case path == "/api/v1/admin/dashboard" && method == http.MethodGet:
				adminMW(http.HandlerFunc(api.AdminDashboard)).ServeHTTP(w, r)
			case path == "/api/v1/admin/users" && method == http.MethodGet:
				adminMW(http.HandlerFunc(api.AdminListUsers)).ServeHTTP(w, r)
			case strings.HasPrefix(path, "/api/v1/admin/users/") && strings.HasSuffix(path, "/status") && method == http.MethodPatch:
				adminMW(http.HandlerFunc(api.AdminUpdateUserStatus)).ServeHTTP(w, r)
			case path == "/api/v1/admin/reports" && method == http.MethodGet:
				adminMW(http.HandlerFunc(api.AdminListReports)).ServeHTTP(w, r)
			case strings.HasPrefix(path, "/api/v1/admin/moderation/cases/") && strings.HasSuffix(path, "/decide") && method == http.MethodPost:
				adminMW(http.HandlerFunc(api.AdminResolveCase)).ServeHTTP(w, r)
			case path == "/api/v1/admin/appeals" && method == http.MethodGet:
				adminMW(http.HandlerFunc(api.AdminListAppeals)).ServeHTTP(w, r)
			case strings.HasPrefix(path, "/api/v1/admin/businesses/") && strings.HasSuffix(path, "/verify") && method == http.MethodPost:
				adminMW(http.HandlerFunc(api.AdminVerifyBusiness)).ServeHTTP(w, r)
			case path == "/api/v1/admin/audit-logs" && method == http.MethodGet:
				adminMW(http.HandlerFunc(api.AdminListAuditLogs)).ServeHTTP(w, r)
			default:
				api.RespondError(w, http.StatusNotFound, "ENDPOINT_NOT_FOUND", "No matching endpoint for "+method+" "+path)
			}
		})

		authMW(authHandler).ServeHTTP(w, r)
	})

	// Global Middlewares: CORS, Recoverer
	handler := middleware.CORS(mux)
	handler = middleware.Recoverer(handler)
	return handler
}

