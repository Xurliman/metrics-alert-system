package models

import (
	"database/sql"
	"github.com/google/uuid"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
}

type DBMetrics struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	MetricType string          `json:"type"`
	Value      sql.NullFloat64 `json:"value,omitempty"`
	Delta      sql.NullInt64   `json:"delta,omitempty"`
}

func (m *Metrics) Equals(someMetric *Metrics) bool {
	if m == nil || someMetric == nil {
		return false
	}

	if m.ID != someMetric.ID || m.MType != someMetric.MType {
		return false
	}

	if (m.Value == nil && someMetric.Value != nil) || (m.Value != nil && someMetric.Value == nil) || (m.Value != nil && someMetric.Value != nil && *m.Value != *someMetric.Value) {
		return false
	}

	if (m.Delta == nil && someMetric.Delta != nil) || (m.Delta != nil && someMetric.Delta == nil) || (m.Delta != nil && someMetric.Delta != nil && *m.Delta != *someMetric.Delta) {
		return false
	}

	return true
}

func (dm *DBMetrics) ToModel() *Metrics {
	if dm.Value.Valid {
		return &Metrics{
			ID:    dm.Name,
			MType: dm.MetricType,
			Value: &dm.Value.Float64,
		}
	}
	if dm.Delta.Valid {
		return &Metrics{
			ID:    dm.Name,
			MType: dm.MetricType,
			Delta: &dm.Delta.Int64,
		}
	}
	return nil
}
