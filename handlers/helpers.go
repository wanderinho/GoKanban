package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"kanban/src"
	"net/http"
	"strconv"
	"context"
	"errors"
)

func writeError(w http.ResponseWriter, status int, message string) {
	errDTO := dto.ErrorDTO{Message: message}
	writeJSON(w, status, errDTO)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func parseIntPath(r *http.Request, path string) (int, error) {
	val := r.PathValue(path)
	return strconv.Atoi(val)
}

func columnToDTO(ctx context.Context, handler *Handler, column *src.Column) (*dto.ColumnDTO, error) {
	tasks, err := handler.httpStorage.GetTasksByColumn(ctx, column.ID)
	if err != nil {
		return nil, err
	}
	tasksSlice := make([]src.Task, 0, len(tasks))
	for _, t := range tasks {
		tasksSlice = append(tasksSlice, *t)
	}
	return &dto.ColumnDTO{
		ID:      column.ID,
		Title:   column.Title,
		Tasks:   tasksSlice,
		BoardID: column.BoardID,
	}, nil
}

func columnsToDTO(ctx context.Context, handler *Handler, columns []*src.Column) ([]dto.ColumnDTO, error) {
	columnsDTO := make([]dto.ColumnDTO, 0, len(columns))
	for _, col := range columns {
		colDTO, err := columnToDTO(ctx, handler, col)
		if err != nil {
			return nil, err
		}
		columnsDTO = append(columnsDTO, *colDTO)
	}
	return columnsDTO, nil
}

func boardToDTO(ctx context.Context, handler *Handler, board *src.Board) (*dto.BoardDTO, error) {
	columns, err := handler.httpStorage.GetColumnsByBoard(ctx, board.ID)
	if err != nil {
		return nil, err
	}
	columnsDTO, err := columnsToDTO(ctx, handler, columns)
	if err != nil {
		return nil, err
	}
	return &dto.BoardDTO{
		ID:      board.ID,
		Title:   board.Title,
		Columns: columnsDTO,
	}, nil
}

func errorStatus(err error) int {
	switch {
	case errors.Is(err, src.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, src.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, src.ErrInvalidInput):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}