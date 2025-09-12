package db

import (
	"L0/internal/cache"
	"L0/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var OrderCache *cache.Cache

// InsertOrder вставляет заказ в БД и кэш
func InsertOrder(ctx context.Context, db *pgxpool.Pool, order models.Order) error {
	// начинаем транзакцию
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Ошибка начала транзакции: %w", err)
	}
	// откат в случае паники или ошибки
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Сначала вставляем запись в orders
	_, err = tx.Exec(ctx, `
		INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("Ошибка вставки заказа: %w", err)
	}

	// Вставляем delivery
	_, err = tx.Exec(ctx, `
		INSERT INTO delivery(order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("Ошибка вставки delivery: %w", err)
	}

	// Вставляем payment
	_, err = tx.Exec(ctx, `
		INSERT INTO payment(order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("Ошибка вставки payment: %w", err)
	}

	// Вставляем items
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO items(order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("Ошибка вставки item: %w", err)
		}
	}

	// Коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("Ошибка коммита транзакции: %w", err)
	}

	// Добавляем в кэш
	if OrderCache != nil {
		OrderCache.Set(order.OrderUID, order)
	}
	log.Println("Заказ успешно добавлен в БД и кэш:", order.OrderUID)
	return nil
}

// достаёт заказ из кэша или БД
func GetOrderByUID(ctx context.Context, db *pgxpool.Pool, uid string) (models.Order, error) {
	// Проверка кэша
	if OrderCache != nil {
		if order, ok := OrderCache.Get(uid); ok {
			log.Println("[CACHE] Найден заказ:", uid)
			return order, nil
		}
	}

	var order models.Order

	// Получаем заказ + delivery + payment через JOIN
	err := db.QueryRow(ctx, `
		SELECT 
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id,
			o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
			p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		JOIN delivery d ON o.order_uid = d.order_uid
		JOIN payment p ON o.order_uid = p.order_uid
		WHERE o.order_uid = $1
	`, uid).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID,
		&order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
		&order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)
	if err != nil {
		return models.Order{}, fmt.Errorf("Ошибка получения заказа: %w", err)
	}

	// Получаем items
	rows, err := db.Query(ctx, `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1
	`, uid)
	if err != nil {
		return models.Order{}, fmt.Errorf("Ошибка запроса items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status); err != nil {
			return models.Order{}, fmt.Errorf("Ошибка сканирования item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	// Кэшируем заказ
	if OrderCache != nil {
		OrderCache.Set(uid, order)
		log.Println("[DB] Заказ загружен и добавлен в кэш:", uid)
	}

	return order, nil
}
