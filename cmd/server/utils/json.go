package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DefaultResponse struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSONSuccess(ctx *gin.Context, data interface{}) {
	Logger.Info(logSuccessFormat(data))
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	//ctx.JSON(http.StatusOK, &DefaultResponse{
	//	Success: true,
	//	Status:  http.StatusOK,
	//	Message: "OK",
	//	Data:    data,
	//})
	ctx.JSON(http.StatusOK, data)
	ctx.Abort()
}

func JSONError(ctx *gin.Context, err error) {
	Logger.Info(logErrorFormat(err))
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.JSON(http.StatusBadRequest, DefaultResponse{
		Success: false,
		Status:  http.StatusBadRequest,
		Message: err.Error(),
		Data:    nil,
	})
	ctx.Abort()
}

func JSONInternalServerError(ctx *gin.Context, err error) {
	Logger.Info(logErrorFormat(err))
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.JSON(http.StatusInternalServerError, DefaultResponse{
		Success: false,
		Status:  http.StatusInternalServerError,
		Message: err.Error(),
		Data:    nil,
	})
	ctx.Abort()
}

func JSONValidationError(ctx *gin.Context, err error) {
	Logger.Info(logErrorFormat(err))
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.JSON(http.StatusUnprocessableEntity, DefaultResponse{
		Success: false,
		Status:  http.StatusUnprocessableEntity,
		Message: err.Error(),
		Data:    nil,
	})
	ctx.Abort()
}

func JSONNotFound(ctx *gin.Context, err error) {
	Logger.Info(logErrorFormat(err))
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.JSON(http.StatusNotFound, DefaultResponse{
		Success: false,
		Status:  http.StatusNotFound,
		Message: err.Error(),
		Data:    nil,
	})
}

func logSuccessFormat(data interface{}) string {
	jsonData, _ := json.MarshalIndent(data, "", "    ")
	return "✅" + string(jsonData)
}

func logErrorFormat(err error) string {
	return "❌ " + err.Error()
}
