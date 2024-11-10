package services

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_counterConverterService_ParamsToMetric(t *testing.T) {
	var (
		counterMetricValue       int64
		resultCounterValue       = counterMetricValue + 1
		existingMetric           *models.Metrics
		secondCounterMetricValue int64 = 3293
	)
	type args struct {
		existingMetric *models.Metrics
		metricName     string
		metricValue    string
	}
	tests := []struct {
		name       string
		args       args
		wantMetric *models.Metrics
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			args: args{
				existingMetric: &models.Metrics{
					ID:    "PollCount",
					MType: constants.CounterMetricType,
					Delta: &counterMetricValue,
				},
				metricName:  "PollCount",
				metricValue: "1",
			},
			wantMetric: &models.Metrics{
				ID:    "PollCount",
				MType: constants.CounterMetricType,
				Delta: &resultCounterValue,
			},
			wantErr: assert.NoError,
		},
		{
			name: "second",
			args: args{
				existingMetric: existingMetric,
				metricName:     "NewMetric",
				metricValue:    "3293",
			},
			wantMetric: &models.Metrics{
				ID:    "NewMetric",
				MType: constants.CounterMetricType,
				Delta: &secondCounterMetricValue,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := counterConverterService{}
			gotMetric, err := s.ParamsToMetric(tt.args.existingMetric, tt.args.metricName, tt.args.metricValue)
			if !tt.wantErr(t, err, fmt.Sprintf("ParamsToMetric(%v, %v, %v)", tt.args.existingMetric, tt.args.metricName, tt.args.metricValue)) {
				return
			}
			assert.Equalf(t, tt.wantMetric, gotMetric, "ParamsToMetric(%v, %v, %v)", tt.args.existingMetric, tt.args.metricName, tt.args.metricValue)
		})
	}
}

func Test_counterConverterService_RequestToMetric(t *testing.T) {
	var (
		counterMetricValue int64 = 989
		existingMetric     *models.Metrics
	)

	type args struct {
		existingMetric *models.Metrics
		metricRequest  requests.MetricsSaveRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Metrics
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			args: args{
				existingMetric: existingMetric,
				metricRequest: requests.MetricsSaveRequest{
					ID:    "SomeMetric",
					MType: constants.CounterMetricType,
					Delta: &counterMetricValue,
				},
			},
			want: &models.Metrics{
				ID:    "SomeMetric",
				MType: constants.CounterMetricType,
				Delta: &counterMetricValue,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := counterConverterService{}
			got, err := s.RequestToMetric(tt.args.existingMetric, tt.args.metricRequest)
			if !tt.wantErr(t, err, fmt.Sprintf("RequestToMetric(%v, %v)", tt.args.existingMetric, tt.args.metricRequest)) {
				return
			}
			assert.Equalf(t, tt.want, got, "RequestToMetric(%v, %v)", tt.args.existingMetric, tt.args.metricRequest)
		})
	}
}
