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
	"go.uber.org/zap"
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
			utils.JSONNotFound(ctx, err)
			return
		}
		utils.JSONError(ctx, err)
		return
	}

	utils.JSONSuccess(ctx, nil)
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

func (c *MetricsController) SaveMany(ctx *gin.Context) {
	var request []requests.MetricsSaveRequest
	err := ctx.ShouldBindWith(&request, binding.JSON)
	if err != nil && err != io.EOF {
		utils.JSONValidationError(ctx, err)
		return
	}

	for _, req := range request {
		if err = req.Validate(); err != nil {
			utils.JSONValidationError(ctx, err)
			return
		}
	}
	utils.Logger.Error("REQUEST", zap.Any("error", request))

	err = c.service.SaveMany(ctx.Request.Context(), request)
	if err != nil {
		utils.JSONInternalServerError(ctx, err)
		return
	}
	utils.JSONSuccess(ctx, "Success")
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

	response, err := resources.ToResponse(metric)
	if err != nil {
		utils.JSONInternalServerError(ctx, err)
		return
	}
	
	utils.JSONSuccess(ctx, *response)
}

func (c *MetricsController) Ping(ctx *gin.Context) {
	err := c.service.Ping(ctx.Request.Context())
	if err != nil {
		utils.JSONInternalServerError(ctx, err)
		return
	}
	utils.JSONSuccess(ctx, "successfully connected to the database")
}
