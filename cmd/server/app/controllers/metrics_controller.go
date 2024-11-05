package controllers

import (
	"errors"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"net/http"
	"time"
)

type MetricsController struct {
	service         interfaces.MetricsServiceInterface
	storeTicker     *time.Ticker
	fileStoragePath string
}

func NewMetricsController(service interfaces.MetricsServiceInterface) interfaces.MetricsControllerInterface {
	return &MetricsController{
		service: service,
	}
}

func (c *MetricsController) List(ctx *gin.Context) {
	data := c.service.List()
	ctx.HTML(http.StatusOK, "metrics-all.html", data)
}

func (c *MetricsController) Save(ctx *gin.Context) {
	var (
		metricType  = ctx.Param("type")
		metricName  = ctx.Param("name")
		metricValue = ctx.Param("value")
		err         error
	)
	switch metricType {
	case constants.GaugeMetricType:
		err = c.service.SaveWhenParams(services.Gauge, metricName, metricValue)
	case constants.CounterMetricType:
		err = c.service.SaveWhenParams(services.Counter, metricName, metricValue)
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

func (c *MetricsController) SaveBody(ctx *gin.Context) {
	metricRequest := new(requests.MetricsSaveRequest)
	err := ctx.ShouldBindWith(&metricRequest, binding.JSON)
	if err != nil && err != io.EOF {
		utils.JSONValidationError(ctx, err)
		return
	}

	if err = metricRequest.Validate(); err != nil {
		utils.JSONValidationError(ctx, err)
		return
	}

	var metric *models.Metrics
	switch metricRequest.MType {
	case constants.GaugeMetricType:
		metric, err = c.service.SaveWhenBody(services.Gauge, *metricRequest)
	case constants.CounterMetricType:
		metric, err = c.service.SaveWhenBody(services.Counter, *metricRequest)
	default:
		utils.JSONError(ctx, constants.ErrInvalidMetricType)
		return
	}

	if err != nil {
		utils.JSONInternalServerError(ctx, err)
		return
	}

	utils.JSONSuccess(ctx, metric)
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
		metricsValue, err = c.service.GetMetricValue(services.Gauge, metricsName)
	case constants.CounterMetricType:
		metricsValue, err = c.service.GetMetricValue(services.Counter, metricsName)
	}

	defer func(err error) {
		if err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.String(http.StatusOK, metricsValue)

	}(err)
}

func (c *MetricsController) ShowBody(ctx *gin.Context) {
	metricRequest := new(requests.MetricsShowRequest)
	err := ctx.ShouldBindWith(&metricRequest, binding.JSON)
	if err != nil && err != io.EOF {
		utils.JSONValidationError(ctx, err)
		return
	}

	if err = metricRequest.Validate(); err != nil {
		utils.JSONValidationError(ctx, err)
		return
	}

	var metric *models.Metrics
	switch metricRequest.MType {
	case constants.GaugeMetricType:
		metric, err = c.service.Show(services.Gauge, metricRequest.ID)
	case constants.CounterMetricType:
		metric, err = c.service.Show(services.Counter, metricRequest.ID)
	default:
		utils.JSONError(ctx, constants.ErrInvalidMetricType)
		return
	}

	if err != nil {
		utils.JSONNotFound(ctx, err)
		return
	}

	utils.JSONSuccess(ctx, metric)
}
