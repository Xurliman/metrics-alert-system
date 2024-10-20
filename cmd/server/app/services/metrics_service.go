package services

import (
	"errors"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"strconv"
)

type MetricsService struct {
}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

var metrics = models.ExistingMetrics()

func (s *MetricsService) FindByName(metricsType string, metricsName string) (metricsValue string, err error) {
	if metricsType == "gauge" {
		value, ok := metrics.Gauge[metricsName]
		if !ok {
			return metricsValue, fmt.Errorf("metrics %s not found", metricsName)
		}
		metricsValue = strconv.FormatFloat(value, 'f', -1, 64)
	} else if metricsType == "counter" {
		value, ok := metrics.Counter[metricsName]
		if !ok {
			return metricsValue, fmt.Errorf("metrics %s not found", metricsName)
		}
		metricsValue = strconv.FormatInt(value, 10)
	}
	return metricsValue, nil
}

func (s *MetricsService) GetAll() map[string]string {
	data := make(map[string]string)
	for metricName, metricValue := range metrics.Gauge {
		data[metricName] = strconv.FormatFloat(metricValue, 'f', -1, 64)
	}
	for metricName, metricValue := range metrics.Counter {
		data[metricName] = strconv.FormatInt(metricValue, 10)
	}
	return data
}

func (s *MetricsService) Save(metricsType string, metricsName string, metricsValue string) error {
	if metricsType == "counter" {
		metricVal, err := strconv.ParseInt(metricsValue, 10, 64)
		if err != nil {
			return errors.New("invalid metrics value for counter type")
		}
		metrics.Counter[metricsName] += metricVal
	} else if metricsType == "gauge" {
		metricVal, err := strconv.ParseFloat(metricsValue, 64)
		if err != nil {
			return errors.New("invalid metrics value for counter type")
		}
		metrics.Gauge[metricsName] = metricVal
	} else {
		return errors.New("invalid metrics type")
	}
	return nil
}
