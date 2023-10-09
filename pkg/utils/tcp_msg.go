package utils

import (
	"io"
	"net"
	"world-of-wisdom/internal/tcp_message"
)

func WriteConn(msg tcp_message.Message, clientConn net.Conn) error {
	marshaledMsg, err := msg.Marshal()
	if err != nil {
		return err
	}

	_, err = clientConn.Write(marshaledMsg)
	if err != nil {
		return err
	}
	return nil
}

func ReadFromConn(conn net.Conn) ([]byte, error) {
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
