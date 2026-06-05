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

    var req dto.UpdateTaskTitleDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
       	writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
        return
	}
	
	if task, err := h.httpStorage.UpdateTaskTitle(r.Context(), boardID, columnID, taskID, req.Title); err != nil {
		writeError(w, errorStatus(err), err.Error())
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

    var req dto.UpdateTaskDescriptionDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
        return
    }

    task, err := h.httpStorage.UpdateTaskDescription(r.Context(), boardID, columnID, taskID, req.Description)
    if err != nil {
       	writeError(w, errorStatus(err), err.Error())
        	return
	}

    writeJSON(w, http.StatusOK, task)
}
