package server

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"world-of-wisdom/internal/hashcash"
	"world-of-wisdom/internal/quotes"
	"world-of-wisdom/internal/tcp_message"
	"world-of-wisdom/pkg/utils"
)

type Server struct {
	powService   hashcash.Service
	quoteService quotes.QuoteService
}

func NewServer(hashcashService hashcash.Service, quoteService quotes.QuoteService) *Server {
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
		req, err := s.readFromConn(clientConn)
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
			err = utils.WriteConn(*response, clientConn)
			if err != nil {
				log.Printf("error sending tcp message: %s", err.Error())
			}
		}
	}
}

func (s *Server) readFromConn(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 1024)
	var data []byte

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		data = append(data, buffer[:n]...)
		if n < len(buffer) {
			break
		}
	}
	return data, nil
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

	hash := hashcash.NewHashcash(5, req.Data)
	log.Printf("adding hash %++v", hash)

	s.powService.AddHashIndicator(binary.BigEndian.Uint64(hash.Rand))
	marshaledStamp, err := json.Marshal(hash)
	if err != nil {
		return nil, ErrFailedToMarshal
	}

	return tcp_message.NewMessage(tcp_message.ChallengeResp, string(marshaledStamp)), nil
}

func (s *Server) handleQuoteRequest(parsedRequest tcp_message.Message) (*tcp_message.Message, error) {
	var stamp hashcash.Hashcash
	err := json.Unmarshal([]byte(parsedRequest.Data), &stamp)
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

	responseMessage := tcp_message.NewMessage(tcp_message.QuoteResp, s.quoteService.GetQuote().QuoteText)
	err = s.powService.DeleteIndicator(randNum)
	if err != nil {
		return nil, ErrFailedToRemoveIndicator
	}

	return responseMessage, nil
}
