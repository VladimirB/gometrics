package service

import (
	model "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/repository"
)

type MetricsService struct {
	storage *repository.MemStorage
}

func NewMetricsService(storage *repository.MemStorage) *MetricsService {
	return &MetricsService{
		storage: storage,
	}
}

func (s *MetricsService) Save(metric model.Metrics) error {
	if err := metric.Validate(); err != nil {
		return err
	}

	// Значения для счетчика необходимо сохранять в Delta
	if metric.MType == model.Counter {
		if metric.Delta == nil {
			metric.Delta = new(int64)
		}
		*metric.Delta = int64(*metric.Value)
		metric.Value = nil
	}

	s.storage.Save(metric.ID, metric)

	return nil
}

func (s *MetricsService) SaveByFields(metricType string, metricId string, value float64) error {
	metric := model.Metrics{
		ID:    metricId,
		MType: metricType,
		Value: new(float64),
	}
	metric.Value = &value

	return s.Save(metric)
}

func (s *MetricsService) Get(metricId string) (model.Metrics, error) {
	return s.storage.Get(metricId)
}
