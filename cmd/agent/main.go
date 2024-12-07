package main

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/controllers"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.InitLogger(os.Getenv("APP_ENV"), constants.LogFilePath)
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

	metricRepository := repositories.NewMetricsRepository()
	metricsService := services.NewMetricsService(metricRepository, cfg)
	metricController := controllers.NewMetricsController(metricsService, cfg)
	metricController.Run(ctx)
}
