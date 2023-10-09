package server

import "errors"

var ErrEmptyPort = errors.New("empty server port")

type Config struct {
	Port int64
}

func (c *Config) Validate() error {
	if c.Port == 0 {
		return ErrEmptyPort
	}
	return nil
}
