package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo, err := initializeRepository(cfg, archiveService)
	if err != nil {
		log.Error("error when initializing repository", zap.Error(err))
		return
	}
	go archiveToFile(ctx, cfg, repo, archiveService)

	r := routes.SetupRoutes(repo, cfg.Key)
	err = r.Run(cfg.GetPort())
	if err != nil {
		log.Fatal("error when starting server", zap.Error(err))
	}
}

func archiveToFile(ctx context.Context, cfg *config.Config, repository interfaces.MetricsRepositoryInterface, archiveService interfaces.ArchiveServiceInterface) {
	storeTicker := time.NewTicker(cfg.GetStoreInterval())
	defer storeTicker.Stop()

	for {
		select {
		case <-storeTicker.C:
			metrics, err := repository.List()
			if err != nil {
				log.Error("error getting list of metrics", zap.Error(err))
			}

			err = archiveService.Archive(metrics)
			if err != nil {
				log.Error("archiving data went wrong", zap.Error(err))
			}
		case <-ctx.Done():
			return
		}
	}
}

func initializeRepository(cfg *config.Config, archiveService interfaces.ArchiveServiceInterface) (interfaces.MetricsRepositoryInterface, error) {
	if cfg.DatabaseDSN == "" {
		if cfg.Restore {
			metrics, err := archiveService.Load()
			if err != nil {
				return nil, constants.ErrLoadingMetricsFromArchive
			}

			return repositories.NewMetricsRepository(metrics), nil
		}
		return nil, constants.ErrDatabaseDSNEmpty
	}

	if err := database.OpenDB(cfg.DatabaseDSN); err != nil {
		return nil, constants.ErrConnectingDatabase
	}
	return repositories.NewDBMetricsRepository(), nil
}
