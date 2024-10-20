package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"github.com/Xurliman/metrics-alert-system/cmd/server/routes"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	r := routes.SetupRoutes()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := config.GetPort()
	if err != nil {
		port = utils.GetPort()
	}

	log.Printf("Starting server on port %v", port)
	err = r.Run(port)
	if err != nil {
		return
	}
}
