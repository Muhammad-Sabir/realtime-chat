package models

import (
	"net"
	"sync"
)

type ConnectionsStore struct {
	connections map[*User]net.Addr
	mu          sync.RWMutex
}

func NewConnectionsStore() *ConnectionsStore {
	return &ConnectionsStore{
		connections: make(map[*User]net.Addr),
	}
}

func (cs *ConnectionsStore) AddConnection(user *User, addr net.Addr) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.connections[user] = addr
}

func (cs *ConnectionsStore) RemoveConnection(user *User) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.connections, user)
}
