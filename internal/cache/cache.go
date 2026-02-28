package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

type Cache struct {
	mu       sync.RWMutex
	items    map[string]CacheItem
	ttl      time.Duration
	maxItems int
}

func NewCache(ttl time.Duration, maxItems int) *Cache {
	c := &Cache{
		items:    make(map[string]CacheItem),
		ttl:      ttl,
		maxItems: maxItems,
	}
	go c.startCleanup()
	return c
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) >= c.maxItems {
		for k := range c.items {
			delete(c.items, k)
			break
		}
	}

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok || time.Now().After(item.Expiration) {
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.items {
			if time.Now().After(v.Expiration) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
