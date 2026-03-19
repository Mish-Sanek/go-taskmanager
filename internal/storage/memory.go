package storage

import (
	"errors"
	"practice/internal/models"
)

type MemoryStore struct {
	tasks  map[int]models.Task
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (m *MemoryStore) GetAllTasks() ([]models.Task, error) {
	result := []models.Task{}
	for _, task := range m.tasks {
		result = append(result, task)
	}
	return result, nil
}

func (m *MemoryStore) CreateTask(title string) (models.Task, error) {
	task := models.Task{ID: m.nextID, Title: title, Completed: false}
	m.tasks[m.nextID] = task
	m.nextID++

	return task, nil
}

func (m *MemoryStore) GetTaskByID(id int) (models.Task, error) {
	task, exist := m.tasks[id]

	if !exist {
		return models.Task{}, errors.New("task not found")
	}

	return task, nil
}

func (m *MemoryStore) DeleteTask(id int) error {
	_, exist := m.tasks[id]
	if !exist {
		return errors.New("task not found")
	}

	delete(m.tasks, id)
	return nil
}

func (m *MemoryStore) UpdateTask(id int, title string, completed bool) (models.Task, error) {
	_, exist := m.tasks[id]
	if !exist {
		return models.Task{}, errors.New("task not found")
	}

	m.tasks[id] = models.Task{
		ID:        id,
		Title:     title,
		Completed: completed,
	}

	return m.tasks[id], nil
}
