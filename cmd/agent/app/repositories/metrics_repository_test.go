package repositories

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/requests"
	"github.com/stretchr/testify/assert"
	"runtime"
	"strconv"
	"testing"
)

func TestMetricsRepository_GetRequestURL(t *testing.T) {
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
			want:    "/gauge/Alloc/" + strconv.FormatFloat(allocValue, 'f', -1, 64),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricsRepository{}
			got, err := r.GetRequestURL(tt.metric)
			if !tt.wantErr(t, err, fmt.Sprintf("GetRequestURL(%v)", tt.metric)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetRequestURL(%v)", tt.metric)
		})
	}
}

func TestMetricsRepository_GetPlainRequest(t *testing.T) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	allocValue := float64(memStats.Alloc)
	type args struct {
		metric *models.Metrics
	}
	tests := []struct {
		name        string
		args        args
		wantRequest *requests.MetricsRequest
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			args: args{
				metric: &models.Metrics{
					ID:    "Alloc",
					MType: constants.GaugeMetricType,
					Value: &allocValue,
				},
			},
			wantRequest: &requests.MetricsRequest{
				ID:    "Alloc",
				MType: constants.GaugeMetricType,
				Value: &allocValue,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricsRepository{}
			gotRequest, err := r.GetPlainRequest(tt.args.metric)
			if !tt.wantErr(t, err, fmt.Sprintf("GetPlainRequest(%v)", tt.args.metric)) {
				return
			}
			assert.Equalf(t, tt.wantRequest, gotRequest, "GetPlainRequest(%v)", tt.args.metric)
		})
	}
}

func BenchmarkMetricsRepository_GetAll(b *testing.B) {
	repo := NewMetricsRepository()
	for i := 0; i < b.N; i++ {
		_ = repo.GetAll()
	}
}
