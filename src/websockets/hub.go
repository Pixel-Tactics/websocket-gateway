package websockets

import (
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
			hub.ClientList.Store(client, true)
		case client := <-hub.CloseChannel:
			hub.closeClient(client)
		case request := <-hub.UserIdChannel:
			hub.setUserId(request.UserId, request.Client)
		case request := <-hub.MessageChannel:
			go hub.ControlRouter.RouteMessage(request.Message, request.Client)
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

func (hub *ClientHub) setUserId(userId string, client *Client) {
	oldClient, ok := hub.UserIdToClient.Load(userId)
	if ok {
		oldClient.Close()
	}

	hub.UserIdToClient.Store(userId, client)
	hub.ClientToUserId.Store(client, userId)
}

func (hub *ClientHub) closeClient(client *Client) {
	userId, ok := hub.ClientToUserId.Load(client)
	if ok {
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
	router router.ControlRouter,
) *ClientHub {
	return &ClientHub{
		ControlRouter:  router,
		UserIdToClient: datastructures.NewSyncMap[string, *Client](),
		ClientToUserId: datastructures.NewSyncMap[*Client, string](),
		ClientList:     datastructures.NewSyncMap[*Client, bool](),
		AddChannel:     make(chan *Client, 1024),
		UserIdChannel:  make(chan *UserIdRequest, 1024),
		CloseChannel:   make(chan *Client, 1024),
		MessageChannel: make(chan *MessageRequest, 1024),
	}
}
