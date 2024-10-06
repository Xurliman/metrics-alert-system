package interfaces

import "github.com/gin-gonic/gin"

type MetricsControllerInterface interface {
	Index(ctx *gin.Context)
	Show(ctx *gin.Context)
	Validate(ctx *gin.Context)
}

type MetricsServiceInterface interface {
	FindByName(metricsType string, metricsName string) (metricsValue string, err error)
	GetAll() map[string]string
}
