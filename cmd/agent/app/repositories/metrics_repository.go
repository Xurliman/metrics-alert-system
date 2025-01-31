package repositories

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
)

type MetricsRepository struct {
	metricsCollection map[string]*models.Metrics
	mu                sync.RWMutex
}

func NewMetricsRepository() interfaces.MetricsRepository {
	return &MetricsRepository{
		metricsCollection: make(map[string]*models.Metrics),
	}
}

func (r *MetricsRepository) GetAll() map[string]*models.Metrics {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copiedMetrics := make(map[string]*models.Metrics, len(r.metricsCollection))
	for k, v := range r.metricsCollection {
		copiedMetrics[k] = v
	}
	return copiedMetrics
}

func (r *MetricsRepository) SaveAll(metrics []*models.Metrics) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, m := range metrics {
		if err := r.save(m); err != nil {
			return err
		}
	}
	return nil
}

func (r *MetricsRepository) Save(metric *models.Metrics) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.save(metric); err != nil {
		return err
	}

	return nil
}

func (r *MetricsRepository) save(metric *models.Metrics) error {
	if metric == nil {
		return constants.ErrInvalidMetric
	}

	m, exists := r.metricsCollection[metric.ID]
	if !exists {
		r.metricsCollection[metric.ID] = metric
		return nil
	}

	switch metric.MType {
	case constants.GaugeMetricType:
		newVal := metric.GetValue() + m.GetValue()
		r.metricsCollection[metric.ID] = models.NewGaugeMetric(metric.ID, newVal)
	case constants.CounterMetricType:
		newVal := metric.GetDelta() + m.GetDelta()
		r.metricsCollection[metric.ID] = models.NewCounterMetric(metric.ID, newVal)
	default:
		return constants.ErrInvalidMetricType
	}

	return nil
}

func (r *MetricsRepository) GetRequestURL(metric *models.Metrics) (value string, err error) {
	switch metric.MType {
	case constants.GaugeMetricType:
		value = metric.GetValueString()
	case constants.CounterMetricType:
		value = metric.GetDeltaString()
	default:
		return "", constants.ErrInvalidMetricType
	}

	return fmt.Sprintf("/%s/%s/%v",
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

func (r *MetricsRepository) GetRequestBody(metric *models.Metrics) ([]byte, error) {
	switch metric.MType {
	case constants.GaugeMetricType:
		request, err := metric.ToGaugeRequest()
		if err != nil {
			return nil, err
		}
		requestBytes, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}

		return requestBytes, nil
	case constants.CounterMetricType:
		request, err := metric.ToCounterRequest()
		if err != nil {
			return nil, err
		}

		requestBytes, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}

		return requestBytes, nil
	default:
		return nil, constants.ErrInvalidMetricType
	}
}
