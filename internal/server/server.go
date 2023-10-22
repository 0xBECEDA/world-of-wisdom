package server

import (
	"encoding/binary"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"net"
	"time"
	"world-of-wisdom/internal/message"
	"world-of-wisdom/internal/pow"
	"world-of-wisdom/internal/quotes"
	"world-of-wisdom/internal/utils"
)

type Server struct {
	powService   pow.Repository
	quoteService quotes.QuoteRepository

	WriteDeadline time.Duration
	ReadDeadline  time.Duration
}

func NewServer(
	hr pow.Repository,
	quoteService quotes.QuoteRepository,
	writeDeadline time.Duration,
	readDeadline time.Duration) *Server {
	return &Server{
		powService:    hr,
		quoteService:  quoteService,
		WriteDeadline: writeDeadline,
		ReadDeadline:  readDeadline,
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

		go s.handleConn(clientConn)
	}
}

func (s *Server) handleConn(clientConn net.Conn) {
	defer clientConn.Close()

	for {
		req, err := utils.ReadFromConn(clientConn, s.ReadDeadline)
		if err != nil {
			log.Printf("error reading request: %s", err.Error())
			return
		}

		if len(req) == 0 {
			continue
		}

		response, err := s.processRequest(req)
		if err != nil {
			log.Printf("error processing request: %s", err.Error())
			continue
		}

		if response != nil {
			err = utils.WriteConn(*response, clientConn, s.WriteDeadline)
			if err != nil {
				log.Printf("error sending tcp message: %s", err.Error())
			}
		}
	}
}

func (s *Server) processRequest(clientRequest []byte) (*message.Message, error) {
	parsedRequest, err := message.Parse(clientRequest)
	if err != nil {
		return nil, err
	}

	switch parsedRequest.Type {
	case message.ChallengeReq:
		return s.challengeRequestHandler(parsedRequest)
	case message.QuoteReq:
		return s.handleQuoteRequest(*parsedRequest)
	default:
		return nil, ErrUnknownRequest
	}
}

func (s *Server) challengeRequestHandler(req *message.Message) (*message.Message, error) {
	if req == nil {
		return nil, ErrEmptyMessage
	}

	hash := pow.NewHashcash(5, req.Data)
	log.Printf("adding hash %++v", hash)

	s.powService.AddIndicator(binary.BigEndian.Uint64(hash.Rand))
	marshaledStamp, err := msgpack.Marshal(hash)
	if err != nil {
		return nil, ErrFailedToMarshal
	}

	return message.NewMessage(message.ChallengeResp, string(marshaledStamp)), nil
}

func (s *Server) handleQuoteRequest(parsedRequest message.Message) (*message.Message, error) {
	var stamp pow.Hashcash
	err := msgpack.Unmarshal([]byte(parsedRequest.Data), &stamp)
	if err != nil {
		return nil, ErrFailedToUnmarshal
	}

	log.Printf("received hash %++v", stamp)

	randNum := binary.BigEndian.Uint64(stamp.Rand)
	ok := s.powService.IndicatorExists(randNum)
	if !ok {
		return nil, ErrFailedToGetRand
	}

	if !stamp.Verify() {
		return nil, ErrChallengeUnsolved
	}

	responseMessage := message.NewMessage(message.QuoteResp, s.quoteService.GetQuote().QuoteText)
	s.powService.DeleteIndicator(randNum)

	return responseMessage, nil
}
