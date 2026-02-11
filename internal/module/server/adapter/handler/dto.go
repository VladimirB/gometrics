package handler

import (
	"fmt"

	"github.com/VladimirB/gometrics/internal/domain"
)

type metricRequest struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (r metricRequest) ToDomain() domain.Metric {
	return domain.Metric{
		ID:    r.ID,
		MType: r.MType,
		Delta: r.Delta,
		Value: r.Value,
	}
}

func (r metricRequest) String() string {
	return fmt.Sprintf("MetricRequest(ID: %s, MType: %s, Delta: %p, Value: %p)", r.ID, r.MType, r.Delta, r.Value)
}

type MetricResponse struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`
}

func mapToMetricResponse(m domain.Metric) MetricResponse {
	return MetricResponse{
		ID:    m.ID,
		MType: m.MType,
		Delta: m.Delta,
		Value: m.Value,
	}
}

func (r MetricResponse) String() string {
	return fmt.Sprintf("MetricResponse(ID: %s, MType: %s, Delta: %p, Value: %p)", r.ID, r.MType, r.Delta, r.Value)
}
