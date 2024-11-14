package services

import (
	"context"
	"errors"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"strconv"
)

type MetricsService struct {
	sw         interfaces.Switcher
	repository interfaces.MetricsRepositoryInterface
}

func NewMetricsService(repository interfaces.MetricsRepositoryInterface, switcher interfaces.Switcher) *MetricsService {
	return &MetricsService{
		repository: repository,
		sw:         switcher,
	}
}

func (s *MetricsService) List() (map[string]string, error) {
	data := make(map[string]string)
	metricsCollection, err := s.repository.List()
	if err != nil {
		return nil, err
	}
	for metricName, metric := range metricsCollection {
		switch metric.MType {
		case constants.GaugeMetricType:
			data[metricName] = gaugeStrVal(metric)
		case constants.CounterMetricType:
			data[metricName] = counterStrVal(metric)
		default:
			continue
		}
	}
	return data, nil
}

func (s *MetricsService) GetMetricValue(metricType, metricName string) (metricValue string, err error) {
	metric, err := s.repository.FindByName(metricName)
	if err != nil {
		return "", err
	}
	switch metricType {
	case constants.GaugeMetricType:
		metricValue = gaugeStrVal(metric)
	case constants.CounterMetricType:
		metricValue = counterStrVal(metric)
	}
	return metricValue, nil
}

func (s *MetricsService) SaveWhenParams(metricType, metricName, metricValue string) (err error) {
	if metricName == "" {
		return constants.ErrEmptyMetricName
	}

	if metricType != constants.GaugeMetricType && metricType != constants.CounterMetricType {
		return constants.ErrInvalidMetricType
	}

	existingMetric, err := s.repository.FindByName(metricName)
	if err != nil && !errors.Is(err, constants.ErrMetricNotFound) {
		return err
	}

	var metric *models.Metrics
	if metricType == constants.GaugeMetricType {
		metric, err = s.sw.ConvertParams(GaugeConverter, existingMetric, metricName, metricValue)
		if err != nil {
			return err
		}
	}

	if metricType == constants.CounterMetricType {
		metric, err = s.sw.ConvertParams(CounterConverter, existingMetric, metricName, metricValue)
		if err != nil {
			return err
		}
	}

	_, err = s.repository.Save(metric)
	if err != nil {
		return err
	}

	return nil
}

func (s *MetricsService) SaveWhenBody(metricRequest requests.MetricsSaveRequest) (entry *models.Metrics, err error) {
	existingMetric, err := s.repository.FindByName(metricRequest.ID)
	if err != nil && !errors.Is(err, constants.ErrMetricNotFound) {
		return nil, err
	}
	var metric *models.Metrics

	switch metricRequest.MType {
	case constants.GaugeMetricType:
		metric, err = s.sw.ConvertRequest(GaugeConverter, existingMetric, metricRequest)
		if err != nil {
			return nil, err
		}
	case constants.CounterMetricType:
		metric, err = s.sw.ConvertRequest(CounterConverter, existingMetric, metricRequest)
		if err != nil {
			return nil, err
		}
	default:
		return nil, constants.ErrInvalidMetricType
	}

	result, err := s.repository.Save(metric)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MetricsService) Show(metricName string) (metric *models.Metrics, err error) {
	return s.repository.FindByName(metricName)
}

func (s *MetricsService) Ping(ctx context.Context) error {
	return s.repository.Ping(ctx)
}

func (s *MetricsService) SaveMany(ctx context.Context, request []requests.MetricsSaveRequest) (err error) {
	metrics := make(map[string]*models.Metrics)
	for _, metricRequest := range request {
		var metric *models.Metrics

		switch metricRequest.MType {
		case constants.GaugeMetricType:
			metric, err = s.sw.ConvertRequest(GaugeConverter, nil, metricRequest)
			if err != nil {
				return err
			}
		case constants.CounterMetricType:
			if existingMetric, ok := metrics[metricRequest.ID]; ok {
				metric, err = s.sw.ConvertRequest(CounterConverter, existingMetric, metricRequest)
				if err != nil {
					return err
				}
			} else {
				metric, err = s.sw.ConvertRequest(CounterConverter, nil, metricRequest)
				if err != nil {
					return err
				}
			}
		default:
			return constants.ErrInvalidMetricType
		}

		metrics[metricRequest.ID] = metric
	}

	var metricsArr []*models.Metrics
	for _, metric := range metrics {
		metricsArr = append(metricsArr, metric)
	}
	err = s.repository.InsertMany(ctx, metricsArr)
	if err != nil {
		return err
	}

	return nil
}

func gaugeStrVal(metric *models.Metrics) (val string) {
	return strconv.FormatFloat(*metric.Value, 'f', -1, 64)
}

func counterStrVal(metric *models.Metrics) (val string) {
	return strconv.FormatInt(*metric.Delta, 10)
}
