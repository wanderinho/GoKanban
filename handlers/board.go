package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"kanban/src"
	"net/http"
)

type Handler struct {
	httpStorage src.StorageInterface
}

func NewHandler(storage src.StorageInterface) *Handler {
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

	board, err := h.httpStorage.NewBoard(r.Context(), req.Title)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
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

	var req dto.CreateColumnDTO
	var res dto.ColumnDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}

	column, err := h.httpStorage.AddColumn(r.Context(), boardID, req.Title)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
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

	var res *dto.BoardDTO
	board, err := h.httpStorage.GetBoard(r.Context(), boardID)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}
	
	res, err = boardToDTO(r.Context(), h, board)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}
	writeJSON(w, http.StatusOK, res)

}

func (h *Handler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	boardID, err := parseIntPath(r, "boardID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	if err := h.httpStorage.RemoveBoard(r.Context(), boardID); err != nil {
		writeError(w, errorStatus(err), err.Error())
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
	var req dto.UpdateBoardTitleDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}
	board, err := h.httpStorage.UpdateBoardTitle(r.Context(), boardID, req.Title)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}
	boardDTO, err := boardToDTO(r.Context(), h, board)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}
	writeJSON(w, http.StatusOK, boardDTO)
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
	column, err := h.httpStorage.GetColumn(r.Context(), boardID, columnID)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}
	res, err := columnToDTO(r.Context(), h, column)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}
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

	if err := h.httpStorage.RemoveColumn(r.Context(), boardID, columnID); err != nil {
		writeError(w, errorStatus(err), err.Error())
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
	var req dto.MoveTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}
	if err := h.httpStorage.MoveTask(r.Context(), boardID, req.TaskID, req.FromColumnID, req.ToColumnID); err != nil {
		writeError(w, errorStatus(err), err.Error())
		return
	}
	necessaryBoard, err := h.httpStorage.GetBoard(r.Context(), boardID)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}
	boardDTO, err := boardToDTO(r.Context(), h, necessaryBoard)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}
	res := dto.MoveTaskResponseDTO{
		BoardID: boardDTO.ID,
		Columns: boardDTO.Columns,
	}
	writeJSON(w, http.StatusOK, res)
}
