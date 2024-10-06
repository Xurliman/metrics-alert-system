package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/routes"
)

func main() {
	r := routes.SetupRoutes()
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
