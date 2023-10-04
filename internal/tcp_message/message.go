package tcp_message

import (
	"encoding/json"
)

type Message struct {
	Type Type   `json:"type"`
	Data string `json:"data"`
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func Parse(req []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(req, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func NewMessage(messageType Type, data string) *Message {
	return &Message{
		Type: messageType,
		Data: data,
	}
}
