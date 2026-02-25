package api

import (
	"L0/internal/service"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// обрабатывает запрос GET /order/{order_uid}
func GetOrderHandler(w http.ResponseWriter, r *http.Request, ctx context.Context, service *service.OrderService) {
	vars := mux.Vars(r)
	uid := vars["order_uid"]

	order, err := service.GetOrder(ctx, uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
