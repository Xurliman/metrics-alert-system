package middlewares

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Middleware interface {
	Handle(next gin.HandlerFunc) gin.HandlerFunc
}

type Request struct {
	URI      string
	Method   string
	Duration time.Duration
}

type Response struct {
	StatusCode int
	Size       int
}
