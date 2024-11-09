package interfaces

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
)

type Switch interface {
	ConvertParams(converter Converter, existingMetric *models.Metrics, metricName, metricValue string) (metric *models.Metrics, err error)
	ConvertRequest(converter Converter, existingMetric *models.Metrics, metricRequest requests.MetricsSaveRequest) (metric *models.Metrics, err error)
}
