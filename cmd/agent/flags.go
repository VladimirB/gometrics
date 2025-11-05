package main

import (
	"flag"

	"github.com/VladimirB/gometrics/internal/config"
)

var agentConfig config.AgentConfig = config.AgentConfig{
	Address: config.NetAddress{
		Host: "localhost",
		Port: 8080,
	},
	PollInterval: 2,
	ReportInterval: 10,
}

func parseFlags() {
	flag.Var(&agentConfig.Address, "a", "string value in host:port format")
	flag.IntVar(&agentConfig.PollInterval, "p", 2, "metric poll interval in seconds")
	flag.IntVar(&agentConfig.ReportInterval, "r", 10, "metric send interval in seconds")
	flag.Parse()
}