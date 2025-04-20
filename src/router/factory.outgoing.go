package router

import (
	"pixeltactics.com/websocket-gateway/src/config"
	"pixeltactics.com/websocket-gateway/src/events"
	"pixeltactics.com/websocket-gateway/src/integrations/communication"
	"pixeltactics.com/websocket-gateway/src/messages"
)

type OutgoingRouterFactory interface {
	Generate(config *config.Route) (OutgoingRouter, error)
}

type OutgoingRouterFactoryImpl struct {
	RMQManager   *communication.RMQManager
	EventManager events.EventManager
	Messager     messages.WebSocketHub
}

func (factory *OutgoingRouterFactoryImpl) Generate(config *config.Route) (OutgoingRouter, error) {
	if config.Type == "queue" {
		return NewOutgoingQueue(config, factory.RMQManager, factory.EventManager, factory.Messager), nil
	}
	return nil, ErrInvalidRouterType
}

func NewOutgoingRouterFactory(
	rmqManager *communication.RMQManager,
	eventManager events.EventManager,
	messager messages.WebSocketHub,
) OutgoingRouterFactory {
	return &OutgoingRouterFactoryImpl{
		RMQManager:   rmqManager,
		EventManager: eventManager,
		Messager:     messager,
	}
}
