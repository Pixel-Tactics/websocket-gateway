package messages

import (
	"encoding/json"
	"errors"
)

var ErrInvalidJson = errors.New("invalid json")

type Messager interface {
	Send(clientId string, message *Message)
}

type ClientMessager interface {
	SendBack(message *Message)
	Messager
}

type WebSocketMessager struct {
	Message     *Message
	ClientId    *string
	Send        func(clientId string, message *Message)
	SendBack    func(message *Message)
	SetClientId func(clientId string)
}

type Message struct {
	Type       string                 `json:"type"`
	Identifier string                 `json:"identifier"`
	Body       map[string]interface{} `json:"body"`
}

func JsonBytesToMessage(jsonBytes []byte) (*Message, error) {
	var raw map[string]json.RawMessage
	err := json.Unmarshal(jsonBytes, &raw)
	if err != nil {
		return nil, ErrInvalidJson
	}

	var messageType string
	err = json.Unmarshal(raw["type"], &messageType)
	if err != nil {
		return nil, ErrInvalidJson
	}

	var identifier string
	err = json.Unmarshal(raw["identifier"], &identifier)
	if err != nil {
		return nil, ErrInvalidJson
	}

	var body map[string]interface{}
	err = json.Unmarshal(raw["body"], &body)
	if err != nil {
		return nil, ErrInvalidJson
	}

	return &Message{Type: messageType, Identifier: identifier, Body: body}, nil
}

func MessageToJsonBytes(message *Message) ([]byte, error) {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return nil, errors.New("message is invalid")
	}

	return jsonBytes, nil
}
