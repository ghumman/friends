package ttlcacheaggregator

import (
	"sync"
	"time"
)

type entry struct {
	value interface{}
	expiresAt time.Time
}

type TTLCache struct {
	mu sync.RWMutex
	entries map[string]entry
	ttl time.Duration
}

func NewTTLCache(ttl time.Duration) *TTLCache {
	return &TTLCache{
		entries: make(map[string] entry),
		ttl: ttl,
	}
}

func (c *TTLCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = entry{
		value: value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	e, exists := c.entries[key]
	defer c.mu.RUnlock()

	if !exists {
		return nil, false
	}

	if time.Now().After(e.expiresAt) {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()

		return nil, false
	}
	return e.value, true
}