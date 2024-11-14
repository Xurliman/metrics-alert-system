package services

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_gaugeConverterService_ParamsToMetric(t *testing.T) {
	var (
		existingMetric        *models.Metrics
		firstGaugeMetricValue float64 = 3293
		gaugeMetricValue      float64 = 98
		resultGaugeValue      float64 = 32
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
				existingMetric: existingMetric,
				metricName:     "SomeNewGaugeMetric",
				metricValue:    "3293",
			},
			wantMetric: &models.Metrics{
				ID:    "SomeNewGaugeMetric",
				MType: constants.GaugeMetricType,
				Value: &firstGaugeMetricValue,
			},
			wantErr: assert.NoError,
		},
		{
			name: "second",
			args: args{
				existingMetric: &models.Metrics{
					ID:    "SomeExistingGaugeMetric",
					MType: constants.GaugeMetricType,
					Value: &gaugeMetricValue,
				},
				metricName:  "SomeExistingGaugeMetric",
				metricValue: "32",
			},
			wantMetric: &models.Metrics{
				ID:    "SomeExistingGaugeMetric",
				MType: constants.GaugeMetricType,
				Value: &resultGaugeValue,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := gaugeConverterService{}
			gotMetric, err := s.ParamsToMetric(tt.args.existingMetric, tt.args.metricName, tt.args.metricValue)
			if !tt.wantErr(t, err, fmt.Sprintf("ParamsToMetric(%v, %v, %v)", tt.args.existingMetric, tt.args.metricName, tt.args.metricValue)) {
				return
			}
			assert.Equalf(t, tt.wantMetric, gotMetric, "ParamsToMetric(%v, %v, %v)", tt.args.existingMetric, tt.args.metricName, tt.args.metricValue)
		})
	}
}

func Test_gaugeConverterService_RequestToMetric(t *testing.T) {
	type args struct {
		existingMetric *models.Metrics
		metricRequest  requests.MetricsSaveRequest
	}
	tests := []struct {
		name       string
		args       args
		wantMetric *models.Metrics
		wantErr    assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := gaugeConverterService{}
			gotMetric, err := s.RequestToMetric(tt.args.existingMetric, tt.args.metricRequest)
			if !tt.wantErr(t, err, fmt.Sprintf("RequestToMetric(%v, %v)", tt.args.existingMetric, tt.args.metricRequest)) {
				return
			}
			assert.Equalf(t, tt.wantMetric, gotMetric, "RequestToMetric(%v, %v)", tt.args.existingMetric, tt.args.metricRequest)
		})
	}
}
