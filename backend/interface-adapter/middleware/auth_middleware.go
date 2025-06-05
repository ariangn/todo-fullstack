package middleware

import (
	"context"
	"net/http"

	"github.com/ariangn/todo-fullstack/backend/infrastructure/auth"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// AuthMiddleware validates the token from HTTP-only cookie and stores the user ID in context.
func AuthMiddleware(authClient auth.AuthClientInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read the "token" cookie
			cookie, err := r.Cookie("token")
			if err != nil || cookie.Value == "" {
				http.Error(w, "missing or invalid token", http.StatusUnauthorized)
				return
			}

			// Validate the JWT and extract claims
			claims, err := authClient.ValidateToken(cookie.Value)
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
