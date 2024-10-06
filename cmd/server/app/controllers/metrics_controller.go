package controllers

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MetricsController struct {
	service interfaces.MetricsServiceInterface
}

func NewMetricsController(service interfaces.MetricsServiceInterface) *MetricsController {
	return &MetricsController{service: service}
}

func (c *MetricsController) Validate(ctx *gin.Context) {
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
	if metricsType == "counter" {
		_, err := strconv.ParseInt(metricsValue, 10, 64)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
	}
	if metricsType == "gauge" {
		_, err := strconv.ParseFloat(metricsValue, 64)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
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
	}
	ctx.String(http.StatusOK, metricsValue)
}
