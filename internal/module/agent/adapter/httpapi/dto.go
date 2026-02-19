package httpapi

import "github.com/VladimirB/gometrics/internal/domain"

type metricDTO struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`
}

func mapToMetricDTO(m domain.Metric) metricDTO {
	return metricDTO{
		ID:    m.ID,
		MType: m.MType,
		Delta: m.Delta,
		Value: m.Value,
	}
}
