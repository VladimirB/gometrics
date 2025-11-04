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
	// Выглядит как кусок бизнес-логики. Поискать другое место под это
	if m, ok := ms.storage[metricName]; ok && metric.MType == models.Counter {
		*metric.Value += *m.Value
	}
	ms.storage[metricName] = metric
}

func (ms *MemStorage) Get(metricName string) (models.Metrics, error) {
	if metric, ok := ms.storage[metricName]; !ok {
		return models.Metrics{}, fmt.Errorf("metric %q not found", metricName)
	} else {
		return metric, nil
	}
}

func (ms *MemStorage) GetAll() map[string]models.Metrics {
	return ms.storage
}
