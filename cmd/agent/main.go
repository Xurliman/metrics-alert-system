package main

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/controllers"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func main() {
	utils.Logger = utils.NewLogger()

	err := godotenv.Load(constants.EnvFilePath)
	if err != nil {
		log.Println(constants.ErrLoadingEnv)
	}

	flagCfg := utils.NewOptions()
	envCfg := config.NewConfig()

	address, err := flagCfg.GetHost()
	if err != nil {
		address, _ = envCfg.GetHost()
	}

	pollInterval, err := flagCfg.GetPollInterval()
	if err != nil {
		pollInterval, _ = envCfg.GetPollInterval()
	}

	reportInterval, err := flagCfg.GetReportInterval()
	if err != nil {
		reportInterval, _ = envCfg.GetReportInterval()
	}

	key, err := flagCfg.GetKey()
	if err != nil {
		key, _ = envCfg.GetKey()
	}

	rateLimit, err := flagCfg.GetRateLimit()
	if err != nil {
		rateLimit, _ = envCfg.GetRateLimit()
	}

	client := http.Client{Timeout: 10 * time.Second}
	metricRepository := repositories.NewMetricsRepository()
	metricsService := services.NewMetricsService(metricRepository)
	metricController := controllers.NewMetricsController(client, metricsService, address, key)

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defer pollTicker.Stop()
	defer reportTicker.Stop()
	for i := 0; i < rateLimit; i++ {
		select {
		case <-pollTicker.C:
			metricController.CollectMetrics()
		case <-reportTicker.C:

			go handleError(metricController.SendMetricsWithParams(ctx))
			go handleError(metricController.SendMetrics(ctx))
			go handleError(metricController.SendCompressedMetrics(ctx))
			go handleError(metricController.SendCompressedMetricsWithParams(ctx))
			go handleError(metricController.SendBatchMetrics())
		}
	}
}

func handleError(err error) {
	if err != nil {
		utils.Logger.Error("error while sending metrics", zap.Error(err))
	}
}
