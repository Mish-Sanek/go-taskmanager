package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"practice/internal/storage"
)

type createTaskRequest struct {
	Title string `json:"title"`
}

type updateTaskRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("error decoding json: %s", err)
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}

	task := storage.CreateTask(req.Title)
	respondJSON(w, http.StatusOK, task)
}

func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := storage.GetAllTasks()
	respondJSON(w, http.StatusOK, tasks)
}

func GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := storage.GetTaskByID(id)

	if err != nil {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = storage.DeleteTask(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var data updateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	oldTask, err := storage.GetTaskByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}

	title := oldTask.Title
	if data.Title != nil {
		title = *data.Title
	}

	completed := oldTask.Completed
	if data.Completed != nil {
		completed = *data.Completed
	}

	updatedTask, err := storage.UpdateTask(id, title, completed)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, updatedTask)
}
