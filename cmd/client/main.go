package main

import (
	"log"

	"world-of-wisdom/internal/client"
)

func main() {
	cfg := &client.Config{}
	if err := cfg.Load(); err != nil {
		log.Fatalf("failed to start client: %s", err.Error())
	}

	tcpClient := client.New(cfg)
	err := tcpClient.Run()
	if err != nil {
		log.Fatalf("failed to start client: %s", err.Error())
	}
}
