package cache

import (
	"sync"
	"time"
)

// Cache is a thread-safe in-memory cache with TTL support
type Cache[K comparable, V any] struct {
	items map[K]cacheItem[V]
	mu    sync.RWMutex
	ttl   time.Duration
}

type cacheItem[V any] struct {
	value      V
	expiration time.Time
}

// New creates a new cache with the specified TTL
func New[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]cacheItem[V]),
		ttl:   ttl,
	}
}

// Set stores a value in the cache with the configured TTL
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem[V]{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}

// Get retrieves a value from the cache if it exists and hasn't expired
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		var zero V
		return zero, false
	}

	if time.Now().After(item.expiration) {
		var zero V
		return zero, false
	}

	return item.value, true
}

// Delete removes a value from the cache
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[K]cacheItem[V])
}

// Cleanup removes expired items from the cache
func (c *Cache[K, V]) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.expiration) {
			delete(c.items, key)
		}
	}
}

// StartCleanup starts a goroutine that periodically cleans up expired items
func (c *Cache[K, V]) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.Cleanup()
		}
	}()
}
