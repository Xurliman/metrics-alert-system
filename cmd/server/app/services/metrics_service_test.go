package services

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"

	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/mocks/servicemocks"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestMetricsService_List(t *testing.T) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	gaugeMetricExample := float64(memStats.Alloc)
	var counterMetricExample int64
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
			repo := servicemocks.NewMetricsRepositoryInterface(t)
			sw := servicemocks.NewSwitcher(t)
			repo.On("List").Return(map[string]*models.Metrics{
				"Alloc": {
					ID:    "Alloc",
					MType: constants.GaugeMetricType,
					Value: &gaugeMetricExample,
				},
				"PollCount": {
					ID:    "PollCount",
					MType: constants.CounterMetricType,
					Delta: &counterMetricExample,
				},
			}, nil)
			s := NewMetricsService(repo, sw)
			metrics, err := s.List()
			assert.NoError(t, err)
			assert.Equal(t, len(metrics), 2)
		})
	}
}

func TestMetricsService_GetMetricValue(t *testing.T) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	wantGaugeMetric := float64(memStats.Alloc)
	var counterMetricExample int64
	tests := []struct {
		name            string
		metricName      string
		metricType      string
		metric          *models.Metrics
		wantMetricValue string
		wantErr         bool
	}{
		{
			name:       "first",
			metricName: "Alloc",
			metricType: constants.GaugeMetricType,
			metric: &models.Metrics{
				ID:    "Alloc",
				MType: constants.GaugeMetricType,
				Value: &wantGaugeMetric,
			},
			wantMetricValue: strconv.FormatFloat(wantGaugeMetric, 'f', -1, 64),
			wantErr:         false,
		},
		{
			name:            "second",
			metricName:      "HeapAlloc",
			metricType:      constants.GaugeMetricType,
			metric:          &models.Metrics{},
			wantMetricValue: "",
			wantErr:         true,
		},
		{
			name:       "third",
			metricName: "PollCount",
			metricType: constants.CounterMetricType,
			metric: &models.Metrics{
				ID:    "PollCount",
				MType: constants.CounterMetricType,
				Delta: &counterMetricExample,
			},
			wantMetricValue: "0",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := servicemocks.NewMetricsRepositoryInterface(t)
			sw := servicemocks.NewSwitcher(t)
			s := NewMetricsService(repo, sw)
			if tt.wantErr {
				repo.On("FindByName", tt.metricName).Return(nil, constants.ErrMetricNotFound)
			} else {
				repo.On("FindByName", tt.metricName).Return(tt.metric, nil)
			}

			gotMetricValue, err := s.GetMetricValue(tt.metricType, tt.metricName)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equalf(t, tt.wantMetricValue, gotMetricValue, "GetMetricValue(%v, %v)", tt.metric, tt.metricName)
		})
	}
}

func TestMetricsService_SaveWhenParams(t *testing.T) {
	someGaugeMetric := float64(232039203)
	someCounterMetric := int64(54854954)
	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name    string
		args    args
		metric  *models.Metrics
		conv    interfaces.Converter
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "first",
			args: args{
				metricType:  constants.GaugeMetricType,
				metricName:  "Foo",
				metricValue: strconv.FormatFloat(someGaugeMetric, 'f', -1, 64),
			},
			metric: &models.Metrics{
				ID:    "Foo",
				MType: constants.GaugeMetricType,
				Value: &someGaugeMetric,
			},
			conv:    GaugeConverter,
			wantErr: assert.NoError,
		},
		{
			name: "second",
			args: args{
				metricType:  constants.CounterMetricType,
				metricName:  "Bar",
				metricValue: strconv.FormatInt(someCounterMetric, 10),
			},
			metric: &models.Metrics{
				ID:    "Bar",
				MType: constants.CounterMetricType,
				Delta: &someCounterMetric,
			},
			conv:    CounterConverter,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := servicemocks.NewMetricsRepositoryInterface(t)
			sw := servicemocks.NewSwitcher(t)
			s := NewMetricsService(repo, sw)
			var existingMetric *models.Metrics
			repo.On("FindByName", tt.args.metricName).Return(existingMetric, constants.ErrMetricNotFound)
			sw.On("ConvertParams", tt.conv, existingMetric, tt.args.metricName, tt.args.metricValue).Return(tt.metric, nil)
			repo.On("Save", tt.metric).Return(tt.metric, nil)
			tt.wantErr(t, s.SaveWhenParams(tt.args.metricType, tt.args.metricName, tt.args.metricValue), fmt.Sprintf("SaveWhenParams(%v, %v, %v)", tt.args.metricType, tt.args.metricName, tt.args.metricValue))
		})
	}
}

