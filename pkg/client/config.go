package client

import "errors"

var ErrEmptyHostName = errors.New("empty host or port, check your envs")

type Config struct {
	Hostname string `json:"HOSTNAME"`
	Port     int64  `json:"PORT"`
	Resource string `json:"RESOURCE"`
}

func (c *Config) Validate() error {
	if c.Hostname == "" || c.Port == 0 {
		return ErrEmptyHostName
	}
	return nil
}
