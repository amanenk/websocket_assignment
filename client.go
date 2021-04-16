package main

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/client"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	err := logger.Load()
	if err != nil {
		log.Fatalf("Failed to laod logger: %v", err)
	}

	args := os.Args[1:]

	connections := 5000
	url := "ws://localhost:5000/ws"

	if len(args) == 1 {
		url = args[0]
	}
	if len(args) >= 2 {
		parsed, err := strconv.Atoi(args[0])
		if err != nil {
			logger.Get().Fatal("failed to parse argument. try again")
		}
		if parsed > 20000 {
			logger.Get().Fatal("number is too big")
		}
		connections = parsed
	}

	var wg sync.WaitGroup

	logger.Get().Info("connecting...", zap.Int("connections", connections), zap.String("url", url))

	for i := 0; i < connections; i++ {
		time.Sleep(time.Millisecond * 1)
		go func() {
			wg.Add(1)
			client.Start(url)
			wg.Done()
		}()
		fmt.Printf("\r%d connections", i+1)
	}

	wg.Wait()
}
