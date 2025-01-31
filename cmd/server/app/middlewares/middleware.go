package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	//Handle(next gin.HandlerFunc) gin.HandlerFunc
	Handle(ctx *gin.Context)
}

type Request struct {
	URI      string
	Method   string
	Duration time.Duration
	Body     map[string]interface{}
	Header   http.Header
	Error    error
}

func (request *Request) Handle(ctx *gin.Context) (size int64) {
	var (
		buf bytes.Buffer
		err error
	)

	request.URI = ctx.Request.RequestURI
	request.Method = ctx.Request.Method
	request.Header = ctx.Request.Header

	_, err = io.Copy(&buf, ctx.Request.Body)
	if err != nil {
		request.Error = err
	}

	requestBodyBytes := buf.Bytes()
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

	err = json.Unmarshal(buf.Bytes(), &request.Body)
	if err != nil {
		errors.Join(request.Error, err)
	}

	return size
}

type Response struct {
	StatusCode int
	Size       int64
	Body       map[string]interface{}
}

type ResponseCapture struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (r *ResponseCapture) Write(b []byte) (int, error) {
	r.Body.Write(b)                  // Capture the response body
	return r.ResponseWriter.Write(b) // Handle it to the actual response
}
