package kafka

import (
	"L0/internal/cache"
	"L0/internal/models"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

func StartConsumer(ctx context.Context, brokers []string, topic string, groupID string, repo interface {
	Get(ctx context.Context, uid string) (*models.Order, error)
}, cache *cache.Cache) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Kafka consumer shutting down...")
				return
			}
			log.Println("Error reading message:", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		if order.OrderUID == "" {
			log.Println("Empty order UID, skipping")
			continue
		}

		cache.Set(order.OrderUID, order)
		log.Println("Order cached:", order.OrderUID)
	}
}
