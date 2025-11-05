package config

import (
	"errors"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	args := strings.Split(s, ":")
	if len(args) != 2 {
		return errors.New("need address in a form host:port")
	}

	port, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	a.Host = args[0]
	a.Port = port

	return nil
}