package repository

import (
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

func (ms *MemStorage) Save(metricName string, metric domain.Metric) {
	ms.storage[metricName] = metric
}

func (ms *MemStorage) Get(metricName string) (domain.Metric, error) {
	if metric, ok := ms.storage[metricName]; !ok {
		return domain.Metric{}, fmt.Errorf("metric %q not found", metricName)
	} else {
		return metric, nil
	}
}

func (ms *MemStorage) GetAll() map[string]domain.Metric {
	return ms.storage
}
