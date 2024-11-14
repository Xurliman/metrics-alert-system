package controllers

import (
	"encoding/json"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/middlewares"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/mocks/servicemocks"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func callOnceToChangeDir() {
	err := os.Chdir("../../../../")
	if err != nil {
		log.Fatalf("Failed to set working directory: %v", err)
	}
}

func setupRoutes(service *servicemocks.MetricsServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLFiles("cmd/server/public/templates/metrics-all.html")
	utils.Logger = utils.NewLogger(gin.TestMode)
	logging := middlewares.NewLoggingMiddleware()
	controller := NewMetricsController(service)

	r.GET("/", logging.Handle(controller.List))
	r.POST("/update/:type/:name/:value", logging.Handle(controller.Save))
	r.POST("/update/", logging.Handle(controller.SaveBody))
	r.GET("/value/:type/:name/", logging.Handle(controller.Show))
	r.POST("/value/", logging.Handle(controller.ShowBody))
	return r
}

func TestMetricsController_List(t *testing.T) {
	callOnceToChangeDir()
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myMap := make(map[string]string)
			service.On("List").Return(myMap, nil)
			req := httptest.NewRequest(tt.method, tt.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.expected.statusCode, resp.Code)
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
		name        string
		metricType  string
		metricName  string
		metricValue string
		wantErr     bool
		method      string
		url         string
		expected    expected
	}{

		{
			name:        "third",
			metricType:  "gauge",
			metricName:  "someMetric",
			metricValue: "527",
			wantErr:     false,
			method:      http.MethodPost,
			url:         "http://localhost:8080/update/gauge/someMetric/527",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "first",
			metricType:  "counter",
			metricName:  "someMetric",
			metricValue: "527",
			wantErr:     false,
			method:      http.MethodPost,
			url:         "http://localhost:8080/update/counter/someMetric/527",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "second",
			metricType:  "gauge",
			metricName:  "someMetric",
			metricValue: "527",
			wantErr:     true,
			method:      http.MethodGet,
			url:         "http://localhost:8080/update/gauge/someMetric/527",
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:        "fourth",
			metricType:  "counter",
			metricName:  "someMetric",
			metricValue: "unknown",
			wantErr:     true,
			method:      http.MethodPost,
			url:         "http://localhost:8080/update/counter/someMetric/unknown",
			expected: expected{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				service.On("SaveWhenParams", tt.metricType, tt.metricName, tt.metricValue).Return(constants.ErrInvalidGaugeMetricValue)
			} else {
				service.On("SaveWhenParams", tt.metricType, tt.metricName, tt.metricValue).Return(nil)
			}

			req := httptest.NewRequest(tt.method, tt.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.expected.statusCode, resp.Code)
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
			service.On("SaveWhenBody", mock.AnythingOfType("requests.MetricsSaveRequest")).Return(nil, constants.ErrInvalidGaugeMetricValue)

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

func TestMetricsController_Show(t *testing.T) {
	service := servicemocks.NewMetricsServiceInterface(t)
	router := setupRoutes(service)
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name       string
		metricType string
		metricName string
		method     string
		wantErr    bool
		url        string
		expected   expected
	}{
		{
			name:       "first",
			metricType: "counter",
			metricName: "someMetric",
			method:     http.MethodGet,
			wantErr:    true,
			url:        "http://localhost:8080/value/counter/someMetric/",
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:       "second",
			metricType: "gauge",
			metricName: "someMetric",
			method:     http.MethodGet,
			wantErr:    false,
			url:        "http://localhost:8080/value/gauge/someMetric/",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:       "third",
			metricType: "gauge",
			metricName: "GCCPUFraction",
			method:     http.MethodGet,
			wantErr:    false,
			url:        "http://localhost:8080/value/gauge/GCCPUFraction/",
			expected: expected{
				statusCode: http.StatusOK,
			},
		},
		{
			name:       "fourth",
			metricType: "counter",
			metricName: "HeapObjects",
			method:     http.MethodGet,
			wantErr:    true,
			url:        "http://localhost:8080/value/counter/HeapObjects/",
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				service.On("GetMetricValue", tt.metricType, tt.metricName).Return("89", constants.ErrMetricNotFound)
			} else {
				service.On("GetMetricValue", tt.metricType, tt.metricName).Return("", nil)
			}

			req := httptest.NewRequest(tt.method, tt.url, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.expected.statusCode, resp.Code)
		})
	}
}

func TestMetricsController_ShowBody(t *testing.T) {
	someGaugeMetric := float64(203)
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
		wantValue  *float64
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
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:       "second",
			url:        "http://localhost:8080/value/",
			metricType: constants.GaugeMetricType,
			metricName: "Alloc",
			body: `{
				"id" : "Alloc",
				"type" : "gauge"
			}`,
			wantBody: `{
	"id": "Alloc",
	"type": "gauge"
	"value": 0
}`,
			want: map[string]interface{}{
				"id":    "Alloc",
				"type":  "gauge",
				"value": someGaugeMetric,
			},
			wantValue: &someGaugeMetric,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				service.On("Show", tt.metricName).Return(nil, constants.ErrMetricNotFound)
			} else {
				service.On("Show", tt.metricName).Return(&models.Metrics{
					ID:    tt.metricName,
					MType: tt.metricType,
					Value: tt.wantValue,
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
