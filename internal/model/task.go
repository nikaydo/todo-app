package model

import (
	"errors"
	"time"

	date "github.com/nikaydo/final/internal/date"
)

var (
	ErrDateFormat = errors.New("ошибка при проверке даты")
)

type Response struct {
	T []Task `json:"tasks"`
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func (task *Task) CheckDate() error {
	now := time.Now()
	today := now.Format("20060102")
	if task.Date == "" {
		task.Date = today
		return nil
	}
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return ErrDateFormat
	}
	if task.Repeat == "" {
		if t.Before(now) {
			task.Date = today
			return nil
		}
		return nil
	}
	d, err := date.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}
	task.Date = d
	return nil
}
