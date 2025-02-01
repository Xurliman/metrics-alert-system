package services

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkMetricsService_SendBatchMetrics(b *testing.B) {
	ctx := context.Background()
	var err error

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
