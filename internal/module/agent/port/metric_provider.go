package port

import (
	"context"

	models "github.com/VladimirB/gometrics/internal/model"
)

type MetricProvider interface {
	Send(ctx context.Context, metric models.Metrics) error
	SendAsJSON(ctx context.Context, metric models.Metrics) error
}