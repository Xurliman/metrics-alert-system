package routes

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/controllers"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/middlewares"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.New()
	r.LoadHTMLFiles("./cmd/server/public/templates/metrics-all.html")

	logging := middlewares.NewLoggingMiddleware()
	metricsService := services.NewMetricsService()
	metricsController := controllers.NewMetricsController(metricsService)

	r.GET("/", logging.Handle(metricsController.Index))
	r.GET("/value/:type/:name/", logging.Handle(metricsController.Show))
	r.POST("/update/:type/:name/:value", logging.Handle(metricsController.Update))
	return r
}
