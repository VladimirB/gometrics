package port

import (
	"context"

	"github.com/VladimirB/gometrics/internal/domain"
)

type MetricRepository interface {
	Save(ctx context.Context, metric domain.Metric) error

	SaveAll(ctx context.Context, metrics []domain.Metric) error

	Get(ctx context.Context, metricID string) (domain.Metric, error)

	GetAll(ctx context.Context) (map[string]domain.Metric, error)
}
