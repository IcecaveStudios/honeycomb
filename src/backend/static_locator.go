package backend

import (
	"context"
	"sync"
)

// StaticLocator finds a back-end HTTP server based on the server name in TLS
// requests (SNI).
type StaticLocator struct {
	endpoints map[string]*Endpoint
	mutex     sync.RWMutex
}

// Locate finds the back-end HTTP server for the given server name.
func (locator *StaticLocator) Locate(_ context.Context, serverName string) *Endpoint {
	locator.mutex.RLock()
	endpoint := locator.endpoints[serverName]
	locator.mutex.RUnlock()

	return endpoint
}

// Add creates a new mapping from server name to back-end HTTP server.
func (locator *StaticLocator) Add(serverName string, endpoint *Endpoint) {
	locator.mutex.Lock()
	if locator.endpoints == nil {
		locator.endpoints = map[string]*Endpoint{}
	}
	locator.endpoints[serverName] = endpoint
	locator.mutex.Unlock()
}
