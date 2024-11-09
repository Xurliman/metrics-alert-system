package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"github.com/Xurliman/metrics-alert-system/cmd/server/routes"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"time"
)

func main() {
	utils.Logger = utils.NewLogger(config.GetAppEnv())

	err := godotenv.Load(constants.EnvFilePath)
	if err != nil {
		utils.Logger.Error("", zap.Error(constants.ErrLoadingEnv))
	}

	flagOptions := utils.NewOptions()

	port, err := flagOptions.GetPort()
	if err != nil {
		port, _ = config.GetPort()
	}

	fileStoragePath, err := flagOptions.GetFileStoragePath()
	if err != nil {
		fileStoragePath = config.GetFileStoragePath()
	}

	storeInterval, err := flagOptions.GetStoreInterval()
	if err != nil {
		storeInterval = config.GetStoreInterval()
	}

	shouldRestore := flagOptions.GetShouldRestore() && config.GetShouldRestore()
	archiveService, err := services.NewArchiveService(fileStoragePath)
	if err != nil {
		utils.Logger.Error("error related to archiving", zap.Error(err))
	}

	dsn, err := flagOptions.GetDatabaseDSN()
	if err != nil {
		dsn = config.GetDatabaseDSN()
	}
	config.Open(dsn)

	r, repo := routes.SetupRoutes(shouldRestore, archiveService)
	utils.Logger.Error("Starting server on port ",
		zap.String("port", port),
	)

	go func() {
		storeTicker := time.NewTicker(storeInterval)
		defer storeTicker.Stop()
		for range storeTicker.C {
			err = archiveService.Archive(repo.List())
			if err != nil {
				utils.Logger.Error("archiving data went wrong", zap.Error(err))
			}
		}
	}()
	err = r.Run(port)
	if err != nil {
		return
	}
}
