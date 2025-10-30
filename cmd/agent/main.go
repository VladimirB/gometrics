package main

import (
	"time"

	"github.com/VladimirB/gometrics/internal/agent"
	model "github.com/VladimirB/gometrics/internal/model"
)

const (
	pollInterval = 2
)

func main() {
	metrics := make(map[string]model.Metrics)

	pollCounter := 0

	for {
		pollCounter++

		if pollCounter % pollInterval == 0 {
			agent.ReadMetrics(metrics)
		}

		time.Sleep(1 * time.Second)
	}
}
