package main

import (
	"fmt"
	"log"
	"world-of-wisdom/internal/hashcash"
	"world-of-wisdom/internal/quotes"
	"world-of-wisdom/pkg/config"
	"world-of-wisdom/pkg/server"
	"world-of-wisdom/storage"
)

func main() {
	cfg := &server.Config{}
	if err := config.LoadConfig(cfg); err != nil {
		log.Fatalf("failed start server %s", err.Error())
	}

	hashRepository := hashcash.NewHashCashRepository(storage.NewStorage())
	hashService := hashcash.NewService(hashRepository)

	quoteService := quotes.NewService(quotes.NewRepository())

	tcpServer := server.NewServer(hashService, quoteService)

	log.Printf("starting tcp server on port %v", cfg.Port)
	log.Fatal(tcpServer.Listen(fmt.Sprintf("localhost:%v", cfg.Port)))
}
