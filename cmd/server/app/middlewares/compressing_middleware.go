package middlewares

import (
	"compress/gzip"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/Xurliman/metrics-alert-system/internal/compressor"
	"github.com/gin-gonic/gin"
	"strings"
)

type CompressingMiddleware struct {
	Response
}

func NewCompressingMiddleware() Middleware {
	return &CompressingMiddleware{}
}

func (c CompressingMiddleware) Handle(ctx *gin.Context) {
	acceptEncoding := ctx.GetHeader("Accept-Encoding")
	if !strings.Contains(acceptEncoding, "gzip") {
		ctx.Next()
		return
	}

	gz := gzip.NewWriter(ctx.Writer)
	defer func(gz *gzip.Writer) {
		err := gz.Close()
		if err != nil {
			utils.JSONInternalServerError(ctx, fmt.Errorf("error closing gzip writer: %v", err))
		}
	}(gz)

	gzw := &compressor.GzipResponseWriter{ResponseWriter: ctx.Writer, Writer: gz}
	ctx.Writer = gzw
	ctx.Writer.Header().Set("Content-Encoding", "gzip")

	ctx.Next()
}
