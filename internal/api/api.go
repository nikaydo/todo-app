package api

import (
	"net/http"

	"github.com/nikaydo/final/internal/db"
)

func Init(webDir string, Database db.Database) {
	var hander Handlers = Handlers{Database: Database}
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", hander.nextDayHandler)
	http.HandleFunc("/api/signin", hander.signin)
	http.HandleFunc("/api/task", auth(hander.taskHandler))
	http.HandleFunc("/api/tasks", auth(hander.getTasks))
	http.HandleFunc("/api/task/done", auth(hander.taskDone))
}
