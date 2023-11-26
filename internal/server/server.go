package server

import (
	"net"
	"time"
	"world-of-wisdom/internal/pow"
	"world-of-wisdom/internal/quotes"
)

type Server struct {
	powService   pow.Repository
	quoteService quotes.QuoteRepository

	writeDeadline time.Duration
	readDeadline  time.Duration
}

func NewServer(
	hr pow.Repository,
	quoteService quotes.QuoteRepository,
	writeDeadline time.Duration,
	readDeadline time.Duration) *Server {
	return &Server{
		powService:    hr,
		quoteService:  quoteService,
		writeDeadline: writeDeadline,
		readDeadline:  readDeadline,
	}
}

func (s *Server) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		clientConn, err := l.Accept()
		if err != nil {
			continue
		}

		go s.handle(clientConn)
	}
}
