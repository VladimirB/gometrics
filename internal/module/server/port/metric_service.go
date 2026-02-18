package port

import (
	"context"

	"github.com/VladimirB/gometrics/internal/domain"
)

type MetricService interface {
	Save(ctx context.Context, metric domain.Metric) error

	SaveAll(ctx context.Context, metrics []domain.Metric) error

	SaveByFields(ctx context.Context, metricType string, metricID string, value float64) error

	Get(ctx context.Context, metricID string) (domain.Metric, error)

	GetAll(ctx context.Context) (map[string]domain.Metric, error)

	RestoreFromFile(ctx context.Context, fileName string) error

	DumpToFile(ctx context.Context, fileName string) error
}
