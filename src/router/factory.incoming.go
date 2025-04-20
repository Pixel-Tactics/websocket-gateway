package router

import (
	"errors"

	"pixeltactics.com/websocket-gateway/src/config"
	"pixeltactics.com/websocket-gateway/src/integrations/communication"
)

var ErrInvalidRouterType = errors.New("invalid router type")

type IncomingRouterFactory interface {
	Generate(config *config.Route) (IncomingRouter, error)
}

type IncomingRouterFactoryImpl struct {
	RMQManager *communication.RMQManager
}

func (factory *IncomingRouterFactoryImpl) Generate(config *config.Route) (IncomingRouter, error) {
	if config.Type == "queue" {
		return NewIncomingQueue(config, factory.RMQManager), nil
	} else if config.Type == "stream" {
		return NewIncomingStream(config, factory.RMQManager), nil
	}
	return nil, ErrInvalidRouterType
}

func NewIncomingRouterFactory(
	rmqManager *communication.RMQManager,
) IncomingRouterFactory {
	return &IncomingRouterFactoryImpl{
		RMQManager: rmqManager,
	}
}
