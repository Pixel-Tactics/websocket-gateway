package router

import (
	"encoding/json"
	"strings"

	"github.com/rabbitmq/amqp091-go"
	"pixeltactics.com/websocket-gateway/src/config"
	"pixeltactics.com/websocket-gateway/src/integrations/communication"
	"pixeltactics.com/websocket-gateway/src/messages"
)

type IncomingQueue struct {
	Config     *config.Route
	RMQManager *communication.RMQManager
}

func NewIncomingQueue(
	config *config.Route,
	rmqManager *communication.RMQManager,
) IncomingRouter {
	return &IncomingQueue{
		Config:     config,
		RMQManager: rmqManager,
	}
}

func (queue *IncomingQueue) RouteMessage(message *messages.Message, client messages.WebSocketClient) error {
	username, err := client.GetUserId()
	if err != nil {
		return err
	}

	channel, err := queue.RMQManager.GetChannel(INCOMING_CHANNEL)
	if err != nil {
		return err
	}

	message.Route = message.Route[len(queue.Config.UserPrefix):]
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	brokerPath := strings.ReplaceAll(queue.Config.BrokerPath, "{{player}}", username)
	return channel.Publish("", brokerPath, true, false, amqp091.Publishing{
		Body: data,
	})
}

func (queue *IncomingQueue) IsRouterOf(route string) bool {
	return strings.HasPrefix(route, queue.Config.UserPrefix)
}
