package models

import "time"

type Task struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	Completed  bool       `json:"completed"`
	DueDate    *time.Time `json:"due_date"`
	IsFavorite bool       `json:"is_favorite"`
	Color      *string    `json:"color"`
	IsArchived bool       `json:"is_archived"`
	RepeatDays []int      `json:"repeat_days"`
}
