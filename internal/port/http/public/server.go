package public

import (
	"encoding/json"
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
	service Service  // Используем интерфейс сервиса //TO DO: здесь интерфейс Service
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

//TODO: здесь 4 метода всё, -это методы Server-в каждый метод объединить всё что ниже расписано-метод в итоге отдаёт json; всё расписано-определить порядок

// getMaxHandler реализация обработчиков с конвертацией в DTO
//func (s *Server) getMaxHandler(ctx context.Context, titles []string) (dto.CoinDTOList, error) {
//	coins, err := s.service.GetMaxRate(ctx, titles)
//	return convertToDTO(coins), err
//}

//// Обработчики запросов
//func (s *Server) handleGetMax(rw http.ResponseWriter, req *http.Request) { //TODO: сигнатура правильная
//	s.handleRequest(rw, req, s.service.GetMaxRate)
//}

//TODO: rest api

func (s *Server) GetMax(rw http.ResponseWriter, req *http.Request) {
	//TODO: все ошибки через http.Erorr и return
	titles := req.URL.Query()["title"]
	if len(titles) == 0 {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		//json.NewEncoder(rw).Encode(map[string]string{"error": "missing 'title' parameter"})
		return
	}

	ctx := req.Context()
	coins, err := s.service.GetMaxRate(ctx, titles) //TODO:ctx:=req.Context() прописать в других методах
	if err != nil {
		var statusCode int
		errorMessage := "internal server error"

		if errors.Is(err, entities.ErrorInvalidParams) {
			statusCode = http.StatusBadRequest
			errorMessage = err.Error()
		} else {
			statusCode = http.StatusBadRequest // Изменено с StatusInternalServerError
		}

		rw.Header().Set("Content-Type", "application/json")
		//rw.WriteHeader(statusCode)
		//json.NewEncoder(rw).Encode(map[string]string{"error": errorMessage})
		return
	}

	if len(coins) == 0 {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(map[string]string{"error": "no data found"})
		return
	}

	result := make(dto.CoinDTOList, len(coins))
	for i, coin := range coins {
		result[i] = dto.CoinDTO{
			Title: coin.Title,
			Rate:  coin.Rate,
			Date:  coin.Date,
		}
	}

	rw.Header().Set("Content-Type", "application/json") //TODO: вместо Set Add
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(result); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
