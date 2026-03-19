package storage

import "practice/internal/models"

type TaskStorage interface {
	CreateTask(title string) (models.Task, error)
	GetAllTasks() ([]models.Task, error)
	GetTaskByID(id int) (models.Task, error)
	DeleteTask(id int) error
	UpdateTask(id int, title string, completed bool) (models.Task, error)
}
