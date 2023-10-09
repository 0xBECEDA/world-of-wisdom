package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"world-of-wisdom/internal/hashcash"
	"world-of-wisdom/internal/tcp_message"
	"world-of-wisdom/utils"
)

const maxIterations = 10000000

type Client struct {
	Hostname string
	Port     int64
	Resource string
}

func NewClient(config *Config) *Client {
	return &Client{
		Hostname: config.Hostname,
		Port:     config.Port,
		Resource: config.Resource,
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

	resp, err := receiveResponse(conn)
	if err != nil {
		return err
	}

	return r.handleChallengeResponse(resp, conn)
}

func (r *Client) handleChallengeResponse(resp string, conn net.Conn) error {
	quoteRequest, err := handleChallengeResponse(resp)
	if err != nil {
		return err
	}

	if err := utils.SendMessage(*quoteRequest, conn); err != nil {
		return err
	}

	respQuote, err := receiveResponse(conn)
	if err != nil {
		return err
	}

	quote, err := unmarshallQuote(respQuote)
	if err != nil {
		return err
	}

	log.Printf("received quote: %s", quote)
	return nil
}

func (r *Client) requestChallenge(conn net.Conn) error {
	msg := tcp_message.NewMessage(tcp_message.ChallengeReq, "")
	return utils.SendMessage(*msg, conn)
}

func receiveResponse(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	return reader.ReadString('\n')
}

func unmarshallQuote(respQuote string) (string, error) {
	quoteResponseMessage := tcp_message.Message{}
	err := json.Unmarshal([]byte(respQuote), &quoteResponseMessage)
	if err != nil {
		return "", err
	}
	return quoteResponseMessage.Data, nil
}

func handleChallengeResponse(resp string) (*tcp_message.Message, error) {
	hash := &hashcash.Hashcash{}
	if err := unmarshallHash(resp, hash); err != nil {
		return nil, err
	}

	_, err := solveHash(hash)
	if err != nil {
		return nil, err
	}
	return prepareQuoteRequest(hash), nil
}

func unmarshallHash(resp string, hash *hashcash.Hashcash) error {
	challengeResponseMessage := tcp_message.Message{}
	err := json.Unmarshal([]byte(resp), &challengeResponseMessage)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(challengeResponseMessage.Data), hash)
}

func solveHash(hash *hashcash.Hashcash) (*hashcash.Hashcash, error) {
	err := hash.ComputeHash(maxIterations)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func prepareQuoteRequest(solvedHash *hashcash.Hashcash) *tcp_message.Message {
	solvedHashMarshalled, _ := json.Marshal(solvedHash)
	return &tcp_message.Message{Type: tcp_message.QuoteReq, Data: string(solvedHashMarshalled)}
}
