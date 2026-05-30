package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"net/http"
)

func (h *Handler) UpdateTaskTitle(w http.ResponseWriter, r *http.Request) {
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

	taskID, err := parseIntPath(r, "taskID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}
	
	necessaryBoard, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	if _, ok := necessaryBoard.Columns[columnID]; !ok {
		writeError(w, http.StatusNotFound, "такой колонки не существует в доске")
		return
	}

 	column, ok := h.httpStorage.Columns[columnID]
    if !ok {
        writeError(w, http.StatusNotFound, "такой колонки не существует в глобальном хранилище")
        return
    }

    if _, ok := column.Tasks[taskID]; !ok {
        writeError(w, http.StatusNotFound, "такой задачи не существует в колонке")
        return
    }

    var req dto.UpdateTaskTitleDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
       	writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
        return
	}
	
	if task, err := h.httpStorage.UpdateTaskTitle(boardID, columnID, taskID, req.Title); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	} else {
		writeJSON(w, http.StatusOK, task)
	}
}


func (h *Handler) UpdateTaskDescription(w http.ResponseWriter, r *http.Request) {
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

	taskID, err := parseIntPath(r, "taskID")
	if err != nil {
		writeError(w, http.StatusBadRequest, "параметр пути должен быть числом")
		return
	}

	necessaryBoard, ok := h.httpStorage.Boards[boardID]
	if !ok {
		writeError(w, http.StatusNotFound, "такой доски не существует")
		return
	}

	if _, ok := necessaryBoard.Columns[columnID]; !ok {
		writeError(w, http.StatusNotFound, "такой колонки не существует в доске")
		return
	}

 	column, ok := h.httpStorage.Columns[columnID]
    if !ok {
        writeError(w, http.StatusNotFound, "такой колонки не существует в глобальном хранилище")
        return
    }

    if _, ok := column.Tasks[taskID]; !ok {
        writeError(w, http.StatusNotFound, "такой задачи не существует в колонке")
        return
    }

	var req dto.UpdateTaskDescriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    	writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
     	return
	}
	
	task := h.httpStorage.Tasks[taskID].UpdateTaskDescription(req.Description)
	writeJSON(w, http.StatusOK, task)
}
