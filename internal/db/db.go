package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikaydo/final/internal/model"
	_ "modernc.org/sqlite"
)

var (
	ErrNotFoundId = errors.New("not found")
	ErrIncorectId = errors.New("incorrect id for updating task")
)

const (
	schema = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title TEXT,
		comment TEXT,
		repeat VARCHAR(128)
	);`
	idx = `CREATE INDEX IF NOT EXISTS scheduler_date_index ON scheduler(date);`
)

type Database struct {
	db *sql.DB
}

func (d *Database) InitDatabase(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("create table schema: %w", err)
	}
	_, err = db.Exec(idx)
	if err != nil {
		return fmt.Errorf("create idndex: %w", err)
	}
	d.db = db
	return nil
}
func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) AddTask(task model.Task) (int64, error) {

	n, err := d.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, fmt.Errorf("add task Exec: %w", err)
	}
	num, err := n.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add task LastInsertId: %w", err)
	}
	return num, nil
}
func (d *Database) UpdateTask(task model.Task) error {
	n, err := d.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id;",
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return fmt.Errorf("update task Exec: %w", err)
	}
	count, err := n.RowsAffected()
	if err != nil {
		return fmt.Errorf("update task RowsAffected: %w", err)
	}
	if count == 0 {
		return ErrIncorectId
	}
	return nil
}

func (d *Database) UpdateDate(id, date string) error {
	n, err := d.db.Exec("UPDATE scheduler SET date = :date WHERE id = :id;",
		sql.Named("id", id),
		sql.Named("date", date))
	if err != nil {
		return fmt.Errorf("update date Exec: %w", err)
	}
	count, err := n.RowsAffected()
	if err != nil {
		return fmt.Errorf("update date RowsAffected: %w", err)
	}
	if count == 0 {
		return ErrIncorectId
	}
	return nil
}

func (d *Database) GetTask(id int) (model.Task, error) {
	var r model.Task
	list := d.db.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	err := list.Scan(&r.ID, &r.Date, &r.Title, &r.Comment, &r.Repeat)
	if err != nil {
		return r, fmt.Errorf("get task Scan: %w", err)
	}
	return r, nil
}

func (d *Database) GetTasks(limit int) (model.Response, error) {
	var result model.Response = model.Response{T: []model.Task{}}
	list, err := d.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit;",
		sql.Named("limit", limit))
	if err != nil {
		return result, fmt.Errorf("get tasks Query: %w", err)
	}
	defer list.Close()
	for list.Next() {
		var r model.Task
		err := list.Scan(&r.ID, &r.Date, &r.Title, &r.Comment, &r.Repeat)
		if err != nil {
			return result, fmt.Errorf("get tasks Scan: %w", err)
		}
		result.T = append(result.T, r)
	}
	err = list.Err()
	if err != nil {
		return result, fmt.Errorf("get tasks check error after search: %w", err)
	}
	return result, nil
}
func (d *Database) SearchTask(search string, limit int) (model.Response, error) {
	var result model.Response = model.Response{T: []model.Task{}}
	searchPattern := "%" + search + "%"
	list, err := d.db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search OR date LIKE :search ORDER BY date LIMIT :limit;",
		sql.Named("search", searchPattern),
		sql.Named("limit", limit))
	if err != nil {
		return result, fmt.Errorf("search task query: %w", err)
	}
	defer list.Close()
	for list.Next() {
		var r model.Task
		err := list.Scan(&r.ID, &r.Date, &r.Title, &r.Comment, &r.Repeat)
		if err != nil {
			return result, fmt.Errorf("search task Scan: %w", err)
		}
		result.T = append(result.T, r)
	}
	return result, nil
}

func (d *Database) DelTask(id int) error {
	_, err := d.db.Exec("DELETE FROM scheduler WHERE id = :id;",
		sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("delete task Exec: %w", err)
	}
	return nil
}
