package dto

import (
	"kanban/src"
)

type CreateBoardDTO struct {
	Title string `json:"title"`
}

type CreateBoardResponseDTO struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CreateColumnDTO struct {
	Title string `json:"title"`
}

type CreateColumnResponseDTO struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	BoardID int    `json:"board_id"`
}

type GetBoardResponseDTO struct {
	ID      int          `json:"id"`
	Title   string       `json:"title"`
	Columns []src.Column `json:"columns"`
}

type UpdateBoardTitleDTO struct {
	Title string `json:"title"`
}

type UpdateTaskTitleDTO struct {
	Title string `json:"title"`
}

type UpdateBoardTitleResponseDTO struct {
	ID      int          `json:"id"`
	Title   string       `json:"title"`
	Columns []src.Column `json:"columns"`
}

type GetColumnResponseDTO struct {
	ID      int        `json:"id"`
	Title   string     `json:"title"`
	Tasks   []src.Task `json:"tasks"`
	BoardID int        `json:"board_id"`
}

type MoveTaskDTO struct {
	TaskID       int `json:"taskID"`
	FromColumnID int `json:"from_column_id"`
	ToColumnID   int `json:"to_column_id"`
}

type MoveTaskResponseDTO struct {
	BoardID int          `json:"board_id"`
	Columns []src.Column `json:"columns"`
}

type CreateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateColumnTitleDTO struct {
	Title string `json:"title"`
}

type UpdateColumnTitleResponseDTO struct {
	ID      int        `json:"id"`
	Title   string     `json:"title"`
	Tasks   []src.Task `json:"tasks"`
	BoardID int        `json:"board_id"`
}

type UpdateTaskDescriptionDTO struct {
	Description string `json:"description"`
}

type ErrorDTO struct {
	Message string `json:"message"`
}
