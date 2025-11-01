package main

import (
	"log"
	"time"

	"github.com/VladimirB/gometrics/internal/agent"
	"github.com/VladimirB/gometrics/internal/metricapi"
	model "github.com/VladimirB/gometrics/internal/model"
)

const (
	pollInterval = 2
	reportInterval = 10
)

func main() {
	metrics := make(map[string]model.Metrics)
	agent := agent.NewAgent(metricapi.NewClient())

	timer := 0
	for {
		timer++

		if timer % pollInterval == 0 {
			agent.ReadMetrics(metrics)
		}

		if timer % reportInterval == 0 {
			for _, metric := range(metrics) {
				err := agent.Send(metric)
				if err != nil {
					log.Println(err)
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}
