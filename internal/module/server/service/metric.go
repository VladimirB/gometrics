package service

import (
	"context"

	model "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/module/server/port"
	"github.com/VladimirB/gometrics/internal/repository"
)

type metricService struct {
	storage *repository.MemStorage
}

func NewMetricsService(storage *repository.MemStorage) port.MetricService {
	return &metricService{
		storage: storage,
	}
}

func (s *metricService) Save(ctx context.Context, metric model.Metrics) error {
	if err := metric.Validate(); err != nil {
		return err
	}

	// Значения для счетчика необходимо сохранять в Delta
	if metric.MType == model.Counter && metric.Delta != nil {
		if saved, err := s.storage.Get(metric.ID); err == nil {
			*metric.Delta += *saved.Delta
		}
	}

	s.storage.Save(metric.ID, metric)

	return nil
}

func (s *metricService) SaveByFields(ctx context.Context, metricType string, metricID string, value float64) error {
	metric := model.Metrics{
		ID:    metricID,
		MType: metricType,
	}

	switch metricType {
	case model.Counter:
		metric.Delta = new(int64)
		*metric.Delta = int64(value)
	case model.Gauge:
		metric.Value = new(float64)
		*metric.Value = value
	}

	return s.Save(ctx, metric)
}

func (s *metricService) Get(ctx context.Context, metricID string) (model.Metrics, error) {
	return s.storage.Get(metricID)
}

func (s *metricService) GetAll(ctx context.Context) map[string]model.Metrics {
	return s.storage.GetAll()
}
