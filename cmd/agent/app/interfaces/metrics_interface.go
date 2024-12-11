package interfaces

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
)

type MetricsController interface {
	Report(ctx context.Context)
	Poll(ctx context.Context)
	Run(ctx context.Context)
}

type MetricsService interface {
	CollectMetricValues() error
	GetAll() map[string]*models.Metrics
	SendBatchMetrics() (err error)
	SendMetric(ctx context.Context, metric *models.Metrics) (errs error)
	SendCompressedMetric(ctx context.Context, metric *models.Metrics) (errs error)
	SendMetricWithParams(ctx context.Context, metric *models.Metrics) (err error)
}

type MetricsRepository interface {
	GetRequestURL(metric *models.Metrics) (string, error)
	GetPlainRequest(metric *models.Metrics) (*requests.MetricsRequest, error)
	GetRequestBody(metric *models.Metrics) ([]byte, error)
	SaveAll(metrics []*models.Metrics) error
	Save(metric *models.Metrics) error
	GetAll() map[string]*models.Metrics
}
