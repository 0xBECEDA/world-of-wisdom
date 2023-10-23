package client

import (
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"net"
	"time"
	"world-of-wisdom/internal/message"
	"world-of-wisdom/internal/pow"
	"world-of-wisdom/internal/utils"
)

const maxIterations = 10000000

type Client struct {
	Hostname      string
	Port          string
	Resource      string
	WriteDeadline time.Duration
	ReadDeadline  time.Duration
}

func NewClient(config *Config) *Client {
	return &Client{
		Hostname:      config.Hostname,
		Port:          config.Port,
		Resource:      config.Resource,
		WriteDeadline: config.WriteDeadline,
		ReadDeadline:  config.ReadDeadline,
	}
}

func (r *Client) Start() error {
	addr := fmt.Sprintf("%v:%v", r.Hostname, r.Port)
	log.Printf("starting client on %s", addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := r.handleConnection(conn); err != nil {
		return err
	}
	return nil
}

func (r *Client) handleConnection(conn net.Conn) error {
	if err := r.requestChallenge(conn); err != nil {
		return err
	}

	resp, err := utils.ReadConn(conn, r.ReadDeadline)
	if err != nil {
		return err
	}

	return r.handleChallengeResponse(resp, conn)
}

func (r *Client) handleChallengeResponse(resp []byte, conn net.Conn) error {
	quoteRequest, err := handleChallengeResponse(resp)
	if err != nil {
		return err
	}

	if err := utils.WriteConn(*quoteRequest, conn, r.WriteDeadline); err != nil {
		return err
	}
	respQuote, err := utils.ReadConn(conn, r.ReadDeadline)
	if err != nil {
		return err
	}

	quote, err := unmarshallQuote(respQuote)
	if err != nil {
		return err
	}

	log.Printf("received quote: '%s'(c)", quote)
	return nil
}

func (r *Client) requestChallenge(conn net.Conn) error {
	msg := message.NewMessage(message.ChallengeReq, "")

	if err := conn.SetWriteDeadline(time.Now().Add(r.ReadDeadline)); err != nil {
		log.Printf("failed set write deadline %v", err)
	}

	return utils.WriteConn(*msg, conn, r.WriteDeadline)
}

func unmarshallQuote(respQuote []byte) (string, error) {
	quoteResponseMessage := message.Message{}
	err := msgpack.Unmarshal(respQuote, &quoteResponseMessage)
	if err != nil {
		return "", err
	}
	return quoteResponseMessage.Data, nil
}

func handleChallengeResponse(resp []byte) (*message.Message, error) {
	hash := &pow.Hashcash{}
	if err := unmarshallHash(resp, hash); err != nil {
		return nil, err
	}

	_, err := solveHash(hash)
	if err != nil {
		return nil, err
	}
	return prepareQuoteRequest(hash), nil
}

func unmarshallHash(resp []byte, hash *pow.Hashcash) error {
	challengeResponseMessage := message.Message{}
	err := msgpack.Unmarshal(resp, &challengeResponseMessage)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal([]byte(challengeResponseMessage.Data), hash)
}

func solveHash(hash *pow.Hashcash) (*pow.Hashcash, error) {
	err := hash.ComputeHash(maxIterations)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func prepareQuoteRequest(solvedHash *pow.Hashcash) *message.Message {
	solvedHashMarshalled, _ := msgpack.Marshal(solvedHash)
	return &message.Message{Type: message.QuoteReq, Data: string(solvedHashMarshalled)}
}
