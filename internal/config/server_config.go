package config

const(
	defaultServerHost = "localhost"
	defaultServerPort = 8080
)

type ServerConfig struct {
	Address NetAddress
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Address: NetAddress{
			Host: defaultServerHost,
			Port: defaultServerPort,
		},
	}
}