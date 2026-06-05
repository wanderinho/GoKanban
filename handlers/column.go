package handlers

import (
	"encoding/json"
	"kanban/handlers/dto"
	"net/http"
)

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
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
	
	var req dto.CreateTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
		return
	}

	task, err := h.httpStorage.AddTask(r.Context(), boardID, columnID, req.Title, req.Description)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
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

	task, err := h.httpStorage.GetTask(r.Context(), boardID, columnID, taskID)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
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

	if err := h.httpStorage.RemoveTask(r.Context(), boardID, columnID, taskID); err != nil {
		writeError(w, errorStatus(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateColumnTitle(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UpdateColumnTitleDTO
	var res *dto.ColumnDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    	writeError(w, http.StatusBadRequest, "не удалось прочитать тело запроса")
     	return
	}
	
	column, err := h.httpStorage.UpdateColumnTitle(r.Context(), boardID, columnID, req.Title)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}

	res, err = columnToDTO(r.Context(), h, column)
	if err != nil {
    	writeError(w, errorStatus(err), err.Error())
     	return
	}
	
	writeJSON(w, http.StatusOK, res)
}
