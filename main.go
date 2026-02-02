package main

import (
	"log"
	"net/http"
	"os"

	_ "todo_final/pkg/api"
	"todo_final/pkg/db"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal("DB init error:", err)
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	webDir := "web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Server is running on port: %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server running error:", err)
	}
}
