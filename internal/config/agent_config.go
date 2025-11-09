package config

const(
	defaultHost = "localhost"
	defaultPort = 8080
	defaultPollInterval = 2
	defaultReportInterval = 10
)

type AgentConfig struct {
	Address NetAddress
	PollInterval int
	ReportInterval int
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		Address: NetAddress{
			Host: defaultHost,
			Port: defaultPort,
		},
		PollInterval: defaultPollInterval,
		ReportInterval: defaultReportInterval,
	}
}