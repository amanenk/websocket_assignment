package main

import (
	"github.com/fdistorted/websocket-practical/server"
)

func main() {
	//start server
	go func() {
		server.Start()
	}()

	select {}
}
