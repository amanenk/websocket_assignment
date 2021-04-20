package ws

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/models"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"github.com/fdistorted/websocket-practical/websocket/client"
	storage2 "github.com/fdistorted/websocket-practical/websocket/storage"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler struct {
	storage *storage2.Storage
}

func NewWebsocketHandler(storage *storage2.Storage) *WebsocketHandler {
	return &WebsocketHandler{storage: storage}
}

func (wh *WebsocketHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger.Get().Debug("got request")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Get().Error("failed to create conn connection %v", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "not a websocket request")
		return
	}

	client := client.NewClient(conn)
	wh.storage.Add(client)
	logger.Get().Debug("client connected", zap.String("client_id", client.ClientId), zap.Int("client", wh.storage.GetClientsCount()))
	defer func() {
		logger.Get().Debug("closing connection", zap.String("client_id", client.ClientId))
		client.SetSubscribed(false)
		err := client.Close()
		if err != nil {
			logger.Get().Error("failed to close client connection", zap.Error(err))
		}
	}()

	go client.Write()
	// reads the message from client
	client.Read(func(data map[string]interface{}, err error) {
		if err != nil {
			//improve error handling
			wh.storage.Delete(client.ClientId)
			return
		}
		command, ok := data["command"]
		if ok {
			switch command.(string) {
			case models.Subscribe:
				client.SetSubscribed(true)
				break
			case models.Unsubscribe:
				client.SetSubscribed(false)
				break
			case models.NumConnections:
				msg := models.NumConnectionsBody{
					NumConnection: wh.storage.GetClientsCount(),
				}
				client.Send(msg)
				break
			default:
				logger.Get().Warn("unsupported command", zap.String("command", data["command"].(string)))
				break
			}
		}
	})
}
