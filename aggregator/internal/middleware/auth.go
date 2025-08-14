package middleware

import (
	"aggregator/internal/service/auth"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"
const roleKey contextKey = "role"

type AuthMiddleware struct {
	svc auth.AuthService
}

func NewAuthMiddleware(svc auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		svc: svc,
	}
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		resp, err := m.svc.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		parent := context.WithValue(r.Context(), userIDKey, resp.UserID)
		ctx := context.WithValue(parent, roleKey, resp.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey).(string)
	return role, ok
}
