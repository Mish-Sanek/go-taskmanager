package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"practice/internal/models"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{db}, nil
}

func (s *PostgresStore) CreateTask(title string) (models.Task, error) {

	var task models.Task

	query := `
    INSERT INTO tasks (title) VALUES ($1)
    RETURNING id, title, completed
  `

	err := s.db.QueryRow(query, title).Scan(&task.ID, &task.Title, &task.Completed)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (s *PostgresStore) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	rows, err := s.db.Query("SELECT id, title, completed FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("error of db request: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Completed)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return tasks, nil
}

func (s *PostgresStore) GetTaskByID(id int) (models.Task, error) {
	var task models.Task

	err := s.db.QueryRow("SELECT id, title, completed FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.Title, &task.Completed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Task{}, fmt.Errorf("task not found")
		}
		return models.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (s *PostgresStore) DeleteTask(id int) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to ger rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (s *PostgresStore) UpdateTask(id int, title string, completed bool) (models.Task, error) {

	result, err := s.db.Exec(`
    UPDATE tasks
    SET title = $1, completed = $2
    WHERE id = $3
  `, title, completed, id)

	if err != nil {
		return models.Task{}, fmt.Errorf("error of execute: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to ger rows affected: %w", err)
	}

	if rows == 0 {
		return models.Task{}, fmt.Errorf("task was not even exist")
	}

	var task models.Task

	err = s.db.QueryRow(`
    SELECT id, title, completed FROM tasks
    WHERE id = $1
  `, id).Scan(&task.ID, &task.Title, &task.Completed)

	if err != nil {
		return models.Task{}, fmt.Errorf("error of getting updated task from db: %w", err)
	}

	return task, nil
}
