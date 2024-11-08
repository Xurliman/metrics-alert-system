package services

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/repositories"
	"strconv"
)

type counterMetricsService struct {
	repository repositories.MetricsRepository
}

func NewCounterMetricsService() interfaces.MetricsInterface {
	return &counterMetricsService{}
}

func (s counterMetricsService) Save(metricName, metricValue string, existingMetric *models.Metrics) error {
	if metricName == "" {
		return constants.ErrEmptyMetricName
	}

	metricVal, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return constants.ErrInvalidCounterMetricValue
	}

	switch existingMetric {
	case nil:
		MetricsCollection[metricName] = &models.Metrics{
			ID:    metricName,
			MType: constants.CounterMetricType,
			Delta: &metricVal,
		}
	default:
		*existingMetric.Delta += metricVal
	}
	return nil
}

func (s counterMetricsService) SaveBody(metricRequest requests.MetricsSaveRequest, existingMetric *models.Metrics) (entry *models.Metrics, err error) {
	if metricRequest.Delta == nil {
		return nil, constants.ErrInvalidCounterMetricValue
	}
	switch existingMetric {
	case nil:
		entry = &models.Metrics{
			ID:    metricRequest.ID,
			MType: constants.CounterMetricType,
			Delta: metricRequest.Delta,
		}
		MetricsCollection[metricRequest.ID] = entry
	default:
		*existingMetric.Delta += *metricRequest.Delta
		entry = existingMetric
	}

	return entry, nil
}

func (s counterMetricsService) GetMetricValue(metricsName string) (metricsValue string, err error) {
	metric, err := s.FindByName(metricsName)
	if err != nil {
		return metricsValue, fmt.Errorf("type counter|metrics %s not found", metricsName)
	}

	metricsValue = strconv.FormatInt(*metric.Delta, 10)
	return metricsValue, nil
}

func (s counterMetricsService) FindByName(metricName string) (*models.Metrics, error) {
	if m, ok := MetricsCollection[metricName]; ok {
		return m, nil
	}
	return nil, constants.ErrMetricNotFound
}
