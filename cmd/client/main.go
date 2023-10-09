package main

import (
	"log"
	"world-of-wisdom/pkg/client"
	"world-of-wisdom/pkg/config"
)

func main() {
	clientConfig := &client.Config{}
	if err := config.LoadConfig(clientConfig); err != nil {
		log.Fatalf("failed to start client: %s", err.Error())
	}

	tcpClient := client.NewClient(clientConfig)
	err := tcpClient.Start()
	if err != nil {
		log.Fatalf("failed to start client: %s", err.Error())
	}
}
