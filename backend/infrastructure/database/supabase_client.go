package database

import (
	"log"
	"os"

	auth "github.com/supabase-community/auth-go"
	postgrest "github.com/supabase-community/postgrest-go"
	storage "github.com/supabase-community/storage-go"
)

// SupabaseClient holds all Supabase‐related subclients.
type SupabaseClient struct {
	Auth    auth.Client // auth.New returns auth.Client (an interface)
	DB      *postgrest.Client
	Storage *storage.Client
}

// NewSupabaseClient initializes the Supabase Auth, PostgREST, and Storage clients.
// It reads SUPABASE_URL and SUPABASE_KEY from the environment.
func NewSupabaseClient() (*SupabaseClient, error) {
	baseURL := os.Getenv("SUPABASE_URL")
	apiKey := os.Getenv("SUPABASE_KEY")
	if baseURL == "" || apiKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	// 1) Initialize Auth client
	// auth.New takes (projectReference, apiKey) and returns auth.Client (an interface).
	projectRef := baseURL
	if len(projectRef) >= 8 && projectRef[:8] == "https://" {
		projectRef = projectRef[8:]
	}
	clientAuth := auth.New(projectRef, apiKey)

	// 2) Initialize PostgREST client
	restURL := baseURL + "/rest/v1"
	headers := map[string]string{
		"apikey":        apiKey,
		"Authorization": "Bearer " + apiKey,
	}
	dbClient := postgrest.NewClient(restURL, "public", headers)

	// 3) Initialize Storage client
	storageURL := baseURL + "/storage/v1"
	storageClient := storage.NewClient(storageURL, apiKey, headers)

	return &SupabaseClient{
		Auth:    clientAuth,
		DB:      dbClient,
		Storage: storageClient,
	}, nil
}
