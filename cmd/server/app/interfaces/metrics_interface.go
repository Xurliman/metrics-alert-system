package interfaces

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/gin-gonic/gin"
)

type MetricsRepositoryInterface interface {
	Ping(ctx *gin.Context) error
}

type MetricsInterface interface {
	FindByName(metricName string) (*models.Metrics, error)
	GetMetricValue(metricsName string) (metricsValue string, err error)
	Save(metricName, metricValue string, existingMetric *models.Metrics) error
	SaveBody(metricRequest requests.MetricsSaveRequest, existingMetric *models.Metrics) (entry *models.Metrics, err error)
}

type MetricsServiceInterface interface {
	List() map[string]string
	GetMetricValue(metric MetricsInterface, metricName string) (metricValue string, err error)
	Show(metric MetricsInterface, metricName string) (entry *models.Metrics, err error)
	SaveWhenParams(metric MetricsInterface, metricName, metricValue string) error
	SaveWhenBody(metric MetricsInterface, metricRequest requests.MetricsSaveRequest) (entry *models.Metrics, err error)
	Ping(ctx *gin.Context) error
}

type MetricsControllerInterface interface {
	List(ctx *gin.Context)
	Show(ctx *gin.Context)
	Save(ctx *gin.Context)
	SaveBody(ctx *gin.Context)
	ShowBody(ctx *gin.Context)
	Ping(ctx *gin.Context)
}
