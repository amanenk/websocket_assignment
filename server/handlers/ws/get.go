package ws

import (
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"github.com/fdistorted/websocket-practical/server/websocket/clients"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var upgrader = websocket.Upgrader{
	//ReadBufferSize:  0,
	//WriteBufferSize: 0,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Get(w http.ResponseWriter, r *http.Request) {
	logger.Get().Debug("got request")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Get().Error("failed to create conn connection %v", zap.Error(err))
		//todo return error
	}

	client := clients.NewClient(conn)
	clients.StorageObject.Add(client)
	logger.Get().Debug("client connected", zap.String("client_id", client.ClientId), zap.Int("clients", clients.StorageObject.GetClientsCount()))
	defer func() {
		err := client.Close()
		if err != nil {
			logger.Get().Error("failed to close client connection", zap.Error(err))
		}
	}()

	go client.Write()
	// reads the message from client
	client.Read()
}
