package main

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/filelimits"
	"github.com/fdistorted/websocket-practical/models"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"github.com/fdistorted/websocket-practical/websocket/client"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func Start(url string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	client := client.NewClient(conn)

	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Println("failed to close client connection")
		}
	}()

	go client.Write()

	done := make(chan bool)

	go func() {
		defer close(done)
		client.Read(func(data map[string]interface{}, err error) {
			if err != nil {
				// todo handle it somehow
				return
			}
			value, ok := data["num_connections"]
			if ok {
				fmt.Printf("num_connections: %v   \r", value)
			}
		})
	}()

	client.Send(models.CommandBody{Command: models.NumConnections})

	subscribeAfter := time.Duration(rand.Intn(50)) * time.Millisecond //randomise a bit subscription message
	subscribeTimer := time.NewTimer(subscribeAfter)

	for {
		select {
		case <-done:
			return
		case <-subscribeTimer.C:
			client.Send(models.CommandBody{Command: models.Subscribe})
		case <-interrupt:
			//fmt.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			client.SendClose()
			select {
			case <-done:
			case <-time.After(2 * time.Second):
			}
			return
		}
	}
}

func main() {
	filelimits.MaxOpenFiles()

	err := logger.Load()
	if err != nil {
		log.Fatalf("Failed to laod logger: %v", err)
	}

	args := os.Args[1:]

	connections := 5000
	url := "ws://localhost:5000/ws"

	if len(args) >= 1 {
		url = args[0]
	}
	if len(args) >= 2 {
		parsed, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("failed to parse argument. try again")
			return
		}
		if parsed > 20000 {
			fmt.Println("number of connections is too big")
			return
		}
		connections = parsed
	}

	var wg sync.WaitGroup

	fmt.Printf("connecting to: %s with %d connections\n", url, connections)

	for i := 0; i < connections; i++ {
		time.Sleep(time.Millisecond * 1)
		wg.Add(1)
		go func() {
			Start(url)
			wg.Done()
		}()
	}

	wg.Wait()
}
