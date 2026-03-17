package storage

import (
	"errors"
	"practice/internal/models"
)

var (
	tasks      = make(map[int]models.Task)
	nextID int = 1
)

func GetAllTasks() []models.Task {
	result := []models.Task{}
	for _, task := range tasks {
		result = append(result, task)
	}
	return result
}

func CreateTask(title string) models.Task {
	task := models.Task{ID: nextID, Title: title, Completed: false}
	tasks[nextID] = task
	nextID++

	return task
}

func GetTaskByID(id int) (models.Task, error) {
	task, exist := tasks[id]

	if !exist {
		return models.Task{}, errors.New("task not found")
	}

	return task, nil
}

func DeleteTask(id int) error {
	_, exist := tasks[id]
	if !exist {
		return errors.New("task not found")
	}

	delete(tasks, id)
	return nil
}

func UpdateTask(id int, title string, completed bool) (models.Task, error) {
	_, exist := tasks[id]
	if !exist {
		return models.Task{}, errors.New("task not found")
	}

	tasks[id] = models.Task{
		ID:        id,
		Title:     title,
		Completed: completed,
	}

	return tasks[id], nil
}
