// cmd/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/ariangn/todo-fullstack/backend/di"
	custommw "github.com/ariangn/todo-fullstack/backend/interface-adapter/middleware"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env file: %v", err)
	}
	// Initialize DI container
	container, err := di.InitializeContainer()
	if err != nil {
		log.Fatalf("failed to initialize container: %v", err)
	}

	// Set up router
	r := chi.NewRouter()
	clientOrigin := os.Getenv("CLIENT_ORIGIN")
	if clientOrigin == "" {
		log.Fatal("CLIENT_ORIGIN must be set")
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{clientOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	}))

	// Single /api route group
	r.Route("/api", func(r chi.Router) {
		// 1) Public endpoints:
		r.Post("/users/register", container.UserController.Register)
		r.Post("/users/login", container.UserController.Login)

		// 2) Protected endpoints (require JWT)
		r.Group(func(r chi.Router) {
			r.Use(custommw.AuthMiddleware(container.AuthClient))
			r.Get("/auth/me", container.UserController.Me)
			r.Post("/users/logout", container.UserController.Logout)

			// /api/todos/*
			r.Route("/todos", func(r chi.Router) {
				r.Post("/", container.TodoController.Create)
				r.Get("/", container.TodoController.List)
				r.Get("/{id}", container.TodoController.GetByID)
				r.Put("/{id}", container.TodoController.Update)
				r.Put("/{id}/status", container.TodoController.ToggleStatus)
				r.Delete("/{id}", container.TodoController.Delete)
				r.Post("/{id}/duplicate", container.TodoController.Duplicate)
			})

			// /api/categories/*
			r.Route("/categories", func(r chi.Router) {
				r.Post("/", container.CategoryController.Create)
				r.Get("/", container.CategoryController.List)
				r.Put("/{id}", container.CategoryController.Update)
				r.Delete("/{id}", container.CategoryController.Delete)
			})

			// /api/tags/*
			r.Route("/tags", func(r chi.Router) {
				r.Post("/", container.TagController.Create)
				r.Get("/", container.TagController.List)
				r.Put("/{id}", container.TagController.Update)
				r.Delete("/{id}", container.TagController.Delete)
			})
		})
	})

	// HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Server listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt to gracefully shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
