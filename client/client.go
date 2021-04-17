package client

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/models"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"github.com/fdistorted/websocket-practical/server/websocket/clients"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func Start(url string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	logger.Get().Debug("connecting", zap.String("url", url))

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	client := clients.NewClient(conn)

	defer func() {
		err := client.Close()
		if err != nil {
			logger.Get().Error("failed to close client connection", zap.Error(err))
		}
	}()

	go client.Write()

	done := make(chan bool)

	go func() {
		defer close(done)
		client.Read(func(data map[string]interface{}) {
			//fmt.Printf("got data: %+v\n", data)
			value, ok := data["num_connections"]
			if ok {
				fmt.Printf("num_connections: %v\r", value)
			}
		})
	}()

	client.Send(models.CommandBody{Command: models.NumConnections})

	subscribeAfter := time.Duration(rand.Intn(100)) * time.Millisecond //randomise a bit subscription message
	subscribeTimer := time.NewTimer(subscribeAfter)

outer:
	for {
		select {
		case <-done:
			return
		case <-subscribeTimer.C:
			client.Send(models.CommandBody{Command: models.Subscribe})
		case <-interrupt:
			//log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			break outer
		}
	}
	logger.Get().Info("exiting")
}
