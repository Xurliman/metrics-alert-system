package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
	"io"
)

type HashingMiddleware struct {
	key string
}

func (h HashingMiddleware) Handle(ctx *gin.Context) {
	hashMsg := ctx.Request.Header.Get("HashSHA256")
	if hashMsg == "" {
		ctx.Next()
		return
	}

	decodedData, err := hex.DecodeString(hashMsg)
	if err != nil {
		utils.JSONError(ctx, fmt.Errorf("invalid hash format: %v", err))
		return
	}

	// Read request body
	body, err := ctx.GetRawData()
	if err != nil {
		utils.JSONError(ctx, fmt.Errorf("failed to read request body: %v", err))
		return
	}

	// Restore the request body so Gin can read it again later
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	hm := hmac.New(sha256.New, []byte(h.key))
	hm.Write(body) // Hash the actual request body
	computedHash := hm.Sum(nil)

	if !hmac.Equal(computedHash, decodedData) {
		utils.JSONError(ctx, constants.ErrInvalidHash)
		return
	}

	ctx.Next()
}

func NewHashingMiddleware(key string) Middleware {
	return &HashingMiddleware{
		key: key,
	}
}
