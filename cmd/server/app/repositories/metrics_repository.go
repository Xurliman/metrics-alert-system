package repositories

import (
	"database/sql"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/gin-gonic/gin"
)

type MetricsRepository struct {
	metricsCollection map[string]*models.Metrics
	db                *sql.DB
}

func NewMetricsRepository(db *sql.DB) interfaces.MetricsRepositoryInterface {
	return &MetricsRepository{
		metricsCollection: make(map[string]*models.Metrics),
		db:                db,
	}
}

func (r *MetricsRepository) Ping(ctx *gin.Context) error {
	return r.db.PingContext(ctx)
}
