package dto

import "kanban/src"

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
	BoardID int    `json:"boardID"`
}

type GetBoardResponseDTO struct {
	ID      int          `json:"id"`
	Title   string       `json:"title"`
	Columns []src.Column `json:"columns"`
}

type UpdateBoardTitleDTO struct {
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
	BoardID int        `json:"boardID"`
}

type MoveTaskDTO struct {
	TaskID       int `json:"taskID"`
	FromColumnID int `json:"from_column_id"`
	ToColumnID   int `json:"to_column_id"`
}

type MoveTaskResponseDTO struct {
	BoardID int          `json:"boardID"`
	Columns []src.Column `json:"columns"`
}

type ErrorDTO struct {
	Message string `json:"message"`
}
