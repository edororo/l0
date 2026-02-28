package service

import (
	"L0/internal/cache"
	"L0/internal/models"
	"context"
	"fmt"
)

type OrderRepository interface {
	Get(ctx context.Context, uid string) (*models.Order, error)
}

type OrderService struct {
	repo  OrderRepository
	cache *cache.Cache
}

func New(repo OrderRepository, c *cache.Cache) *OrderService {
	return &OrderService{repo: repo, cache: c}
}

func (s *OrderService) Get(ctx context.Context, uid string) (*models.Order, error) {
	if uid == "" {
		return nil, fmt.Errorf("Invalid order UID")
	}

	if val, ok := s.cache.Get(uid); ok {
		o := val.(models.Order)
		return &o, nil
	}

	order, err := s.repo.Get(ctx, uid)
	if err != nil {
		return nil, err
	}

	s.cache.Set(uid, *order)
	return order, nil
}
