package services

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"testing"
)

func TestMetricsService_CollectMetrics(t *testing.T) {
	oldMetrics := models.OldMetrics{
		Gauge: map[string]float64{
			"gauge1": 100.5,
			"gauge2": 200.1,
		},
		Counter: map[string]int64{
			"counter1": 10,
			"counter2": 20,
		},
	}

	// Create service with old metrics
	service := &MetricsService{
		oldMetricsCollection: oldMetrics,
		metricsCollection:    nil,
	}

	service.ConvertToMetrics()

	if len(service.metricsCollection) != 4 {
		t.Errorf("Expected 4 metrics, got %d", len(service.metricsCollection))
	}

	if metric, ok := service.metricsCollection["gauge1"]; ok {
		if metric.MType != constants.GaugeMetricType || *metric.Value != 100.5 {
			t.Errorf("Unexpected metric for gauge1: %+v", metric)
		}
	} else {
		t.Error("Metric 'gauge1' not found")
	}

	if metric, ok := service.metricsCollection["counter1"]; ok {
		if metric.MType != constants.CounterMetricType || *metric.Delta != 10 {
			t.Errorf("Unexpected metric for counter1: %+v", metric)
		}
	} else {
		t.Error("Metric 'counter1' not found")
	}
}
