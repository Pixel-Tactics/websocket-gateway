package router

import "pixeltactics.com/websocket-gateway/src/messages"

type QueueRouter struct {
	// Queue
}

func NewQueueRouter() Router {
	return &QueueRouter{}
}

func (queue *QueueRouter) RouteMessage(message *messages.Message, client messages.WebSocketClient) error {
	return nil
}

func (queue *QueueRouter) IsPrefixOf(route string) bool {
	return false
}

func (queue *QueueRouter) IsIncoming() bool {
	return false
}
