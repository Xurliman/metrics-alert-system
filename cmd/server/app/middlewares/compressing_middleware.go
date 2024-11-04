package middlewares

import (
	"compress/gzip"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

type CompressingMiddleware struct {
	Response
}

func NewCompressingMiddleware() Middleware {
	return &CompressingMiddleware{}
}

func (c CompressingMiddleware) Handle(next gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		acceptEncoding := ctx.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			next(ctx)
			return
		}

		gz := gzip.NewWriter(ctx.Writer)
		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				utils.JSONInternalServerError(ctx, err)
			}
		}(gz)

		gzw := &GzipResponseWriter{ResponseWriter: ctx.Writer, Writer: gz}
		ctx.Writer = gzw
		ctx.Writer.Header().Set("Content-Encoding", "gzip")

		next(ctx)
	}
}

type GzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (gzw *GzipResponseWriter) Write(data []byte) (int, error) {
	return gzw.Writer.Write(data)
}
