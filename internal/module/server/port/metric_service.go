package port

import (
	"context"

	models "github.com/VladimirB/gometrics/internal/model"
)

type MetricService interface {
	Save(ctx context.Context, metric models.Metrics) error

	SaveByFields(ctx context.Context, metricType string, metricID string, value float64) error

	Get(ctx context.Context, metricID string) (models.Metrics, error)

	GetAll(ctx context.Context) map[string]models.Metrics
}
