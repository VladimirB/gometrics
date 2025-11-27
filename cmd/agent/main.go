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

	logger.Log.Info("Running Agent", zap.String("Start time", time.Now().Local().String()))

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
			sendFunc := func() error {
				return agent.SendMetrics(cfg.Address, metrics)
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
