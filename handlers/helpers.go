package handlers

import (
	"net/http"
	"kanban/handlers/dto"
	"kanban/src"
	"encoding/json"
	"strconv"
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

func (h *Handler) boardExists(boardID int) error {
    if _, ok := h.httpStorage.Boards[boardID]; !ok {
        return errors.New("такой доски не существует")
    }
    return nil
}

func tasksMapToSlice(handler *Handler, tasksMap map[int]struct{}) []src.Task {
    tasksSlice := make([]src.Task, 0, len(tasksMap))
    for taskID := range tasksMap {
        if task, ok := handler.httpStorage.Tasks[taskID]; ok {
            tasksSlice = append(tasksSlice, *task)
        }
    }
    return tasksSlice
}

func columnsMapToDTO(handler *Handler, board *src.Board) []dto.ColumnDTO {
    columnsDTO := make([]dto.ColumnDTO, 0, len(board.Columns))
    for colID := range board.Columns {
        col := handler.httpStorage.Columns[colID]
        tasksSlice := tasksMapToSlice(handler, col.Tasks)
        colDTO := dto.ColumnDTO{
            ID:      col.ID,
            Title:   col.Title,
            Tasks:   tasksSlice,
            BoardID: col.BoardID,
        }
        columnsDTO = append(columnsDTO, colDTO)
    }
    return columnsDTO
}