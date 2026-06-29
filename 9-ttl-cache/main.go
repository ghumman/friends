package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	value      interface{}
	expiryTime time.Time
}

type TTLCache struct {
	cache map[string]CacheItem
	mu    sync.RWMutex // 💡 Fixed: Changed to RWMutex to support RLock
}

func NewTTLCache() *TTLCache {
	return &TTLCache{
		cache: make(map[string]CacheItem),
	}
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, ok := c.cache[key]

	if !ok {
		c.mu.RUnlock() // Make sure to unlock before returning!
		return nil, false
	}

	if time.Now().After(item.expiryTime) {
		c.mu.RUnlock() // 💡 CRITICAL: Explicitly release Read Lock BEFORE upgrading to Write Lock!

		c.mu.Lock()
		// Double-check pattern (Perfect!)
		if currentItem, exists := c.cache[key]; exists && time.Now().After(currentItem.expiryTime) {
			delete(c.cache, key)
		}
		c.mu.Unlock()
		return nil, false
	}

	c.mu.RUnlock() // Release read lock for successful path
	return item.value, true
}

func (c *TTLCache) Set(key string, value interface{}, ttl time.Duration) { // 💡 Typo fix: changed ttl type back to time.Duration
	c.mu.Lock()
	defer c.mu.Unlock()

	item := CacheItem{
		value:      value,
		expiryTime: time.Now().Add(ttl), // 💡 Typo fix: Calculate relative expiration
	}

	c.cache[key] = item
}