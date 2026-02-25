package service

import (
	"L0/internal/cache"
	"L0/internal/db"
	"L0/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService struct {
	db    *pgxpool.Pool
	cache *cache.Cache
}

func NewOrderService(db *pgxpool.Pool, cache *cache.Cache) *OrderService {
	return &OrderService{db: db, cache: cache}
}

func (s *OrderService) GetOrder(ctx context.Context, uid string) (models.Order, error) {
	if order, ok := s.cache.Get(uid); ok {
		return order, nil
	}

	order, err := db.GetOrderByUID(ctx, s.db, uid)
	if err != nil {
		return models.Order{}, err
	}

	s.cache.Set(uid, order)
	return order, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, order models.Order) error {
	if err := db.InsertOrder(ctx, s.db, order); err != nil {
		return err
	}
	s.cache.Set(order.OrderUID, order)
	return nil
}
