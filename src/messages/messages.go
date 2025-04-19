package messages

import (
	"encoding/json"
	"errors"
)

var ErrInvalidJson = errors.New("invalid json")
var ErrInvalidMessage = errors.New("invalid message")

type WebSocketClient interface {
	GetUserId() (string, error)
	SendToUserId(userId string, message *Message)
	Send(message *Message)
}

type WebSocketHub interface {
	SendToUserId(userId string, message *Message)
}

type Message struct {
	Route string                 `json:"route"`
	Data  map[string]interface{} `json:"data"`
}

func JsonBytesToMessage(jsonBytes []byte) (*Message, error) {
	var message Message
	err := json.Unmarshal(jsonBytes, &message)
	if err != nil {
		return nil, ErrInvalidJson
	}
	return &message, nil
}

func MessageToJsonBytes(message *Message) ([]byte, error) {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return nil, ErrInvalidMessage
	}
	return jsonBytes, nil
}

func CreateMessage(route string, data map[string]interface{}, message string) *Message {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["message"] = message

	return &Message{
		Route: route,
		Data:  data,
	}
}

func Error(err error) *Message {
	return &Message{
		Route: "ERROR",
		Data: map[string]interface{}{
			"message": err.Error(),
		},
	}
}
