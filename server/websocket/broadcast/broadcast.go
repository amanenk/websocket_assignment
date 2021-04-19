package broadcast

import (
	"github.com/fdistorted/websocket-practical/models"
	"github.com/fdistorted/websocket-practical/server/websocket/storage"
	"time"
)

func InitBroadcast(storage *storage.Storage) chan bool {
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

				storage.Mutex.Lock()
				for _, client := range storage.Clients {
					if client.Sub {
						msg.ClientId = client.ClientId
						client.Send(msg)
					}
				}
				storage.Mutex.Unlock()
			}
		}
	}()
	return stopChannel
}
