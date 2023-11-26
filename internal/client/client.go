package client

import (
	"fmt"
	"log"
	"net"
	"time"
)

const maxIterations = 10000000

type Client struct {
	hostname      string
	port          string
	resource      string
	writeDeadline time.Duration
	readDeadline  time.Duration
}

func New(config *Config) *Client {
	return &Client{
		hostname:      config.Hostname,
		port:          config.Port,
		resource:      config.Resource,
		writeDeadline: config.WriteDeadline,
		readDeadline:  config.ReadDeadline,
	}
}

func (c *Client) Run() error {
	addr := fmt.Sprintf("%v:%v", c.hostname, c.port)
	log.Printf("starting client on %s", addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := c.handleConn(conn); err != nil {
		return err
	}
	return nil
}
