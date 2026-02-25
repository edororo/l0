package kafka

import (
	"L0/internal/models"
	"L0/internal/service"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

func StartConsumer(ctx context.Context, brokers []string, topic string, groupID string, service *service.OrderService) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Printf("Kafka consumer запущен. Topic=%s", topic)

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Consumer остановлен")
				return
			}
			log.Println("Ошибка чтения сообщения:", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Println("Ошибка JSON:", err)
			continue
		}

		if order.OrderUID == "" {
			log.Println("Пустой order_uid")
			continue
		}

		if err := service.CreateOrder(ctx, order); err != nil {
			log.Println("Ошибка сохранения:", err)
			continue
		}

		log.Println("Заказ обработан:", order.OrderUID)
	}
}
