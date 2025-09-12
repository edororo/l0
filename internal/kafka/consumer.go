package kafka

import (
	"L0/internal/cache"
	"L0/internal/db"
	"L0/internal/models"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// StartConsumer запускает Kafka consumer для топика orders
func StartConsumer(ctx context.Context, brokers []string, topic string, groupID string, dbPool *pgxpool.Pool, orderCache *cache.Cache) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	log.Printf("Kafka consumer запущен. Topic=%s", topic)

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			continue
		}

		// Проверяем пустое сообщение
		if len(msg.Value) == 0 {
			log.Println("Пустое сообщение, пропускаем")
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Println("Ошибка парсинга JSON:", err, "Сообщение:", string(msg.Value))
			continue
		}

		if order.OrderUID == "" {
			log.Println("Пустой order_uid, пропускаем")
			continue
		}

		// Проверка дубликата в кэше
		if _, exists := orderCache.Get(order.OrderUID); exists {
			log.Println("Заказ уже существует в кэше, пропускаем:", order.OrderUID)
			continue
		}

		// Сохраняем заказ в БД
		if err := db.InsertOrder(ctx, dbPool, order); err != nil {
			log.Printf("Ошибка вставки заказа %s: %v", order.OrderUID, err)
			continue
		}

		// Добавляем в кэш
		orderCache.Set(order.OrderUID, order)

		log.Println("Заказ успешно обработан и добавлен в БД и кэш:", order.OrderUID)
	}
}
