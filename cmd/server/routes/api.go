package routes

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/controllers"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.LoadHTMLFiles("./cmd/server/public/templates/metrics-all.html")
	metricsService := services.NewMetricsService()
	metricsController := controllers.NewMetricsController(metricsService)
	r.GET("/", metricsController.Index)
	r.GET("/value/:type/:name", metricsController.Show)
	r.POST("/update/:type/:name/:value", metricsController.Validate)
	return r
}
