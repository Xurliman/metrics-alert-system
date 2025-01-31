package repositories

import (
	"context"
	"testing"

	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/stretchr/testify/assert"
)

var (
	someGaugeMetricVal   = 239.232
	someCounterMetricVal = int64(3455)
)

func BenchmarkMetricsRepository_Save(b *testing.B) {
	repo := NewMetricsRepository(make(map[string]*models.Metrics))
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

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := repo.Save(someGaugeMetric)
			assert.NoError(b, err)

			_, err = repo.Save(someCounterMetric)
			assert.NoError(b, err)
		}
	})
}

func BenchmarkMetricsRepository_InsertMany(b *testing.B) {
	metrics := []*models.Metrics{
		{
			ID:    "someGaugeMetric",
			MType: constants.GaugeMetricType,
			Value: &someGaugeMetricVal,
		},
		{
			ID:    "someCounterMetric",
			MType: constants.CounterMetricType,
			Delta: &someCounterMetricVal,
		},
	}
	repo := NewMetricsRepository(make(map[string]*models.Metrics))
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := repo.InsertMany(ctx, metrics)
			assert.NoError(b, err)
		}
	})
}
