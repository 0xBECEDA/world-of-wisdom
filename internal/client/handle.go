package client

import (
	"log"
	"net"
	"time"

	"github.com/vmihailenco/msgpack/v5"

	"world-of-wisdom/internal/message"
	"world-of-wisdom/internal/pow"
	"world-of-wisdom/internal/utils"
)

func (c *Client) handleConn(conn net.Conn) error {
	if err := c.requestChallenge(conn); err != nil {
		return err
	}
	resp, err := utils.ReadConn(conn, c.readDeadline)
	if err != nil {
		return err
	}

	return c.handleChallengeResp(resp, conn)
}

func (c *Client) requestChallenge(conn net.Conn) error {
	msg := message.NewMessage(message.ChallengeReq, "")

	if err := conn.SetWriteDeadline(time.Now().Add(c.readDeadline)); err != nil {
		log.Printf("failed set write deadline %v", err)
	}

	return utils.WriteConn(*msg, conn, c.writeDeadline)
}

func (c *Client) unmarshallQuote(respQuote []byte) (string, error) {
	quoteResponseMessage := message.Message{}
	err := msgpack.Unmarshal(respQuote, &quoteResponseMessage)
	if err != nil {
		return "", err
	}
	return quoteResponseMessage.Data, nil
}

func (c *Client) unmarshallChallenge(resp []byte, hash *pow.Hashcash) error {
	challengeResponseMessage := message.Message{}
	err := msgpack.Unmarshal(resp, &challengeResponseMessage)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal([]byte(challengeResponseMessage.Data), hash)
}

func (c *Client) prepareQuoteRequest(solvedHash *pow.Hashcash) *message.Message {
	solvedHashMarshalled, _ := msgpack.Marshal(solvedHash)
	return &message.Message{Type: message.QuoteReq, Data: string(solvedHashMarshalled)}
}

func (c *Client) solveChallenge(hash *pow.Hashcash) (*pow.Hashcash, error) {
	err := hash.Compute(maxIterations)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (c *Client) handleChallengeResp(resp []byte, conn net.Conn) error {
	quoteRequest, err := c.handleChallengeResponse(resp)
	if err != nil {
		return err
	}

	if err := utils.WriteConn(*quoteRequest, conn, c.writeDeadline); err != nil {
		return err
	}
	respQuote, err := utils.ReadConn(conn, c.readDeadline)
	if err != nil {
		return err
	}

	quote, err := c.unmarshallQuote(respQuote)
	if err != nil {
		return err
	}

	log.Printf("got quote: '%s'(c)", quote)
	return nil
}

func (c *Client) handleChallengeResponse(resp []byte) (*message.Message, error) {
	hash := &pow.Hashcash{}
	if err := c.unmarshallChallenge(resp, hash); err != nil {
		return nil, err
	}

	_, err := c.solveChallenge(hash)
	if err != nil {
		return nil, err
	}
	return c.prepareQuoteRequest(hash), nil
}
