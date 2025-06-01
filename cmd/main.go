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

	"github.com/ariangn/todo-go/di"
	custommw "github.com/ariangn/todo-go/interface-adapter/middleware"
)

func main() {
	// initialize DI container via Wire
	container, err := di.InitializeContainer()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// set up a Chi router with middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// public (unauthenticated) routes: user register & login
	r.Route("/api", func(r chi.Router) {
		r.Post("/users/register", container.UserController.Register)
		r.Post("/users/login", container.UserController.Login)
	})

	// protected routes (require JWT) under /api
	r.Route("/api", func(r chi.Router) {
		// apply JWT‐based AuthMiddleware (injects userID into context)
		r.Use(custommw.AuthMiddleware(container.AuthClient))

		// To-Do endpoints
		r.Route("/todos", func(r chi.Router) {
			r.Post("/", container.TodoController.Create)
			r.Get("/", container.TodoController.List)
			r.Get("/{id}", container.TodoController.GetByID)
			r.Put("/{id}", container.TodoController.Update)
			r.Put("/{id}/status", container.TodoController.ToggleStatus)
			r.Delete("/{id}", container.TodoController.Delete)
			r.Post("/{id}/duplicate", container.TodoController.Duplicate)
		})

		// Category endpoints
		r.Route("/categories", func(r chi.Router) {
			r.Post("/", container.CategoryController.Create)
			r.Get("/", container.CategoryController.List)
			r.Put("/{id}", container.CategoryController.Update)
			r.Delete("/{id}", container.CategoryController.Delete)
		})

		// Tag endpoints
		r.Route("/tags", func(r chi.Router) {
			r.Post("/", container.TagController.Create)
			r.Get("/", container.TagController.List)
			r.Put("/{id}", container.TagController.Update)
			r.Delete("/{id}", container.TagController.Delete)
		})
	})

	// start HTTP Server with graceful shutdown
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
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("⚠️ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped.")
}
