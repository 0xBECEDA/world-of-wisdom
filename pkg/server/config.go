package server

import (
	"errors"
	"os"
)

var (
	ErrEmptyPort = errors.New("empty server port")
)

type Config struct {
	Port string
}

func (c *Config) Load() error {
	port := os.Getenv("PORT")
	if port == "" {
		return ErrEmptyPort
	}
	c.Port = port
	return nil
}
