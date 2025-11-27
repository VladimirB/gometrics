package mapper

import models "github.com/VladimirB/gometrics/internal/model"

type MetricResponse struct {
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

func MetricToResponse(metric models.Metrics) MetricResponse {
	var result = MetricResponse{
		ID: metric.ID,
		Type: metric.MType,
	}

	switch metric.MType {
	case models.Counter:
		result.Value = float64(*metric.Delta)
	case models.Gauge:
		result.Value = *metric.Value
	}

	return result
}

func MetricToValue(metric models.Metrics) float64 {
	if metric.MType == models.Counter {
		return float64(*metric.Delta)
	} else {
		return *metric.Value
	}
}
