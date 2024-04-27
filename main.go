package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {

	log.Println("[INFO] Connecting to database...")
	db, err := sql.Open("sqlite", "scheduler.db")
	checkDbError(err, "connection")
	defer db.Close()

	s := NewStorage(db)

	s.InitDatabase()

	t := NewTaskService(s)

	http.Handle("/", http.FileServer(http.Dir("./web/")))
	http.Handle("/api/nextdate", http.HandlerFunc(t.ExtractParamsDateHandler))
	http.Handle("/api/task", http.HandlerFunc(t.TaskHandler))
	http.Handle("GET /api/tasks", http.HandlerFunc(t.TasksHandler))
	http.Handle("POST /api/task/done", http.HandlerFunc(t.DoneHandler))

	log.Println("[INFO] Starting server on port 7540...")
	log.Fatal(http.ListenAndServe(":7540", nil))
}
