package controllers

import (
	"errors"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/mocks/servicemocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"testing"
)

func setupRoutes(service *servicemocks.MetricsServiceInterface) *gin.Engine {
	controller := NewMetricsController(service)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLFiles("../../public/templates/metrics-all.html")
	r.GET("/", controller.Index)
	r.GET("/value/:type/:name/", controller.Show)
	r.POST("/update/:type/:name/:value", controller.Update)
	return r
}

func TestMetricsController_Index(t *testing.T) {
	service := servicemocks.NewMetricsServiceInterface(t)
	router := setupRoutes(service)
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
			myMap := make(map[string]string)
			service.On("GetAll").Return(myMap)
			req := httptest.NewRequest(test.method, test.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, test.expected.statusCode, resp.Code)
		})
	}
}

func TestMetricsController_Show(t *testing.T) {
	service := servicemocks.NewMetricsServiceInterface(t)
	router := setupRoutes(service)
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name        string
		metricsType string
		metricsName string
		method      string
		wantErr     bool
		url         string
		expected    expected
	}{
		{
			name:        "first",
			metricsType: "counter",
			metricsName: "someMetric",
			method:      http.MethodGet,
			wantErr:     true,
			url:         "http://localhost:8080/value/counter/someMetric/",
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:        "second",
			metricsType: "gauge",
			metricsName: "someMetric",
			method:      http.MethodGet,
			wantErr:     false,
			url:         "http://localhost:8080/value/gauge/someMetric/",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "third",
			metricsType: "gauge",
			metricsName: "GCCPUFraction",
			method:      http.MethodGet,
			wantErr:     false,
			url:         "http://localhost:8080/value/gauge/GCCPUFraction/",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "fourth",
			metricsType: "counter",
			metricsName: "HeapObjects",
			method:      http.MethodGet,
			wantErr:     true,
			url:         "http://localhost:8080/value/counter/HeapObjects/",
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.wantErr {
				service.On("FindByName", test.metricsType, test.metricsName).Return("", errors.New("not found"))
			} else {
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				service.On("FindByName", test.metricsType, test.metricsName).Return(strconv.FormatFloat(memStats.GCCPUFraction, 'f', -1, 64), nil)
			}
			req := httptest.NewRequest(test.method, test.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, test.expected.statusCode, resp.Code)
		})
	}
}

func TestMetricsController_Update(t *testing.T) {
	service := servicemocks.NewMetricsServiceInterface(t)
	router := setupRoutes(service)
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name         string
		metricsType  string
		metricsName  string
		metricsValue string
		wantErr      bool
		method       string
		url          string
		expected     expected
	}{

		{
			name:         "third",
			metricsType:  "gauge",
			metricsName:  "someMetric",
			metricsValue: "527",
			wantErr:      false,
			method:       http.MethodPost,
			url:          "http://localhost:8080/update/gauge/someMetric/527",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:         "first",
			metricsType:  "counter",
			metricsName:  "someMetric",
			metricsValue: "527",
			wantErr:      false,
			method:       http.MethodPost,
			url:          "http://localhost:8080/update/counter/someMetric/527",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:         "second",
			metricsType:  "gauge",
			metricsName:  "someMetric",
			metricsValue: "527",
			wantErr:      true,
			method:       http.MethodGet,
			url:          "http://localhost:8080/update/gauge/someMetric/527",
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:         "fourth",
			metricsType:  "counter",
			metricsName:  "someMetric",
			metricsValue: "unknown",
			wantErr:      true,
			method:       http.MethodPost,
			url:          "http://localhost:8080/update/counter/someMetric/unknown",
			expected: expected{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.wantErr {
				service.On("Save", test.metricsType, test.metricsName, test.metricsValue).Return(errors.New("invalid metrics value for counter type"))
			} else {
				service.On("Save", test.metricsType, test.metricsName, test.metricsValue).Return(nil)
			}
			req := httptest.NewRequest(test.method, test.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, test.expected.statusCode, resp.Code)
		})
	}
}
