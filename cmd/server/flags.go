package main

import (
	"flag"

	"github.com/VladimirB/gometrics/internal/config"
)

var serverConfig config.ServerConfig = config.ServerConfig{
	Address: config.NetAddress{
		Host: "localhost",
		Port: 8080,
	},
}

func parseFlags() {
	flag.Var(&serverConfig.Address, "a", "string value in host:port format")
	flag.Parse()
}