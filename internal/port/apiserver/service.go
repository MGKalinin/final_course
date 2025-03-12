package apiserver

import (
	"context"
	"net/http"
	"time"

	"final_course/internal/cases"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// CoinDTO struct DTO для монет
type CoinDTO struct {
	Title string    `json:"title"`
	Rate  float64   `json:"rate"`
	Date  time.Time `json:"date"`
}

type CoinsRequest struct {
	Titles []string `json:"titles"`
}

type CoinsResponse struct {
	Coins []CoinDTO `json:"coins"`
}

// APIError struct обработка ошибок API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAPIError(code int, message string) error {
	return APIError{Code: code, Message: message}
}

func (e APIError) Error() string {
	return e.Message
}

func (e APIError) ToHTTPCode() int {
	return e.Code
}

// RESTPort определяет интерфейс REST API
type RESTPort interface {
	GetCoins(ctx context.Context, titles []string) (*CoinsResponse, error)
	GetMaxRate(ctx context.Context, titles []string) (*CoinsResponse, error)
	GetMinRate(ctx context.Context, titles []string) (*CoinsResponse, error)
	GetAvgRate(ctx context.Context, titles []string) (*CoinsResponse, error)
}

// RESTHandler реализует обработчики REST API
type RESTHandler struct {
	service cases.Service
}

func NewRESTHandler(service cases.Service) *RESTHandler {
	return &RESTHandler{
		service: service,
	}
}

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

func (h *RESTHandler) getCoins(w http.ResponseWriter, r *http.Request) {
	var req CoinsRequest
	if err := chi.Bind(r, &req); err != nil {
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
	w.JSON(http.StatusOK, dtoResult)
}

func convertToDTO(coins []Coin) *CoinsResponse {
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
