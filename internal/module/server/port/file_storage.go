package port

import (
	"context"

	"github.com/VladimirB/gometrics/internal/domain"
)

type MetricFileStorage interface {
	Read(ctx context.Context, fileName string) ([]domain.Metric, error)

	Write(ctx context.Context, metrics []domain.Metric, fileName string) error
}
