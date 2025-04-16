package router

import (
	"errors"

	"pixeltactics.com/websocket-gateway/src/messages"
)

var ErrBadGateway = errors.New("bad gateway")

type ControlRouter interface {
	RouteMessage(message *messages.Message, client messages.WebSocketClient)
}

type Router interface {
	RouteMessage(message *messages.Message, client messages.WebSocketClient) error
	IsPrefixOf(route string) bool
	IsIncoming() bool
}

type ControlRouterImpl struct {
	Routers []Router
}

func (control *ControlRouterImpl) RouteMessage(message *messages.Message, client messages.WebSocketClient) {
	for _, router := range control.Routers {
		if router.IsPrefixOf(message.Route) && router.IsIncoming() {
			err := router.RouteMessage(message, client)
			if err != nil {
				client.Send(messages.Error(ErrBadGateway))
				return
			}
			return
		}
	}
}

func NewControlRouter(routers []Router) ControlRouter {
	return &ControlRouterImpl{
		Routers: routers,
	}
}
