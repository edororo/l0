package cache

import (
	"L0/internal/models"
	"sync"
	"time"
)

type cacheItem struct {
	order      models.Order
	expiration time.Time
}

type Cache struct {
	mu       sync.RWMutex
	data     map[string]cacheItem
	ttl      time.Duration
	maxItems int
}

func NewCache(ttl time.Duration, maxItems int) *Cache {
	c := &Cache{
		data:     make(map[string]cacheItem),
		ttl:      ttl,
		maxItems: maxItems,
	}

	go c.startCleanup()

	return c
}

func (c *Cache) Set(key string, order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// если превышен лимит — удаляем самый старый элемент
	if len(c.data) >= c.maxItems {
		c.evictOldest()
	}

	c.data[key] = cacheItem{
		order:      order,
		expiration: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(key string) (models.Order, bool) {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return models.Order{}, false
	}

	if time.Now().After(item.expiration) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return models.Order{}, false
	}

	return item.order, true
}

// удаляет самый старый элемент
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for k, v := range c.data {
		if oldestKey == "" || v.expiration.Before(oldestTime) {
			oldestKey = k
			oldestTime = v.expiration
		}
	}
	delete(c.data, oldestKey)
}

// автоочистка раз в минуту
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		now := time.Now()

		c.mu.Lock()
		for k, v := range c.data {
			if now.After(v.expiration) {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}
}
