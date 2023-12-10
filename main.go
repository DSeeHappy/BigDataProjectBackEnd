package main

import (
	"Backend/config"
	"Backend/server"
	"log"
	"os"
)

func main() {
	log.Println("Starting Jobs App")

	log.Println("Initializing configuration")
	initConfig := config.InitConfig(getConfigFileName())

	log.Println("Initializing database")
	dbHandler := server.InitDatabase(initConfig)

	log.Println("Initializing HTTP sever")
	httpServer := server.InitHttpServer(initConfig, dbHandler)

	httpServer.Start()
}

func getConfigFileName() string {
	env := os.Getenv("ENV")

	if env != "" {
		return "jobs-" + env
	}

	return "jobs"
}
