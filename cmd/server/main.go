package main

import (
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

	addr := ":" + cfg.Port
	log.Printf("starting tcp server on addr %v", addr)
	log.Fatal(tcpServer.Listen(addr))
}
