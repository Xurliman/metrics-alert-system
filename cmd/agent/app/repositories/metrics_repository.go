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

func (r *MetricsRepository) GetBytes(metricRequest *requests.MetricsRequest) ([]byte, error) {
	requestBytes, err := json.Marshal(metricRequest)
	if err != nil {
		return nil, err
	}
	return requestBytes, nil
}

func (r *MetricsRepository) GetRequestURL(metric *models.Metrics, address string) (string, error) {
	var (
		value string
		err   error
	)

	switch metric.MType {
	case constants.GaugeMetricType:
		value, err = metric.GetValue()
		if err != nil {
			return "", err
		}
	case constants.CounterMetricType:
		value, err = metric.GetDelta()
		if err != nil {
			return "", err
		}
	default:
		return "", constants.ErrInvalidMetricType
	}

	return fmt.Sprintf("http://%s/update/%s/%s/%v",
		address,
		metric.MType,
		metric.ID,
		value,
	), nil
}

func (r *MetricsRepository) GetPlainRequest(metric *models.Metrics) (request *requests.MetricsRequest, err error) {
	switch metric.MType {
	case constants.GaugeMetricType:
		request, err = metric.ToGaugeRequest()
		if err != nil {
			return nil, err
		}
	case constants.CounterMetricType:
		request, err = metric.ToCounterRequest()
		if err != nil {
			return nil, err
		}
	default:
		return nil, constants.ErrInvalidMetricType
	}

	return request, nil
}
