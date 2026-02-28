package repository

import (
	"L0/internal/models"
	"context"
)

type OrderRepository interface {
	Get(ctx context.Context, uid string) (*models.Order, error)
	Create(ctx context.Context, o *models.Order) error
}
