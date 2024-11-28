package services

import (
	"context"
	"github.com/DataDog/gopsutil/mem"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/utils"
	"log"
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
		oldMetricsCollection: models.OldMetrics{
			Gauge:   make(map[string]float64),
			Counter: make(map[string]int64),
		},
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

	if v, err := mem.VirtualMemory(); err == nil {
		oldCollection.Gauge["TotalMemory"] = float64(v.Total)
		oldCollection.Gauge["FreeMemory"] = float64(v.Free)
		oldCollection.Gauge["CPUutilization1"] = float64(v.Used)
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

func (s *MetricsService) Generator(ctx context.Context) chan *models.Metrics {
	inputCh := make(chan *models.Metrics, len(s.metricsCollection))
	go func() {
		defer close(inputCh)

		for _, metric := range s.metricsCollection {
			log.Println("I have a metric: ", metric) // Debug log
			select {
			case <-ctx.Done():
				log.Println("Ctx done in generator", ctx.Err())
				return
			case inputCh <- metric:
				log.Println("Generated metric: ", metric) // Debug log
			}
		}
	}()
	return inputCh
}

func (s *MetricsService) URLConstructor(ctx context.Context, inputCh chan *models.Metrics, address string) chan models.Result {
	resultCh := make(chan models.Result, len(s.metricsCollection))
	go func() {
		defer close(resultCh)

		for metric := range inputCh {
			request, err := s.repository.GetRequestURL(metric, address)
			result := models.NewURLResult(request, err)
			select {
			case <-ctx.Done():
				return
			case resultCh <- result:
			}
		}
	}()

	return resultCh
}

func (s *MetricsService) RequestConstructor(ctx context.Context, inputCh chan *models.Metrics) chan models.Result {
	resultCh := make(chan models.Result, len(s.metricsCollection))
	go func() {
		defer close(resultCh)

		for metric := range inputCh {
			request, err := s.repository.GetPlainRequest(metric)
			result := models.NewRequestResult(request, err)
			select {
			case <-ctx.Done():
				return
			case resultCh <- result:
			}
		}
	}()

	return resultCh
}

func (s *MetricsService) ByteTransformer(ctx context.Context, inputCh chan models.Result) chan models.Result {
	resultCh := make(chan models.Result, len(s.metricsCollection))
	go func() {
		defer close(resultCh)

		for input := range inputCh {
			if err := input.Error(); err != nil {
				resultCh <- models.NewResult(nil, err)
				return
			}
			request, err := s.repository.GetBytes(input.Request())
			result := models.NewResult(request, err)
			select {
			case <-ctx.Done():
				return
			case resultCh <- result:
			}
		}
	}()

	return resultCh
}

func (s *MetricsService) RequestCompressor(ctx context.Context, inputCh chan models.Result) chan models.Result {
	resultCh := make(chan models.Result, len(s.metricsCollection))
	go func() {
		defer close(resultCh)
		for input := range inputCh {
			if err := input.Error(); err != nil {
				resultCh <- models.NewResult(nil, err)
				return
			}
			request, err := utils.Compress(input.Bytes())
			result := models.NewResult(request, err)
			select {
			case <-ctx.Done():
				log.Println("Ctx done in request compressor", ctx.Err())
				return
			case resultCh <- result:
				log.Println("Request compressor done", result)
			}
		}
	}()

	return resultCh
}
