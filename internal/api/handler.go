package api

import (
	"encoding/json"
	"net/http"

	appErr "L0/internal/errors"
	"L0/internal/service"

	"github.com/gorilla/mux"
)

func GetOrder(s *service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := mux.Vars(r)["order_uid"]
		o, err := s.Get(r.Context(), id)

		if err != nil {
			switch err {
			case appErr.ErrInvalidInput:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case appErr.ErrNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, appErr.ErrInternal.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(o)
	}
}
