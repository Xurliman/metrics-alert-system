package models

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
)

type OldMetrics struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
}

func (m *Metrics) ToGaugeRequest() requests.MetricsRequest {
	return requests.MetricsRequest{
		ID:    m.ID,
		MType: constants.GaugeMetricType,
		Value: m.Value,
	}
}

func (m *Metrics) ToCounterRequest() requests.MetricsRequest {
	return requests.MetricsRequest{
		ID:    m.ID,
		MType: constants.CounterMetricType,
		Delta: m.Delta,
	}
}
