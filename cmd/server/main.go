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
	"github.com/Xurliman/metrics-alert-system/internal/cert"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func getGitCommit() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(output))
}

func getDate() string {
	cmd := exec.Command("date", "+%Y-%m-%d")
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(output))
}

func getValue(val string, fallback func() string) string {
	if val == "" {
		return fallback()
	}
	return val
}

// go build -ldflags "-X 'main.buildVersion=1.0.0' -X 'main.buildDate=2024-04-03' -X 'main.buildCommit=$(git rev-parse HEAD)'"  -o server cmd/server/main.go
func init() {
	log.InitLogger(os.Getenv("APP_ENV"), constants.LogFilePath)
	log.Info("Build:",
		zap.String("version", getValue(buildVersion, func() string { return "N/A" })),
		zap.String("date", getValue(buildDate, getDate)),
		zap.String("commit", getValue(buildCommit, getGitCommit)),
	)
}

func main() {
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

	err = cert.GenerateKeyPair()
	if err != nil {
		log.Fatal("failed to generate key pair", zap.Error(err))
	}
	log.Info("certificate successfully written")

	r := routes.SetupRoutes(repo, cfg)
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
