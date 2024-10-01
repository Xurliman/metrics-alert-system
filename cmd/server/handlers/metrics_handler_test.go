package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
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
			url:    "http://localhost:8080/update/counter/someMetric/527",
			expected: expected{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "third",
			method: http.MethodPost,
			url:    "http://localhost:8080/update/counter//527",
			expected: expected{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:   "fourth",
			method: http.MethodPost,
			url:    "http://localhost:8080/update/unknown/someMetric/527",
			expected: expected{
				statusCode: http.StatusNotImplemented,
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
			client := &http.Client{}
			request, err := http.NewRequest(test.method, test.url, nil)
			assert.NoError(t, err)
			response, err := client.Do(request)
			if assert.NoError(t, err) {
				defer response.Body.Close()
				assert.Equal(t, test.expected.statusCode, response.StatusCode)
			} else {
				assert.Fail(t, err.Error())
			}
		})
	}
}
