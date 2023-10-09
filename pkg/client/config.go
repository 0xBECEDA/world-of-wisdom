package client

import (
	"errors"
	"os"
)

var (
	ErrEmptyPort = errors.New("empty client port")
	ErrEmptyHost = errors.New("empty client host")
)

type Config struct {
	Hostname string
	Port     string
	Resource string
}

func (c *Config) Load() error {
	host := os.Getenv("HOSTNAME")
	if host == "" {
		return ErrEmptyHost
	}
	c.Hostname = host

	port := os.Getenv("PORT")
	if port == "" {
		return ErrEmptyPort
	}
	c.Port = port
	c.Resource = os.Getenv("RESOURCE")
	return nil
}
