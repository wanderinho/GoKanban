package src

import (
	"time"
)

// структура таски
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ColumnID    int       `json:"column_id"`
}
