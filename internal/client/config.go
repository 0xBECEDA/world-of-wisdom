package client

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	defaultReadDeadline  = 10
	defaultWriteDeadline = 10
)

var (
	ErrEmptyPort = errors.New("empty client port")
	ErrEmptyHost = errors.New("empty client host")
)

type Config struct {
	Hostname      string
	Port          string
	Resource      string
	WriteDeadline time.Duration
	ReadDeadline  time.Duration
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

	writeDeadline := os.Getenv("WRITE_DEADLINE")
	if writeDeadline == "" {
		c.WriteDeadline = defaultWriteDeadline * time.Second
	} else {
		dur, err := strconv.Atoi(writeDeadline)
		if err != nil {
			return err
		}
		c.WriteDeadline = time.Duration(dur) * time.Second
	}

	readDeadline := os.Getenv("READ_DEADLINE")
	if readDeadline == "" {
		c.ReadDeadline = defaultReadDeadline * time.Second
	} else {
		dur, err := strconv.Atoi(readDeadline)
		if err != nil {
			return err
		}
		c.ReadDeadline = time.Duration(dur) * time.Second
	}
	return nil
}
