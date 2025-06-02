package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ariangn/todo-fullstack/backend/infrastructure/auth"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func AuthMiddleware(authClient auth.AuthClientInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]
			claims, err := authClient.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			sub, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "invalid token subject", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extracts the userID from context if set
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
