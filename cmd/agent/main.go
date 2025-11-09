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

	timer := 0
	for {
		timer++

		if timer % config.PollInterval == 0 {
			agent.ReadMetrics(metrics)
		}

		if timer % config.ReportInterval == 0 {
			for _, metric := range(metrics) {
				err := agent.Send(config.Address.String(), metric)
				if err != nil {
					log.Println(err)
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}
