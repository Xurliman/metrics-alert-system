package interfaces

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
)

type Converter interface {
	ParamsToMetric(existingMetric *models.Metrics, metricName, metricValue string) (metric *models.Metrics, err error)
	RequestToMetric(existingMetric *models.Metrics, metricRequest requests.MetricsSaveRequest) (metric *models.Metrics, err error)
}
