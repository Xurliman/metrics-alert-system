package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/services"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/utils"
	"net/http"
	"time"
)

func main() {
	options := utils.ParseFlags()
	pollInterval := options.GetPollInterval()
	reportInterval := options.GetReportInterval()
	client := http.Client{Timeout: 10 * time.Second}
	metrics := services.CollectMetrics()
	services.SendMetrics(client, metrics, options.GetAddr())
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	defer pollTicker.Stop()
	defer reportTicker.Stop()
	for {
		select {
		case <-pollTicker.C:
			metrics = services.CollectMetrics()
		case <-reportTicker.C:
			services.SendMetrics(client, metrics, options.GetAddr())
		}
	}
}
