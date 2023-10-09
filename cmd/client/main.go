package main

import (
	"log"
	"world-of-wisdom/pkg/client"
)

func main() {
	cfg := &client.Config{}
	if err := cfg.Load(); err != nil {
		log.Fatalf("failed to start client: %s", err.Error())
	}

	tcpClient := client.NewClient(cfg)
	err := tcpClient.Start()
	if err != nil {
		log.Fatalf("failed to start client: %s", err.Error())
	}
}
