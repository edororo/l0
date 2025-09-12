package cache

import (
	"L0/internal/models"
	"sync"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

// NewCache создаёт новый пустой кэш
func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]models.Order),
	}
}

// Set добавляет/обновляет заказ по UID
func (c *Cache) Set(uid string, order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[uid] = order
}

// Get достаёт заказ по UID
func (c *Cache) Get(uid string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orders[uid]
	return order, ok
}

// GetAll возвращает все заказы (например, для прогрева кэша из БД при старте)
func (c *Cache) GetAll() map[string]models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	copy := make(map[string]models.Order)
	for k, v := range c.orders {
		copy[k] = v
	}
	return copy
}
