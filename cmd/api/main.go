// go run cmd/api/main.go

// postgres://app:secret@localhost:5432/tasks?sslmode=disable
//

package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"practice/internal/config"
	"practice/internal/handlers"
	"practice/internal/logger"
	"practice/internal/storage"
	"time"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.Env)
	mux := http.NewServeMux()
	var store storage.TaskStorage

	if cfg.Env == "production" {
		postgresStore, err := storage.NewPostgresStore(cfg.DatabaseURL)
		if err != nil {
			slog.Error("failed to create postgres store", "error", err)
		}
		store = postgresStore
	} else {
		store = storage.NewMemoryStore()
	}

	taskHandler := handlers.NewTaskHandler(store)

	mux.HandleFunc("GET /tasks", taskHandler.GetAllTasks)
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		slog.Info("server is opening:", "port", cfg.Port)

		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down server....")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}
