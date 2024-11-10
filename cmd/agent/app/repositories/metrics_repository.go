package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
)

type MetricsRepository struct{}

func NewMetricsRepository() interfaces.MetricsRepository {
	return &MetricsRepository{}
}

func (r *MetricsRepository) GetRequestBody(metric *models.Metrics) ([]byte, error) {
	switch metric.MType {
	case constants.GaugeMetricType:
		request, err := json.Marshal(metric.ToGaugeRequest())
		if err != nil {
			return nil, err
		}
		return request, err
	case constants.CounterMetricType:
		request, err := json.Marshal(metric.ToCounterRequest())
		if err != nil {
			return nil, err
		}
		return request, nil
	default:
		return nil, constants.ErrInvalidMetricType
	}
}

func (r *MetricsRepository) GetRequestURL(metric *models.Metrics, address string) (string, error) {
	switch metric.MType {
	case constants.GaugeMetricType:
		return fmt.Sprintf("http://%s/update/%s/%s/%f",
			address,
			constants.GaugeMetricType,
			metric.ID,
			*metric.Value,
		), nil
	case constants.CounterMetricType:
		return fmt.Sprintf("http://%s/update/%s/%s/%v",
			address,
			constants.CounterMetricType,
			metric.ID,
			*metric.Delta,
		), nil
	default:
		return "", constants.ErrInvalidMetricType
	}
}

func (r *MetricsRepository) GetPlainRequest(metric *models.Metrics) requests.MetricsRequest {
	var request requests.MetricsRequest

	switch metric.MType {
	case constants.GaugeMetricType:
		request = metric.ToGaugeRequest()
	case constants.CounterMetricType:
		request = metric.ToCounterRequest()
	}

	return request
}
