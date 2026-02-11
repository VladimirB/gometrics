package port

import (
	"context"

	"github.com/VladimirB/gometrics/internal/domain"
)

type MetricRepository interface {
	Save(ctx context.Context, metric domain.Metric) error
}
