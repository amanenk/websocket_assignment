package websocket

import (
	"encoding/json"
	"github.com/fdistorted/websocket-practical/models"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"github.com/fdistorted/websocket-practical/server/storage"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var upgrader = websocket.Upgrader{
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

	// reads the message from client
	go func() {
		defer conn.Close()
	outer:
		for {
			var cmd models.CommandBody
			mt, body, err := conn.ReadMessage()
			if err != nil {
				handleClientError(conn, "failed to read from client", err)
				break outer
			}

			logger.Get().Info("got message", zap.Int("mt", mt), zap.String("body", string(body)))

			if mt == websocket.TextMessage {
				// decode the json request to comment
				err = json.Unmarshal(body, &cmd)
				if err != nil {
					logger.Get().Error("failed to parse command JSON", zap.Error(err))
					err := conn.WriteMessage(websocket.TextMessage, []byte("failed to parse json")) // todo replace with json message
					if err != nil {
						handleClientError(conn, "failed to send num connections", err)
						break outer
					}
				} else {
					switch cmd.Command {
					case "SUBSCRIBE":
						storage.ClientStorage.Subscribe(conn)
						break
					case "UNSUBSCRIBE":
						storage.ClientStorage.Unsubscribe(conn)
						break
					case "NUM_CONNECTIONS":
						msg := models.NumConnectionsBody{
							NumConnection: storage.ClientStorage.GetClientsCount(),
						}
						err := conn.WriteJSON(msg)
						if err != nil {
							handleClientError(conn, "failed to send num connections", err)
							break outer
						}
						break
					}
				}

			}
		}
	}()

	storage.ClientStorage.Add(conn)
	logger.Get().Info("client connected", zap.String("client_id", storage.ClientStorage.Clients[conn].ClientId), zap.Int("clients", storage.ClientStorage.GetClientsCount()))
}

func handleClientError(conn *websocket.Conn, msg string, err error) {
	clientId := storage.ClientStorage.Clients[conn].ClientId
	storage.ClientStorage.Delete(conn)
	logger.Get().Error("websocket error, disconnecting",
		zap.Int("clients",
			storage.ClientStorage.GetClientsCount()),
		zap.String("client_id", clientId),
		zap.String("error_msg", msg),
		zap.Error(err))

}
