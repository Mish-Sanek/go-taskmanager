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

type TaskHandler struct {
	store storage.TaskStorage
}

func NewTaskHandler(store storage.TaskStorage) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
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

	task, err := h.store.CreateTask(req.Title)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create task")
		return
	}
	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAllTasks()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get tasks")
		return
	}
	respondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.store.GetTaskByID(id)

	if err != nil {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.store.DeleteTask(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "task not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
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

	oldTask, err := h.store.GetTaskByID(id)
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

	updatedTask, err := h.store.UpdateTask(id, title, completed)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, updatedTask)
}
