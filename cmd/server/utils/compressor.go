package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
)

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w, err := gzip.NewWriterLevel(&buf, flate.BestSpeed)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type GzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (gzw *GzipResponseWriter) Write(data []byte) (int, error) {
	return gzw.Writer.Write(data)
}

func Decompress(data []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(data))
	defer func(r io.ReadCloser) {
		err := r.Close()
		if err != nil {
			Logger.Error("decompressing err", zap.Error(err))
		}
	}(r)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
