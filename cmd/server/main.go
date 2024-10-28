package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"github.com/Xurliman/metrics-alert-system/cmd/server/routes"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
)

func main() {
	utils.Logger = utils.NewLogger(config.GetAppEnv())

	err := godotenv.Load(constants.EnvFilePath)
	if err != nil {
		log.Fatal(constants.ErrLoadingEnv)
	}

	r := routes.SetupRoutes()

	port, err := utils.GetPort()
	if err != nil {
		port, err = config.GetPort()
	}

	utils.Logger.Info("Starting server on port %v",
		zap.String("port", port),
	)
	err = r.Run(port)
	if err != nil {
		return
	}
}
