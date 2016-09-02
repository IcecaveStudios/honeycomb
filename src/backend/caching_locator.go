package backend

import (
	"context"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// CachingLocator wraps another locator to provide basic caching functionality.
type CachingLocator struct {
	// The actual locator used to find endpoints
	Inner Locator

	// How long a positive (successful) "locate" response should be cached.
	PositiveTTL time.Duration

	// How long a negative (unsuccessful) "locate" response should be cached.
	NegativeTTL time.Duration

	// The maximum number of cached endpoints.
	MaxSize uint

	// An optional logger for information about the locator.
	Logger *log.Logger

	cache atomic.Value
	mutex sync.Mutex
}

// Locate finds the back-end HTTP server for the given server name.
func (locator *CachingLocator) Locate(ctx context.Context, serverName string) *Endpoint {
	serverName = strings.ToLower(serverName)

	// Look in the cache ...
	cache, _ := locator.cache.Load().(cacheEntries)
	if endpoint, ok := cache.fetch(serverName); ok {
		return endpoint
	}

	// The server name is not in the cache, forward it to the inner locator ...
	if locator.Logger != nil {
		locator.Logger.Printf("backend: Cache miss for '%s'", serverName)
	}

	return locator.forward(ctx, serverName)
}

func (locator *CachingLocator) forward(ctx context.Context, serverName string) *Endpoint {
	locator.mutex.Lock()
	defer locator.mutex.Unlock()

	// Check if another goroutine already added this server name to the cache ...
	cache, _ := locator.cache.Load().(cacheEntries)
	if endpoint, ok := cache.fetch(serverName); ok {
		return endpoint
	}

	// Query the inner locator ...
	endpoint := locator.Inner.Locate(ctx, serverName)
	ttl := locator.computeTTL(endpoint)

	if locator.Logger != nil {
		if endpoint == nil {
			locator.Logger.Printf(
				"backend: Caching unresolvable route from '%s' for %s",
				serverName,
				ttl,
			)
		} else {
			locator.Logger.Printf(
				"backend: Caching route from '%s' to '%s' (%s) for %s",
				serverName,
				endpoint.Address,
				endpoint.Description,
				ttl,
			)
		}
	}

	entry := cacheEntry{
		ExpiresAt: time.Now().Add(ttl),
		Endpoint:  endpoint,
	}

	// And store the result in the cache ...
	cache = cache.update(serverName, entry)
	locator.cache.Store(cache)

	return endpoint
}

func (locator *CachingLocator) computeTTL(endpoint *Endpoint) time.Duration {
	var ttl time.Duration

	if endpoint == nil {
		ttl = locator.NegativeTTL
	} else {
		ttl = locator.PositiveTTL
	}

	if ttl == 0 {
		return time.Duration(15 * time.Second)
	}

	return ttl
}

type cacheEntries map[string]cacheEntry
type cacheEntry struct {
	ExpiresAt time.Time
	Endpoint  *Endpoint
}

func (entries cacheEntries) fetch(serverName string) (*Endpoint, bool) {
	entry, hasEntry := entries[serverName]

	if hasEntry && entry.ExpiresAt.After(time.Now()) {
		return entry.Endpoint, true
	}

	return nil, false
}

func (entries cacheEntries) update(
	serverName string,
	entry cacheEntry,
) cacheEntries {
	now := time.Now()
	updated := cacheEntries{serverName: entry}

	for name, entry := range entries {
		if entry.ExpiresAt.After(now) && name != serverName {
			updated[name] = entry
		}
	}

	return updated
}
