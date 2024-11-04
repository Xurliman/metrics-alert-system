package middlewares

import (
	"bytes"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
)

type DecompressingMiddleware struct{}

func NewDecompressingMiddleware() Middleware {
	return &DecompressingMiddleware{}
}

func (d DecompressingMiddleware) Handle(next gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contentEncoding := ctx.GetHeader("Content-Encoding")
		if !strings.Contains(contentEncoding, "gzip") {
			next(ctx)
			return
		}

		var buf bytes.Buffer
		_, err := io.Copy(&buf, ctx.Request.Body)
		if err != nil {
			utils.JSONInternalServerError(ctx, fmt.Errorf("error reading request body: %v", err))
			return
		}

		decompressedBody, err := utils.Decompress(buf.Bytes())
		if err != nil {
			utils.JSONInternalServerError(ctx, fmt.Errorf("error decompressing request body: %v", err))
			return
		}

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(decompressedBody))
		next(ctx)
	}
}
