package services

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"strconv"
)

type counterConverterService struct{}

func (s counterConverterService) ParamsToMetric(existingMetric *models.Metrics, metricName, metricValue string) (metric *models.Metrics, err error) {
	metricVal, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return nil, constants.ErrInvalidCounterMetricValue
	}

	switch existingMetric {
	case nil:
		return defaultCounterMetric(metricName, &metricVal)
	default:
		if existingMetric.MType != constants.CounterMetricType {
			return nil, constants.ErrMetricExists
		}

		if existingMetric.Delta == nil {
			return defaultCounterMetric(metricName, &metricVal)
		}

		*existingMetric.Delta += metricVal
		return existingMetric, nil
	}
}

func (s counterConverterService) RequestToMetric(existingMetric *models.Metrics, metricRequest requests.MetricsSaveRequest) (*models.Metrics, error) {
	if metricRequest.Delta == nil {
		return nil, constants.ErrInvalidCounterMetricValue
	}

	switch existingMetric {
	case nil:
		return defaultCounterMetric(metricRequest.ID, metricRequest.Delta)
	default:
		if existingMetric.MType != constants.CounterMetricType {
			return nil, constants.ErrMetricExists
		}

		if existingMetric.Delta == nil {
			return defaultCounterMetric(metricRequest.ID, metricRequest.Delta)
		}

		*existingMetric.Delta += *metricRequest.Delta
		return existingMetric, nil
	}
}

func defaultCounterMetric(name string, delta *int64) (*models.Metrics, error) {
	return &models.Metrics{
		ID:    name,
		MType: constants.CounterMetricType,
		Delta: delta,
	}, nil
}
