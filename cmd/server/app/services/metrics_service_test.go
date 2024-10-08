package services

import (
	"testing"
)

func TestMetricsService_FindByName(t *testing.T) {
	type args struct {
		metricsType string
		metricsName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "first",
			args: args{
				metricsType: "gauge",
				metricsName: "test",
			},
			wantErr: true,
		},
		{
			name: "second",
			args: args{
				metricsType: "gauge",
				metricsName: "HeapIdle",
			},
			wantErr: false,
		},
		{
			name: "third",
			args: args{
				metricsType: "gauge",
				metricsName: "TotalAlloc",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MetricsService{}
			_, err := s.FindByName(tt.args.metricsType, tt.args.metricsName)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
