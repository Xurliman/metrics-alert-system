package services

import (
	"github.com/stretchr/testify/assert"
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
		//{
		//	name: "second",
		//	args: args{
		//		metricsType: "gauge",
		//		metricsName: "HeapIdle",
		//	},
		//	wantErr: false,
		//},
		//{
		//	name: "third",
		//	args: args{
		//		metricsType: "gauge",
		//		metricsName: "TotalAlloc",
		//	},
		//	wantErr: false,
		//},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &MetricsService{}
			_, err := s.FindByName(test.args.metricsType, test.args.metricsName)
			if (err != nil) != test.wantErr {
				assert.Equal(t, test.wantErr, err != nil, "FindByName() error = %v, wantErr %v", err, test.wantErr)
				return
			}
		})
	}
}
