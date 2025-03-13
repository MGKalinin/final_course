package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type APIHandler struct {
	service ServiceInterface
}

func NewAPIHandler(service ServiceInterface) *APIHandler {
	return &APIHandler{service: service}
}

func (h *APIHandler) SetupRoutes(r *chi.Mux) {
	r.Get("/coins", h.GetCoins)
}

func (h *APIHandler) GetCoins(w http.ResponseWriter, r *http.Request) {
	coins, err := h.service.GetCoins()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(coins)
}
