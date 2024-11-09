package controllers

import (
	"errors"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/resources"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"net/http"
)

type MetricsController struct {
	service interfaces.MetricsServiceInterface
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
	)

	if err := c.service.SaveWhenParams(metricType, metricName, metricValue); err != nil {
		if errors.Is(err, constants.ErrEmptyMetricName) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusBadRequest)
		return
	}

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

	metric, err := c.service.SaveWhenBody(*metricRequest)
	if err != nil {
		utils.JSONInternalServerError(ctx, err)
		return
	}

	utils.JSONSuccess(ctx, metric)
}

func (c *MetricsController) Show(ctx *gin.Context) {
	var (
		metricType  = ctx.Param("type")
		metricName  = ctx.Param("name")
		metricValue string
		err         error
	)

	metricValue, err = c.service.GetMetricValue(metricType, metricName)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.String(http.StatusOK, metricValue)
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

	metric, err := c.service.Show(metricRequest.ID)
	if err != nil {
		utils.JSONNotFound(ctx, err)
		return
	}

	utils.JSONSuccess(ctx, resources.ToResponse(metric))
}

func (c *MetricsController) Ping(ctx *gin.Context) {
	err := c.service.Ping(ctx)
	if err != nil {
		utils.JSONInternalServerError(ctx, err)
		return
	}
	utils.JSONSuccess(ctx, "successfully connected to the database")
}
