package websocket

import (
	"github.com/fdistorted/websocket-practical/server/clients-storage"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
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

	client := clients_storage.NewClient(conn)
	clients_storage.ClientStorage.Add(client)
	logger.Get().Debug("client connected", zap.String("client_id", clients_storage.ClientStorage.Clients[conn].ClientId), zap.Int("clients", clients_storage.ClientStorage.GetClientsCount()))
	defer client.Close()

	go client.Write()
	// reads the message from client
	client.Read()
}
