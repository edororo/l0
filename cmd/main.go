package main

import (
	"L0/internal/api"
	"L0/internal/cache"
	"L0/internal/kafka"
	"L0/internal/service"
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

	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE")
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer pool.Close()
	log.Println("Подключение к PostgreSQL успешно")

	orderCache := cache.NewCache(5*time.Minute, 1000)
	orderService := service.NewOrderService(pool, orderCache)

	go func() {
		api.StartServer(ctx, ":8081", orderService)
	}()

	go func() {
		brokers := []string{os.Getenv("BROKERS")}
		topic := os.Getenv("TOPIC")
		groupID := os.Getenv("GROUP_ID")
		kafka.StartConsumer(ctx, brokers, topic, groupID, orderService)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Завершение работы...")
	cancel()
	time.Sleep(time.Second)
	log.Println("Приложение остановлено")
}
