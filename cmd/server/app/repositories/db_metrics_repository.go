package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type DBMetricsRepository struct {
	db *sql.DB
}

func NewDBMetricsRepository() interfaces.MetricsRepositoryInterface {
	return &DBMetricsRepository{
		db: DB,
	}
}

func (r *DBMetricsRepository) Save(metric *models.Metrics) *models.Metrics {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dbMetric, err := r.Find(ctx, metric.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if dbMetric != nil {
		err = r.Update(ctx, dbMetric.Id, metric)
	} else {
		err = r.Insert(ctx, metric)
	}
	if err != nil {
		return nil
	}

	return metric
}

func (r *DBMetricsRepository) FindByName(metricName string) (*models.Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dbMetric, err := r.Find(ctx, metricName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrMetricNotFound
		}
		return nil, err
	}

	return dbMetric.ToModel(), nil
}

func (r *DBMetricsRepository) List() map[string]*models.Metrics {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	metrics := make(map[string]*models.Metrics)

	dbMetrics := r.FindAll(ctx)
	for _, dbMetric := range dbMetrics {
		if dbMetric != nil {
			metrics[dbMetric.Name] = dbMetric.ToModel()
		}
	}

	return metrics
}

func (r *DBMetricsRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *DBMetricsRepository) Find(ctx context.Context, metricName string) (*models.DBMetrics, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, metric_type, value, delta FROM metrics WHERE name = $1`, metricName)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var dbMetric models.DBMetrics
	err := row.Scan(&dbMetric.Id, &dbMetric.Name, &dbMetric.MetricType, &dbMetric.Value, &dbMetric.Delta)
	if err != nil {
		return nil, err
	}

	return &dbMetric, nil
}

func (r *DBMetricsRepository) FindAll(ctx context.Context) []*models.DBMetrics {
	var dbMetrics []*models.DBMetrics
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, metric_type, value, delta FROM metrics")
	if err != nil {
		return nil
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			utils.Logger.Error("error closing sql.rows", zap.Error(err))
		}
	}(rows)

	for rows.Next() {
		var dbMetric models.DBMetrics
		err = rows.Scan(&dbMetric.Id, &dbMetric.Name, &dbMetric.MetricType, &dbMetric.Value, &dbMetric.Delta)
		if err != nil {
			return nil
		}

		dbMetrics = append(dbMetrics, &dbMetric)
	}

	if err = rows.Err(); err != nil {
		return nil
	}

	return dbMetrics
}

func (r *DBMetricsRepository) Insert(ctx context.Context, metric *models.Metrics) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO metrics (name, metric_type, value, delta) VALUES($1, $2, $3, $4)`,
		metric.ID,
		metric.MType,
		metric.Value,
		metric.Delta,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBMetricsRepository) Update(ctx context.Context, metricId uuid.UUID, metric *models.Metrics) error {
	_, err := r.db.ExecContext(ctx, `UPDATE metrics SET name = $1, metric_type = $2, value = $3, delta = $4 WHERE id = $5`,
		metric.ID,
		metric.MType,
		metric.Value,
		metric.Delta,
		metricId.String(),
	)
	if err != nil {
		return err
	}
	return nil
}
