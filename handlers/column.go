package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"kanban/src"
	"net/http"
	"strconv"
)

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	boardIdStr := r.PathValue("boardID")
	boardID, err := strconv.Atoi(boardIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	columnIdStr := r.PathValue("columnID")
	columnID, err := strconv.Atoi(columnIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	_, ok = necessaryBoard.Columns[columnID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой колонки не существует")
		return
	}
	
	var req dto.CreateTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}

	task, err := h.httpStorage.AddTask(boardID, columnID, req.Title, req.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	boardIdStr := r.PathValue("boardID")
	boardID, err := strconv.Atoi(boardIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	columnIdStr := r.PathValue("columnID")
	columnID, err := strconv.Atoi(columnIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	taskIdStr := r.PathValue("taskID")
	taskID, err := strconv.Atoi(taskIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	_, ok = necessaryBoard.Columns[columnID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой колонки не существует")
		return
	}

	task, err := h.httpStorage.GetTask(boardID, columnID, taskID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	boardIdStr := r.PathValue("boardID")
	boardID, err := strconv.Atoi(boardIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	columnIdStr := r.PathValue("columnID")
	columnID, err := strconv.Atoi(columnIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	taskIdStr := r.PathValue("taskID")
	taskID, err := strconv.Atoi(taskIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	_, ok = necessaryBoard.Columns[columnID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой колонки не существует")
		return
	}

	if err := h.httpStorage.RemoveTask(boardID, columnID, taskID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateColumnTitle(w http.ResponseWriter, r *http.Request) {
	boardIdStr := r.PathValue("boardID")
	boardID, err := strconv.Atoi(boardIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	columnIdStr := r.PathValue("columnID")
	columnID, err := strconv.Atoi(columnIdStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	_, ok = necessaryBoard.Columns[columnID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой колонки не существует")
		return
	}

	var req dto.UpdateColumnTitleDTO
	var res dto.ColumnDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    	writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
     	return
	}
	
	column, err := h.httpStorage.UpdateColumnTitle(boardID, columnID, req.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res.ID = column.ID
	res.Title = column.Title
	
	taskSlice := make([]src.Task, 0, len(column.Tasks))
	for taskID := range column.Tasks {
		task := h.httpStorage.Tasks[taskID]

		taskSlice = append(taskSlice, *task)	
	}
	res.BoardID = column.BoardID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
