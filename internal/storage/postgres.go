package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"practice/internal/models"
	"time"
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

func (s *PostgresStore) CreateTask(title string, dueDate *time.Time, isFavorite bool, color *string, repeatDays []int) (models.Task, error) {

	var task models.Task

	query := `
    INSERT INTO tasks (title, due_date, is_favorite, color, repeat_days) VALUES ($1, $2, $3, $4, $5)
    RETURNING id, title, completed, due_date, is_favorite, color, is_archived, repeat_days
  `

	err := s.db.QueryRow(query, title, dueDate, isFavorite, color, repeatDays).Scan(&task.ID, &task.Title, &task.Completed, &task.DueDate, &task.IsFavorite, &task.Color, &task.IsArchived, &task.RepeatDays)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (s *PostgresStore) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	rows, err := s.db.Query("SELECT id, title, completed, due_date, is_favorite, color, is_archived, repeat_days FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("error of db request: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.DueDate, &task.IsFavorite, &task.Color, &task.IsArchived, &task.RepeatDays)
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

	err := s.db.QueryRow("SELECT id, title, completed, due_date, is_favorite, color, is_archived, repeat_days FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.Title, &task.Completed, &task.DueDate, &task.IsFavorite, &task.Color, &task.IsArchived, &task.RepeatDays)
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

func (s *PostgresStore) UpdateTask(id int, title string, completed bool, dueDate *time.Time, isFavorite bool, color *string, isArchived bool, repeatDays []int) (models.Task, error) {

	result, err := s.db.Exec(`
    UPDATE tasks
    SET title = $1, completed = $2, due_date = $3, is_favorite = $4, color = $5, is_archived = $6, repeat_days = $7
    WHERE id = $8
  `, title, completed, dueDate, isFavorite, color, isArchived, repeatDays, id)

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
    SELECT id, title, completed, due_date, is_favorite, color, is_archived, repeat_days FROM tasks
    WHERE id = $1
  `, id).Scan(&task.ID, &task.Title, &task.Completed, &task.DueDate, &task.IsFavorite, &task.Color, &task.IsArchived, &task.RepeatDays)

	if err != nil {
		return models.Task{}, fmt.Errorf("error of getting updated task from db: %w", err)
	}

	return task, nil
}
