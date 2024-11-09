package services

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/utils"
	"math/rand"
	"runtime"
)

type MetricsService struct {
	repository           interfaces.MetricsRepository
	metricsCollection    map[string]*models.Metrics
	oldMetricsCollection models.OldMetrics
}

func NewMetricsService(repository interfaces.MetricsRepository) *MetricsService {

	return &MetricsService{
		repository:        repository,
		metricsCollection: make(map[string]*models.Metrics),
	}
}

var pollCount int64

func (s *MetricsService) CollectMetricValues() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	oldCollection := models.OldMetrics{
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
	s.oldMetricsCollection = oldCollection
	s.ConvertToMetrics()
}

func (s *MetricsService) ConvertToMetrics() {
	metrics := make(map[string]*models.Metrics)
	for metric, value := range s.oldMetricsCollection.Gauge {
		metrics[metric] = &models.Metrics{
			ID:    metric,
			MType: constants.GaugeMetricType,
			Value: &value,
		}
	}
	for metric, value := range s.oldMetricsCollection.Counter {
		metrics[metric] = &models.Metrics{
			ID:    metric,
			MType: constants.CounterMetricType,
			Delta: &value,
		}
	}
	s.metricsCollection = metrics
}

func (s *MetricsService) GetRequestBodies() ([][]byte, error) {
	var requestsToSend [][]byte
	for _, metric := range s.metricsCollection {
		request, err := s.repository.GetRequestBody(metric)
		if err != nil {
			return nil, err
		}
		requestsToSend = append(requestsToSend, request)
	}
	return requestsToSend, nil
}

func (s *MetricsService) GetCompressedRequestBodies() ([][]byte, error) {
	var requestsToSend [][]byte
	for _, metric := range s.metricsCollection {
		request, err := s.repository.GetRequestBody(metric)
		if err != nil {
			return nil, err
		}

		compressedRequest, err := utils.Compress(request)
		if err != nil {
			return nil, err
		}

		requestsToSend = append(requestsToSend, compressedRequest)
	}
	return requestsToSend, nil
}

func (s *MetricsService) GetRequestURLs(address string) ([]string, error) {
	var urls []string
	for _, metric := range s.metricsCollection {
		url, err := s.repository.GetRequestURL(metric, address)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
