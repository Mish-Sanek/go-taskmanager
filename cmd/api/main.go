// go run cmd/api/main.go

package main

import (
	"log"
	"net/http"
	"practice/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", handlers.GetAllTasksHandler)
	mux.HandleFunc("POST /tasks", handlers.CreateTaskHandler)
	mux.HandleFunc("GET /tasks/{id}", handlers.GetTaskByIDHandler)
	mux.HandleFunc("DELETE /tasks/{id}", handlers.DeleteTaskHandler)
	mux.HandleFunc("PUT /tasks/{id}", handlers.UpdateTaskHandler)

	log.Println("Сервер запущен на :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Ощибка запуска сервера:", err)
	}
}
