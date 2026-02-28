package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"
)

var (
	DatabaseURL   string
	Brokers       string
	Topic         string
	GroupID       string
	CacheTTL      time.Duration
	CacheMaxItems int
	HTTPPort      string
)

func LoadEnv() error {
	_ = godotenv.Load()

	DatabaseURL = os.Getenv("DATABASE")
	Brokers = os.Getenv("BROKERS")
	Topic = os.Getenv("TOPIC")
	GroupID = os.Getenv("GROUP_ID")
	HTTPPort = os.Getenv("HTTP_PORT")

	var err error
	CacheTTL, err = time.ParseDuration(os.Getenv("CACHE_TTL"))
	if err != nil {
		CacheTTL = 5 * time.Minute
	}
	CacheMaxItems = 1000
	if v := os.Getenv("CACHE_MAX_ITEMS"); v != "" {
		fmt.Sscanf(v, "%d", &CacheMaxItems)
	}

	return nil
}
