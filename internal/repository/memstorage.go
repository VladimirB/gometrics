package repository

import (
	"fmt"

	models "github.com/VladimirB/gometrics/internal/model"
)

type MemStorage struct {
	storage map[string]models.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]models.Metrics),
	}
}

func (ms *MemStorage) Save(metricName string, metric models.Metrics) {
	ms.storage[metricName] = metric
}

func (ms *MemStorage) Get(metricName string) (models.Metrics, error) {
	if metric, ok := ms.storage[metricName]; !ok {
		return models.Metrics{}, fmt.Errorf("metric %q not found", metricName)
	} else {
		return metric, nil
	}
}
