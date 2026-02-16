package service

import (
	"context"
	"fmt"

	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/module/server/port"
)

type metricService struct {
	repo        port.MetricRepository
	fileStorage port.MetricFileStorage
}

func NewMetricsService(repo port.MetricRepository, fileStorage port.MetricFileStorage) port.MetricService {
	return &metricService{
		repo:        repo,
		fileStorage: fileStorage,
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

	return s.repo.Save(ctx, metric)
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

func (s *metricService) GetAll(ctx context.Context) (map[string]domain.Metric, error) {
	return s.repo.GetAll(ctx)
}

func (s *metricService) RestoreFromFile(ctx context.Context, fileName string) error {
	metrics, err := s.fileStorage.Read(ctx, fileName)
	if err != nil {
		return fmt.Errorf("error on reading from file: %w", err)
	}

	for _, metric := range metrics {
		if err := s.repo.Save(ctx, metric); err != nil {
			return fmt.Errorf("error on save metric in repo: %w", err)
		}
	}

	return nil
}

func (s *metricService) DumpToFile(ctx context.Context, fielName string) error {
	metrics, err := s.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("cant get metrics from repo: %w", err)
	}

	var temp []domain.Metric
	for _, metric := range metrics {
		temp = append(temp, metric)
	}

	if err := s.fileStorage.Write(ctx, temp, fielName); err != nil {
		return fmt.Errorf("error on dump metrics to file: %w", err)
	}

	return nil
}
