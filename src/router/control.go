package router

import (
	"log"

	"pixeltactics.com/websocket-gateway/src/messages"
)

type ControlRouter interface {
	AddIncomingRouter(router IncomingRouter)
	AddOutgoingRouter(router OutgoingRouter)

	// Runs outgoing routers in separate goroutine.
	Run()

	RouteMessage(message *messages.Message, client messages.WebSocketClient)
}

type ControlRouterImpl struct {
	IncomingRouters []IncomingRouter
	OutgoingRouters []OutgoingRouter
}

func NewControlRouter() ControlRouter {
	return &ControlRouterImpl{
		IncomingRouters: make([]IncomingRouter, 0),
		OutgoingRouters: make([]OutgoingRouter, 0),
	}
}

func (control *ControlRouterImpl) AddIncomingRouter(router IncomingRouter) {
	control.IncomingRouters = append(control.IncomingRouters, router)
}

func (control *ControlRouterImpl) AddOutgoingRouter(router OutgoingRouter) {
	control.OutgoingRouters = append(control.OutgoingRouters, router)
}

// Runs outgoing routers in separate goroutine.
func (control *ControlRouterImpl) Run() {
	for _, router := range control.OutgoingRouters {
		go router.Run()
	}
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
