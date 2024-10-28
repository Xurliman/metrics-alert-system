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
		log.Fatal(constants.ErrLoadingEnv)
	}

	envCfg := utils.NewOptions()
	cfg := config.NewConfig()

	address, err := envCfg.GetHost()
	if err != nil {
		address, _ = cfg.GetHost()
	}

	pollInterval, err := envCfg.GetPollInterval()
	if err != nil {
		pollInterval, _ = cfg.GetPollInterval()
	}

	reportInterval, err := envCfg.GetReportInterval()
	if err != nil {
		reportInterval, _ = cfg.GetReportInterval()
	}

	client := http.Client{Timeout: 10 * time.Second}
	metrics := services.CollectMetrics()
	services.SendMetrics(client, metrics, address)

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
