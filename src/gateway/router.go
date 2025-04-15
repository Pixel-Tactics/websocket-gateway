package gateway

import (
	"errors"

	"pixeltactics.com/websocket-gateway/src/messages"
)

type Router interface {
	RouteMessage(message *messages.Message, client messages.WebSocketClient)
}

type RouterImpl struct {
}

func (router *RouterImpl) RouteMessage(message *messages.Message, client messages.WebSocketClient) {
	switch message.Route {
	default:
		client.Send(messages.Error(errors.New("invalid type")))
	}
}

func NewRouter() Router {
	return &RouterImpl{}
}
