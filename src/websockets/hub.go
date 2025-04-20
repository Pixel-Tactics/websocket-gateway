package websockets

import (
	"log"

	"pixeltactics.com/websocket-gateway/src/events"
	"pixeltactics.com/websocket-gateway/src/messages"
	"pixeltactics.com/websocket-gateway/src/router"
	"pixeltactics.com/websocket-gateway/src/utils/datastructures"
)

type MessageRequest struct {
	Message *messages.Message
	Client  *Client
}

type UserIdRequest struct {
	UserId string
	Client *Client
}

type ClientHub struct {
	EventManager  events.EventManager
	ControlRouter router.ControlRouter

	UserIdToClient *datastructures.SyncMap[string, *Client]
	ClientToUserId *datastructures.SyncMap[*Client, string]
	ClientList     *datastructures.SyncMap[*Client, bool]

	AddChannel     chan *Client
	CloseChannel   chan *Client
	UserIdChannel  chan *UserIdRequest
	MessageChannel chan *MessageRequest
}

func (hub *ClientHub) Run() {
	for {
		select {
		case client := <-hub.AddChannel:
			log.Println("[DEBUG] Storing client...")
			hub.ClientList.Store(client, true)
			log.Println("[DEBUG] Stored client.")
		case client := <-hub.CloseChannel:
			log.Println("[DEBUG] Closing client...")
			hub.closeClient(client)
			log.Println("[DEBUG] Closed client.")
		case request := <-hub.UserIdChannel:
			log.Println("[DEBUG] Requesting for username...")
			hub.setUserId(request.UserId, request.Client)
			log.Println("[DEBUG] Requested for username.")
		case request := <-hub.MessageChannel:
			log.Println("[DEBUG] Routing message...")
			hub.ControlRouter.RouteMessage(request.Message, request.Client)
			log.Println("[DEBUG] Routed message.")
		}
	}
}

func (hub *ClientHub) GetAllUserId() []string {
	return hub.UserIdToClient.Keys()
}

func (hub *ClientHub) GetClientFromUserId(userId string) (*Client, bool) {
	client, ok := hub.UserIdToClient.Load(userId)
	if !ok {
		return nil, false
	}
	return client, true
}

func (hub *ClientHub) GetUserIdFromClient(client *Client) (string, bool) {
	userId, ok := hub.ClientToUserId.Load(client)
	if !ok {
		return "", false
	}
	return userId, true
}

func (hub *ClientHub) SendToUserId(userId string, message *messages.Message) {
	otherClient, ok := hub.GetClientFromUserId(userId)
	if ok {
		otherClient.Receive <- message
	}
}

func (hub *ClientHub) setUserId(userId string, client *Client) {
	oldClient, ok := hub.UserIdToClient.Load(userId)
	if ok {
		oldClient.Close()
	}

	hub.UserIdToClient.Store(userId, client)
	hub.ClientToUserId.Store(client, userId)

	client.Receive <- messages.CreateMessage("AUTH", nil, "successfully authenticated as "+userId)
	log.Println("[DEBUG] Setting user id...")
	hub.EventManager.Emit("user-connect", userId)
}

func (hub *ClientHub) closeClient(client *Client) {
	log.Println("[DEBUG] Deleting user id...")
	userId, ok := hub.ClientToUserId.Load(client)
	if ok {
		hub.EventManager.Emit("user-disconnect", userId)
		hub.ClientToUserId.Delete(client)
	}
	_, ok = hub.UserIdToClient.Load(userId)
	if ok {
		hub.UserIdToClient.Delete(userId)
	}
	_, ok = hub.ClientList.Load(client)
	if ok {
		hub.ClientList.Delete(client)
		close(client.Receive)
	}
}

func NewClientHub(
	controlRouter router.ControlRouter,
	eventManager events.EventManager,
) *ClientHub {
	return &ClientHub{
		EventManager:   eventManager,
		ControlRouter:  controlRouter,
		UserIdToClient: datastructures.NewSyncMap[string, *Client](),
		ClientToUserId: datastructures.NewSyncMap[*Client, string](),
		ClientList:     datastructures.NewSyncMap[*Client, bool](),
		AddChannel:     make(chan *Client, 1024),
		UserIdChannel:  make(chan *UserIdRequest, 1024),
		CloseChannel:   make(chan *Client, 1024),
		MessageChannel: make(chan *MessageRequest, 1024),
	}
}
