package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"practice/internal/storage"
	"time"
)

type createTaskRequest struct {
	Title      string     `json:"title"`
	DueDate    *time.Time `json:"due_date,omitempty"`
	IsFavorite bool       `json:"is_favorite"`
	Color      *string    `json:"color"`
	RepeatDays []int      `json:"repeat_days"`
}

type updateTaskRequest struct {
	Title      *string    `json:"title"`
	Completed  *bool      `json:"completed"`
	DueDate    *time.Time `json:"due_date"`
	IsFavorite *bool      `json:"is_favorite"`
	Color      *string    `json:"color"`
	IsArchived *bool      `json:"is_archived"`
	RepeatDays *[]int     `json:"repeat_days"`
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

	task, err := h.store.CreateTask(req.Title, req.DueDate, req.IsFavorite, req.Color, req.RepeatDays)
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

	dueDate := oldTask.DueDate
	if data.DueDate != nil {
		dueDate = data.DueDate
	}

	isFavorite := oldTask.IsFavorite
	if data.IsFavorite != nil {
		isFavorite = *data.IsFavorite
	}

	color := oldTask.Color
	if data.Color != nil {
		color = data.Color
	}

	isArchived := oldTask.IsArchived
	if data.IsArchived != nil {
		isArchived = *data.IsArchived
	}

	repeatDays := oldTask.RepeatDays
	if data.RepeatDays != nil {
		repeatDays = *data.RepeatDays
	}

	updatedTask, err := h.store.UpdateTask(id, title, completed, dueDate, isFavorite, color, isArchived, repeatDays)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, updatedTask)
}
