package server

import (
	"fmt"
	"net/http"

	"github.com/nikaydo/final/internal/api"
	"github.com/nikaydo/final/internal/db"
)

func Run(port, webDir string, Database db.Database) error {
	api.Init(webDir, Database)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
