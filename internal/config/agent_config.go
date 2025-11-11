package config

import "flag"

const (
	defaultAgentHost      = "localhost"
	defaultAgentPort      = 8080
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type AgentConfig struct {
	Address        NetAddress
	PollInterval   int
	ReportInterval int
}

func GetAgentConfig() AgentConfig {
	config := AgentConfig{
		Address: NetAddress{
			Host: defaultAgentHost,
			Port: defaultAgentPort,
		},
		PollInterval:   defaultPollInterval,
		ReportInterval: defaultReportInterval,
	}

	flag.Var(&config.Address, "a", "string value in host:port format")
	flag.IntVar(&config.PollInterval, "p", defaultPollInterval, "metric poll interval in seconds")
	flag.IntVar(&config.ReportInterval, "r", defaultReportInterval, "metric send interval in seconds")
	flag.Parse()

	return config
}