func TestMetricsService_SaveWhenBody(t *testing.T) {
	someGaugeMetric := float64(290)
	someCounterMetric := int64(54854954)
	type args struct {
		metricRequest requests.MetricsSaveRequest
	}
	tests := []struct {
		name      string
		args      args
		conv      interfaces.Converter
		metric    *models.Metrics
		wantEntry *models.Metrics
		wantErr   assert.ErrorAssertionFunc
		exists    bool
	}{
		{
			name: "first",
			args: args{
				metricRequest: requests.MetricsSaveRequest{
					ID:    "SomeGaugeMetric",
					MType: constants.GaugeMetricType,
					Value: &someGaugeMetric,
				},
			},
			conv: GaugeConverter,
			wantEntry: &models.Metrics{
				ID:    "SomeGaugeMetric",
				MType: constants.GaugeMetricType,
				Value: &someGaugeMetric,
			},
			metric: &models.Metrics{
				ID:    "SomeGaugeMetric",
				MType: constants.GaugeMetricType,
				Value: &someGaugeMetric,
			},
			wantErr: assert.NoError,
		},
		{
			name: "second",
			args: args{
				metricRequest: requests.MetricsSaveRequest{
					ID:    "SomeCounterMetric",
					MType: constants.CounterMetricType,
					Delta: &someCounterMetric,
				},
			},
			conv: CounterConverter,
			wantEntry: &models.Metrics{
				ID:    "SomeCounterMetric",
				MType: constants.CounterMetricType,
				Delta: &someCounterMetric,
			},
			metric: &models.Metrics{
				ID:    "SomeCounterMetric",
				MType: constants.CounterMetricType,
				Delta: &someCounterMetric,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := servicemocks.NewMetricsRepositoryInterface(t)
			sw := servicemocks.NewSwitcher(t)
			s := NewMetricsService(repo, sw)
			var existingMetric *models.Metrics
			repo.On("FindByName", tt.args.metricRequest.ID).Return(nil, constants.ErrMetricNotFound)
			sw.On("ConvertRequest", tt.conv, existingMetric, tt.args.metricRequest).Return(tt.metric, nil)
			repo.On("Save", tt.metric).Return(tt.metric, nil)
			gotEntry, err := s.SaveWhenBody(tt.args.metricRequest)
			if !tt.wantErr(t, err, fmt.Sprintf("SaveWhenBody(%v)", tt.args.metricRequest)) {
				return
			}
			assert.Equalf(t, tt.wantEntry, gotEntry, "SaveWhenBody(%v)", tt.args.metricRequest)
		})
	}
}

func TestMetricsService_Show(t *testing.T) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	wantAllocValue := float64(memStats.Alloc)

	type args struct {
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
				metricName: "Foo",
			},
			wantEntry: nil,
			wantErr:   assert.Error,
		},
		{
			name: "second",
			args: args{
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
			repo := servicemocks.NewMetricsRepositoryInterface(t)
			sw := servicemocks.NewSwitcher(t)
			s := NewMetricsService(repo, sw)
			if tt.wantEntry == nil {
				repo.On("FindByName", tt.args.metricName).Return(nil, constants.ErrMetricNotFound)
			} else {
				repo.On("FindByName", tt.args.metricName).Return(tt.wantEntry, nil)
			}
			gotEntry, err := s.Show(tt.args.metricName)
			if !tt.wantErr(t, err, fmt.Sprintf("FindByName(%v)", tt.args.metricName)) {
				return
			}
			assert.Equalf(t, tt.wantEntry, gotEntry, "FindByName(%v)", tt.args.metricName)
		})
	}
}
