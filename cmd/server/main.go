package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"github.com/Xurliman/metrics-alert-system/cmd/server/database"
	"github.com/Xurliman/metrics-alert-system/cmd/server/routes"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"time"
)

func main() {
	log.InitLogger(os.Getenv("APP_ENV"), constants.LogFilePath)
	err := godotenv.Load(constants.EnvFilePath)
	if err != nil {
		log.Warn("error when loading env", zap.Error(constants.ErrLoadingEnv))
	}

	cfg, err := config.Setup()
	if err != nil {
		log.Error("error when setting up configuration", zap.Error(err))
	}

	archiveService := services.NewArchiveService(cfg.FileStoragePath)
	var repo interfaces.MetricsRepositoryInterface
	if err = database.OpenDB(cfg.DatabaseDSN); err != nil {
		log.Warn("error connecting to database", zap.Error(err))
		repo = repositories.NewMetricsRepository(cfg.Restore, archiveService)
		go func() {
			storeTicker := time.NewTicker(cfg.GetStoreInterval())
			defer storeTicker.Stop()
			for range storeTicker.C {
				metrics, err := repo.List()
				if err != nil {
					log.Error("error getting list of metrics", zap.Error(err))
				}

				err = archiveService.Archive(metrics)
				if err != nil {
					log.Error("archiving data went wrong", zap.Error(err))
				}
			}
		}()
	} else {
		repo = repositories.NewDBMetricsRepository()
	}

	r := routes.SetupRoutes(repo, cfg.Key)
	log.Info("Starting server on port ",
		zap.String("port", cfg.GetPort()),
	)
	err = r.Run(cfg.GetPort())
	if err != nil {
		return
	}
}
