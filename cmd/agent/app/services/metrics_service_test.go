package services

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	someGaugeMetricVal   = 239.232
	someCounterMetricVal = int64(3455)
)

func TestMetricsService_CollectMetrics(t *testing.T) {
	metrics := map[string]*models.Metrics{
		"gauge1": models.NewGaugeMetric("gauge1", 100.5),
		"gauge2": models.NewGaugeMetric("gauge2", 200.1),

		"counter1": models.NewCounterMetric("counter1", 10),
		"counter2": models.NewCounterMetric("counter2", 20),
	}

	// Create service with old metrics
	service := &MetricsService{
		metricsCollection: metrics,
	}

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

func BenchmarkMetricsService_SendBatchMetrics(b *testing.B) {
	repo := repositories.NewMetricsRepository()
	cfg, err := config.Setup()
	assert.NoError(b, err)
	ctx := context.Background()

	someGaugeMetric := &models.Metrics{
		ID:    "someGaugeMetric",
		MType: constants.GaugeMetricType,
		Value: &someGaugeMetricVal,
	}
	someCounterMetric := &models.Metrics{
		ID:    "someCounterMetric",
		MType: constants.CounterMetricType,
		Delta: &someCounterMetricVal,
	}

	service := NewMetricsService(repo, cfg)
	b.Run("SendBatchMetrics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = service.SendBatchMetrics()
			assert.NoError(b, err)
		}
	})
	b.Run("SendMetricWithParams", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = service.SendMetricWithParams(ctx, someCounterMetric)
		}
	})
	b.Run("SendCompressedMetric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = service.SendCompressedMetric(ctx, someGaugeMetric)
		}
	})
	b.Run("SendMetric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = service.SendMetric(ctx, someCounterMetric)
		}
	})
}
