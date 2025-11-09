package main

import (
	"flag"

	"github.com/VladimirB/gometrics/internal/config"
)

func parseFlags(config *config.AgentConfig) {
	flag.Var(&config.Address, "a", "string value in host:port format")
	flag.IntVar(&config.PollInterval, "p", 2, "metric poll interval in seconds")
	flag.IntVar(&config.ReportInterval, "r", 10, "metric send interval in seconds")
	flag.Parse()
}