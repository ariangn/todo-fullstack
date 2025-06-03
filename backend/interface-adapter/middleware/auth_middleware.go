package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ariangn/todo-fullstack/backend/infrastructure/auth"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// AuthMiddleware validates the Bearer token and stores the user ID in context.
func AuthMiddleware(authClient auth.AuthClientInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Expect header: Authorization: Bearer <token>
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

			// Validate the JWT and extract claims
			claims, err := authClient.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Expect the “sub” claim to be the user’s ID
			sub, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "invalid token subject", http.StatusUnauthorized)
				return
			}

			// Store userID in context for downstream handlers
			ctx := context.WithValue(r.Context(), userIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts the userID (string) from context, if present.
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
