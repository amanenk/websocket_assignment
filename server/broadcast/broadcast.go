package broadcast

import (
	"github.com/fdistorted/websocket-practical/models"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"github.com/fdistorted/websocket-practical/server/storage"
	"go.uber.org/zap"
	"time"
)

func InitBroadcast() {
	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for conn, client := range storage.ClientStorage.Clients {
					if client.Sub {
						msg := models.BeaconBody{
							ClientId:  client.ClientId,
							Timestamp: time.Now().Unix(),
						}
						err := conn.WriteJSON(msg)
						if err != nil {
							logger.Get().Info("client disconnected", zap.String("client_id", client.ClientId), zap.Int("clients", storage.ClientStorage.GetClientsCount()))
							storage.ClientStorage.Delete(conn)
							break
						}
					}
				}
			}
		}
	}()
}
