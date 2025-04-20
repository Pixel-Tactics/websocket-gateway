package router

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"pixeltactics.com/websocket-gateway/src/config"
	"pixeltactics.com/websocket-gateway/src/events"
	"pixeltactics.com/websocket-gateway/src/integrations/communication"
	"pixeltactics.com/websocket-gateway/src/messages"
)

type OutgoingQueue struct {
	Config     *config.Route
	RMQManager *communication.RMQManager
	Messager   messages.WebSocketHub

	Users map[string]*UserQueue

	AddUser    chan string
	DeleteUser chan string
	Messages   chan UsernameMessage
}

type UsernameMessage struct {
	Username string
	Message  *messages.Message
}

func NewOutgoingQueue(
	config *config.Route,
	rmqManager *communication.RMQManager,
	eventManager events.EventManager,
	messager messages.WebSocketHub,
) *OutgoingQueue {
	queue := &OutgoingQueue{
		Config:     config,
		RMQManager: rmqManager,
		Messager:   messager,
		Users:      make(map[string]*UserQueue),
		AddUser:    make(chan string, 256),
		DeleteUser: make(chan string, 256),
		Messages:   make(chan UsernameMessage, 256),
	}
	eventManager.On("user-connect", func(event interface{}) error {
		username, ok := event.(string)
		if !ok {
			log.Println("[ERROR] cannot convert event to string")
			return nil
		}
		queue.AddUser <- username
		return nil
	})
	eventManager.On("user-disconnect", func(event interface{}) error {
		username, ok := event.(string)
		if !ok {
			log.Println("[ERROR] cannot convert event to string")
			return nil
		}
		queue.DeleteUser <- username
		return nil
	})
	return queue
}

func (queue *OutgoingQueue) Run() {
	for {
		select {
		case username := <-queue.AddUser:
			log.Println("Got add user")
			_, exists := queue.Users[username]
			if exists {
				continue
			}
			userQueue := NewUserQueue(username, queue)
			queue.Users[username] = userQueue
			go userQueue.Run()
		case username := <-queue.DeleteUser:
			log.Println("Got delete user")
			userQueue, exists := queue.Users[username]
			if !exists {
				continue
			}
			userQueue.Close <- true
			delete(queue.Users, username)
		case msg := <-queue.Messages:
			log.Println("Got message")
			queue.Messager.SendToUserId(msg.Username, msg.Message)
			// default:
			// 	log.Println("Nothing")
		}
	}
}

type UserQueue struct {
	Username string
	Parent   *OutgoingQueue
	Close    chan bool
}

func NewUserQueue(username string, parent *OutgoingQueue) *UserQueue {
	return &UserQueue{
		Username: username,
		Parent:   parent,
		Close:    make(chan bool, 256),
	}
}

// Runs the user queue.
func (queue *UserQueue) Run() {
	retrying := false
	for {
		if retrying {
			log.Println("channel closed, retrying...")
			time.Sleep(3 * time.Second)
		}
		retrying = true

		channelId := OUTGOING_CHANNEL + "session" + ":" + queue.Username
		channel, err := queue.Parent.RMQManager.GetChannel(channelId)
		if err != nil {
			log.Println(err)
			continue
		}

		queueName := strings.ReplaceAll(queue.Parent.Config.BrokerPath, "{{player}}", queue.Username)
		_, err = channel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			amqp091.Table{
				"x-expires": int32(60000),
			},
		)
		if err != nil {
			log.Println(err)
			continue
		}

		stop := queue.consume(channelId, channel, queueName)
		if stop {
			break
		}
	}
}

// Listens for RabbitMQ queue and close signal. Returns false if lost connection, or true if intended to be closed.
func (queue *UserQueue) consume(channelId string, channel *amqp091.Channel, queueName string) bool {
	consumer, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		return false
	}
	for {
		select {
		case <-queue.Close:
			queue.Parent.RMQManager.CloseChannel(channelId)
			return true
		case delivery, ok := <-consumer:
			if !ok {
				return false
			}

			var message *messages.Message
			errBody := json.Unmarshal(delivery.Body, &message)
			err = delivery.Ack(false)
			if err != nil {
				return false
			}
			if errBody != nil {
				log.Println(errBody)
				continue
			}

			queue.Parent.Messages <- UsernameMessage{
				Username: queue.Username,
				Message:  message,
			}
		}
	}
}
