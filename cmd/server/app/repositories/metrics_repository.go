package repositories

import (
	"context"
	"database/sql"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"go.uber.org/zap"
	"runtime"
	"sync"
)

type MetricsRepository struct {
	metricsCollection map[string]*models.Metrics
	db                *sql.DB
	mu                sync.RWMutex
}

func NewMetricsRepository(shouldRestore bool, archiveService interfaces.ArchiveServiceInterface) interfaces.MetricsRepositoryInterface {
	if !shouldRestore {
		return &MetricsRepository{
			metricsCollection: defaultMetrics(),
			db:                DB,
			mu:                sync.RWMutex{},
		}
	}

	metrics, err := archiveService.Load()
	if err != nil {
		utils.Logger.Error("error loading metrics from file", zap.Error(err))
		return &MetricsRepository{
			metricsCollection: defaultMetrics(),
			db:                DB,
			mu:                sync.RWMutex{},
		}
	}
	return &MetricsRepository{
		metricsCollection: metrics,
		db:                DB,
		mu:                sync.RWMutex{},
	}
}

func (r *MetricsRepository) Save(metric *models.Metrics) *models.Metrics {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metricsCollection[metric.ID] = metric
	return metric
}

func (r *MetricsRepository) FindByName(metricName string) (*models.Metrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if m, ok := r.metricsCollection[metricName]; ok {
		return m, nil
	}
	return nil, constants.ErrMetricNotFound
}

func (r *MetricsRepository) List() map[string]*models.Metrics {
	return r.metricsCollection
}

func (r *MetricsRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func defaultMetrics() map[string]*models.Metrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	gaugeMetricExample := float64(memStats.Alloc)
	var counterMetricExample int64
	defaultMap := make(map[string]*models.Metrics)
	defaultMap["Alloc"] = &models.Metrics{
		ID:    "Alloc",
		MType: constants.GaugeMetricType,
		Value: &gaugeMetricExample,
	}
	defaultMap["PollCount"] = &models.Metrics{
		ID:    "PollCount",
		MType: constants.CounterMetricType,
		Delta: &counterMetricExample,
	}
	return defaultMap
}
