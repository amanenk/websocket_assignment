package broadcast

import (
	"github.com/fdistorted/websocket-practical/models"
	"github.com/fdistorted/websocket-practical/server/websocket/clients"
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

				clients.StorageObject.Mutex.Lock()
				for _, client := range clients.StorageObject.Clients {
					if client.Sub {
						msg.ClientId = client.ClientId
						client.Send(msg)
					}
				}
				clients.StorageObject.Mutex.Unlock()
			}
		}
	}()
	return stopChannel
}
