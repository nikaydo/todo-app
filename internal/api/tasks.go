package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/nikaydo/final/internal/db"
	"github.com/nikaydo/final/internal/model"
)

var (
	ErrInvalidID    = errors.New("invalid or missing identifier")
	ErrInvalidTitle = errors.New("invalid title")
)

type Handlers struct {
	Database db.Database
}

func (h *Handlers) taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.addTask(w, r)
	case http.MethodGet:
		h.getTask(w, r)
	case http.MethodPut:
		h.UpdateTask(w, r)
	case http.MethodDelete:
		h.delTask(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) addTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJSONResponse(w, MsgError{Error: ErrInvalidTitle.Error()}, http.StatusBadRequest)
		return
	}

	if err := validateAndAdjustDate(&task); err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := h.Database.AddTask(task)
	if err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, MsgId{Id: id}, http.StatusCreated)
}

func (h *Handlers) getTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSONResponse(w, MsgError{Error: ErrInvalidID.Error()}, http.StatusBadRequest)
		return
	}
	num, err := strconv.Atoi(id)
	if err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	n, err := h.Database.GetTask(num)
	if err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, n, http.StatusOK)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJSONResponse(w, MsgError{Error: ErrInvalidTitle.Error()}, http.StatusBadRequest)
		return
	}
	if err := validateAndAdjustDate(&task); err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := h.Database.UpdateTask(task); err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, struct{}{}, http.StatusOK)
}
func (h *Handlers) delTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSONResponse(w, MsgError{Error: ErrInvalidID.Error()}, http.StatusBadRequest)
		return
	}
	num, err := strconv.Atoi(id)
	if err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := h.Database.DelTask(num); err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, struct{}{}, http.StatusOK)
}

func (h *Handlers) getTasks(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			search = t.Format("20060102")
		}
		tasks, err := h.Database.SearchTask(search, 100)
		if err != nil {
			writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, tasks, http.StatusOK)
		return
	}

	tasks, err := h.Database.GetTasks(100)
	if err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, tasks, http.StatusOK)
}

func (h *Handlers) taskDone(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSONResponse(w, MsgError{Error: ErrInvalidID.Error()}, http.StatusBadRequest)
		return
	}
	num, err := strconv.Atoi(id)
	if err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	task, err := h.Database.GetTask(num)
	if err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
		return
	}
	if task.Repeat == "" {
		if err := h.Database.DelTask(num); err != nil {
			writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, struct{}{}, http.StatusOK)
		return
	}
	if err := task.CheckDate(); err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	if err := h.Database.UpdateTask(task); err != nil {
		writeJSONResponse(w, MsgError{Error: err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, struct{}{}, http.StatusOK)
}
