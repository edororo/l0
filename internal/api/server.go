package api

import (
	"L0/internal/cache"
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"path/filepath"
)

// StartServer запускает HTTP сервер с API и фронтом
func StartServer(ctx context.Context, addr string, dbPool *pgxpool.Pool, orderCache *cache.Cache) {
	router := mux.NewRouter()

	// Эндпоинт API для получения заказа
	router.HandleFunc("/order/{order_uid}", func(w http.ResponseWriter, r *http.Request) {
		GetOrderHandler(w, r, ctx, dbPool, orderCache)
	}).Methods("GET")

	// Раздача frontend.html при заходе на /
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Определяем путь к файлу относительно проекта
		path := filepath.Join("internal", "api", "frontend.html")
		http.ServeFile(w, r, path)
	}).Methods("GET")

	log.Println("HTTP сервер запущен на", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Ошибка запуска HTTP сервера:", err)
	}
}
