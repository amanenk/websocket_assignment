package storage

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	ClientId string
	Sub      bool
}

// SafeCounter is safe to use concurrently.
type ClientsStorage struct {
	mu      sync.Mutex
	Clients map[*websocket.Conn]*Client
}

func (c *ClientsStorage) Subscribe(ws *websocket.Conn) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.Clients[ws].Sub = true
	c.mu.Unlock()
}

func (c *ClientsStorage) Unsubscribe(ws *websocket.Conn) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.Clients[ws].Sub = false
	c.mu.Unlock()
}

func (c *ClientsStorage) Delete(ws *websocket.Conn) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	delete(c.Clients, ws)
	c.mu.Unlock()
}

func (c *ClientsStorage) Add(ws *websocket.Conn) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.Clients[ws] = &Client{
		uuid.New().String(),
		false,
	}
	c.mu.Unlock()
}

func (c *ClientsStorage) GetClientsCount() int {
	return len(c.Clients)
}

var ClientStorage = ClientsStorage{Clients: make(map[*websocket.Conn]*Client)}
