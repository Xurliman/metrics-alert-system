package controllers

import (
	"encoding/json"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/middlewares"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/mocks/servicemocks"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/services"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func setupRoutes(service *servicemocks.MetricsServiceInterface) *gin.Engine {
	utils.Logger = utils.NewLogger(gin.TestMode)
	logging := middlewares.NewLoggingMiddleware()
	controller := NewMetricsController(service)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLFiles("../../public/templates/metrics-all.html")
	r.GET("/", logging.Handle(controller.List))
	r.POST("/update/:type/:name/:value", logging.Handle(controller.Save))
	r.POST("/update/", logging.Handle(controller.SaveBody))
	r.GET("/value/:type/:name/", logging.Handle(controller.Show))
	r.POST("/value/", logging.Handle(controller.ShowBody))
	return r
}

func TestMetricsController_List(t *testing.T) {
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
			service.On("List").Return(myMap)
			req := httptest.NewRequest(test.method, test.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, test.expected.statusCode, resp.Code)
		})
	}
}

func TestMetricsController_ShowBody(t *testing.T) {
	service := servicemocks.NewMetricsServiceInterface(t)
	router := setupRoutes(service)
	tests := []struct {
		name       string
		url        string
		metricType string
		metricName string
		body       string
		wantBody   string
		want       map[string]interface{}
		wantErr    bool
	}{
		{
			name:       "first",
			url:        "http://localhost:8080/value/",
			metricType: constants.GaugeMetricType,
			metricName: "HeapIdle",
			body: `{
				"id" : "HeapIdle",
				"type": "gauge"
			}`,
			wantBody: `{
    "success": false,
    "status": 404,
    "message": "metric not found",
    "data": null
}`,
			want: map[string]interface{}{
				"success": false,
				"status":  float64(404),
				"message": "metric not found",
				"data":    nil,
			},
			wantErr: true,
		},
		{
			name:       "second",
			url:        "http://localhost:8080/value/",
			metricType: constants.CounterMetricType,
			metricName: "PollCount",
			body: `{
				"id" : "PollCount",
				"type" : "counter"
			}`,
			wantBody: `{
	"id": "PollCount",
	"type": "counter"
}`,
			want: map[string]interface{}{
				"id":   "PollCount",
				"type": "counter",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.metricType {
			case constants.GaugeMetricType:
				service.On("Show", services.Gauge, tt.metricName).Return(nil, constants.ErrMetricNotFound)
			case constants.CounterMetricType:
				service.On("Show", services.Counter, tt.metricName).Return(&models.Metrics{
					ID:    tt.metricName,
					MType: tt.metricType,
				}, nil)
			}
			body := strings.NewReader(tt.body)
			req := httptest.NewRequest(http.MethodPost, tt.url, body)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			var got map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &got)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, got)
			service.AssertExpectations(t)
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
				switch test.metricsType {
				case constants.GaugeMetricType:
					service.On("GetMetricValue", services.Gauge, test.metricsName).Return("", constants.ErrInvalidMetricType)
				case constants.CounterMetricType:
					service.On("GetMetricValue", services.Counter, test.metricsName).Return("", constants.ErrInvalidMetricType)
				}
			} else {
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				switch test.metricsType {
				case constants.GaugeMetricType:
					service.On("GetMetricValue", services.Gauge, test.metricsName).Return(strconv.FormatFloat(memStats.GCCPUFraction, 'f', -1, 64), nil)
				case constants.CounterMetricType:
					service.On("GetMetricValue", services.Counter, test.metricsName).Return(strconv.FormatUint(memStats.HeapObjects, 10), nil)
				}
			}
			req := httptest.NewRequest(test.method, test.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, test.expected.statusCode, resp.Code)
		})
	}
}

func TestMetricsController_SaveBody(t *testing.T) {
	service := servicemocks.NewMetricsServiceInterface(t)
	router := setupRoutes(service)
	tests := []struct {
		name            string
		url             string
		body            string
		wantBody        string
		wantComplexBody string
		wantComplexMap  map[string]interface{}
		wantErr         bool
	}{
		{
			name: "first",
			url:  "http://localhost:8080/update/",
			body: `{
				"id" : "HeapIdle",
				"type": "gauge"
			}`,
			wantBody: `{}`,
			wantComplexBody: `{
	"success": false,
	"status": 500,
	"message": "invalid metrics value for gauge type",
	"data": null
}`,
			wantComplexMap: map[string]interface{}{
				"success": false,
				"status":  float64(500),
				"message": "invalid metrics value for gauge type",
				"data":    nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service.On("SaveWhenBody", services.Gauge, mock.AnythingOfType("requests.MetricsSaveRequest")).Return(nil, constants.ErrInvalidGaugeMetricValue)

			body := strings.NewReader(tt.body)
			req := httptest.NewRequest(http.MethodPost, tt.url, body)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			var got map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &got)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantComplexMap, got)
			service.AssertExpectations(t)
		})
	}
}

func TestMetricsController_Save(t *testing.T) {
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
				switch test.metricsType {
				case constants.GaugeMetricType:
					service.On("SaveWhenParams", services.Gauge, test.metricsName, test.metricsValue).Return(constants.ErrInvalidGaugeMetricValue)
				case constants.CounterMetricType:
					service.On("SaveWhenParams", services.Counter, test.metricsName, test.metricsValue).Return(constants.ErrInvalidCounterMetricValue)
				}
			} else {
				switch test.metricsType {
				case constants.GaugeMetricType:
					service.On("SaveWhenParams", services.Gauge, test.metricsName, test.metricsValue).Return(nil)
				case constants.CounterMetricType:
					service.On("SaveWhenParams", services.Counter, test.metricsName, test.metricsValue).Return(nil)
				}
			}

			req := httptest.NewRequest(test.method, test.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, test.expected.statusCode, resp.Code)
		})
	}
}
