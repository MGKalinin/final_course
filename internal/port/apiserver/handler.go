package apiserver

import (
	"encoding/json"
	"net/http"

	"final_course/internal/cases"
)

// RESTHandler - обработчик REST API
type RESTHandler struct {
	service cases.Service
}

// NewRESTHandler создает новый RESTHandler
func NewRESTHandler(service cases.Service) *RESTHandler {
	return &RESTHandler{
		service: service,
	}
}

func (h *RESTHandler) getCoins(w http.ResponseWriter, r *http.Request) {
	var req CoinsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	result, err := h.service.GetLastRates(ctx, req.Titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dtoResult := convertToDTO(result)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoResult)
}

func convertToDTO(coins []cases.Coin) *CoinsResponse {
	dtoCoins := make([]CoinDTO, len(coins))
	for i, coin := range coins {
		dtoCoins[i] = CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
			Date:  coin.Date,
		}
	}
	return &CoinsResponse{Coins: dtoCoins}
}
