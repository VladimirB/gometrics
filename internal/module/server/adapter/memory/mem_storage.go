package memory

import (
	"context"
	"fmt"

	"github.com/VladimirB/gometrics/internal/domain"
)

type MemStorage struct {
	storage map[string]domain.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]domain.Metric),
	}
}

func (ms *MemStorage) Save(ctx context.Context, metric domain.Metric) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("repo save aborted: %w", err)
	}

	ms.storage[metric.ID] = metric
	return nil
}

func (ms *MemStorage) Get(ctx context.Context, metricID string) (domain.Metric, error) {
	if err := ctx.Err(); err != nil {
		return domain.Metric{}, fmt.Errorf("repo save aborted: %w", err)
	}

	if metric, ok := ms.storage[metricID]; !ok {
		return domain.Metric{}, fmt.Errorf("metric %q not found", metricID)
	} else {
		return metric, nil
	}
}

func (ms *MemStorage) GetAll(ctx context.Context) map[string]domain.Metric {
	return ms.storage
}
