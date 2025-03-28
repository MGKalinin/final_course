package public

import (
	"encoding/json"
	"final_course/internal/entities"
	"final_course/pkg/dto"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Server структура реализующая интерфейс ServerInterface.
type Server struct {
	service Service  // Используем интерфейс сервиса
	router  *chi.Mux // Роутер для обработки HTTP-запросов
}

// NewServer конструктор
func NewServer(service Service) (*Server, error) {
	if service == nil || service == Service(nil) {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "service is nil")
	}
	r := chi.NewRouter()
	return &Server{
		service: service,
		router:  r,
	}, nil
}

// Run реализация метода  ServerInterface
func (s *Server) Run() {
	s.router.Get("/v1/max", s.GetMax)
	s.router.Get("/v1/min", s.GetMin)
	s.router.Get("/v1/avg", s.GetAverage)
	s.router.Get("/v1/last", s.GetLastRate)

	log.Printf("Server starting on :8080")
	http.ListenAndServe(":8080", s.router)
}

// GetMax обрабатывает GET-запросы к эндпоинту /v1/max, возвращая максимальные ставки для указанных монет.
func (s *Server) GetMax(rw http.ResponseWriter, req *http.Request) {
	// Извлекаем параметр titles из строки запроса
	titlesStr := req.URL.Query().Get("titles")
	if titlesStr == "" {
		// Если параметр отсутствует, возвращаем ошибку 400 (Bad Request)
		http.Error(rw, "missing 'titles' parameter", http.StatusBadRequest)
		return
	}

	titles := strings.Split(titlesStr, ",")
	ctx := req.Context()
	coins, err := s.service.GetMaxRate(ctx, titles)
	if err != nil {
		// Если произошла ошибка в сервисе, возвращаем ошибку 400 (Bad Request)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if coins == nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Преобразуем в CoinDTOList для JSON
	var dtoList dto.CoinDTOList
	for _, coin := range coins {
		dtoList = append(dtoList, dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
			Date:  coin.Date, // TODO: может быть ебота со временем -см на тестах
		})
	}

	// Устанавливаем заголовок и кодируем ответ
	rw.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(dtoList)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

//TODO:  документация swagger; конфиг; миграция бд;
//доделать остальные три метода

// GetMin обрабатывает GET-запросы к эндпоинту /v1/min
func (s *Server) GetMin(rw http.ResponseWriter, req *http.Request) {
	titlesStr := req.URL.Query().Get("titles")
	if titlesStr == "" {
		http.Error(rw, "missing 'titles' parameter", http.StatusBadRequest)
		return
	}
	titles := strings.Split(titlesStr, ",")
	ctx := req.Context()
	coins, err := s.service.GetMinRate(ctx, titles)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if coins == nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	var dtoList dto.CoinDTOList
	for _, coin := range coins {
		dtoList = append(dtoList, dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
			Date:  coin.Date,
		})
	}
	rw.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(dtoList)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// GetAverage обрабатывает GET-запросы к эндпоинту /v1/avg
func (s *Server) GetAverage(rw http.ResponseWriter, req *http.Request) {
	titlesStr := req.URL.Query().Get("titles")
	if titlesStr == "" {
		http.Error(rw, "missing 'titles' parameter", http.StatusBadRequest)
		return
	}
	titles := strings.Split(titlesStr, ",")
	ctx := req.Context()
	coins, err := s.service.GetAvgRate(ctx, titles)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if coins == nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	var dtoList dto.CoinDTOList
	for _, coin := range coins {
		dtoList = append(dtoList, dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
			Date:  coin.Date,
		})
	}
	rw.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(dtoList)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// GetLastRate обрабатывает GET-запросы к эндпоинту /v1/last
func (s *Server) GetLastRate(rw http.ResponseWriter, req *http.Request) {
	titlesStr := req.URL.Query().Get("titles")
	if titlesStr == "" {
		http.Error(rw, "missing 'titles' parameter", http.StatusBadRequest)
		return
	}
	titles := strings.Split(titlesStr, ",")
	ctx := req.Context()
	coins, err := s.service.GetLastRates(ctx, titles)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if coins == nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	var dtoList dto.CoinDTOList
	for _, coin := range coins {
		dtoList = append(dtoList, dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
			Date:  coin.Date,
		})
	}
	rw.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(dtoList)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
