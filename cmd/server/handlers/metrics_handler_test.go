package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	type expected struct {
		statusCode int
	}
	tests := []struct {
		name     string
		method   string
		url      string
		expected expected
	}{
		//{
		//	name:   "first",
		//	method: http.MethodPost,
		//	url:    "http://localhost:8080/update/counter/someMetric/527",
		//	expected: expected{
		//		statusCode: http.StatusOK,
		//	},
		//},
		{
			name:   "second",
			method: http.MethodGet,
			url:    "http://localhost:8080/update/gauge/someMetric/527",
			expected: expected{
				statusCode: http.StatusBadRequest,
			},
		},
		//{
		//	name:   "third",
		//	method: http.MethodPost,
		//	url:    "http://localhost:8080/update/gauge/someMetric/527",
		//	expected: expected{
		//		statusCode: http.StatusOK,
		//	},
		//},
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
			r := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			UpdateMetrics(w, r)
			result := w.Result()
			assert.Equal(t, test.expected.statusCode, result.StatusCode)
			defer result.Body.Close()
		})
	}
}
