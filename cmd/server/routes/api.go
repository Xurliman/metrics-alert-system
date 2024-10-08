package routes

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/controllers"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func SetupRoutes() *gin.Engine {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Using env from machine")
	}
	mode := os.Getenv("GIN_MODE")
	gin.SetMode(mode)
	r := gin.New()
	r.LoadHTMLFiles("./cmd/server/public/templates/metrics-all.html")
	metricsService := services.NewMetricsService()
	metricsController := controllers.NewMetricsController(metricsService)
	r.GET("/", metricsController.Index)
	r.GET("/value/:type/:name", metricsController.Show)
	r.POST("/update/:type/:name/:value", metricsController.Validate)
	return r
}
