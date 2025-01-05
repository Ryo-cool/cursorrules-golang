package cache

import (
	"math"
	"sync"
	"time"
)

// Item represents a cached item
type Item struct {
	Value      interface{}
	Expiration int64
	LastAccess int64
}

// Cache represents an in-memory cache with LRU eviction
type Cache struct {
	items     map[string]Item
	mu        sync.RWMutex
	maxItems  int
	hitCount  uint64
	missCount uint64
}

// Config represents cache configuration
type Config struct {
	MaxItems        int
	CleanupInterval time.Duration
}

// New creates a new cache instance with configuration
func New(config Config) *Cache {
	if config.MaxItems <= 0 {
		config.MaxItems = 1000 // default size
	}
	if config.CleanupInterval <= 0 {
		config.CleanupInterval = time.Minute
	}

	cache := &Cache{
		items:    make(map[string]Item),
		maxItems: config.MaxItems,
	}

	go cache.cleanup(config.CleanupInterval)
	return cache
}

// Set adds an item to the cache
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict items
	if len(c.items) >= c.maxItems {
		c.evictLRU()
	}

	c.items[key] = Item{
		Value:      value,
		Expiration: time.Now().Add(duration).UnixNano(),
		LastAccess: time.Now().UnixNano(),
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if !exists {
		c.missCount++
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		delete(c.items, key)
		c.missCount++
		return nil, false
	}

	// Update last access time
	item.LastAccess = time.Now().UnixNano()
	c.items[key] = item
	c.hitCount++

	return item.Value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// GetStats returns cache statistics
func (c *Cache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := float64(c.hitCount + c.missCount)
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(c.hitCount) / total
	}

	return map[string]interface{}{
		"size":       len(c.items),
		"max_size":   c.maxItems,
		"hit_count":  c.hitCount,
		"miss_count": c.missCount,
		"hit_rate":   hitRate,
	}
}

// cleanup removes expired items from the cache
func (c *Cache) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for key, item := range c.items {
			if now > item.Expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// evictLRU removes the least recently used item
func (c *Cache) evictLRU() {
	var oldestKey string
	var oldestAccess int64 = math.MaxInt64

	for key, item := range c.items {
		if item.LastAccess < oldestAccess {
			oldestAccess = item.LastAccess
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}
