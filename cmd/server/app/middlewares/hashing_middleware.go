package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
)

type HashingMiddleware struct {
	key string
}

func (h HashingMiddleware) Handle(next gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hashMsg := ctx.Request.Header.Get("HashSHA256")
		if hashMsg == "" {
			next(ctx)
			return
		}
		decodedData, err := hex.DecodeString(hashMsg)
		if err != nil {
			utils.JSONError(ctx, err)
			return
		}

		hm := hmac.New(sha256.New, []byte(h.key))
		hm.Write(decodedData)
		sign := hm.Sum(nil)

		if !hmac.Equal(sign, decodedData) {
			utils.JSONError(ctx, errors.New("invalid hash"))
			return
		}
	}
}

func NewHashingMiddleware(key string) Middleware {
	return &HashingMiddleware{
		key: key,
	}
}
