package interfaces

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/gin-gonic/gin"
)

type MetricsRepositoryInterface interface {
	Ping(ctx *gin.Context) error
	Save(metric *models.Metrics) *models.Metrics
	FindByName(metricName string) (*models.Metrics, error)
	List() map[string]*models.Metrics
}
