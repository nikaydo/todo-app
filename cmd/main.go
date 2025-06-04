package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	db "github.com/nikaydo/final/internal/db"
	"github.com/nikaydo/final/internal/server"
)

func checkEnv(s string) string {
	e := os.Getenv(s)
	if e == "" {
		log.Println(s + " not set")
		return ""
	}
	log.Println(s+":", e)
	return e
}

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	todo_port := checkEnv("TODO_PORT")
	todo_dbfile := checkEnv("TODO_DBFILE")
	var Database db.Database
	err = Database.InitDatabase(todo_dbfile)
	if err != nil {
		log.Fatalln(err)
	}
	defer Database.Close()
	webDir := "./web"
	if err := server.Run(todo_port, webDir, Database); err != nil {
		log.Fatalln(err)
		return
	}
}
