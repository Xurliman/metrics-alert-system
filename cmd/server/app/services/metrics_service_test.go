package services

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricsService_FindCounterMetric(t *testing.T) {
	type args struct {
		metricsName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "first",
			args:    args{metricsName: "test"},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{}
			_, err := s.FindCounterMetric(tt.args.metricsName)
			if !tt.wantErr(t, err, fmt.Sprintf("FindCounterMetric(%v)", tt.args.metricsName)) {
				return
			}
		})
	}
}

func TestMetricsService_FindGaugeMetric(t *testing.T) {
	type args struct {
		metricsName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "first",
			args:    args{metricsName: "HeapIdle"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{}
			_, err := s.FindGaugeMetric(tt.args.metricsName)
			if !tt.wantErr(t, err, fmt.Sprintf("FindGaugeMetric(%v)", tt.args.metricsName)) {
				return
			}
		})
	}
}

func TestMetricsService_FindMetricByName(t *testing.T) {
	type args struct {
		metricsType string
		metricsName string
	}
	tests := []struct {
		name             string
		args             args
		wantMetricsValue string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name:             "first",
			args:             args{metricsType: "gauge", metricsName: "test"},
			wantMetricsValue: "test",
			wantErr:          assert.Error,
		},
		{
			name:             "second",
			args:             args{metricsType: "gauge", metricsName: "HeapIdle"},
			wantMetricsValue: "test",
			wantErr:          assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{}
			_, err := s.FindMetricByName(tt.args.metricsType, tt.args.metricsName)
			if !tt.wantErr(t, err, fmt.Sprintf("FindMetricByName(%v, %v)", tt.args.metricsType, tt.args.metricsName)) {
				return
			}
		})
	}
}

func TestMetricsService_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Get All",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{}
			output := s.GetAll()
			if !tt.wantErr {
				assert.Equal(t, len(output), 0)
			}
		})
	}
}

func TestMetricsService_SaveCounterMetric(t *testing.T) {
	type args struct {
		metricsName  string
		metricsValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Save Counter Metric",
			args: args{
				metricsName:  "CounterMetric",
				metricsValue: "0980",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{}
			tt.wantErr(t, s.SaveCounterMetric(tt.args.metricsName, tt.args.metricsValue), fmt.Sprintf("SaveCounterMetric(%v, %v)", tt.args.metricsName, tt.args.metricsValue))
		})
	}
}

func TestMetricsService_SaveGaugeMetric(t *testing.T) {
	type args struct {
		metricsName  string
		metricsValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Save Gauge Metric",
			args: args{
				metricsName:  "HeapIdle",
				metricsValue: "23109",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{}
			tt.wantErr(t, s.SaveGaugeMetric(tt.args.metricsName, tt.args.metricsValue), fmt.Sprintf("SaveGaugeMetric(%v, %v)", tt.args.metricsName, tt.args.metricsValue))
		})
	}
}

func TestMetricsService_SaveMetric(t *testing.T) {
	type args struct {
		metricsType  string
		metricsName  string
		metricsValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "saveGauge",
			args: args{
				metricsName:  "test1",
				metricsValue: "2302930",
				metricsType:  "gauge",
			},
			wantErr: assert.NoError,
		},
		{
			name: "saveCounter",
			args: args{
				metricsName:  "test2",
				metricsValue: "2302930",
				metricsType:  "counter",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{}
			tt.wantErr(t, s.SaveMetric(tt.args.metricsType, tt.args.metricsName, tt.args.metricsValue), fmt.Sprintf("SaveMetric(%v, %v, %v)", tt.args.metricsType, tt.args.metricsName, tt.args.metricsValue))
		})
	}
}
