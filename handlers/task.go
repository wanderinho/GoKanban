package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"kanban/src"
	"net/http"
	"strconv"
)

func UpdateTaskDescription(w http.ResponseWriter, r *http.Request) {
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

	necessaryTask, check := necessaryColumn.Tasks[taskID]
	if !check {
		writeError(w, http.StatusNotFound, "такой задачи не существует")
		return
	}

	var req dto.UpdateTaskDescriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    	writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
     	return
	}
	
	task := necessaryTask.UpdateTaskDescription(req.Description)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
