// Package services define business logic for agent to send and collect metrics from the server
package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"

	"github.com/DataDog/gopsutil/mem"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"github.com/Xurliman/metrics-alert-system/internal/compressor"
)

type MetricsService struct {
	repository        interfaces.MetricsRepository
	client            http.Client
	cfg               *config.Config
	metricsCollection map[string]*models.Metrics
}

func NewMetricsService(repository interfaces.MetricsRepository, cfg *config.Config) *MetricsService {
	return &MetricsService{
		repository:        repository,
		client:            http.Client{},
		cfg:               cfg,
		metricsCollection: make(map[string]*models.Metrics),
	}
}

var pollCount int64

func (s *MetricsService) CollectMetricValues() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metricsCollection := []*models.Metrics{
		models.NewGaugeMetric("Alloc", float64(memStats.Alloc)),
		models.NewGaugeMetric("BuckHashSys", float64(memStats.BuckHashSys)),
		models.NewGaugeMetric("Frees", float64(memStats.Frees)),
		models.NewGaugeMetric("GCCPUFraction", memStats.GCCPUFraction),
		models.NewGaugeMetric("GCSys", float64(memStats.GCSys)),
		models.NewGaugeMetric("HeapAlloc", float64(memStats.HeapAlloc)),
		models.NewGaugeMetric("HeapIdle", float64(memStats.HeapIdle)),
		models.NewGaugeMetric("HeapInuse", float64(memStats.HeapInuse)),
		models.NewGaugeMetric("HeapObjects", float64(memStats.HeapObjects)),
		models.NewGaugeMetric("HeapReleased", float64(memStats.HeapReleased)),
		models.NewGaugeMetric("HeapSys", float64(memStats.HeapSys)),
		models.NewGaugeMetric("LastGC", float64(memStats.LastGC)),
		models.NewGaugeMetric("Lookups", float64(memStats.Lookups)),
		models.NewGaugeMetric("MCacheInuse", float64(memStats.MCacheInuse)),
		models.NewGaugeMetric("MCacheSys", float64(memStats.MCacheSys)),
		models.NewGaugeMetric("MSpanInuse", float64(memStats.MSpanInuse)),
		models.NewGaugeMetric("MSpanSys", float64(memStats.MSpanSys)),
		models.NewGaugeMetric("Mallocs", float64(memStats.Mallocs)),
		models.NewGaugeMetric("NextGC", float64(memStats.NextGC)),
		models.NewGaugeMetric("NumForcedGC", float64(memStats.NumForcedGC)),
		models.NewGaugeMetric("NumGC", float64(memStats.NumGC)),
		models.NewGaugeMetric("OtherSys", float64(memStats.OtherSys)),
		models.NewGaugeMetric("PauseTotalNs", float64(memStats.PauseTotalNs)),
		models.NewGaugeMetric("StackInuse", float64(memStats.StackInuse)),
		models.NewGaugeMetric("StackSys", float64(memStats.StackSys)),
		models.NewGaugeMetric("Sys", float64(memStats.Sys)),
		models.NewGaugeMetric("TotalAlloc", float64(memStats.TotalAlloc)),
		models.NewGaugeMetric("RandomValue", rand.Float64()),

		models.NewCounterMetric("PollCount", pollCount),
	}

	if v, err := mem.VirtualMemory(); err == nil {
		metricsCollection = append(metricsCollection,
			models.NewGaugeMetric("TotalMemory", float64(v.Total)),
			models.NewGaugeMetric("FreeMemory", float64(v.Free)),
			models.NewGaugeMetric("CPUutilization1", float64(v.Used)),
		)
	}
	pollCount++
	err := s.repository.SaveAll(metricsCollection)
	if err != nil {
		return err
	}
	return nil
}

func (s *MetricsService) GetAll() map[string]*models.Metrics {
	return s.repository.GetAll()
}

func (s *MetricsService) SendBatchMetrics() (err error) {
	url := fmt.Sprintf("%s/updates/", s.cfg.GetHost())
	var requestsToSend []requests.MetricsRequest

	for _, metric := range s.metricsCollection {
		request, err := s.repository.GetPlainRequest(metric)
		if err != nil {
			return err
		}
		requestsToSend = append(requestsToSend, *request)
	}

	marshalledRequest, err := json.Marshal(requestsToSend)
	if err != nil {
		return err
	}

	compressedRequest, err := compressor.Compress(marshalledRequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(compressedRequest))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Encoding", "gzip")
	return s.makeRequest(request)

}

func (s *MetricsService) SendMetricWithParams(ctx context.Context, metric *models.Metrics) (err error) {
	urlParams, err := s.repository.GetRequestURL(metric)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", s.cfg.GetHost()+urlParams, nil)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "text/plain")
	if err = s.makeRequest(request); err != nil {
		return err
	}
	return nil
}

func (s *MetricsService) SendCompressedMetric(ctx context.Context, metric *models.Metrics) (errs error) {
	url := fmt.Sprintf("%s/update/", s.cfg.GetHost())
	requestBody, err := s.repository.GetRequestBody(metric)
	if err != nil {
		return err
	}

	compressedRequest, err := compressor.Compress(requestBody)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(compressedRequest))
	if err != nil {
		return errors.Join(errs, err)
	}

	dst, err := s.hashRequest(compressedRequest)
	if err != nil && !errors.Is(err, constants.ErrKeyMissing) {
		errs = errors.Join(errs, err)
	} else {
		request.Header.Set("HashSHA256", dst)
	}

	request.Header.Set("Content-Encoding", "gzip")
	if err = s.makeRequest(request); err != nil {
		return errors.Join(errs, err)
	}
	return errs
}

func (s *MetricsService) SendMetric(ctx context.Context, metric *models.Metrics) error {
	url := fmt.Sprintf("%s/update/", s.cfg.GetHost())
	requestBody, err := s.repository.GetRequestBody(metric)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	dst, err := s.hashRequest(requestBody)
	if err != nil && !errors.Is(err, constants.ErrKeyMissing) {
		return err
	} else {
		request.Header.Set("HashSHA256", dst)
	}

	if err = s.makeRequest(request); err != nil {
		return err
	}
	return nil
}

func (s *MetricsService) makeRequest(request *http.Request) (err error) {
	request.Header.Add("Content-Type", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return err
	}

	err = response.Body.Close()
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return constants.ErrStatusNotOK
	}
	return nil
}

func (s *MetricsService) hashRequest(requestBody []byte) (string, error) {
	if s.cfg.Key == "" {
		return "", constants.ErrKeyMissing
	}
	h := hmac.New(sha256.New, []byte(s.cfg.Key))
	h.Write(requestBody)
	dst := h.Sum(nil)
	return hex.EncodeToString(dst), nil
}
