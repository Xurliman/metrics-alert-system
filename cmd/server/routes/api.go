package routes

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/controllers"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/middlewares"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(metricsRepository interfaces.MetricsRepositoryInterface) *gin.Engine {
	decompression := middlewares.NewDecompressingMiddleware()
	logging := middlewares.NewLoggingMiddleware()
	compression := middlewares.NewCompressingMiddleware()
	r := gin.New()
	r.LoadHTMLFiles("./cmd/server/public/templates/metrics-all.html")

	metricsService := services.NewMetricsService(metricsRepository, services.NewSwitchService())
	metricsController := controllers.NewMetricsController(metricsService)

	r.GET("/", decompression.Handle(logging.Handle(compression.Handle(metricsController.List))))
	r.GET("/value/:type/:name", decompression.Handle(logging.Handle(compression.Handle(metricsController.Show))))
	r.GET("/ping", decompression.Handle(compression.Handle(metricsController.Ping)))
	r.POST("/value/", decompression.Handle(logging.Handle(compression.Handle(metricsController.ShowBody))))
	r.POST("/update/:type/:name/:value", decompression.Handle(logging.Handle(compression.Handle(metricsController.Save))))
	r.POST("/update/", decompression.Handle(logging.Handle(compression.Handle(metricsController.SaveBody))))
	return r
}
