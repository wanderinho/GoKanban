package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"kanban/src"
	"net/http"
	"strconv"
)

type Handler struct {
	httpStorage *src.Storage
}

func NewHandler(storage *src.Storage) *Handler {
	return &Handler{
		httpStorage: storage,
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	var err dto.ErrorDTO
	err.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBoardDTO
	var res dto.CreateBoardResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}

	board, err := h.httpStorage.NewBoard(req.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res.ID = board.ID
	res.Title = board.Title
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	_, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	var req dto.CreateColumnDTO
	var res dto.CreateColumnResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}

	column, err := h.httpStorage.AddColumn(boardID, req.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res.ID = column.ID
	res.Title = column.Title
	res.BoardID = column.BoardID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetBoard(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	var res dto.GetBoardResponseDTO
	board, err := h.httpStorage.GetBoard(boardID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	res.ID = board.ID
	res.Title = board.Title
	for k, _ := range board.Columns {
		res.Columns = append(res.Columns, *h.httpStorage.Columns[k])
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	_, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	if err := h.httpStorage.RemoveBoard(boardID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateBoardTitle(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	_, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	var req dto.UpdateBoardTitleDTO
	var res dto.UpdateBoardTitleResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}

	board, err := h.httpStorage.UpdateBoardTitle(boardID, req.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res.ID = board.ID
	res.Title = board.Title
	for k, _ := range board.Columns {
		res.Columns = append(res.Columns, *h.httpStorage.Columns[k])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) GetColumn(w http.ResponseWriter, r *http.Request) {
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

	_, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	var res dto.GetColumnResponseDTO
	column, err := h.httpStorage.GetColumn(boardID, columnID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	res.ID = column.ID
	res.Title = column.Title
	for k, _ := range column.Tasks {
		res.Tasks = append(res.Tasks, *h.httpStorage.Tasks[k])
	}
	res.BoardID = column.BoardID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
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

	_, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	if err := h.httpStorage.RemoveColumn(boardID, columnID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MoveTask(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	var req dto.MoveTaskDTO
	var res dto.MoveTaskResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}

	if err := h.httpStorage.MoveTask(boardID, req.TaskID, req.FromColumnID, req.ToColumnID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res.BoardID = necessaryBoard.ID
	for k, _ := range necessaryBoard.Columns {
		res.Columns = append(res.Columns, *h.httpStorage.Columns[k])
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
