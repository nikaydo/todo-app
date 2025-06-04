package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/nikaydo/final/internal/date"
	"github.com/nikaydo/final/internal/model"
)

type MsgId struct {
	Id int64 `json:"id"`
}

type MsgError struct {
	Error string `json:"error"`
}

func writeJSONResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		log.Println("error writing response", err)
	}
}

func validateAndAdjustDate(task *model.Task) error {
	layout := "20060102"
	nowDate := time.Now().Format(layout)
	if task.Date == "" {
		task.Date = nowDate
	}
	t, err := time.Parse(layout, task.Date)
	if err != nil {
		return err
	}
	taskDate := t.Format(layout)
	if task.Repeat != "" {
		if _, err := date.NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			return err
		}
		if taskDate < nowDate {
			nextDate, _ := date.NextDate(time.Now(), task.Date, task.Repeat)
			task.Date = nextDate
		}
	}
	if taskDate < nowDate {
		task.Date = nowDate
	}
	return nil
}
