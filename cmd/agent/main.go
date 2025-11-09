package main

import (
	"fmt"
	"log"
	"time"

	"github.com/VladimirB/gometrics/internal/agent"
	"github.com/VladimirB/gometrics/internal/config"
	"github.com/VladimirB/gometrics/internal/metricapi"
	model "github.com/VladimirB/gometrics/internal/model"
)

func main() {
	metrics := make(map[string]model.Metrics)
	agent := agent.NewAgent(metricapi.NewClient())

	config := config.NewAgentConfig()
	parseFlags(config)
	fmt.Println("run agent with config:", config)

	pollTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			agent.ReadMetrics(metrics)
		case <-reportTicker.C:
			for _, metric := range metrics {
				err := agent.Send(config.Address.String(), metric)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
