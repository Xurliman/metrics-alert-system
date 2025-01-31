package repositories

import (
	"context"
	"database/sql"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	dsn := "host=localhost port=5432 user=postgres password=kali dbname=metrics sslmode=disable"
	testDB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}

	testDB.SetMaxIdleConns(5)
	testDB.SetMaxOpenConns(10)
	testDB.SetConnMaxLifetime(5 * time.Minute)
	testDB.SetConnMaxIdleTime(2 * time.Minute)
	DB = testDB

	// Run tests
	code := m.Run()

	// Cleanup
	_, _ = testDB.ExecContext(context.Background(), "TRUNCATE metrics RESTART IDENTITY")
	_ = testDB.Close()

	os.Exit(code)
}

func BenchmarkDBMetricsRepository_Save(b *testing.B) {
	repo := NewDBMetricsRepository()
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

func BenchmarkDBMetricsRepository_InsertMany(b *testing.B) {
	repo := NewDBMetricsRepository()
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
	ctx := context.Background()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := repo.InsertMany(ctx, metrics)
			assert.NoError(b, err)
		}
	})
}
