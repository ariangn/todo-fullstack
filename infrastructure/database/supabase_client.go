package database

import (
    "log"
    "os"

    postgrest "github.com/supabase-community/postgrest-go"
    auth "github.com/supabase-community/auth-go"
    storage "github.com/supabase-community/storage-go"
)

type SupabaseClient struct {
    Auth    *auth.Client
    DB      *postgrest.Client
    Storage *storage.Client
}

func NewSupabaseClient() (*SupabaseClient, error) {
    url := os.Getenv("SUPABASE_URL")
    key := os.Getenv("SUPABASE_KEY")
    if url == "" || key == "" {
        log.Fatal("SUPABASE_URL and SUPABASE_KEY must be set in environment")
    }

    // Initialize Supabase Auth 
    authClient := auth.NewClient(url+"/auth/v1", auth.WithAPIKey(key))

    // Initialize PostgREST (database) client
    dbClient := postgrest.New(url+"/rest/v1", postgrest.WithAPIKey(key))

    // Initialize Supabase Storage client
    storageClient := storage.NewClient(url, key)

    return &SupabaseClient{
        Auth:    authClient,
        DB:      dbClient,
        Storage: storageClient,
    }, nil
}
