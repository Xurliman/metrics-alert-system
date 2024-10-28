package services

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"strconv"
)

type MetricsService struct {
}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

var metrics = models.ExistingMetrics()

func (s *MetricsService) GetAll() map[string]string {
	data := make(map[string]string)
	for metricName, metricValue := range metrics.Gauge {
		data[metricName] = strconv.FormatFloat(metricValue, 'f', -1, 64)
	}
	for metricName, metricValue := range metrics.Counter {
		data[metricName] = strconv.FormatInt(metricValue, 10)
	}
	return data
}

func (s *MetricsService) SaveGaugeMetric(metricsName string, metricsValue string) error {
	if metricsName == "" {
		return constants.ErrEmptyMetricName
	}
	metricVal, err := strconv.ParseFloat(metricsValue, 64)
	if err != nil {
		return constants.ErrInvalidGaugeMetricValue
	}

	metrics.Gauge[metricsName] = metricVal
	return nil
}

func (s *MetricsService) SaveCounterMetric(metricsName string, metricsValue string) error {
	metricVal, err := strconv.ParseInt(metricsValue, 10, 64)
	if err != nil {
		return constants.ErrInvalidCounterMetricValue
	}

	metrics.Counter[metricsName] += metricVal
	return nil
}

func (s *MetricsService) SaveMetric(metricsType string, metricsName string, metricsValue string) error {
	switch metricsType {
	case constants.GaugeMetricType:
		return s.SaveGaugeMetric(metricsName, metricsValue)
	case constants.CounterMetricType:
		return s.SaveCounterMetric(metricsName, metricsValue)
	default:
		return constants.ErrInvalidMetricType
	}
}

func (s *MetricsService) FindGaugeMetric(metricsName string) (metricsValue string, err error) {
	value, ok := metrics.Gauge[metricsName]
	if !ok {
		return metricsValue, fmt.Errorf("type gauge|metrics %s not found", metricsName)
	}

	metricsValue = strconv.FormatFloat(value, 'f', -1, 64)
	return metricsValue, nil
}

func (s *MetricsService) FindCounterMetric(metricsName string) (metricsValue string, err error) {
	value, ok := metrics.Counter[metricsName]
	if !ok {
		return metricsValue, fmt.Errorf("type counter|metrics %s not found", metricsName)
	}

	metricsValue = strconv.FormatInt(value, 10)
	return metricsValue, nil
}

func (s *MetricsService) FindMetricByName(metricsType string, metricsName string) (metricsValue string, err error) {
	switch metricsType {
	case constants.GaugeMetricType:
		return s.FindGaugeMetric(metricsName)
	case constants.CounterMetricType:
		return s.FindCounterMetric(metricsName)
	default:
		return metricsValue, constants.ErrInvalidMetricType
	}
}
