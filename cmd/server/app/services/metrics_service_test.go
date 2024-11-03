package services

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/mocks/servicemocks"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/stretchr/testify/assert"
	"runtime"
	"strconv"
	"testing"
)

func TestMetricsService_GetMetricValue(t *testing.T) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	wantGaugeMetric := memStats.Alloc
	type args struct {
		metric     interfaces.MetricsInterface
		metricName string
	}
	tests := []struct {
		name            string
		args            args
		wantMetricValue string
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			args: args{
				metric:     Gauge,
				metricName: "Alloc",
			},
			wantMetricValue: strconv.FormatUint(wantGaugeMetric, 10),
			wantErr:         assert.NoError,
		},
		{
			name: "second",
			args: args{
				metric:     Gauge,
				metricName: "HeapAlloc",
			},
			wantMetricValue: "",
			wantErr:         assert.Error,
		},
		{
			name: "third",
			args: args{
				metric:     Counter,
				metricName: "PollCount",
			},
			wantMetricValue: "0",
			wantErr:         assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMetricsService()
			metricsInterface := servicemocks.NewMetricsInterface(t)

			if tt.wantMetricValue != "" {
				metricsInterface.On("GetMetricValue", tt.args.metricName).Return(tt.wantMetricValue, nil)
			} else {
				metricsInterface.On("GetMetricValue", tt.args.metricName).Return(tt.wantMetricValue, constants.ErrMetricNotFound)
			}

			gotMetricValue, err := s.GetMetricValue(metricsInterface, tt.args.metricName)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMetricValue(%v, %v)", tt.args.metric, tt.args.metricName)) {
				return
			}
			assert.Equalf(t, tt.wantMetricValue, gotMetricValue, "GetMetricValue(%v, %v)", tt.args.metric, tt.args.metricName)
		})
	}
}

func TestMetricsService_List(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "List All",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMetricsService()
			assert.Equal(t, len(s.List()), 2)
		})
	}
}

func TestMetricsService_SaveWhenBody(t *testing.T) {
	someGaugeMetric := float64(290)

	type args struct {
		metric        interfaces.MetricsInterface
		metricRequest requests.MetricsSaveRequest
	}
	tests := []struct {
		name      string
		args      args
		wantEntry *models.Metrics
		wantErr   assert.ErrorAssertionFunc
		exists    bool
	}{
		{
			name: "first",
			args: args{
				metric: Gauge,
				metricRequest: requests.MetricsSaveRequest{
					ID:    "SomeMetric",
					MType: constants.GaugeMetricType,
					Value: &someGaugeMetric,
				},
			},
			wantEntry: &models.Metrics{
				ID:    "SomeMetric",
				MType: constants.GaugeMetricType,
				Value: &someGaugeMetric,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMetricsService()
			metricsInterface := servicemocks.NewMetricsInterface(t)

			if tt.exists {
				metricsInterface.On("FindByName", tt.args.metricRequest.ID).Return(tt.wantEntry, nil)
				metricsInterface.On("SaveBody", tt.args.metricRequest, tt.wantEntry).Return(&models.Metrics{
					ID:    tt.args.metricRequest.ID,
					MType: constants.GaugeMetricType,
					Value: &someGaugeMetric,
				}, nil)

			} else {
				var existingMetric *models.Metrics
				metricsInterface.On("FindByName", tt.args.metricRequest.ID).Return(nil, constants.ErrMetricNotFound)
				metricsInterface.On("SaveBody", tt.args.metricRequest, existingMetric).Return(&models.Metrics{
					ID:    tt.args.metricRequest.ID,
					MType: constants.GaugeMetricType,
					Value: &someGaugeMetric,
				}, nil)
			}

			gotEntry, err := s.SaveWhenBody(metricsInterface, tt.args.metricRequest)
			if !tt.wantErr(t, err, fmt.Sprintf("SaveWhenBody(%v, %v)", tt.args.metric, tt.args.metricRequest)) {
				return
			}
			assert.Equalf(t, tt.wantEntry, gotEntry, "SaveWhenBody(%v, %v)", tt.args.metric, tt.args.metricRequest)
		})
	}
}

func TestMetricsService_SaveWhenParams(t *testing.T) {
	type args struct {
		metric      interfaces.MetricsInterface
		metricName  string
		metricValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			args: args{
				metric:      Gauge,
				metricName:  "Foo",
				metricValue: "232039203",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMetricsService()
			tt.wantErr(t, s.SaveWhenParams(tt.args.metric, tt.args.metricName, tt.args.metricValue), fmt.Sprintf("SaveWhenParams(%v, %v, %v)", tt.args.metric, tt.args.metricName, tt.args.metricValue))
		})
	}
}

func TestMetricsService_Show(t *testing.T) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	wantAllocValue := float64(memStats.Alloc)

	type args struct {
		metric     interfaces.MetricsInterface
		metricName string
	}
	tests := []struct {
		name      string
		args      args
		wantEntry *models.Metrics
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			args: args{
				metric:     Gauge,
				metricName: "Foo",
			},
			wantEntry: nil,
			wantErr:   assert.Error,
		},
		{
			name: "second",
			args: args{
				metric:     Gauge,
				metricName: "Alloc",
			},
			wantEntry: &models.Metrics{
				ID:    "Alloc",
				MType: constants.GaugeMetricType,
				Value: &wantAllocValue,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMetricsService()
			metricsInterface := servicemocks.NewMetricsInterface(t)

			if tt.wantEntry != nil {
				metricsInterface.On("FindByName", tt.args.metricName).Return(tt.wantEntry, nil)
			} else {
				metricsInterface.On("FindByName", tt.args.metricName).Return(nil, constants.ErrMetricNotFound)
			}

			gotEntry, err := s.Show(metricsInterface, tt.args.metricName)
			if !tt.wantErr(t, err, fmt.Sprintf("Show(%v, %v)", metricsInterface, tt.args.metricName)) {
				return
			}
			assert.Equalf(t, tt.wantEntry, gotEntry, "Show(%v, %v)", metricsInterface, tt.args.metricName)
		})
	}
}
