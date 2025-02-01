// Package requests define a struct which is used to wrap metrics while sending a request to the server
package requests

type MetricsRequest struct {
	ID    string   `json:"id" validate:"required"`
	MType string   `json:"type" validate:"required"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
