package main

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/controllers"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"github.com/Xurliman/metrics-alert-system/internal/rsa"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
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

// go build -ldflags "-X 'main.buildVersion=1.0.0' -X 'main.buildDate=2024-04-03' -X 'main.buildCommit=$(git rev-parse HEAD)'"  -o agent cmd/agent/main.go
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
		log.Error("error when loading env", zap.Error(constants.ErrLoadingEnv))
	}

	cfg, err := config.Setup()
	if err != nil {
		log.Fatal("error parsing config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	if cfg.CryptoKey != "" {
		err = rsa.LoadPublicKey(cfg.CryptoKey)
		if err != nil {
			log.Fatal("error loading crypto key", zap.Error(err))
		}
	}

	metricRepository := repositories.NewMetricsRepository()
	metricsService := services.NewMetricsService(metricRepository, cfg)
	metricController := controllers.NewMetricsController(metricsService, cfg)
	metricController.Run(ctx)
}
