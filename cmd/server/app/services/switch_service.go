package services

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/http/requests"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/interfaces"
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
)

type SwitchService struct{}

func NewSwitchService() interfaces.Switcher {
	return &SwitchService{}
}

func (s SwitchService) ConvertParams(converter interfaces.Converter, existingMetric *models.Metrics, metricName, metricValue string) (metric *models.Metrics, err error) {
	return converter.ParamsToMetric(existingMetric, metricName, metricValue)
}

func (s SwitchService) ConvertRequest(converter interfaces.Converter, existingMetric *models.Metrics, metricRequest requests.MetricsSaveRequest) (metric *models.Metrics, err error) {
	return converter.RequestToMetric(existingMetric, metricRequest)
}

var (
	GaugeConverter   = gaugeConverterService{}
	CounterConverter = counterConverterService{}
)
