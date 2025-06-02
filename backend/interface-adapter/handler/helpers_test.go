// interface-adapter/handler/helpers_test.go
package handler

import (
	"os"
	"testing"
	"time"

	"github.com/ariangn/todo-fullstack/backend/infrastructure/auth"
	"github.com/ariangn/todo-fullstack/backend/interface-adapter/middleware"
	"github.com/go-chi/chi/v5"
)

// ptrString returns a pointer to the given string.
func ptrString(s string) *string {
	return &s
}

// ptrTime returns a pointer to the given time.Time.
func ptrTime(t time.Time) *time.Time {
	return &t
}

// generateAuthToken sets JWT_SECRET and uses auth.NewAuthClient() to produce a valid JWT.
func generateAuthToken(t *testing.T, secret, userID string) string {
	t.Helper() // mark as helper so failures point to the test, not this function

	// Ensure the code under test sees JWT_SECRET
	os.Setenv("JWT_SECRET", secret)

	authClient := auth.NewAuthClient()
	token, err := authClient.GenerateToken(userID, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return token
}

// mountProtectedRouter installs the real AuthMiddleware(authClient) on the given router.
func mountProtectedRouter(r *chi.Mux, secret string) {
	// Ensure JWT_SECRET is set
	os.Setenv("JWT_SECRET", secret)

	authClient := auth.NewAuthClient()
	r.Use(middleware.AuthMiddleware(authClient))
}
