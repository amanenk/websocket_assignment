package main

import (
	"github.com/fdistorted/websocket-practical/server"
	"time"
)

func main() {
	//start server
	go func() {
		server.Start()
	}()

	time.Sleep(time.Millisecond * 100)

	//todo start client

	select {}
}
