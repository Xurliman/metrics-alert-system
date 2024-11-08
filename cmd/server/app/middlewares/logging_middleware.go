package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type LoggingMiddleware struct {
	Request
	Response
}

func NewLoggingMiddleware() Middleware {
	return LoggingMiddleware{}
}

func (l LoggingMiddleware) Handle(next gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		size := l.Request.Handle(ctx)

		var respBody bytes.Buffer
		respCapture := &ResponseCapture{
			ResponseWriter: ctx.Writer,
			Body:           &respBody,
		}
		ctx.Writer = respCapture

		next(ctx)
		
		l.Request.Duration = time.Since(start)
		l.Response.Size = size
		l.Response.StatusCode = ctx.Writer.Status()
		err := json.Unmarshal(respCapture.Body.Bytes(), &l.Response.Body)
		if err != nil {
			l.Response.Error = err
		}
		utils.Logger.Info(
			"info",
			zap.String("uri", l.Request.URI),
			zap.String("method", l.Request.Method),
			zap.Duration("duration", l.Request.Duration),
			zap.Reflect("request_body", l.Request.Body),
			zap.NamedError("request_err", l.Request.Error),

			zap.Int("status", l.Response.StatusCode),
			zap.Int64("size", l.Response.Size),
			zap.Reflect("response_body", l.Response.Body),
			zap.NamedError("response_err", l.Response.Error),
		)
	}
}
