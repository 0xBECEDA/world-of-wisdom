package main

import (
	"fmt"
	"log"
	"world-of-wisdom/internal/hashcash"
	"world-of-wisdom/internal/quotes"
	"world-of-wisdom/internal/storage"
	"world-of-wisdom/pkg/server"
)

func main() {
	cfg := &server.Config{}
	if err := cfg.Load(); err != nil {
		log.Fatalf("failed start server %s", err.Error())
	}

	hashRepository := hashcash.NewHashCashRepository(storage.NewStorage())
	hashService := hashcash.NewService(hashRepository)

	quoteService := quotes.NewService(quotes.NewRepository())

	tcpServer := server.NewServer(hashService, quoteService)

	log.Printf("starting tcp server on port %v", cfg.Port)
	log.Fatal(tcpServer.Listen(fmt.Sprintf("localhost:%v", cfg.Port)))
}
