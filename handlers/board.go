package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"kanban/src"
	"net/http"
	"strconv"
)

func writeError(w http.ResponseWriter, status int, message string) {
	var err dto.ErrorDTO
	err.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func CreateBoard(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBoardDTO
	var res dto.CreateBoardResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}

	board, err := src.NewBoard(req.Title)
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

func CreateColumn(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := src.BoardMap[boardID]
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

	column, err := necessaryBoard.AddColumn(req.Title)
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

func GetBoard(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	var res dto.GetBoardResponseDTO
	board, err := src.GetBoard(boardID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	res.ID = board.ID
	res.Title = board.Title
	for _, v := range board.Columns {
		res.Columns = append(res.Columns, *v)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func DeleteBoard(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := src.BoardMap[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	if err := necessaryBoard.RemoveBoard(); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateBoardTitle(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := src.BoardMap[boardID]
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

	board, err := necessaryBoard.UpdateBoardTitle(req.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
	}

	res.ID = board.ID
	res.Title = board.Title
	for _, v := range board.Columns {
		res.Columns = append(res.Columns, *v)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func GetColumn(w http.ResponseWriter, r *http.Request) {
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

	necessaryBoard, ok := src.BoardMap[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	var res dto.GetColumnResponseDTO
	column, err := necessaryBoard.GetColumn(columnID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	res.ID = column.ID
	res.Title = column.Title
	for _, v := range column.Tasks {
		res.Tasks = append(res.Tasks, *v)
	}
	res.BoardID = column.BoardID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func DeleteColumn(w http.ResponseWriter, r *http.Request) {
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

	necessaryBoard, ok := src.BoardMap[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	if err := necessaryBoard.RemoveColumn(columnID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func MoveTask(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("boardID")

	boardID, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := src.BoardMap[boardID]
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

	if err := necessaryBoard.MoveTask(req.TaskID, req.FromColumnID, req.ToColumnID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res.BoardID = necessaryBoard.ID
	for _, v := range necessaryBoard.Columns {
		res.Columns = append(res.Columns, *v)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
