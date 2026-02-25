package main

import (
	"context"
	"log"
	"time"

	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/domain"
	"github.com/VladimirB/gometrics/internal/module/agent/adapter/httpapi"
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

	metricAPI := httpapi.NewMetricAPI(cfg.Address)
	agentService := service.NewAgentService(metricAPI, cfg.ReportInterval)

	pollTicker := time.NewTicker(cfg.PollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer reportTicker.Stop()

	metrics := make(map[string]domain.Metric)
	for {
		select {
		case <-pollTicker.C:
			agentService.ReadMetrics(metrics)
		case <-reportTicker.C:
			agentService.SendMetrics(context.Background(), metrics)
		}
	}
}
