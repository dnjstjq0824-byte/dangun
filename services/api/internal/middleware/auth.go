package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/booking-network/bn-near/services/api/internal/models"
	"github.com/booking-network/bn-near/services/api/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	JwtSecretKey = []byte("bn-near-super-secure-production-master-secret-key-32b")
)

type contextKey string

const (
	UserContextKey      contextKey = "auth_user"
	RequestIDContextKey contextKey = "request_id"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func GenerateTokens(user *models.User) (accessToken, refreshToken string, expiresAt time.Time, err error) {
	expiresAt = time.Now().Add(24 * time.Hour)
	isAdmin := user.Email == "admin@bn-near.com" || strings.HasPrefix(user.Nickname, "NearAdmin")

	accessClaims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bn-near",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = token.SignedString(JwtSecretKey)
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshToken = uuid.NewString()
	return accessToken, refreshToken, expiresAt, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func GetUserFromContext(ctx context.Context) *models.User {
	val := ctx.Value(UserContextKey)
	if val == nil {
		return nil
	}
	if u, ok := val.(*models.User); ok {
		return u
	}
	return nil
}

func RequireAuth(store storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenStr := ""

			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				tokenStr = r.URL.Query().Get("token")
			}

			if tokenStr == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]string{
						"code":    "UNAUTHORIZED",
						"message": "Authorization header or token parameter required",
					},
				})
				return
			}

			claims, err := ValidateToken(tokenStr)
			if err != nil {
				respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]string{
						"code":    "INVALID_TOKEN",
						"message": "Session token invalid or expired",
					},
				})
				return
			}

			user, err := store.GetUserByID(claims.UserID)
			if err != nil {
				respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]string{
						"code":    "USER_NOT_FOUND",
						"message": "User not found or deleted",
					},
				})
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalAuth(store storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenStr := ""
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				tokenStr = r.URL.Query().Get("token")
			}

			if tokenStr != "" {
				if claims, err := ValidateToken(tokenStr); err == nil {
					if user, err := store.GetUserByID(claims.UserID); err == nil {
						r = r.WithContext(context.WithValue(r.Context(), UserContextKey, user))
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAdmin(store storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"error": map[string]string{"code": "UNAUTHORIZED", "message": "Authentication required"},
				})
				return
			}

			// Admin check
			if user.Email != "admin@bn-near.com" && !strings.HasPrefix(user.Nickname, "NearAdmin") {
				respondJSON(w, http.StatusForbidden, map[string]interface{}{
					"error": map[string]string{"code": "FORBIDDEN", "message": "Admin privileges required"},
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Request-ID")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), RequestIDContextKey, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
					"error": map[string]string{
						"code":    "INTERNAL_SERVER_ERROR",
						"message": "An unexpected server error occurred",
					},
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
