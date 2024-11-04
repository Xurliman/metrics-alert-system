package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/utils"
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

func SendMetrics(client http.Client, metrics models.Metrics, address string) {
	for metric, value := range metrics.Gauge {
		url := fmt.Sprintf("http://%s/update/", address)
		request, err := json.Marshal(requests.MetricsRequest{
			ID:    metric,
			MType: constants.GaugeMetricType,
			Value: &value,
		})
		if err != nil {
			return
		}

		response, err := client.Post(url, "application/json", bytes.NewReader(request))
		if err != nil {
			return
		}

		err = response.Body.Close()
		if err != nil {
			return
		}
	}

	for metric, value := range metrics.Counter {
		url := fmt.Sprintf("http://%s/update/", address)
		request, err := json.Marshal(requests.MetricsRequest{
			ID:    metric,
			MType: constants.CounterMetricType,
			Delta: &value,
		})
		if err != nil {
			return
		}

		response, err := client.Post(url, "application/json", bytes.NewReader(request))
		if err != nil {
			return
		}

		err = response.Body.Close()
		if err != nil {
			return
		}
	}
}

func SendMetricsWithParam(client http.Client, metrics models.Metrics, address string) {
	for metric, value := range metrics.Gauge {
		url := fmt.Sprintf("http://%s/update/%s/%s/%f", address, constants.GaugeMetricType, metric, value)
		response, err := client.Post(url, "text/plain", nil)
		if err != nil {
			return
		}
		err = response.Body.Close()
		if err != nil {
			return
		}
	}

	for metric, value := range metrics.Counter {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v", address, constants.CounterMetricType, metric, value)
		response, err := client.Post(url, "text/plain", nil)
		if err != nil {
			return
		}
		err = response.Body.Close()
		if err != nil {
			return
		}
	}
}

func SendCompressedMetrics(client http.Client, metrics models.Metrics, address string) {
	for metric, value := range metrics.Gauge {
		url := fmt.Sprintf("http://%s/update/", address)
		body, err := json.Marshal(requests.MetricsRequest{
			ID:    metric,
			MType: constants.GaugeMetricType,
			Value: &value,
		})
		if err != nil {
			return
		}

		request, err := compress(body, url)
		if err != nil {
			return
		}

		response, err := client.Do(request)
		if err != nil {
			return
		}

		err = response.Body.Close()
		if err != nil {
			return
		}
	}

	for metric, value := range metrics.Counter {
		url := fmt.Sprintf("http://%s/update/", address)
		body, err := json.Marshal(requests.MetricsRequest{
			ID:    metric,
			MType: constants.CounterMetricType,
			Delta: &value,
		})
		if err != nil {
			return
		}

		request, err := compress(body, url)
		if err != nil {
			return
		}

		response, err := client.Do(request)
		if err != nil {
			return
		}

		err = response.Body.Close()
		if err != nil {
			return
		}
	}
}

func SendCompressedMetricsWithParam(client http.Client, metrics models.Metrics, address string) {
	for metric, value := range metrics.Gauge {
		url := fmt.Sprintf("http://%s/update/%s/%s/%f", address, constants.GaugeMetricType, metric, value)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		response, err := client.Do(req)
		if err != nil {
			return
		}
		err = response.Body.Close()
		if err != nil {
			return
		}
	}

	for metric, value := range metrics.Counter {
		url := fmt.Sprintf("http://%s/update/%s/%s/%v", address, constants.CounterMetricType, metric, value)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		response, err := client.Do(req)
		if err != nil {
			return
		}

		err = response.Body.Close()
		if err != nil {
			return
		}
	}
}

func compress(body []byte, url string) (*http.Request, error) {
	compressedRequest, err := utils.Compress(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(compressedRequest))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")

	return request, nil
}
