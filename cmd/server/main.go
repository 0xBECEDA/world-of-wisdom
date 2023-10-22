package main

import (
	"log"
	"world-of-wisdom/internal/pow"
	"world-of-wisdom/internal/quotes"
	"world-of-wisdom/internal/server"
	"world-of-wisdom/internal/storage"
)

func main() {
	cfg := &server.Config{}
	if err := cfg.Load(); err != nil {
		log.Fatalf("failed start server %s", err.Error())
	}

	tcpServer := server.NewServer(pow.NewHashCashRepository(storage.NewStorage()), quotes.NewRepository(), cfg.WriteDeadline, cfg.ReadDeadline)
	addr := ":" + cfg.Port

	log.Printf("starting tcp server on addr %v", addr)
	log.Fatal(tcpServer.Listen(addr))
}
