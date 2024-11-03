package services

import (
	"errors"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"runtime"
	"strconv"
)

var MetricsCollection map[string]*models.Metrics

type MetricsService struct{}

func NewMetricsService() *MetricsService {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	gaugeMetricExample := float64(memStats.Alloc)
	var counterMetricExample int64
	MetricsCollection = map[string]*models.Metrics{
		"Alloc": {
			ID:    "Alloc",
			MType: constants.GaugeMetricType,
			Value: &gaugeMetricExample,
		},
		"PollCount": {
			ID:    "PollCount",
			MType: constants.CounterMetricType,
			Delta: &counterMetricExample,
		},
	}
	return &MetricsService{}
}

func (s *MetricsService) List() map[string]string {
	data := make(map[string]string)
	for metricName, metric := range MetricsCollection {
		switch metric.MType {
		case constants.GaugeMetricType:
			data[metricName] = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		case constants.CounterMetricType:
			data[metricName] = strconv.FormatInt(*metric.Delta, 10)
		}
	}
	return data
}

func (s *MetricsService) GetMetricValue(metric interfaces.MetricsInterface, metricName string) (metricValue string, err error) {
	return metric.GetMetricValue(metricName)
}

func (s *MetricsService) SaveWhenParams(metric interfaces.MetricsInterface, metricName, metricValue string) error {
	existingMetric, err := metric.FindByName(metricName)
	if err != nil && !errors.Is(err, constants.ErrMetricNotFound) {
		return err
	}

	return metric.Save(metricName, metricValue, existingMetric)
}

func (s *MetricsService) SaveWhenBody(metric interfaces.MetricsInterface, metricRequest requests.MetricsSaveRequest) (entry *models.Metrics, err error) {
	existingMetric, err := metric.FindByName(metricRequest.ID)
	if err != nil && !errors.Is(err, constants.ErrMetricNotFound) {
		return nil, err
	}

	return metric.SaveBody(metricRequest, existingMetric)
}

func (s *MetricsService) Show(metric interfaces.MetricsInterface, metricName string) (entry *models.Metrics, err error) {
	return metric.FindByName(metricName)
}

var (
	Counter = counterMetricsService{}
	Gauge   = gaugeMetricsService{}
)
