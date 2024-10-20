package controllers

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MetricsController struct {
	service interfaces.MetricsServiceInterface
}

func NewMetricsController(service interfaces.MetricsServiceInterface) *MetricsController {
	return &MetricsController{service: service}
}

func (c *MetricsController) Update(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodPost {
		ctx.Status(http.StatusBadRequest)
	}
	metricsType := ctx.Param("type")
	if metricsType != "gauge" && metricsType != "counter" {
		ctx.Status(http.StatusBadRequest)
		return
	}
	metricsName := ctx.Param("name")
	if metricsName == "" {
		ctx.Status(http.StatusNotFound)
		return
	}
	metricsValue := ctx.Param("value")
	err := c.service.Save(metricsType, metricsName, metricsValue)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *MetricsController) Index(ctx *gin.Context) {
	data := c.service.GetAll()
	ctx.HTML(http.StatusOK, "metrics-all.html", data)
}

func (c *MetricsController) Show(ctx *gin.Context) {
	metricsType := ctx.Param("type")
	metricsName := ctx.Param("name")
	metricsValue, err := c.service.FindByName(metricsType, metricsName)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.String(http.StatusOK, metricsValue)
}
