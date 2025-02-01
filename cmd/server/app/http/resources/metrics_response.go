// Package resources define structs used to wrap responses
package resources

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
)

type MetricsResponse struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func ToResponse(metric *models.Metrics) (*MetricsResponse, error) {
	switch metric.MType {
	case constants.GaugeMetricType:
		return &MetricsResponse{
			ID:    metric.ID,
			MType: metric.MType,
			Value: metric.Value,
		}, nil
	case constants.CounterMetricType:
		return &MetricsResponse{
			ID:    metric.ID,
			MType: metric.MType,
			Delta: metric.Delta,
		}, nil
	default:
		return nil, constants.ErrInvalidMetricType
	}
}
