package controllers

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/resources"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RPCMetricsController struct {
	service interfaces.MetricsServiceInterface
	pb.UnimplementedMetricsServiceServer
}

func NewRPCMetricsController(service interfaces.MetricsServiceInterface) *RPCMetricsController {
	return &RPCMetricsController{
		service: service,
	}
}

func (c *RPCMetricsController) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	var out pb.AddResponse
	var metricRequest requests.MetricsSaveRequest

	metricRequest.ID = in.Metrics.Id
	metricRequest.MType = in.Metrics.Type
	metricRequest.Value = in.Metrics.Value
	metricRequest.Delta = in.Metrics.Delta

	metric, err := c.service.SaveWhenBody(metricRequest)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response, err := resources.ToResponse(metric)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	out.Metrics = &pb.Metrics{
		Id:    response.ID,
		Type:  response.MType,
		Value: response.Value,
		Delta: response.Delta,
	}

	return &out, nil
}

func (c *RPCMetricsController) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	var out pb.GetResponse

	metric, err := c.service.Show(in.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response, err := resources.ToResponse(metric)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	out.Metrics = &pb.Metrics{
		Id:    response.ID,
		Type:  response.MType,
		Value: response.Value,
		Delta: response.Delta,
	}

	return &out, nil
}
