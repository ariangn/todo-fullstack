package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/ariangn/todo-fullstack/backend/di"
	custommw "github.com/ariangn/todo-fullstack/backend/interface-adapter/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env file: %v", err)
	}

	// Initialize DI container
	container, err := di.InitializeContainer()
	if err != nil {
		log.Fatalf("failed to initialize container: %v", err)
	}

	clientOrigin := os.Getenv("CLIENT_ORIGIN")
	if clientOrigin == "" {
		log.Fatal("CLIENT_ORIGIN must be set in your .env file")
	}

	// Set up router with common middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // logs every request
	r.Use(middleware.Recoverer) // prevents panics from crashing server
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{clientOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	}))

	// Define /api routes
	r.Route("/api", func(r chi.Router) {
		// Public routes
		r.Post("/users/register", container.UserController.Register)
		r.Post("/users/login", container.UserController.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(custommw.AuthMiddleware(container.AuthClient))

			// Auth info
			r.Get("/auth/me", container.UserController.Me)
			r.Post("/users/logout", container.UserController.Logout)

			// Todos
			r.Route("/todos", func(r chi.Router) {
				r.Post("/", container.TodoController.Create)
				r.Get("/", container.TodoController.List)
				r.Get("/{id}", container.TodoController.GetByID)
				r.Put("/{id}", container.TodoController.Update)
				r.Patch("/{id}/status", container.TodoController.ToggleStatus)
				r.Delete("/{id}", container.TodoController.Delete)
				r.Post("/{id}/duplicate", container.TodoController.Duplicate)
			})

			// Categories
			r.Route("/categories", func(r chi.Router) {
				r.Post("/", container.CategoryController.Create)
				r.Get("/", container.CategoryController.List)
				r.Put("/{id}", container.CategoryController.Update)
				r.Delete("/{id}", container.CategoryController.Delete)
			})

			// Tags
			r.Route("/tags", func(r chi.Router) {
				r.Post("/", container.TagController.Create)
				r.Get("/", container.TagController.List)
				r.Put("/{id}", container.TagController.Update)
				r.Delete("/{id}", container.TagController.Delete)
			})
		})
	})

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// Handle shutdown signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced to shutdown: %v", err)
	}

	log.Println("server exited cleanly")
}
