package storage

import (
	"github.com/fdistorted/websocket-practical/websocket/client"
	"sync"
)

// SafeCounter is safe to use concurrently.
type Storage struct {
	Mutex   sync.RWMutex
	Clients map[string]*client.Client
}

func (c *Storage) Delete(clientId string) {
	c.Mutex.Lock()
	delete(c.Clients, clientId)
	c.Mutex.Unlock()
}

func (c *Storage) Add(client *client.Client) {
	c.Mutex.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.Clients[client.ClientId] = client
	c.Mutex.Unlock()
}

func (c *Storage) GetClientsCount() int {
	c.Mutex.RLock()
	clientsCount := len(c.Clients)
	c.Mutex.RUnlock()
	return clientsCount
}

func (c *Storage) GetClientId() int {
	c.Mutex.RLock()
	clientId := len(c.Clients)
	c.Mutex.RUnlock()
	return clientId
}

//var StorageObject = &Storage{Clients: make(map[string]*Client)}

func NewStorage() *Storage {
	return &Storage{Clients: make(map[string]*client.Client)}
}
