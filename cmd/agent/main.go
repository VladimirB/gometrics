package main

import (
	"log"
	"time"

	"github.com/VladimirB/gometrics/internal/agent"
	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/logger"
	"github.com/VladimirB/gometrics/internal/metricapi"
	model "github.com/VladimirB/gometrics/internal/model"
	"go.uber.org/zap"
)

const logLevel = "info"

func main() {
	if err := logger.Initialize(logLevel); err != nil {
		log.Fatal(err)
	}
	defer logger.Log.Sync()

	logger.Log.Info("Running Agent", zap.Time("Start time", time.Now()))

	metrics := make(map[string]model.Metrics)
	agent := agent.NewAgent(metricapi.NewClient())

	cfg := config.GetAgentConfig()

	pollTicker := time.NewTicker(cfg.PollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			agent.ReadMetrics(metrics)
		case <-reportTicker.C:
			agent.SendMetrics(cfg.Address, metrics)
		}
	}
}
