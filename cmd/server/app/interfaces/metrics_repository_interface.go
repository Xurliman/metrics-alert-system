package interfaces

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
)

type MetricsRepositoryInterface interface {
	Ping(ctx context.Context) error
	Save(metric *models.Metrics) *models.Metrics
	FindByName(metricName string) (*models.Metrics, error)
	List() map[string]*models.Metrics
}
