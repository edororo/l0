package db

import (
	"L0/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

// WarmUpCache загружает все заказы из БД в кэш при старте сервиса
func WarmUpCache(ctx context.Context, db *pgxpool.Pool) error {
	// временное хранилище, чтобы собрать все заказы
	orders := make(map[string]*models.Order)

	// 1. Достаём заказы вместе с delivery и payment
	rows, err := db.Query(ctx, `
		SELECT 
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
			
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
			
			p.transaction, p.request_id, p.currency, p.provider, p.amount, 
			p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		JOIN delivery d ON o.order_uid = d.order_uid
		JOIN payments p ON o.order_uid = p.order_uid
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ord models.Order
		err := rows.Scan(
			&ord.OrderUID, &ord.TrackNumber, &ord.Entry, &ord.Locale, &ord.InternalSignature,
			&ord.CustomerID, &ord.DeliveryService, &ord.ShardKey, &ord.SmID, &ord.DateCreated, &ord.OofShard,
			&ord.Delivery.Name, &ord.Delivery.Phone, &ord.Delivery.Zip, &ord.Delivery.City,
			&ord.Delivery.Address, &ord.Delivery.Region, &ord.Delivery.Email,
			&ord.Payment.Transaction, &ord.Payment.RequestID, &ord.Payment.Currency,
			&ord.Payment.Provider, &ord.Payment.Amount, &ord.Payment.PaymentDT,
			&ord.Payment.Bank, &ord.Payment.DeliveryCost, &ord.Payment.GoodsTotal, &ord.Payment.CustomFee,
		)
		if err != nil {
			return err
		}
		orders[ord.OrderUID] = &ord
	}

	// 2. Достаём все items и добавляем их к заказам
	itemRows, err := db.Query(ctx, `
		SELECT order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items
	`)
	if err != nil {
		return err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var it models.Item
		var orderUID string
		err := itemRows.Scan(
			&orderUID, &it.ChrtID, &it.TrackNumber, &it.Price, &it.RID,
			&it.Name, &it.Sale, &it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status,
		)
		if err != nil {
			return err
		}

		// находим заказ и добавляем к нему item
		if ord, ok := orders[orderUID]; ok {
			ord.Items = append(ord.Items, it)
		}
	}

	// 3. Переносим всё в глобальный кэш
	for _, ord := range orders {
		OrderCache.Set(ord.OrderUID, *ord)
	}

	log.Printf("Кэш прогрет: загружено %d заказов\n", len(orders))
	return nil
}
