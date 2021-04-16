package broadcast

import (
	"github.com/fdistorted/websocket-practical/models"
	"github.com/fdistorted/websocket-practical/server/clients-storage"
	"time"
)

func InitBroadcast() chan bool {
	ticker := time.NewTicker(100 * time.Millisecond)
	stopChannel := make(chan bool)

	go func() {
		for {
			select {
			case <-stopChannel:
				ticker.Stop()
				return
			case <-ticker.C:
				msg := models.BeaconBody{
					Timestamp: time.Now().Unix(),
				}

				clients_storage.ClientStorage.Mutex.Lock()
				for _, client := range clients_storage.ClientStorage.Clients {
					if client.Sub {
						go func() {
							msg.ClientId = client.ClientId
							client.Send(msg)
						}()
					}
				}
				clients_storage.ClientStorage.Mutex.Unlock()
			}
		}
	}()
	return stopChannel
}
