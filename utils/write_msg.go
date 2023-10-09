package utils

import (
	"net"
	"world-of-wisdom/internal/tcp_message"
)

func SendMessage(msg tcp_message.Message, clientConn net.Conn) error {
	marshaledMsg, err := msg.Marshal()
	if err != nil {
		return err
	}

	// TODO возможно добавить делиметер
	_, err = clientConn.Write(marshaledMsg)
	if err != nil {
		return err
	}
	return nil
}
