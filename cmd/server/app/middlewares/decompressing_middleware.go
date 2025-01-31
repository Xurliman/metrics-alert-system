package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/Xurliman/metrics-alert-system/internal/compressor"
	"github.com/gin-gonic/gin"
)

type DecompressingMiddleware struct{}

func NewDecompressingMiddleware() Middleware {
	return &DecompressingMiddleware{}
}

func (d DecompressingMiddleware) Handle(ctx *gin.Context) {
	contentEncoding := ctx.GetHeader("Content-Encoding")
	if !strings.Contains(contentEncoding, "gzip") {
		ctx.Next()
		return
	}

	var buf bytes.Buffer
	_, err := io.Copy(&buf, ctx.Request.Body)
	if err != nil {
		utils.JSONInternalServerError(ctx, fmt.Errorf("error reading request body: %v", err))
		return
	}

	decompressedBody, err := compressor.Decompress(buf.Bytes())
	if err != nil {
		utils.JSONInternalServerError(ctx, fmt.Errorf("error decompressing request body: %v", err))
		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(decompressedBody))

	ctx.Next()
}
