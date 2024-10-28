package interfaces

import "github.com/gin-gonic/gin"

type MetricsControllerInterface interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Update(ctx *gin.Context)
}

type MetricsServiceInterface interface {
	GetAll() map[string]string
	FindMetricByName(metricsType string, metricsName string) (metricsValue string, err error)
	FindCounterMetric(metricsName string) (metricsValue string, err error)
	FindGaugeMetric(metricsName string) (metricsValue string, err error)
	SaveGaugeMetric(metricsName string, metricsValue string) error
	SaveCounterMetric(metricsName string, metricsValue string) error
	SaveMetric(metricsType string, metricsName string, metricsValue string) error
}
