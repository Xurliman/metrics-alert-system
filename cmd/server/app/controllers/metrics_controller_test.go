package controllers

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsController_Index(t *testing.T) {
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name     string
		method   string
		url      string
		expected expected
	}{
		{
			name:   "first",
			method: http.MethodGet,
			url:    "http://localhost:8080/",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := routes.SetupRoutes()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(test.method, test.url, nil)
			router.ServeHTTP(w, r)

			result := w.Result()
			assert.Equal(t, test.expected.statusCode, result.StatusCode)
			defer result.Body.Close()
		})
	}
}

func TestMetricsController_Show(t *testing.T) {
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name     string
		method   string
		url      string
		expected expected
	}{
		{
			name:   "first",
			method: http.MethodGet,
			url:    "http://localhost:8080/value/gauge/LastGC",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:   "second",
			method: http.MethodGet,
			url:    "http://localhost:8080/update/gauge/someMetric",
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.Default()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(test.method, test.url, nil)
			router.ServeHTTP(w, r)

			result := w.Result()
			assert.Equal(t, test.expected.statusCode, result.StatusCode)
			defer result.Body.Close()
		})
	}
}

func TestMetricsController_Validate(t *testing.T) {
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name     string
		method   string
		url      string
		expected expected
	}{
		{
			name:   "first",
			method: http.MethodPost,
			url:    "http://localhost:8080/update/counter/someMetric/527",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:   "second",
			method: http.MethodGet,
			url:    "http://localhost:8080/update/gauge/someMetric/527",
			expected: expected{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "third",
			method: http.MethodPost,
			url:    "http://localhost:8080/update/gauge/someMetric/527",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:   "fourth",
			method: http.MethodPost,
			url:    "http://localhost:8080/update/unknown/someMetric/527",
			expected: expected{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "fifth",
			method: http.MethodPost,
			url:    "http://localhost:8080/update/counter/someMetric/unknown",
			expected: expected{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.Default()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(test.method, test.url, nil)
			router.ServeHTTP(w, r)

			result := w.Result()
			assert.Equal(t, test.expected.statusCode, result.StatusCode)
			defer result.Body.Close()
		})
	}
}
