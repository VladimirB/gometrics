package config

const(
	defaultAgentHost = "localhost"
	defaultAgentPort = 8080
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
			Host: defaultAgentHost,
			Port: defaultAgentPort,
		},
		PollInterval: defaultPollInterval,
		ReportInterval: defaultReportInterval,
	}
}