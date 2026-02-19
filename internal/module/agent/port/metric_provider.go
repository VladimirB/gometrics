package port

import (
	"context"

	"github.com/VladimirB/gometrics/internal/domain"
)

type MetricProvider interface {
	Send(ctx context.Context, metric domain.Metric) error

	SendAll(ctx context.Context, metrics []domain.Metric) error
}
