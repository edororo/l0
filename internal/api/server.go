package api

import (
	"L0/internal/service"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
)

func StartServer(ctx context.Context, addr string, svc *service.OrderService) {
	router := mux.NewRouter()

	router.HandleFunc("/order/{order_uid}", func(w http.ResponseWriter, r *http.Request) {
		GetOrderHandler(w, r, ctx, svc)
	}).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("internal", "api", "frontend.html")
		http.ServeFile(w, r, path)
	}).Methods("GET")

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("internal/api/static"))),
	)

	log.Println("HTTP сервер запущен на", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Ошибка запуска HTTP сервера:", err)
	}
}
