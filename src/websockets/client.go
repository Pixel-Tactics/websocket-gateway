package websockets

import (
	"errors"
	"log"

	"pixeltactics.com/websocket-gateway/src/messages"
	jwt_auth "pixeltactics.com/websocket-gateway/src/utils/jwt"

	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{}

type Client struct {
	Id      string
	UserId  string
	Hub     *ClientHub
	Conn    *websocket.Conn
	Receive chan *messages.Message
}

func (client *Client) GetUserId() (string, error) {
	if client.UserId == "" {
		return "", errors.New("user unauthenticated")
	}
	return client.UserId, nil
}

func (client *Client) Send(message *messages.Message) {
	client.Receive <- message
}

func (client *Client) SendToUserId(userId string, message *messages.Message) {
	client.Hub.SendToUserId(userId, message)
}

func (client *Client) Close() {
	client.Hub.CloseChannel <- client
	client.Conn.Close()
}

func (client *Client) handleReceive() {
	defer func() {
		client.Close()
	}()
	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	acceptingAuth := true
	for {
		_, jsonBytes, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Authenticate
		if acceptingAuth {
			userId, err := jwt_auth.Validate(string(jsonBytes))
			if err != nil {
				log.Println(err)
				client.Receive <- messages.Error(err)
				continue
			}
			client.UserId = userId
			client.Hub.UserIdChannel <- &UserIdRequest{
				UserId: userId,
				Client: client,
			}
			acceptingAuth = false
			continue
		}

		message, err := messages.JsonBytesToMessage(jsonBytes)
		if err != nil {
			log.Println(err)
			client.Receive <- messages.Error(err)
			continue
		}
		client.Hub.MessageChannel <- &MessageRequest{
			Message: message,
			Client:  client,
		}
	}
}

func (client *Client) handleSend() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.Receive:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			jsonBytes, err := messages.MessageToJsonBytes(message)
			if err != nil {
				return
			}
			w.Write(jsonBytes)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWebSocket(hub *ClientHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Id:      uuid.NewString(),
		Hub:     hub,
		Conn:    conn,
		Receive: make(chan *messages.Message, 256),
	}

	client.Hub.AddChannel <- client

	go client.handleSend()
	go client.handleReceive()
}
