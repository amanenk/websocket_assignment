package clients

import (
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
)

type Client struct {
	mu           sync.Mutex
	socket       *websocket.Conn
	ClientId     string
	Sub          bool
	writeChannel chan interface{}
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ClientId:     uuid.New().String(),
		Sub:          false,
		socket:       conn,
		writeChannel: make(chan interface{}),
	}
}

func (c *Client) Send(data interface{}) {
	c.writeChannel <- data
}

func (c *Client) Close() error {
	return c.socket.Close()
}

func (c *Client) SetSubscribed(sub bool) {
	c.Sub = sub
}

func (c *Client) Read(callback func(data map[string]interface{})) {
	for {
		var read map[string]interface{}
		err := c.socket.ReadJSON(&read)
		if err != nil {
			handleClientError("failed to read from client", c.ClientId, err)
			return
		}
		logger.Get().Debug("got message")
		callback(read)
	}
}

func (c *Client) Write() {
	for {
		select {
		case toWrite := <-c.writeChannel:
			c.mu.Lock()
			err := c.socket.WriteJSON(toWrite)
			if err != nil {
				handleClientError("failed to send to connection", c.ClientId, err)
			}
			c.mu.Unlock()
		}
	}
}

func handleClientError(msg, clientId string, err error) {
	StorageObject.Delete(clientId)
	logger.Get().Error("ws error, disconnecting",
		zap.Int("clients", StorageObject.GetClientsCount()),
		zap.String("client_id", clientId),
		zap.String("error_msg", msg),
		zap.Error(err))
}
