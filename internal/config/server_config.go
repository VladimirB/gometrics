package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerAddress = "localhost:8080"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func GetServerConfig() ServerConfig {
	config := ServerConfig{}
	if err := env.Parse(&config); err != nil {
		log.Println(err)
	}

	var addressFlag string
	flag.StringVar(&addressFlag, "a", defaultServerAddress, "server address")
	flag.Parse()

	if config.Address == "" {
		config.Address = addressFlag
	}

	return config
}
