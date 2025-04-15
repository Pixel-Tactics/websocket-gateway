package gateway

import (
	"errors"

	"pixeltactics.com/websocket-gateway/src/messages"
)

type Router interface {
	RouteMessage(client *messages.WebSocketMessager)
}

type RouterImpl struct {
}

func (router *RouterImpl) RouteMessage(client *messages.WebSocketMessager) {
	switch client.Message.Type {
	default:
		client.SendBack(Error(errors.New("invalid type")))
	}
}

func NewRouter() Router {
	return &RouterImpl{}
}
