package client

import (
	"github.com/fdistorted/websocket-practical/models"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
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

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func() {
		err := c.Close()
		if err != nil {
			logger.Get().Error("failed to close client connection", zap.Error(err))
		}
	}()

	done := make(chan bool)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logger.Get().Error("read:", zap.Error(err))
				return
			}
			logger.Get().Debug("recv:", zap.String("msg", string(message)))

		}
	}()

	subscribeAfter := time.Duration(rand.Intn(100)) * time.Millisecond //randomise a bit subscription message
	unsubscribeAfter := subscribeAfter + 120*time.Second
	subscribeTimer := time.NewTimer(subscribeAfter)
	unsubscribeTimer := time.NewTimer(unsubscribeAfter)
	getConnectionsTicker := time.NewTicker(1 * time.Second)

outer:
	for {
		select {
		case <-done:
			return
		case <-subscribeTimer.C:
			cmd := models.CommandBody{Command: models.Subscribe}
			err := c.WriteJSON(cmd)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-unsubscribeTimer.C:
			cmd := models.CommandBody{Command: models.Unsubscribe}
			err := c.WriteJSON(cmd)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-getConnectionsTicker.C:
			cmd := models.CommandBody{Command: models.NumConnections}
			err := c.WriteJSON(cmd)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
