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

func (c *Client) Read(callback func(data map[string]interface{}, err error)) {
	for {
		var read map[string]interface{}
		err := c.socket.ReadJSON(&read)
		if err != nil {
			handleClientError("failed to read from client", c.ClientId, err)
			callback(nil, err)
			return
		}
		logger.Get().Debug("got message")
		callback(read, nil)
	}
}

func (c *Client) Write() {
	for {
		select {
		case toWrite := <-c.writeChannel:
			c.mu.Lock()
			err := c.socket.WriteJSON(toWrite)
			if err != nil {
				handleClientError("failed to write to client", c.ClientId, err)
				// todo decide if we need to close the connection
				return
			}
			c.mu.Unlock()
		}
	}
}

func handleClientError(msg, clientId string, err error) {
	if websocket.IsUnexpectedCloseError(err,
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
		websocket.CloseNoStatusReceived) {
		logger.Get().Error("ws error",
			zap.String("client_id", clientId),
			zap.String("error_msg", msg),
			zap.Error(err))
	}

}
