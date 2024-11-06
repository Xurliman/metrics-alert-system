package requests

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
)

// MetricsShowRequest Could you please give instructions on how to evaluate my custom validator func "required_if"
// I wanted to make smth like omitempty if required_if=MType:counter is not true,
// but the following is not even getting to the logic where I've implemented this custom func...
// another option was writing omitempty OR required_if but this doesn't work either:
// panicking Undefined validation function 'omitempty' on field 'Value'
type MetricsShowRequest struct {
	ID    string   `json:"id" validate:"required"`
	MType string   `json:"type" validate:"required,oneof=gauge counter"`
	Value *float64 `json:"value,omitempty" validate:"omitempty,required_if=MType:gauge"`
	Delta *int64   `json:"delta,omitempty" validate:"omitempty,required_if=MType:counter"`
}

func (r *MetricsShowRequest) Validate() error {
	return utils.ExtractValidationErrors(r)
}
func (r *MetricsShowRequest) ToModel() models.Metrics {
	return models.Metrics{
		ID:    r.ID,
		MType: r.MType,
		Value: r.Value,
		Delta: r.Delta,
	}
}

type MetricsSaveRequest struct {
	ID    string   `json:"id" validate:"required"`
	MType string   `json:"type" validate:"required,oneof=gauge counter"`
	Value *float64 `json:"value,omitempty" validate:"omitempty,required_if=MType:gauge"`
	Delta *int64   `json:"delta,omitempty" validate:"omitempty,required_if=MType:counter"`
}

func (r *MetricsSaveRequest) Validate() error {
	return utils.ExtractValidationErrors(r)
}
func (r *MetricsSaveRequest) ToModel() models.Metrics {
	return models.Metrics{
		ID:    r.ID,
		MType: r.MType,
		Value: r.Value,
		Delta: r.Delta,
	}
}
