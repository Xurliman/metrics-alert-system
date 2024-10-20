package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/services"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/utils"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

func main() {
	address, pollInterval, reportInterval := handleStartupParameters()
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

func handleStartupParameters() (address string, pollInterval, reportInterval time.Duration) {
	options := utils.ParseFlags()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	address, err = config.GetHost()
	if err != nil {
		address = options.GetAddr()
	}
	pollInterval, err = config.GetPollInterval()
	if err != nil {
		pollInterval = options.GetPollInterval()
	}
	reportInterval, err = config.GetReportInterval()
	if err != nil {
		reportInterval = options.GetReportInterval()
	}
	return address, pollInterval, reportInterval
}
