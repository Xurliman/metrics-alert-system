package services

import (
	"context"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/models"
)

func ExampleMetricsService_SendMetric() {
	someGaugeMetric := &models.Metrics{
		ID:    "someGaugeMetric",
		MType: constants.GaugeMetricType,
		Value: &someGaugeMetricVal,
	}
	ctx := context.Background()

	_ = service.SendMetric(ctx, someGaugeMetric)
	//	Output:
	//	nil

}

func ExampleMetricsService_SendCompressedMetric() {
	someCounterMetric := &models.Metrics{
		ID:    "someCounterMetric",
		MType: constants.GaugeMetricType,
		Delta: &someCounterMetricVal,
	}
	ctx := context.Background()

	_ = service.SendCompressedMetric(ctx, someCounterMetric)
	//  Output:
	//  invalid metrics value for gauge type
}

func ExampleMetricsService_SendMetricWithParams() {
	someGaugeMetric := &models.Metrics{
		ID:    "someGaugeMetric",
		MType: constants.CounterMetricType,
		Value: &someGaugeMetricVal,
	}
	ctx := context.Background()

	_ = service.SendMetricWithParams(ctx, someGaugeMetric)
	//  Output:
	//  invalid metrics value for counter type
}

func ExampleMetricsService_SendBatchMetrics() {
	_ = service.SendBatchMetrics()
	//	Output:
	//	nil
}
