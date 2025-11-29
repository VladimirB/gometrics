package agent

import (
	"errors"
	"math/rand"
	"runtime"

	"github.com/VladimirB/gometrics/internal/logger"
	model "github.com/VladimirB/gometrics/internal/model"
	"go.uber.org/zap"
)

type Agent struct {
	metricSender MetricSender
}

type MetricSender interface {
	Send(destination string, metric model.Metrics) error
	SendAsJSON(destination string, metric model.Metrics) error
}

func NewAgent(sender MetricSender) *Agent {
	return &Agent{
		metricSender: sender,
	}
}

func (Agent) ReadMetrics(metrics map[string]model.Metrics) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	fill(metrics, model.Gauge, model.Alloc, float64(stats.Alloc))
	fill(metrics, model.Gauge, model.BuckHashSys, float64(stats.BuckHashSys))
	fill(metrics, model.Gauge, model.Frees, float64(stats.Frees))
	fill(metrics, model.Gauge, model.GCCPUFraction, stats.GCCPUFraction)
	fill(metrics, model.Gauge, model.GCSys, float64(stats.GCSys))
	fill(metrics, model.Gauge, model.HeapAlloc, float64(stats.HeapAlloc))
	fill(metrics, model.Gauge, model.HeapIdle, float64(stats.HeapIdle))
	fill(metrics, model.Gauge, model.HeapInuse, float64(stats.HeapInuse))
	fill(metrics, model.Gauge, model.HeapObjects, float64(stats.HeapObjects))
	fill(metrics, model.Gauge, model.HeapReleased, float64(stats.HeapReleased))
	fill(metrics, model.Gauge, model.HeapSys, float64(stats.HeapSys))
	fill(metrics, model.Gauge, model.LastGC, float64(stats.LastGC))
	fill(metrics, model.Gauge, model.Lookups, float64(stats.Lookups))
	fill(metrics, model.Gauge, model.MCacheInuse, float64(stats.MCacheInuse))
	fill(metrics, model.Gauge, model.MCacheSys, float64(stats.MCacheSys))
	fill(metrics, model.Gauge, model.MSpanInuse, float64(stats.MSpanInuse))
	fill(metrics, model.Gauge, model.MSpanSys, float64(stats.MSpanSys))
	fill(metrics, model.Gauge, model.Mallocs, float64(stats.Mallocs))
	fill(metrics, model.Gauge, model.NextGC, float64(stats.NextGC))
	fill(metrics, model.Gauge, model.NumForcedGC, float64(stats.NumForcedGC))
	fill(metrics, model.Gauge, model.NumGC, float64(stats.NumGC))
	fill(metrics, model.Gauge, model.OtherSys, float64(stats.OtherSys))
	fill(metrics, model.Gauge, model.PauseTotalNs, float64(stats.PauseTotalNs))
	fill(metrics, model.Gauge, model.StackInuse, float64(stats.StackInuse))
	fill(metrics, model.Gauge, model.StackSys, float64(stats.StackSys))
	fill(metrics, model.Gauge, model.Sys, float64(stats.Sys))
	fill(metrics, model.Gauge, model.TotalAlloc, float64(stats.TotalAlloc))
	fill(metrics, model.Gauge, model.RandomValue, rand.Float64())

	var counter = 0
	pollCount, ok := metrics[model.PollCount]
	if ok {
		counter = int(*pollCount.Delta)
	}
	counter++
	fill(metrics, model.Counter, model.PollCount, float64(counter))
}

func fill(metrics map[string]model.Metrics, metricType string, metricName string, value float64) {
	metric, ok := metrics[metricName]
	if !ok {
		metric = model.Metrics{
			ID:    metricName,
			MType: metricType,
		}

		if metricType == model.Counter {
			metric.Delta = new(int64)
		} else {
			metric.Value = new(float64)
		}

		metrics[metricName] = metric
	}

	if metricType == model.Counter {
		*metric.Delta = int64(value)
	} else {
		*metric.Value = value
	}
}

func (a Agent) SendMetrics(server string, metrics map[string]model.Metrics) error {
	for _, metric := range metrics {
		if err := a.metricSender.SendAsJSON(server, metric); err != nil {
			logger.Log.Error("error on metric send", zap.Error(err), zap.String("Metric", metric.String()))
			fill(metrics, model.Counter, model.PollCount, 0)
			return errors.New("error on metric send")
		}
	}

	return nil
}
