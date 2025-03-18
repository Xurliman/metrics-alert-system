package services

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"github.com/Xurliman/metrics-alert-system/internal/pb"
	"go.uber.org/zap"
)

type RPCMetricsService struct {
	client pb.MetricsServiceClient
}

func NewRPCMetricsService(client pb.MetricsServiceClient) *RPCMetricsService {
	return &RPCMetricsService{
		client: client,
	}
}

func (s *RPCMetricsService) TestMetrics() {
	gaugeVal := float64(239)
	counterVal := int64(223)

	metrics := []*pb.Metrics{
		{Id: "rpcMetric1", Type: constants.GaugeMetricType, Value: &gaugeVal},
		{Id: "rpcMetric2", Type: constants.CounterMetricType, Delta: &counterVal},
	}

	for _, m := range metrics {
		_, err := s.client.Add(context.Background(), &pb.AddRequest{Metrics: m})
		if err != nil {
			log.Error("error adding metric", zap.Error(err))
		}
	}

	for _, metricName := range []string{"rpcMetric2", "someFakeMetric"} {
		resp, err := s.client.Get(context.Background(), &pb.GetRequest{Name: metricName})
		if err != nil {
			log.Error("failed to get metrics", zap.Error(err))
		} else {
			log.Info("got metrics", zap.Any(metricName, resp.GetMetrics()))
		}
	}
}
