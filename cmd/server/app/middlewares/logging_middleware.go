package middlewares

import (
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
	return func(c *gin.Context) {
		start := time.Now()

		next(c)

		l.Request.URI = c.Request.RequestURI
		l.Request.Method = c.Request.Method
		l.Request.Duration = time.Since(start)
		l.Response.Size = c.Writer.Size()
		l.Response.StatusCode = c.Writer.Status()
		utils.Logger.Info(
			"info",
			zap.String("uri", l.Request.URI),
			zap.String("method", l.Request.Method),
			zap.Duration("duration", l.Request.Duration),
			zap.Int("status", l.Response.StatusCode),
			zap.Int("size", l.Response.Size),
		)
	}
}
