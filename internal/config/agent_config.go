package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	defaultAgentAddress   = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
)

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

type agentFlags struct {
	address string
	poll    int
	report  int
}

func GetAgentConfig() AgentConfig {
	config := AgentConfig{}
	if err := env.Parse(&config); err != nil {
		log.Println(err)
	}

	flags := parseAgentFlags()

	if config.Address == "" {
		config.Address = flags.address
	}

	if config.PollInterval == 0 {
		config.PollInterval = time.Duration(flags.poll) * time.Second
	}

	if config.ReportInterval == 0 {
		config.ReportInterval = time.Duration(flags.report) * time.Second
	}

	return config
}

func parseAgentFlags() agentFlags {
	flags := agentFlags{}
	flag.StringVar(&flags.address, "a", defaultAgentAddress, "server address")
	flag.IntVar(&flags.poll, "p", defaultPollInterval, "metric poll interval in seconds")
	flag.IntVar(&flags.report, "r", defaultReportInterval, "metric send interval in seconds")
	flag.Parse()
	return flags
}
