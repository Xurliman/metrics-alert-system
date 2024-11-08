package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/stretchr/testify/assert"
	"runtime"
	"strconv"
	"testing"
)

func TestMetricsRepository_GetRequestBody(t *testing.T) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	allocValue := float64(memStats.Alloc)

	tests := []struct {
		name    string
		metric  *models.Metrics
		wantMap map[string]interface{}
		wantErr bool
	}{
		{
			name: "first",
			metric: &models.Metrics{
				ID:    "Alloc",
				MType: constants.GaugeMetricType,
				Value: &allocValue,
			},
			wantMap: map[string]interface{}{
				"id":    "Alloc",
				"type":  "gauge",
				"value": &allocValue,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewMetricsRepository()
			got, err := r.GetRequestBody(tt.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRequestBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wantBody, err := json.Marshal(tt.wantMap)
			assert.NoError(t, err)
			assert.Equal(t, string(wantBody), string(got))
		})
	}
}

func TestMetricsRepository_GetRequestUrl(t *testing.T) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	allocValue := float64(memStats.Alloc)

	tests := []struct {
		name    string
		metric  *models.Metrics
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			metric: &models.Metrics{
				ID:    "Alloc",
				MType: constants.GaugeMetricType,
				Value: &allocValue,
			},
			want:    "http://localhost:8080/update/gauge/Alloc/" + strconv.FormatFloat(allocValue, 'f', 6, 64),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricsRepository{}
			address := "localhost:8080"
			got, err := r.GetRequestUrl(tt.metric, address)
			if !tt.wantErr(t, err, fmt.Sprintf("GetRequestUrl(%v, %v)", tt.metric, address)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetRequestUrl(%v, %v)", tt.metric, address)
		})
	}
}
