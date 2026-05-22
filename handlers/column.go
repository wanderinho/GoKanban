package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"kanban/src"
	"net/http"
	"strconv"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
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

	necessaryColumn, flag := necessaryBoard.Columns[columnID]
	if !flag {
		writeError(w, http.StatusNotFound, "такой колонки не существует")
		return
	}
	
	var req dto.CreateTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}

	task, err := necessaryColumn.AddTask(req.Title, req.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
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

	necessaryBoard, ok := src.BoardMap[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	necessaryColumn, flag := necessaryBoard.Columns[columnID]
	if !flag {
		writeError(w, http.StatusNotFound, "такой колонки не существует")
		return
	}

	task, err := necessaryColumn.GetTask(taskID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
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

	necessaryBoard, ok := src.BoardMap[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	necessaryColumn, flag := necessaryBoard.Columns[columnID]
	if !flag {
		writeError(w, http.StatusNotFound, "такой колонки не существует")
		return
	}

	if err := necessaryColumn.RemoveTask(taskID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateColumnTitle(w http.ResponseWriter, r *http.Request) {
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

	necessaryColumn, flag := necessaryBoard.Columns[columnID]
	if !flag {
		writeError(w, http.StatusNotFound, "такой колонки не существует")
		return
	}

	var req dto.UpdateColumnTitleDTO
	var res dto.UpdateColumnTitleResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    	writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
     	return
	}
	
	column, err := necessaryColumn.UpdateColumnTitle(req.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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
