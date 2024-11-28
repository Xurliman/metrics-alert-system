package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"github.com/Xurliman/metrics-alert-system/cmd/server/database"
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
		utils.Logger.Error("error when loading env", zap.Error(constants.ErrLoadingEnv))
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
	archiveService := services.NewArchiveService(fileStoragePath)

	dsn, err := flagOptions.GetDatabaseDSN()
	if err != nil {
		dsn = config.GetDatabaseDSN()
	}

	key, err := flagOptions.GetKey()
	if err != nil {
		key = config.GetKey()
	}

	utils.Logger.Debug("SERVER: ",
		zap.String("address", port),
		zap.String("store_interval", storeInterval.String()),
		zap.Bool("should_restore", shouldRestore),
		zap.String("key", key),
	)
	var repo interfaces.MetricsRepositoryInterface
	if err = database.OpenDB(dsn); err != nil {
		utils.Logger.Error("error connecting to database", zap.Error(err))
		repo = repositories.NewMetricsRepository(shouldRestore, archiveService)
		go func() {
			storeTicker := time.NewTicker(storeInterval)
			defer storeTicker.Stop()
			for range storeTicker.C {
				metrics, err := repo.List()
				if err != nil {
					utils.Logger.Error("error getting list of metrics", zap.Error(err))
				}

				err = archiveService.Archive(metrics)
				if err != nil {
					utils.Logger.Error("archiving data went wrong", zap.Error(err))
				}
			}
		}()
	} else {
		repo = repositories.NewDBMetricsRepository()
	}

	r := routes.SetupRoutes(repo, key)

	err = r.Run(port)
	if err != nil {
		return
	}

	utils.Logger.Error("Starting server on port ",
		zap.String("port", port),
	)
}
