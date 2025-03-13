package apiserver

import (
	"context"
	"encoding/json"
	"final_course/internal/cases"
	"final_course/internal/entities"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// RESTHandler - обработчик REST API
type RESTHandler struct {
	service *cases.Service
}

// NewRESTHandler создает новый RESTHandler
func NewRESTHandler(service cases.Service) *RESTHandler {
	return &RESTHandler{
		service: service,
	}
}

// getCoins обрабатывает запрос на получение последних курсов валют
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

// convertToDTO преобразует массив entities.Coin в DTO-формат
func convertToDTO(coins []entities.Coin) *CoinsResponse {
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

// SetupRoutes регистрирует маршруты API
func (h *RESTHandler) SetupRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			middleware.Recoverer,
			middleware.RequestID,
		)

		r.Get("/coins", h.getCoins)
		r.Get("/coins/max", h.getMaxRate)
		r.Get("/coins/min", h.getMinRate)
		r.Get("/coins/avg", h.getAvgRate)
	})
}

func (h *RESTHandler) getMaxRate(w http.ResponseWriter, r *http.Request) {
	titles := []string{"title1", "title2"} // Замените на реальные данные
	coins, err := h.service.GetMaxRate(context.Background(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Вернуть данные пользователю
}

func (h *RESTHandler) getMinRate(w http.ResponseWriter, r *http.Request) {
	titles := []string{"title1", "title2"} // Замените на реальные данные
	coins, err := h.service.GetMinRate(context.Background(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Вернуть данные пользователю
}

func (h *RESTHandler) getAvgRate(w http.ResponseWriter, r *http.Request) {
	titles := []string{"title1", "title2"} // Замените на реальные данные
	coins, err := h.service.GetAvgRate(context.Background(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Вернуть данные пользователю
}
