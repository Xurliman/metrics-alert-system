package models

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
}

func (m *Metrics) Equals(someMetric *Metrics) bool {
	flag := false
	if m.ID == someMetric.ID {
		flag = true
	}
	if m.MType == someMetric.MType {
		flag = true
	}
	if m.Value == someMetric.Value {
		flag = true
	}
	if m.Delta == someMetric.Delta {
		flag = true
	}
	return flag
}
