package main

import (
	"L0/internal/api"
	"L0/internal/cache"
	"L0/internal/db"
	"L0/internal/kafka"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настройки подключения к PostgreSQL
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment")
	}
	databaseURL := os.Getenv("DATABASE")
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}
	defer pool.Close()
	log.Println("Подключение к PostgreSQL успешно")

	//  Инициализация кэша
	db.OrderCache = cache.NewCache()
	log.Println("Кэш заказов создан")

	//  HTTP сервер
	go func() {
		api.StartServer(ctx, ":8081", pool, db.OrderCache)
	}()

	//  Kafka consumer
	go func() {
		brokers := []string{os.Getenv("BROKERS")}
		topic := os.Getenv("TOPIC")
		groupID := os.Getenv("GROUP_ID")
		kafka.StartConsumer(ctx, brokers, topic, groupID, pool, db.OrderCache)
	}()

	//  Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Завершение работы приложения...")
	cancel()
	time.Sleep(1 * time.Second) // даём горутинам завершиться
	log.Println("Приложение остановлено")
}
