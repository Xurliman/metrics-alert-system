package interfaces

import (
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
)

type MetricsController interface {
	SendMetrics() (err error)
	SendCompressedMetrics() (err error)
	SendMetricsWithParams() (err error)
	SendCompressedMetricsWithParams() (err error)
	SendBatchMetrics() (err error)
	CollectMetrics()
}

type MetricsService interface {
	CollectMetricValues()
	GetRequestURLs(address string) ([]string, error)
	GetRequestBodies() ([][]byte, error)
	GetCompressedRequestBodies() ([][]byte, error)
	GetCompressedRequestBody() ([]byte, error)
}

type MetricsRepository interface {
	GetRequestBody(metric *models.Metrics) ([]byte, error)
	GetRequestURL(metric *models.Metrics, address string) (string, error)
	GetPlainRequest(metric *models.Metrics) (*requests.MetricsRequest, error)
}
