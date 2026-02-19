package service

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/module/agent/port"
	"github.com/VladimirB/gometrics/internal/shared/logger"
	"go.uber.org/zap"
)

type AgentService struct {
	metricProvider      port.MetricProvider
	reportMetricTimeout time.Duration
}

func NewAgentService(metricProvider port.MetricProvider, reportMetricTimeout time.Duration) *AgentService {
	return &AgentService{
		metricProvider:      metricProvider,
		reportMetricTimeout: reportMetricTimeout,
	}
}

func (AgentService) ReadMetrics(metrics map[string]domain.Metric) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	fill(metrics, domain.Gauge, domain.Alloc, float64(stats.Alloc))
	fill(metrics, domain.Gauge, domain.BuckHashSys, float64(stats.BuckHashSys))
	fill(metrics, domain.Gauge, domain.Frees, float64(stats.Frees))
	fill(metrics, domain.Gauge, domain.GCCPUFraction, stats.GCCPUFraction)
	fill(metrics, domain.Gauge, domain.GCSys, float64(stats.GCSys))
	fill(metrics, domain.Gauge, domain.HeapAlloc, float64(stats.HeapAlloc))
	fill(metrics, domain.Gauge, domain.HeapIdle, float64(stats.HeapIdle))
	fill(metrics, domain.Gauge, domain.HeapInuse, float64(stats.HeapInuse))
	fill(metrics, domain.Gauge, domain.HeapObjects, float64(stats.HeapObjects))
	fill(metrics, domain.Gauge, domain.HeapReleased, float64(stats.HeapReleased))
	fill(metrics, domain.Gauge, domain.HeapSys, float64(stats.HeapSys))
	fill(metrics, domain.Gauge, domain.LastGC, float64(stats.LastGC))
	fill(metrics, domain.Gauge, domain.Lookups, float64(stats.Lookups))
	fill(metrics, domain.Gauge, domain.MCacheInuse, float64(stats.MCacheInuse))
	fill(metrics, domain.Gauge, domain.MCacheSys, float64(stats.MCacheSys))
	fill(metrics, domain.Gauge, domain.MSpanInuse, float64(stats.MSpanInuse))
	fill(metrics, domain.Gauge, domain.MSpanSys, float64(stats.MSpanSys))
	fill(metrics, domain.Gauge, domain.Mallocs, float64(stats.Mallocs))
	fill(metrics, domain.Gauge, domain.NextGC, float64(stats.NextGC))
	fill(metrics, domain.Gauge, domain.NumForcedGC, float64(stats.NumForcedGC))
	fill(metrics, domain.Gauge, domain.NumGC, float64(stats.NumGC))
	fill(metrics, domain.Gauge, domain.OtherSys, float64(stats.OtherSys))
	fill(metrics, domain.Gauge, domain.PauseTotalNs, float64(stats.PauseTotalNs))
	fill(metrics, domain.Gauge, domain.StackInuse, float64(stats.StackInuse))
	fill(metrics, domain.Gauge, domain.StackSys, float64(stats.StackSys))
	fill(metrics, domain.Gauge, domain.Sys, float64(stats.Sys))
	fill(metrics, domain.Gauge, domain.TotalAlloc, float64(stats.TotalAlloc))
	fill(metrics, domain.Gauge, domain.RandomValue, rand.Float64())

	var counter = 0
	pollCount, ok := metrics[domain.PollCount]
	if ok {
		counter = int(*pollCount.Delta)
	}
	counter++
	fill(metrics, domain.Counter, domain.PollCount, float64(counter))
}

func fill(metrics map[string]domain.Metric, metricType string, metricName string, value float64) {
	metric, ok := metrics[metricName]
	if !ok {
		metric = domain.Metric{
			ID:    metricName,
			MType: metricType,
		}

		if metricType == domain.Counter {
			metric.Delta = new(int64)
		} else {
			metric.Value = new(float64)
		}

		metrics[metricName] = metric
	}

	if metricType == domain.Counter {
		*metric.Delta = int64(value)
	} else {
		*metric.Value = value
	}
}

func (a AgentService) SendMetrics(ctx context.Context, metrics map[string]domain.Metric) error {
	ctx, cancel := context.WithTimeout(ctx, a.reportMetricTimeout)
	defer cancel()

	var bunch []domain.Metric
	for _, v := range metrics {
		bunch = append(bunch, v)
	}

	if err := a.metricProvider.SendAll(ctx, bunch); err != nil {
		fill(metrics, domain.Counter, domain.PollCount, 0)
		logger.Log.Error("send metrics failed", zap.Error(err))
		return fmt.Errorf("send metrics failed: %w", err)
	}

	return nil
}
