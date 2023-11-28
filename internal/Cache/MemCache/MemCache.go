package MemCache

import (
	"short-link/internal/Config"
	cache_interface "short-link/internal/Core/Ports"
	"sync"
	"time"
)

// CacheItem represents a single cache item
type CacheItem struct {
	Content    []interface{}
	Expiration int64
}

// Cache represents an in-memory cache
type Cache struct {
	Config *Config.Config
	items  map[string]CacheItem
	mutex  sync.RWMutex
}

// Set adds an item to the cache
func (c *Cache) SetSlice(key string, content []interface{}, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = CacheItem{
		Content:    content,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

// Get retrieves an item from the cache
func (c *Cache) GetSlice(key string) ([]interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.items[key]
	if !found || (item.Expiration > 0 && time.Now().UnixNano() > item.Expiration) {
		return nil, false
	}
	return item.Content, true
}

// CreateMemCache creates a new Cache object
func CreateMemCache(cfg *Config.Config) cache_interface.MemCacheInterface {
	return &Cache{
		Config: cfg,
		items:  make(map[string]CacheItem),
	}
}
