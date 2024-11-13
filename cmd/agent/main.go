package main

import (
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

	client := http.Client{Timeout: 10 * time.Second}
	metricRepository := repositories.NewMetricsRepository()
	metricsService := services.NewMetricsService(metricRepository)
	metricController := controllers.NewMetricsController(client, metricsService, address)

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	defer pollTicker.Stop()
	defer reportTicker.Stop()
	for {
		select {
		case <-pollTicker.C:
			metricController.CollectMetrics()
		case <-reportTicker.C:
			handleError(metricController.SendMetricsWithParams())
			handleError(metricController.SendMetrics())
			handleError(metricController.SendCompressedMetrics())
			handleError(metricController.SendCompressedMetricsWithParams())
			handleError(metricController.SendBatchMetrics())
		}
	}
}

func handleError(err error) {
	if err != nil {
		utils.Logger.Error("send metrics with params error", zap.Error(err))
	}
}
