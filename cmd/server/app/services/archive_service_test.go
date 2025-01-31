package services

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/stretchr/testify/assert"
)

func callOnceToChangeDir() {
	err := os.Chdir("../../../../")
	if err != nil {
		log.Fatalf("Failed to set working directory: %v", err)
	}
}

func TestArchiveService_Archive(t *testing.T) {
	callOnceToChangeDir()
	var counterMetricExample int64
	var gaugeMetricExample float64
	type fields struct {
		path string
	}
	type args struct {
		metrics map[string]*models.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			fields: fields{
				path: constants.DefaultFileStoragePath,
			},
			args: args{
				metrics: map[string]*models.Metrics{
					"PollCount": {
						ID:    "PollCount",
						MType: constants.CounterMetricType,
						Delta: &counterMetricExample,
					},
					"SomeGaugeMetric": {
						ID:    "SomeGaugeMetric",
						MType: constants.GaugeMetricType,
						Value: &gaugeMetricExample,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArchiveService(tt.fields.path)
			tt.wantErr(t, a.Archive(tt.args.metrics), fmt.Sprintf("Archive(%v)", tt.args.metrics))
		})
	}
}

func TestArchiveService_Load(t *testing.T) {
	var counterMetricExample int64
	var gaugeMetricExample float64

	type fields struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]*models.Metrics
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			fields: fields{
				path: constants.DefaultFileStoragePath,
			},
			want: map[string]*models.Metrics{
				"PollCount": {
					ID:    "PollCount",
					MType: constants.CounterMetricType,
					Delta: &counterMetricExample,
				},
				"SomeGaugeMetric": {
					ID:    "SomeGaugeMetric",
					MType: constants.GaugeMetricType,
					Value: &gaugeMetricExample,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := NewArchiveService(tt.fields.path)
			got, err := a.Load()
			if !tt.wantErr(t, err, fmt.Sprintf("Load() %v %v", got, err)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Load()")
		})
	}
}
