package interfaces

import (
	"context"

	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
)

type MetricsServiceInterface interface {
	List() (map[string]string, error)
	GetMetricValue(metricType, metricName string) (metricValue string, err error)
	Show(metricName string) (entry *models.Metrics, err error)
	SaveWhenParams(metricType, metricName, metricValue string) error
	SaveWhenBody(metricRequest requests.MetricsSaveRequest) (entry *models.Metrics, err error)
	Ping(ctx context.Context) error
	SaveMany(ctx context.Context, request []requests.MetricsSaveRequest) error
}
