package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/utils"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

func main() {
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
	metrics := services.CollectMetrics()

	services.SendMetrics(client, metrics, address)
	services.SendMetricsWithParam(client, metrics, address)
	services.SendCompressedMetrics(client, metrics, address)
	services.SendCompressedMetricsWithParam(client, metrics, address)

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	defer pollTicker.Stop()
	defer reportTicker.Stop()
	for {
		select {
		case <-pollTicker.C:
			metrics = services.CollectMetrics()
		case <-reportTicker.C:
			services.SendMetrics(client, metrics, address)
		}
	}
}
