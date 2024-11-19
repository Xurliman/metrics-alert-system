package requests

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/models"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
)

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
