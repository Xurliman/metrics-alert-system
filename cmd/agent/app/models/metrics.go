// Package models describe metrics structures
package models

import (
	"strconv"

	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
)

type OldMetrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type Metrics struct {
	ID    string
	MType string
	Value *float64
	Delta *int64
}

func NewCounterMetric(name string, delta int64) *Metrics {
	return &Metrics{
		ID:    name,
		MType: constants.CounterMetricType,
		Delta: &delta,
	}
}

func NewGaugeMetric(name string, value float64) *Metrics {
	return &Metrics{
		ID:    name,
		MType: constants.GaugeMetricType,
		Value: &value,
	}
}

func (m *Metrics) ToGaugeRequest() (*requests.MetricsRequest, error) {
	if m.MType != constants.GaugeMetricType {
		return nil, constants.ErrInvalidMetricType
	}
	return &requests.MetricsRequest{
		ID:    m.ID,
		MType: constants.GaugeMetricType,
		Value: m.Value,
	}, nil
}

func (m *Metrics) ToCounterRequest() (*requests.MetricsRequest, error) {
	if m.MType != constants.CounterMetricType {
		return nil, constants.ErrInvalidMetricType
	}
	return &requests.MetricsRequest{
		ID:    m.ID,
		MType: constants.CounterMetricType,
		Delta: m.Delta,
	}, nil
}

func (m *Metrics) GetValueString() string {
	return strconv.FormatFloat(m.GetValue(), 'f', -1, 64)
}

func (m *Metrics) GetDeltaString() string {
	return strconv.FormatInt(m.GetDelta(), 10)
}

func (m *Metrics) GetValue() float64 {
	if m.Value == nil {
		return 0
	}
	return *m.Value
}

func (m *Metrics) GetDelta() int64 {
	if m.Delta == nil {
		return 0
	}
	return *m.Delta
}
