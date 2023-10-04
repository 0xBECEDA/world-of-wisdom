package main

import (
	"go/internal/hashcash"
	"go/internal/quotes"
	"go/pkg/server"
	"go/storage"
	"log"
)

func main() {
	hashRepository := hashcash.NewHashCashRepository(storage.NewStorage())
	hashService := hashcash.NewService(hashRepository)

	quoteService := quotes.NewService(quotes.NewRepository())

	tcpServer := server.NewServer(hashService, quoteService)
	log.Fatal(tcpServer.Listen(":8080"))
}
