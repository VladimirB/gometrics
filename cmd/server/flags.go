package main

import (
	"flag"

	"github.com/VladimirB/gometrics/internal/config"
)

var serverConfig config.ServerConfig = config.ServerConfig{
	Host: "localhost",
	Port: 8080,
}

func parseFlags() {
	flag.Var(&serverConfig, "a", "string value in host:port format")
	flag.Parse()
}