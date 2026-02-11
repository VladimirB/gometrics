package main

import (
	"context"
	"log"
	"time"

	"github.com/VladimirB/gometrics/internal/config"
	model "github.com/VladimirB/gometrics/internal/model"
	"github.com/VladimirB/gometrics/internal/module/agent/adapter/http_api"
	"github.com/VladimirB/gometrics/internal/module/agent/service"
	"github.com/VladimirB/gometrics/internal/shared/logger"
	"go.uber.org/zap"
)

const logLevel = "info"

func main() {
	if err := logger.Initialize(logLevel); err != nil {
		log.Fatal(err)
	}
	defer logger.Log.Sync()

	logger.Log.Info("Running Agent", zap.String("Start time", time.Now().Local().String()))

	cfg := config.GetAgentConfig()

	metricApi := http_api.NewMetricApi(cfg.Address)
	agentService := service.NewAgentService(metricApi, cfg.ReportInterval)

	pollTicker := time.NewTicker(cfg.PollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer reportTicker.Stop()

	metrics := make(map[string]model.Metrics)
	for {
		select {
		case <-pollTicker.C:
			agentService.ReadMetrics(metrics)
		case <-reportTicker.C:
			sendFunc := func() error {
				return agentService.SendMetrics(context.Background(), metrics)
			}

			doWithRetry(sendFunc, 5, 3*time.Second)
		}
	}
}

func doWithRetry(fn func() error, maxRetries int, delay time.Duration) {
	for attempts := 0; attempts < maxRetries; attempts++ {
		if err := fn(); err == nil {
			break
		}

		if attempts < maxRetries-1 {
			time.Sleep(delay)
		}
	}
}
