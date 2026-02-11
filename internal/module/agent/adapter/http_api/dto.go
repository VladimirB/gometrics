package http_api

import models "github.com/VladimirB/gometrics/internal/model"

type metricRequest struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`
}

func mapToMetricRequest(m models.Metrics) metricRequest {
	return metricRequest{
		ID: m.ID,
		MType: m.MType,
		Delta: m.Delta,
		Value: m.Value,
	}
}
