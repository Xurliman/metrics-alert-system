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
	GetRequestBodies() ([][]byte, error)
	GetCompressedRequestBodies() ([][]byte, error)
	GetRequestUrls(address string) ([]string, error)
}

type MetricsRepository interface {
	GetRequestBody(metric *models.Metrics) ([]byte, error)
	GetRequestUrl(metric *models.Metrics, address string) (string, error)
}
