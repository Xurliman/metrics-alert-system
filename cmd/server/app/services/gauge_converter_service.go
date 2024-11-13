package services

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"strconv"
)

type gaugeConverterService struct{}

func (s gaugeConverterService) ParamsToMetric(existingMetric *models.Metrics, metricName, metricValue string) (metric *models.Metrics, err error) {
	metricVal, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return nil, constants.ErrInvalidGaugeMetricValue
	}

	switch existingMetric {
	case nil:
		return defaultGaugeMetric(metricName, &metricVal)
	default:
		if existingMetric.MType != constants.GaugeMetricType {
			return nil, constants.ErrMetricExists
		}

		if existingMetric.Value == nil {
			return defaultGaugeMetric(metricName, &metricVal)
		}

		*existingMetric.Value = metricVal
		return existingMetric, nil
	}
}

func (s gaugeConverterService) RequestToMetric(existingMetric *models.Metrics, metricRequest requests.MetricsSaveRequest) (metric *models.Metrics, err error) {
	if metricRequest.Value == nil {
		return nil, constants.ErrInvalidGaugeMetricValue
	}
	switch existingMetric {
	case nil:
		return defaultGaugeMetric(metricRequest.ID, metricRequest.Value)
	default:
		if existingMetric.MType != constants.GaugeMetricType {
			return nil, constants.ErrMetricExists
		}

		if existingMetric.Value == nil {
			return defaultGaugeMetric(metricRequest.ID, metricRequest.Value)
		}

		*existingMetric.Value = *metricRequest.Value
		return existingMetric, nil
	}
}

func defaultGaugeMetric(name string, value *float64) (metric *models.Metrics, err error) {
	return &models.Metrics{
		ID:    name,
		MType: constants.GaugeMetricType,
		Value: value,
	}, nil
}
