package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/services"
	"net/http"
	"time"
)

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second
	client := http.Client{Timeout: 10 * time.Second}
	metrics := services.CollectMetrics()
	services.SendMetrics(client, metrics)
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	defer pollTicker.Stop()
	defer reportTicker.Stop()
	for {
		select {
		case <-pollTicker.C:
			metrics = services.CollectMetrics()
		case <-reportTicker.C:
			services.SendMetrics(client, metrics)
		}
	}
}
