package server

import (
	"encoding/binary"
	"encoding/json"
	hashcash2 "go/internal/hashcash"
	"go/internal/quotes"
	"go/internal/tcp_message"
	"go/utils"
	"log"
	"net"
)

type Server struct {
	powService   hashcash2.Service
	quoteService quotes.QuoteService
}

func NewServer(hashcashService hashcash2.Service, quoteService quotes.QuoteService) *Server {
	return &Server{
		powService:   hashcashService,
		quoteService: quoteService,
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
		req := make([]byte, 0)
		_, err := clientConn.Read(req)
		if err != nil {
			return
		}

		response, err := s.processRequest(req)
		if err != nil {
			return
		}

		if response != nil {
			err = utils.SendMessage(*response, clientConn)
			if err != nil {
				log.Printf("error sending tcp message:")
			}
		}
	}
}

func (s *Server) processRequest(clientRequest []byte) (*tcp_message.Message, error) {
	parsedRequest, err := tcp_message.Parse(clientRequest)
	if err != nil {
		return nil, err
	}

	switch parsedRequest.Type {
	case tcp_message.ChallengeReq:
		return s.challengeRequestHandler(parsedRequest)
	case tcp_message.QuoteReq:
		return s.handleQuoteRequest(*parsedRequest)
	default:
		return nil, ErrUnknownRequest
	}
}

func (s *Server) challengeRequestHandler(req *tcp_message.Message) (*tcp_message.Message, error) {
	if req == nil {
		return nil, ErrEmptyMessage
	}

	hash := hashcash2.NewHashcash(5, req.Data)
	log.Printf("adding hash %++v", hash)

	if err := s.powService.AddHashIndicator(binary.BigEndian.Uint64(hash.Rand)); err != nil {
		return nil, ErrFailedToAddIndicator
	}

	marshaledStamp, err := json.Marshal(hash)
	if err != nil {
		return nil, ErrFailedToMarshal
	}

	return tcp_message.NewMessage(tcp_message.ChallengeResp, string(marshaledStamp)), nil
}

func (s *Server) handleQuoteRequest(parsedRequest tcp_message.Message) (*tcp_message.Message, error) {
	var stamp hashcash2.Hashcash
	err := json.Unmarshal([]byte(parsedRequest.Data), &stamp)
	if err != nil {
		return nil, ErrFailedToUnmarshal
	}

	log.Printf("received hash %++v", stamp)

	randNum := binary.BigEndian.Uint64(stamp.Rand)
	_, err = s.powService.GetHashIndicator(randNum)
	if err != nil {
		return nil, ErrFailedToGetRand
	}

	if !stamp.Verify() {
		return nil, ErrChallengeUnsolved
	}

	responseMessage := tcp_message.NewMessage(tcp_message.QuoteResp, s.quoteService.GetQuote().QuoteText)
	err = s.powService.RemoveHashIndicator(randNum)
	if err != nil {
		return nil, ErrFailedToRemoveIndicator
	}

	log.Printf("response tcp_message: %++v", responseMessage)
	return responseMessage, nil
}
