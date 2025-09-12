package db

import (
	"database/sql"
	"os"

	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment")
	}

	connStr := os.Getenv("DATABASE")

	var err error
	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Cannot ping database:", err)
	}

	log.Println("Connected to database successfully")
}
