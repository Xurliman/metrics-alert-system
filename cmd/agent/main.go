package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/services"
	"time"
)

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	metrics := services.CollectMetrics()
	for {
		select {
		case <-time.Tick(pollInterval):
			metrics = services.CollectMetrics()
		case <-time.Tick(reportInterval):
			services.SendMetrics(metrics)
		}
	}
}
