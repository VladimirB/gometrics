package config

import (
	"errors"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Host string
	Port int
}

func (c ServerConfig) String() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (c *ServerConfig) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}

	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}

	c.Host = hp[0]
	c.Port = port
	return nil
}