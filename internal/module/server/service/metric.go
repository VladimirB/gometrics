package service

import (
	"context"

	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/module/server/port"
)

type metricService struct {
	repo port.MetricRepository
}

func NewMetricsService(repo port.MetricRepository) port.MetricService {
	return &metricService{
		repo: repo,
	}
}

func (s *metricService) Save(ctx context.Context, metric domain.Metric) error {
	if err := metric.Validate(); err != nil {
		return err
	}

	// Значения для счетчика необходимо сохранять в Delta
	if metric.MType == domain.Counter && metric.Delta != nil {
		if saved, err := s.repo.Get(ctx, metric.ID); err == nil {
			*metric.Delta += *saved.Delta
		}
	}

	s.repo.Save(ctx, metric)

	return nil
}

func (s *metricService) SaveByFields(ctx context.Context, metricType string, metricID string, value float64) error {
	metric := domain.Metric{
		ID:    metricID,
		MType: metricType,
	}

	switch metricType {
	case domain.Counter:
		metric.Delta = new(int64)
		*metric.Delta = int64(value)
	case domain.Gauge:
		metric.Value = new(float64)
		*metric.Value = value
	}

	return s.Save(ctx, metric)
}

func (s *metricService) Get(ctx context.Context, metricID string) (domain.Metric, error) {
	return s.repo.Get(ctx, metricID)
}

func (s *metricService) GetAll(ctx context.Context) map[string]domain.Metric {
	return s.repo.GetAll(ctx)
}
