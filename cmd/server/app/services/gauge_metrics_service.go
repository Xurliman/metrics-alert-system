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

type gaugeMetricsService struct {
	repository repositories.MetricsRepository
}

func NewGaugeMetricsService() interfaces.MetricsInterface {
	return &gaugeMetricsService{}
}

func (s gaugeMetricsService) Save(metricName, metricValue string, existingMetric *models.Metrics) error {
	if metricName == "" {
		return constants.ErrEmptyMetricName
	}

	metricVal, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return constants.ErrInvalidGaugeMetricValue
	}

	switch existingMetric {
	case nil:
		MetricsCollection[metricName] = &models.Metrics{
			ID:    metricName,
			MType: constants.GaugeMetricType,
			Value: &metricVal,
		}
	default:
		*existingMetric.Value = metricVal
	}
	return nil
}

func (s gaugeMetricsService) GetMetricValue(metricsName string) (metricsValue string, err error) {
	metric, err := s.FindByName(metricsName)
	if err != nil {
		return metricsValue, fmt.Errorf("type gauge|metrics %s not found", metricsName)
	}

	metricsValue = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	return metricsValue, nil
}

func (s gaugeMetricsService) SaveBody(metricRequest requests.MetricsSaveRequest, existingMetric *models.Metrics) (entry *models.Metrics, err error) {
	if metricRequest.Value == nil {
		return nil, constants.ErrInvalidGaugeMetricValue
	}

	switch existingMetric {
	case nil:
		entry = &models.Metrics{
			ID:    metricRequest.ID,
			MType: constants.GaugeMetricType,
			Value: metricRequest.Value,
		}
		MetricsCollection[metricRequest.ID] = entry
	default:
		*existingMetric.Value = *metricRequest.Value
		entry = existingMetric
	}

	return entry, nil
}

func (s gaugeMetricsService) FindByName(metricName string) (*models.Metrics, error) {
	if m, ok := MetricsCollection[metricName]; ok {
		return m, nil
	}
	return nil, constants.ErrMetricNotFound
}
