package services

import (
	"errors"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type MetricsService struct {
	sw         interfaces.Switch
	repository interfaces.MetricsRepositoryInterface
}

func NewMetricsService(repository interfaces.MetricsRepositoryInterface, switcher interfaces.Switch) *MetricsService {
	return &MetricsService{
		repository: repository,
		sw:         switcher,
	}
}

func (s *MetricsService) List() map[string]string {
	data := make(map[string]string)
	metricsCollection := s.repository.List()
	for metricName, metric := range metricsCollection {
		switch metric.MType {
		case constants.GaugeMetricType:
			data[metricName] = gaugeStrVal(metric)
		case constants.CounterMetricType:
			data[metricName] = counterStrVal(metric)
		}
	}
	return data
}

func (s *MetricsService) GetMetricValue(metricType, metricName string) (metricValue string, err error) {
	metric, err := s.repository.FindByName(metricName)
	if err != nil {
		return "", err
	}
	switch metricType {
	case constants.GaugeMetricType:
		metricValue = gaugeStrVal(metric)
	case constants.CounterMetricType:
		metricValue = counterStrVal(metric)
	}
	return metricValue, nil
}

func (s *MetricsService) SaveWhenParams(metricType, metricName, metricValue string) error {
	if metricName == "" {
		return constants.ErrEmptyMetricName
	}

	existingMetric, err := s.repository.FindByName(metricName)
	if err != nil && !errors.Is(err, constants.ErrMetricNotFound) {
		return err
	}

	var metric *models.Metrics
	switch metricType {
	case constants.GaugeMetricType:
		metric, err = s.sw.ConvertParams(GaugeConverter, existingMetric, metricName, metricValue)
	case constants.CounterMetricType:
		metric, err = s.sw.ConvertParams(CounterConverter, existingMetric, metricName, metricValue)
	default:
		return constants.ErrInvalidMetricType
	}
	if err != nil {
		return err
	}
	_ = s.repository.Save(metric)
	return nil
}

func (s *MetricsService) SaveWhenBody(metricRequest requests.MetricsSaveRequest) (entry *models.Metrics, err error) {
	existingMetric, err := s.repository.FindByName(metricRequest.ID)
	if err != nil && !errors.Is(err, constants.ErrMetricNotFound) {
		return nil, err
	}

	var metric *models.Metrics
	switch metricRequest.MType {
	case constants.GaugeMetricType:
		metric, err = s.sw.ConvertRequest(GaugeConverter, existingMetric, metricRequest)
	case constants.CounterMetricType:
		metric, err = s.sw.ConvertRequest(CounterConverter, existingMetric, metricRequest)
	default:
		return nil, constants.ErrInvalidMetricType
	}
	if err != nil {
		return nil, err
	}
	return s.repository.Save(metric), nil
}

func (s *MetricsService) Show(metricName string) (metric *models.Metrics, err error) {
	return s.repository.FindByName(metricName)
}

func (s *MetricsService) Ping(ctx *gin.Context) error {
	return s.repository.Ping(ctx)
}

func gaugeStrVal(metric *models.Metrics) (val string) {
	return strconv.FormatFloat(*metric.Value, 'f', -1, 64)
}

func counterStrVal(metric *models.Metrics) (val string) {
	return strconv.FormatInt(*metric.Delta, 10)
}
