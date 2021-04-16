package clients_storage

import (
	"encoding/json"
	"github.com/fdistorted/websocket-practical/models"
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

func (c *Client) Read() error {
	for {
		var cmd models.CommandBody
		mt, body, err := c.socket.ReadMessage()
		if err != nil {
			handleClientError(c.socket, "failed to read from client", c.ClientId, err)
			return err
		}

		logger.Get().Debug("got message", zap.Int("mt", mt), zap.String("body", string(body)))

		if mt == websocket.TextMessage {
			// decode the json request to comment
			err = json.Unmarshal(body, &cmd)
			if err != nil {
				logger.Get().Error("failed to parse command JSON", zap.Error(err))
				c.writeChannel <- "failed to parse json"
			} else {
				switch cmd.Command {
				case "SUBSCRIBE":
					ClientStorage.Subscribe(c.socket)
					break
				case "UNSUBSCRIBE":
					ClientStorage.Unsubscribe(c.socket)
					break
				case "NUM_CONNECTIONS":
					msg := models.NumConnectionsBody{
						NumConnection: ClientStorage.GetClientsCount(),
					}
					c.writeChannel <- msg
					break
				}
			}

		}
	}
}

func (c *Client) Write() error {
	for {
		select {
		case toWrite := <-c.writeChannel:
			c.mu.Lock()
			err := c.socket.WriteJSON(toWrite)
			if err != nil {
				handleClientError(c.socket, "failed to send to connection", c.ClientId, err)
				return err
			}
			c.mu.Unlock()
		}
	}
}

func handleClientError(conn *websocket.Conn, msg, clientId string, err error) {
	ClientStorage.Delete(conn)
	logger.Get().Error("websocket error, disconnecting",
		zap.Int("clients", ClientStorage.GetClientsCount()),
		zap.String("client_id", clientId),
		zap.String("error_msg", msg),
		zap.Error(err))
}
