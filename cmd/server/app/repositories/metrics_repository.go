package repositories

import (
	"context"
	"database/sql"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"sync"
)

type MetricsRepository struct {
	metricsCollection map[string]*models.Metrics
	db                *sql.DB
	mu                sync.RWMutex
}

func NewMetricsRepository(metrics map[string]*models.Metrics) interfaces.MetricsRepositoryInterface {
	return &MetricsRepository{
		metricsCollection: metrics,
		db:                DB,
		mu:                sync.RWMutex{},
	}
}

func (r *MetricsRepository) Save(metric *models.Metrics) (*models.Metrics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metricsCollection[metric.ID] = metric
	return metric, nil
}

func (r *MetricsRepository) FindByName(metricName string) (*models.Metrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if m, ok := r.metricsCollection[metricName]; ok {
		return m, nil
	}
	return nil, constants.ErrMetricNotFound
}

func (r *MetricsRepository) List() (map[string]*models.Metrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.metricsCollection, nil
}

func (r *MetricsRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *MetricsRepository) InsertMany(ctx context.Context, metrics []*models.Metrics) (err error) {
	for _, metric := range metrics {
		_, err = r.Save(metric)
		if err != nil {
			return err
		}
	}
	return nil
}
