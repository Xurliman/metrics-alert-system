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
	memStats             runtime.MemStats
	oldMetricsCollection models.OldMetrics
}

func NewMetricsService(repository interfaces.MetricsRepository) *MetricsService {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &MetricsService{
		repository:        repository,
		metricsCollection: make(map[string]*models.Metrics),
		memStats:          memStats,
	}
}

func (s *MetricsService) CollectMetricValues() {
	oldCollection := models.OldMetrics{
		Gauge: map[string]float64{
			"Alloc":       float64(s.memStats.Alloc),
			"BuckHashSys": float64(s.memStats.BuckHashSys),
			//"Frees":         float64(s.memStats.Frees),
			"GCCPUFraction": s.memStats.GCCPUFraction,
			"GCSys":         float64(s.memStats.GCSys),
			//"HeapAlloc":     float64(s.memStats.HeapAlloc),
			//"HeapIdle":     float64(s.memStats.HeapIdle),
			//"HeapInuse":    float64(s.memStats.HeapInuse),
			//"HeapObjects":  float64(s.memStats.HeapObjects),
			"HeapReleased": float64(s.memStats.HeapReleased),
			"HeapSys":      float64(s.memStats.HeapSys),
			"LastGC":       float64(s.memStats.LastGC),
			"Lookups":      float64(s.memStats.Lookups),
			"MCacheInuse":  float64(s.memStats.MCacheInuse),
			"MCacheSys":    float64(s.memStats.MCacheSys),
			"MSpanInuse":   float64(s.memStats.MSpanInuse),
			"MSpanSys":     float64(s.memStats.MSpanSys),
			//"Mallocs":      float64(s.memStats.Mallocs),
			"NextGC":       float64(s.memStats.NextGC),
			"NumForcedGC":  float64(s.memStats.NumForcedGC),
			"NumGC":        float64(s.memStats.NumGC),
			"OtherSys":     float64(s.memStats.OtherSys),
			"PauseTotalNs": float64(s.memStats.PauseTotalNs),
			"StackInuse":   float64(s.memStats.StackInuse),
			"StackSys":     float64(s.memStats.StackSys),
			"Sys":          float64(s.memStats.Sys),
			//"TotalAlloc":   float64(s.memStats.TotalAlloc),
			"RandomValue": rand.Float64(),
		},
		Counter: map[string]int64{
			//"PollCount": pollCount,
		},
	}
	s.oldMetricsCollection = oldCollection
	s.CollectMetrics()
}

func (s *MetricsService) CollectMetrics() {
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
