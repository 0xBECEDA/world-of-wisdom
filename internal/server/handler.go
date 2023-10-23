package server

import (
	"encoding/binary"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"net"
	"world-of-wisdom/internal/message"
	"world-of-wisdom/internal/pow"
	"world-of-wisdom/internal/utils"
)

func (s *Server) handle(clientConn net.Conn) {
	defer clientConn.Close()

	for {
		req, err := utils.ReadConn(clientConn, s.ReadDeadline)
		if err != nil {
			log.Printf("error reading request: %s", err.Error())
			return
		}

		if len(req) == 0 {
			continue
		}

		response, err := s.processClientRequest(req)
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

func (s *Server) processClientRequest(clientRequest []byte) (*message.Message, error) {
	parsedRequest, err := message.Unmarshal(clientRequest)
	if err != nil {
		return nil, err
	}

	switch parsedRequest.Type {
	case message.ChallengeReq:
		return s.challengeHandler(parsedRequest)
	case message.QuoteReq:
		return s.quoteHandler(*parsedRequest)
	default:
		return nil, ErrUnknownRequest
	}
}

func (s *Server) challengeHandler(req *message.Message) (*message.Message, error) {
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

func (s *Server) quoteHandler(parsedRequest message.Message) (*message.Message, error) {
	var hash pow.Hashcash
	err := msgpack.Unmarshal([]byte(parsedRequest.Data), &hash)
	if err != nil {
		return nil, ErrFailedToUnmarshal
	}

	randNum := binary.BigEndian.Uint64(hash.Rand)
	ok := s.powService.IndicatorExists(randNum)
	if !ok {
		return nil, ErrFailedToGetRand
	}

	if !hash.Verify() {
		return nil, ErrChallengeUnsolved
	}

	responseMessage := message.NewMessage(message.QuoteResp, s.quoteService.GetAnyQuote().QuoteText)
	s.powService.DeleteIndicator(randNum)

	return responseMessage, nil
}
