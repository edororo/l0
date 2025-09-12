package api

import (
	"L0/internal/cache"
	"L0/internal/db"
	"L0/internal/models"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

// обрабатывает запрос GET /order/{order_uid}
func GetOrderHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, dbPool *pgxpool.Pool, orderCache *cache.Cache) {
	// Достаём order_uid из URL
	vars := mux.Vars(r)
	uid := vars["order_uid"]

	// Сначала ищем в кэше
	if order, ok := orderCache.Get(uid); ok {
		log.Println("Заказ найден в кэше:", uid)
		writeJSON(w, order)
		return
	}

	// Если нет в кэш ищем в БД
	order, err := db.GetOrderByUID(ctx, dbPool, uid)
	if err != nil {
		http.Error(w, "Ошибка при получении заказа", http.StatusInternalServerError)
		return
	}

	// Если заказ не найден
	if order.OrderUID == "" {
		http.Error(w, "Заказ не найден", http.StatusNotFound)
		return
	}

	// Добавляем найденный заказ в кэш
	orderCache.Set(uid, order)

	log.Println("Заказ загружен из БД и добавлен в кэш:", uid)
	writeJSON(w, order)
}

// Вспомогательная функция для отправки JSON
func writeJSON(w http.ResponseWriter, data models.Order) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ") // красиво форматируем JSON
	if err := enc.Encode(data); err != nil {
		http.Error(w, "Ошибка сериализации JSON", http.StatusInternalServerError)
	}
}
