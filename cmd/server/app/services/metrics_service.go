package services

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"math/rand"
	"runtime"
	"strconv"
)

var pollCount int64

type MetricsService struct {
}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

func (s *MetricsService) CollectMetrics() models.Metrics {
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

func (s *MetricsService) FindByName(metricsType string, metricsName string) (metricsValue string, err error) {
	metrics := s.CollectMetrics()
	if metricsType == "gauge" {
		value, ok := metrics.Gauge[metricsName]
		if !ok {
			return metricsValue, fmt.Errorf("metrics %s not found", metricsName)
		}
		metricsValue = strconv.FormatFloat(value, 'f', -1, 64)
	} else if metricsType == "counter" {
		value, ok := metrics.Counter[metricsName]
		if !ok {
			return metricsValue, fmt.Errorf("metrics %s not found", metricsName)
		}
		metricsValue = strconv.FormatInt(value, 10)
	}
	return metricsValue, nil
}

func (s *MetricsService) GetAll() map[string]string {
	metrics := s.CollectMetrics()
	data := make(map[string]string)
	for metricName, metricValue := range metrics.Gauge {
		data[metricName] = strconv.FormatFloat(metricValue, 'f', -1, 64)
	}
	for metricName, metricValue := range metrics.Counter {
		data[metricName] = strconv.FormatInt(metricValue, 10)
	}
	return data
}
