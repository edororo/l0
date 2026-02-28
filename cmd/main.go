package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"L0/internal/api"
	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/kafka"
	"L0/internal/repository/postgres"
	"L0/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := config.LoadEnv(); err != nil {
		log.Println("Warning: .env not found")
	}

	pool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	defer pool.Close()
	log.Println("PostgreSQL connected")

	ttl, maxItems := config.CacheTTL, config.CacheMaxItems
	orderCache := cache.NewCache(ttl, maxItems)

	repo := postgres.NewRepo(pool)
	orderService := service.New(repo, orderCache)

	go api.StartServer(ctx, config.HTTPPort, orderService)
	go kafka.StartConsumer(ctx, []string{config.Brokers}, config.Topic, config.GroupID, repo, orderCache)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down...")
	cancel()
	time.Sleep(1 * time.Second)
	log.Println("Application stopped")
}
