package server

import (
	"log"
	"net/http"
	"todo_final/pkg/api"
	cfg "todo_final/pkg/config"
)

func StartServ() error {
	serverAddr := cfg.CfgStruct.Server.IP + ":" + cfg.CfgStruct.Server.Port

	api.Init()
	log.Printf("Starting server at: %s:%s", cfg.CfgStruct.Server.IP, cfg.CfgStruct.Server.Port)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Printf("Error starting server: %s", err)
		return err
	}

	return nil
}
