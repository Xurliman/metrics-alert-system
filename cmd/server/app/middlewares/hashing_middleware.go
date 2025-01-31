package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
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

	hm := hmac.New(sha256.New, []byte(h.key))
	hm.Write(decodedData)
	sign := hm.Sum(nil)

	if !hmac.Equal(sign, decodedData) {
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
