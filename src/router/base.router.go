package router

import (
	"errors"

	"pixeltactics.com/websocket-gateway/src/messages"
)

const INCOMING_CHANNEL = "incoming:"
const OUTGOING_CHANNEL = "outgoing:"

var ErrBadGateway = errors.New("bad gateway")
var ErrInternalServer = errors.New("internal server error")

type IncomingRouter interface {
	RouteMessage(message *messages.Message, client messages.WebSocketClient) error
	IsRouterOf(route string) bool
}

type OutgoingRouter interface {
	Run()
}
