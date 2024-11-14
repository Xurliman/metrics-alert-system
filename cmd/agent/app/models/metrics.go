package models

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
	"strconv"
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

func (m *Metrics) GetValue() (string, error) {
	if m.Value == nil {
		return "", constants.ErrInvalidGaugeMetricValue
	}
	return strconv.FormatFloat(*m.Value, 'f', -1, 64), nil
}

func (m *Metrics) GetDelta() (string, error) {
	if m.Delta == nil {
		return "", constants.ErrInvalidCounterMetricValue
	}
	return strconv.FormatInt(*m.Delta, 10), nil
}
