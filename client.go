package main

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/client"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		time.Sleep(time.Millisecond * 1)
		go func() {
			wg.Add(1)
			client.Start()
			wg.Done()
		}()
		fmt.Printf("\r%d clients", i+1)
	}

	wg.Wait()
}
