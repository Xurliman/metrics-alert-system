package services

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/models"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
)

var pollCount int64

func CollectMetrics() models.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	metrics := models.Metrics{
		Gauge: map[string]float64{
			"Alloc":         float64(memStats.Alloc),
			"BuckHashSys":   float64(memStats.BuckHashSys),
			"Frees":         float64(memStats.Frees),
			"GCCPUFraction": memStats.GCCPUFraction,
			"GCSys":         float64(memStats.GCSys),
			"HeapAlloc":     float64(memStats.HeapAlloc),
			"HeapIdle":      float64(memStats.HeapIdle),
			"HeapInuse":     float64(memStats.HeapInuse),
			"HeapObjects":   float64(memStats.HeapObjects),
			"HeapReleased":  float64(memStats.HeapReleased),
			"HeapSys":       float64(memStats.HeapSys),
			"LastGC":        float64(memStats.LastGC),
			"Lookups":       float64(memStats.Lookups),
			"MCacheInuse":   float64(memStats.MCacheInuse),
			"MCacheSys":     float64(memStats.MCacheSys),
			"MSpanInuse":    float64(memStats.MSpanInuse),
			"MSpanSys":      float64(memStats.MSpanSys),
			"Mallocs":       float64(memStats.Mallocs),
			"NextGC":        float64(memStats.NextGC),
			"NumForcedGC":   float64(memStats.NumForcedGC),
			"NumGC":         float64(memStats.NumGC),
			"OtherSys":      float64(memStats.OtherSys),
			"PauseTotalNs":  float64(memStats.PauseTotalNs),
			"StackInuse":    float64(memStats.StackInuse),
			"StackSys":      float64(memStats.StackSys),
			"Sys":           float64(memStats.Sys),
			"TotalAlloc":    float64(memStats.TotalAlloc),
			"RandomValue":   rand.Float64(),
		},
		Counter: map[string]int64{
			"PollCount": pollCount,
		},
	}
	pollCount++
	return metrics
}

func SendMetrics(metrics models.Metrics) {
	client := &http.Client{}
	for metric, value := range metrics.Gauge {
		url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%f", metric, value)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			log.Printf("Error creating request for gauge metrics: %s\n", err)
			continue
		}
		request.Header.Set("Content-Type", "text/plain")
		response, err := client.Do(request)
		if err != nil {
			log.Printf("Error sending request for gauge metrics: %s\n", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing body for gauge metrics: %s\n", err)
			}
		}(response.Body)
		log.Printf("Sent gauge %s with value %f, response: %s\n", metric, value, response.Status)
	}

	for metric, value := range metrics.Counter {
		url := fmt.Sprintf("http://localhost:8080/update/counter/%s/%v", metric, value)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			log.Printf("Error creating request for counter metrics: %s\n", err)
		}
		request.Header.Set("Content-Type", "text/plain")
		response, err := client.Do(request)
		if err != nil {
			log.Printf("Error sending request for counter metrics: %s\n", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing body for counter metrics: %s\n", err)
			}
		}(response.Body)
		log.Printf("Sent counter %s with value %f, response: %s\n", metric, value, response.Status)
	}
}