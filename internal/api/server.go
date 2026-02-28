package api

import (
	"L0/internal/service"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func StartServer(ctx context.Context, addr string, svc *service.OrderService) {
	router := mux.NewRouter()

	router.HandleFunc("/order/{order_uid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["order_uid"]

		order, err := svc.Get(ctx, uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			log.Printf("Error encoding order: %v", err)
		}
	}).Methods("GET")

	log.Printf("HTTP server running on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
