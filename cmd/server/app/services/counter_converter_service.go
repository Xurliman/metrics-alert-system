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
		return &models.Metrics{
			ID:    metricName,
			MType: constants.CounterMetricType,
			Delta: &metricVal,
		}, nil
	default:
		if existingMetric.MType != constants.CounterMetricType {
			return nil, constants.ErrMetricExists
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
		return &models.Metrics{
			ID:    metricRequest.ID,
			MType: constants.CounterMetricType,
			Delta: metricRequest.Delta,
		}, nil
	default:
		if existingMetric.MType != constants.CounterMetricType {
			return nil, constants.ErrMetricExists
		}
		*existingMetric.Delta += *metricRequest.Delta
		return existingMetric, nil
	}
}
