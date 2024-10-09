package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/routes"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"log"
)

func main() {
	r := routes.SetupRoutes()
	port := utils.GetParsedPort()
	log.Printf("Starting server on port %v", port)
	err := r.Run(port)
	if err != nil {
		return
	}
}
