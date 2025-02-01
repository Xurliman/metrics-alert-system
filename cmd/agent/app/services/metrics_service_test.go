package services

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/repositories"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/config"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"go.uber.org/zap"
	"testing"

	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
)

var (
	someGaugeMetricVal   = 239.232
	someCounterMetricVal = int64(3455)
	service              interfaces.MetricsService
)

func TestMain(m *testing.M) {
	repo := repositories.NewMetricsRepository()
	cfg, err := config.Setup()
	if err != nil {
		log.Fatal("failed to setup configuration for testing", zap.Error(err))
	}
	service = NewMetricsService(repo, cfg)
}

func TestMetricsService_CollectMetrics(t *testing.T) {
	metrics := map[string]*models.Metrics{
		"gauge1": models.NewGaugeMetric("gauge1", 100.5),
		"gauge2": models.NewGaugeMetric("gauge2", 200.1),

		"counter1": models.NewCounterMetric("counter1", 10),
		"counter2": models.NewCounterMetric("counter2", 20),
	}

	// Create service with old metrics
	serviceWithMetrics := &MetricsService{
		metricsCollection: metrics,
	}

	if len(serviceWithMetrics.metricsCollection) != 4 {
		t.Errorf("Expected 4 metrics, got %d", len(serviceWithMetrics.metricsCollection))
	}

	if metric, ok := serviceWithMetrics.metricsCollection["gauge1"]; ok {
		if metric.MType != constants.GaugeMetricType || *metric.Value != 100.5 {
			t.Errorf("Unexpected metric for gauge1: %+v", metric)
		}
	} else {
		t.Error("Metric 'gauge1' not found")
	}

	if metric, ok := serviceWithMetrics.metricsCollection["counter1"]; ok {
		if metric.MType != constants.CounterMetricType || *metric.Delta != 10 {
			t.Errorf("Unexpected metric for counter1: %+v", metric)
		}
	} else {
		t.Error("Metric 'counter1' not found")
	}
}
