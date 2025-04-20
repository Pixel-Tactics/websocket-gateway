package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/rabbitmq/amqp091-go"
	"github.com/xeipuuv/gojsonschema"
	"pixeltactics.com/websocket-gateway/src/config"
	"pixeltactics.com/websocket-gateway/src/integrations/communication"
	"pixeltactics.com/websocket-gateway/src/messages"
)

type IncomingStream struct {
	Config     *config.Route
	RMQManager *communication.RMQManager
}

func NewIncomingStream(
	config *config.Route,
	rmqManager *communication.RMQManager,
) IncomingRouter {
	return &IncomingStream{
		Config:     config,
		RMQManager: rmqManager,
	}
}

func (queue *IncomingStream) RouteMessage(message *messages.Message, client messages.WebSocketClient) error {
	username, err := client.GetUserId()
	if err != nil {
		return err
	}

	channel, err := queue.RMQManager.GetChannel(INCOMING_CHANNEL)
	if err != nil {
		return err
	}

	data, err := json.Marshal(message.Data)
	if err != nil {
		return err
	}

	schema := strings.ReplaceAll(queue.Config.Schema, "{{player}}", username)
	log.Println(schema)
	log.Println(string(data))
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(data)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Println(err)
		return errors.New("invalid body")
	}
	if !result.Valid() {
		var errMsgs []string
		for _, schemaErr := range result.Errors() {
			errMsgs = append(errMsgs, schemaErr.String())
		}
		return fmt.Errorf("invalid body: %s", strings.Join(errMsgs, "; "))
	}

	return channel.Publish("", queue.Config.BrokerPath, true, false, amqp091.Publishing{
		Body: data,
	})
}

func (queue *IncomingStream) IsRouterOf(route string) bool {
	return route == queue.Config.UserPrefix
}
