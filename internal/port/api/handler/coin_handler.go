package handler

import (
	"encoding/json"
	"net/http"

	"final_course/internal/port/api/dto"
	"github.com/go-chi/chi/v5"
)

type CoinService interface {
	GetCoinPrice(symbol string) (float64, error)
}

type CoinHandler struct {
	Service CoinService
}

func NewCoinHandler(service CoinService) *CoinHandler {
	return &CoinHandler{Service: service}
}

func (h *CoinHandler) GetCoinPrice(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")
	price, err := h.Service.GetCoinPrice(symbol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := dto.CoinResponse{
		Symbol: symbol,
		Price:  price,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
