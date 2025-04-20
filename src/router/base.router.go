package router

import (
	"errors"
	"log"

	"pixeltactics.com/websocket-gateway/src/messages"
)

const INCOMING_CHANNEL = "incoming:"
const OUTGOING_CHANNEL = "outgoing:"

var ErrBadGateway = errors.New("bad gateway")
var ErrInternalServer = errors.New("internal server error")

type ControlRouter interface {
	AddIncomingRouter(router IncomingRouter)
	RouteMessage(message *messages.Message, client messages.WebSocketClient)
}

type IncomingRouter interface {
	RouteMessage(message *messages.Message, client messages.WebSocketClient) error
	IsRouterOf(route string) bool
}

type ControlRouterImpl struct {
	IncomingRouters []IncomingRouter
}

func (control *ControlRouterImpl) AddIncomingRouter(router IncomingRouter) {
	control.IncomingRouters = append(control.IncomingRouters, router)
}

func (control *ControlRouterImpl) RouteMessage(message *messages.Message, client messages.WebSocketClient) {
	log.Println("GOT")
	for _, router := range control.IncomingRouters {
		if router.IsRouterOf(message.Route) {
			err := router.RouteMessage(message, client)
			if err != nil {
				log.Println("[ERROR] " + err.Error())
				client.Send(messages.Error(ErrInternalServer))
				return
			}
			return
		}
	}
	client.Send(messages.Error(ErrBadGateway))
}

func NewControlRouter() ControlRouter {
	return &ControlRouterImpl{
		IncomingRouters: make([]IncomingRouter, 0),
	}
}
