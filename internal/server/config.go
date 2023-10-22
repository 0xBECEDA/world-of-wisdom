package server

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	defaultReadDeadline  = 10
	defaultWriteDeadline = 10
	defaultKeyTTL        = 20
)

var (
	ErrEmptyPort = errors.New("empty server port")
)

type Config struct {
	Port          string
	WriteDeadline time.Duration
	ReadDeadline  time.Duration
	KeyTTL        time.Duration
}

func (c *Config) Load() error {
	port := os.Getenv("PORT")
	if port == "" {
		return ErrEmptyPort
	}
	c.Port = port

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

	keyTTL := os.Getenv("KEY_TTL")
	if keyTTL == "" {
		c.KeyTTL = defaultKeyTTL * time.Second
	} else {
		dur, err := strconv.Atoi(keyTTL)
		if err != nil {
			return err
		}
		c.KeyTTL = time.Duration(dur) * time.Second
	}
	return nil
}
