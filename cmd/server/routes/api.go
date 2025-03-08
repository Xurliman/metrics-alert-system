package routes

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/controllers"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/middlewares"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/server/config"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(metricsRepository interfaces.MetricsRepositoryInterface, cfg *config.Config) *gin.Engine {
	decompression := middlewares.NewDecompressingMiddleware()
	logging := middlewares.NewLoggingMiddleware()
	compression := middlewares.NewCompressingMiddleware()
	hashing := middlewares.NewHashingMiddleware(cfg.Key)
	decrypting := middlewares.NewDecryptingMiddleware(cfg.CryptoKey)

	r := gin.New()
	r.LoadHTMLFiles("./cmd/server/public/templates/metrics-all.html")

	metricsService := services.NewMetricsService(metricsRepository, services.NewSwitchService())
	metricsController := controllers.NewMetricsController(metricsService)

	r.Use(decrypting.Handle, decompression.Handle, hashing.Handle, compression.Handle, logging.Handle)
	r.GET("/ping", metricsController.Ping)
	r.GET("/", metricsController.List)
	r.GET("/value/:type/:name", metricsController.Show)
	r.POST("/value/", metricsController.ShowBody)
	r.POST("/update/:type/:name/:value", metricsController.Save)
	r.POST("/update/", metricsController.SaveBody)
	r.POST("/updates/", metricsController.SaveMany)
	return r
}
