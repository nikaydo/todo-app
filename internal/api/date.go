package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	d "github.com/nikaydo/final/internal/date"
)

var (
	ErrValidationDate = errors.New("error validation date")
)

func (h *Handlers) nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	dr, err := time.Parse("20060102", now)
	if err != nil {
		dr = time.Now()
	}
	newDate, err := d.NextDate(dr, date, repeat)
	if err != nil {
		writeJSONResponse(w, MsgError{Error: ErrValidationDate.Error()}, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(newDate))
	if err != nil {
		log.Println("error writing response", err)
	}
}
