package routes

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/controllers"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/middlewares"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	compression := middlewares.NewCompressingMiddleware()
	decompression := middlewares.NewDecompressingMiddleware()
	r := gin.New()
	r.LoadHTMLFiles("./cmd/server/public/templates/metrics-all.html")

	logging := middlewares.NewLoggingMiddleware()
	metricsService := services.NewMetricsService()
	metricsController := controllers.NewMetricsController(metricsService)

	r.GET("/", decompression.Handle(logging.Handle(compression.Handle(metricsController.List))))
	r.GET("/value/:type/:name", decompression.Handle(logging.Handle(compression.Handle(metricsController.Show))))
	r.POST("/value/", decompression.Handle(logging.Handle(compression.Handle(metricsController.ShowBody))))
	r.POST("/update/:type/:name/:value", decompression.Handle(logging.Handle(compression.Handle(metricsController.Save))))
	r.POST("/update/", decompression.Handle(logging.Handle(compression.Handle(metricsController.SaveBody))))
	return r
}
