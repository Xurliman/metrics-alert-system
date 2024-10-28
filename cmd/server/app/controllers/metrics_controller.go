package controllers

import (
	"errors"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
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
	var (
		metricsType  = ctx.Param("type")
		metricsName  = ctx.Param("name")
		metricsValue = ctx.Param("value")
		err          error
	)

	switch metricsType {
	case constants.GaugeMetricType:
		err = c.service.SaveGaugeMetric(metricsName, metricsValue)
	case constants.CounterMetricType:
		err = c.service.SaveCounterMetric(metricsName, metricsValue)
	default:
		ctx.Status(http.StatusBadRequest)
		return
	}

	defer func(err error) {
		if err != nil {
			if errors.Is(err, constants.ErrEmptyMetricName) {
				ctx.Status(http.StatusNotFound)
			}
			ctx.Status(http.StatusBadRequest)
			return
		}
	}(err)

	ctx.Status(http.StatusOK)
}

func (c *MetricsController) Index(ctx *gin.Context) {
	data := c.service.GetAll()
	ctx.HTML(http.StatusOK, "metrics-all.html", data)
}

func (c *MetricsController) Show(ctx *gin.Context) {
	var (
		metricsType  = ctx.Param("type")
		metricsName  = ctx.Param("name")
		metricsValue string
		err          error
	)

	switch metricsType {
	case constants.GaugeMetricType:
		metricsValue, err = c.service.FindGaugeMetric(metricsName)
	case constants.CounterMetricType:
		metricsValue, err = c.service.FindCounterMetric(metricsName)
	}

	defer func(err error) {
		if err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.String(http.StatusOK, metricsValue)

	}(err)
}
