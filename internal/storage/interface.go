package storage

import (
	"practice/internal/models"
	"time"
)

type TaskStorage interface {
	CreateTask(title string, dueDate *time.Time, isFavorite bool, color *string, repeatDays []int) (models.Task, error)
	GetAllTasks() ([]models.Task, error)
	GetTaskByID(id int) (models.Task, error)
	UpdateTask(id int, title string, completed bool, dueDate *time.Time, isFavorite bool, color *string, isArchived bool, repeatDays []int) (models.Task, error)
	DeleteTask(id int) error
}
