package dto

import (
	"kanban/src"
)

type CreateBoardDTO struct {
	Title string `json:"title"`
}

type CreateColumnDTO struct {
	Title string `json:"title"`
}

type BoardDTO struct {
	ID int
	Title string
	Columns []ColumnDTO
}

type ColumnDTO struct {
	ID int
	Title string
	Tasks []src.Task
	BoardID int
}

type UpdateBoardTitleDTO struct {
	Title string `json:"title"`
}

type UpdateTaskTitleDTO struct {
	Title string `json:"title"`
}

type MoveTaskDTO struct {
	TaskID       int `json:"taskID"`
	FromColumnID int `json:"from_column_id"`
	ToColumnID   int `json:"to_column_id"`
}

type MoveTaskResponseDTO struct {
	BoardID int          `json:"board_id"`
	Columns []ColumnDTO `json:"columns"`
}

type CreateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateColumnTitleDTO struct {
	Title string `json:"title"`
}

type UpdateTaskDescriptionDTO struct {
	Description string `json:"description"`
}

type ErrorDTO struct {
	Message string `json:"message"`
}
