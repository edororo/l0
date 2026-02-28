package postgres

import (
	"context"
	"fmt"

	"L0/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Get(ctx context.Context, uid string) (*models.Order, error) {
	var order models.Order
	err := r.pool.QueryRow(ctx,
		`SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
		delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid=$1`, uid).
		Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return nil, fmt.Errorf("Order not found")
	}
	return &order, nil
}
