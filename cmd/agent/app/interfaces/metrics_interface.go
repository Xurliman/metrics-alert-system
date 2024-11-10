package interfaces

import "github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"

type MetricsController interface {
	SendMetrics()
	SendCompressedMetrics()
	SendMetricsWithParams()
	SendCompressedMetricsWithParams()
	CollectMetrics()
}

type MetricsService interface {
	CollectMetricValues()
	GetRequestURLs(address string) ([]string, error)
	GetRequestBodies() ([][]byte, error)
	GetCompressedRequestBodies() ([][]byte, error)
}

type MetricsRepository interface {
	GetRequestBody(metric *models.Metrics) ([]byte, error)
	GetRequestURL(metric *models.Metrics, address string) (string, error)
}
