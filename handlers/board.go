package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"kanban/src"
	"net/http"
)

type Handler struct {
	httpStorage *src.Storage
}

func NewHandler(storage *src.Storage) *Handler {
	return &Handler{
		httpStorage: storage,
	}
}

func (h *Handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBoardDTO
	var res dto.BoardDTO
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
	writeJSON(w, http.StatusCreated, res)
}

func (h *Handler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	boardID, err := parseIntPath(r, "boardID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	if err := h.boardExists(boardID); err != nil {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	var req dto.CreateColumnDTO
	var res dto.ColumnDTO
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
	writeJSON(w, http.StatusCreated, res)
}

func (h *Handler) GetBoard(w http.ResponseWriter, r *http.Request) {
	boardID, err := parseIntPath(r, "boardID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	var res dto.BoardDTO
	board, err := h.httpStorage.GetBoard(boardID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	
	res.ID = board.ID
	res.Title = board.Title
	res.Columns = columnsMapToDTO(h, board)
	writeJSON(w, http.StatusOK, res)

}

func (h *Handler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	boardID, err := parseIntPath(r, "boardID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	if err := h.boardExists(boardID); err != nil {
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
	boardID, err := parseIntPath(r, "boardID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	if err := h.boardExists(boardID); err != nil {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	var req dto.UpdateBoardTitleDTO
	var res dto.BoardDTO
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
	res.Columns = columnsMapToDTO(h, board)

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) GetColumn(w http.ResponseWriter, r *http.Request) {
	boardID, err := parseIntPath(r, "boardID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	columnID, err := parseIntPath(r, "columnID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	if err := h.boardExists(boardID); err != nil {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	var res dto.ColumnDTO
	column, err := h.httpStorage.GetColumn(boardID, columnID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	res.ID = column.ID
	res.Title = column.Title
	res.Tasks = tasksMapToSlice(h, column.Tasks)
	res.BoardID = column.BoardID
	
	writeJSON(w, http.StatusOK, res)

}

func (h *Handler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	boardID, err := parseIntPath(r, "boardID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	columnID, err := parseIntPath(r, "columnID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	if err := h.boardExists(boardID); err != nil{
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
	boardID, err := parseIntPath(r, "boardID")
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
	res.Columns = columnsMapToDTO(h, necessaryBoard)
	
	writeJSON(w, http.StatusOK, res)

}
