package main

import (
	"flag"

	"github.com/VladimirB/gometrics/internal/config"
)

func parseFlags(config *config.ServerConfig) {
	flag.Var(&config.Address, "a", "string value in host:port format")
	flag.Parse()
}
