package public

import (
	"context"
	"encoding/json"
	"final_course/internal/cases"
	"final_course/internal/entities"
	"final_course/pkg/dto"
	"github.com/pkg/errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server структура реализующая интерфейс ServerInterface.
// Содержит поля для сервиса и роутера.
type Server struct {
	service cases.ServiceInterface // Используем интерфейс сервиса //TO DO: здесь интерфейс Service
	router  *chi.Mux               // Роутер для обработки HTTP-запросов
}

// NewServer конструктор
func NewServer(service cases.ServiceInterface) *Server {
	return &Server{
		service: service,
		router:  chi.NewRouter(),
	}
}

// Run реализация метода  ServerInterface
func (s *Server) Run(addr string) error {
	s.router.Get("/v1/max", s.handleGetMax)
	s.router.Get("/v1/min", s.handleGetMin)
	s.router.Get("/v1/avg", s.handleGetAverage)
	s.router.Get("/v1/last", s.handleGetLastRate)

	log.Printf("Server started on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

// getMaxHandler реализация обработчиков с конвертацией в DTO
func (s *Server) getMaxHandler(ctx context.Context, titles []string) (dto.CoinDTOList, error) {
	coins, err := s.service.GetMaxRate(ctx, titles)
	return convertToDTO(coins), err
}

// Обработчики запросов
func (s *Server) handleGetMax(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r, s.service.GetMaxRate)
}

func (s *Server) handleGetMin(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r, s.service.GetMinRate)
}

func (s *Server) handleGetAverage(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r, s.service.GetAvgRate)
}

func (s *Server) handleGetLastRate(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r, s.service.GetLastRates)
}

// Общий обработчик для всех запросов
func (s *Server) handleRequest(
	w http.ResponseWriter,
	r *http.Request,
	handler func(context.Context, []string) ([]entities.Coin, error),
) {
	titles := r.URL.Query()["title"]
	if len(titles) == 0 {
		http.Error(w, "missing 'title' parameter", http.StatusBadRequest)
		return
	}

	coins, err := handler(r.Context(), titles)
	if err != nil {
		s.handleError(w, err)
		return
	}

	if len(coins) == 0 {
		http.Error(w, "no data found", http.StatusNotFound)
		return
	}

	s.sendJSON(w, convertToDTO(coins), http.StatusOK)
}

// Обработка ошибок
func (s *Server) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entities.ErrorInvalidParams):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", http.StatusBadRequest)
	}
}

// Отправка JSON ответа
func (s *Server) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// Конвертация в DTO
func convertToDTO(coins []entities.Coin) dto.CoinDTOList {
	result := make(dto.CoinDTOList, len(coins))
	for i, coin := range coins {
		result[i] = dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
			Date:  coin.Date,
		}
	}
	return result
}
