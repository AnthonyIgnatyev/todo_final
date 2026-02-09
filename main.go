package main

import (
	"log"

	"todo_final/pkg/config"
	database "todo_final/pkg/db"
	"todo_final/pkg/server"
)

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error to load config:", err)
	}

	err = database.Init()
	if err != nil {
		log.Fatal("Error to init database:", err)
	}
	defer database.CloseDB()

	err = server.StartServ()
	if err != nil {
		log.Fatal("Error to start server:", err)
	}
}
