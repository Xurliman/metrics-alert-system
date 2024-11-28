package interfaces

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
)

type MetricsController interface {
	SendMetrics(ctx context.Context) (err error)
	SendCompressedMetrics(ctx context.Context) (err error)
	SendMetricsWithParams(ctx context.Context) (err error)
	SendBatchMetrics(ctx context.Context) (err error)
	CollectMetrics()
}

type MetricsService interface {
	CollectMetricValues()

	Generator(ctx context.Context) chan *models.Metrics
	ByteTransformer(ctx context.Context, inputCh chan models.Result) chan models.Result
	URLConstructor(ctx context.Context, inputCh chan *models.Metrics, address string) chan models.Result
	RequestConstructor(ctx context.Context, inputCh chan *models.Metrics) chan models.Result
	RequestCompressor(ctx context.Context, inputCh chan models.Result) chan models.Result
}

type MetricsRepository interface {
	GetRequestURL(metric *models.Metrics, address string) (string, error)
	GetPlainRequest(metric *models.Metrics) (*requests.MetricsRequest, error)
	GetBytes(metricRequest *requests.MetricsRequest) ([]byte, error)
}
