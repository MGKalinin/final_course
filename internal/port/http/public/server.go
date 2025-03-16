package public

import (
	"context"
	"encoding/json"
	"final_course/internal/cases"
	"final_course/internal/entities"
	"final_course/pkg/dto"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//TO DO: под реализацию методов в сервисе/кейсес создать интерфейс методов макс/мин...структура сервера,
//конструктор-сервер нужно запустить-здесь его запустить и здесь его методы
//TO DO: ручки -пример на 24 строке -4 метода, 4 ручки

// Server структура реализующая интерфейс ServerInterface.
// Содержит поля для сервиса и роутера.
type Server struct {
	service *cases.Service // Сервис для выполнения бизнес-логики //TODO: здесь интерфейс Service
	router  *chi.Mux       // Роутер для обработки HTTP-запросов
}

// ServerInterface описывает интерфейс методов получения результатов, запрашиваемых пользователем.
// Включает методы для получения максимального, минимального, среднего значений и последнего курса.
type ServerInterface interface { //TODO: здесь имплементация методов Service- это вынести в отдельный файл
	GetMax(ctx context.Context, titles []string) (dto.CoinDTOList, error)
	GetMin(ctx context.Context, titles []string) (dto.CoinDTOList, error)
	GetAverage(ctx context.Context, titles []string) (dto.CoinDTOList, error)
	GetLastRate(ctx context.Context, titles []string) (dto.CoinDTOList, error)
	Run(addr string) error // Метод для запуска сервера
}

//TODO: любая функция /метод возвращают ошибку

// NewServer конструктор для создания нового сервера.
// Инициализирует сервер с заданным сервисом и роутером.
func NewServer(service *cases.Service) *Server {
	return &Server{
		service: service,
		router:  chi.NewRouter(),
	}
}

// GetMax реализация метода получения максимального значения.
// Вызывает метод сервиса для получения максимального значения и конвертирует результат в CoinDTOList.
//func (s *Server) GetMax(ctx context.Context, titles []string) (dto.CoinDTOList, error) {
//	coins, err := s.service.GetMaxRate(ctx, titles)
//	if err != nil {
//		return dto.CoinDTOList{}, err
//	}
//	return dto.CoinDTOList{Coins: convertToCoinDTO(coins)}, nil
//}

// GetMin реализация метода получения минимального значения.
// Вызывает метод сервиса для получения минимального значения и конвертирует результат в CoinDTOList.
func (s *Server) GetMin(ctx context.Context, titles []string) (dto.CoinDTOList, error) {
	coins, err := s.service.GetMinRate(ctx, titles)
	if err != nil {
		return dto.CoinDTOList{}, err
	}
	return dto.CoinDTOList{Coins: convertToCoinDTO(coins)}, nil
}

// GetAverage реализация метода получения среднего значения.
// Вызывает метод сервиса для получения среднего значения и конвертирует результат в CoinDTOList.
func (s *Server) GetAverage(ctx context.Context, titles []string) (dto.CoinDTOList, error) {
	coins, err := s.service.GetAvgRate(ctx, titles)
	if err != nil {
		return dto.CoinDTOList{}, err
	}
	return dto.CoinDTOList{Coins: convertToCoinDTO(coins)}, nil
}

// GetLastRate реализация метода получения последнего значения.
// Вызывает метод сервиса для получения последнего значения и конвертирует результат в CoinDTOList.
func (s *Server) GetLastRate(ctx context.Context, titles []string) (dto.CoinDTOList, error) {
	coins, err := s.service.GetLastRates(ctx, titles)
	if err != nil {
		return dto.CoinDTOList{}, err
	}
	return dto.CoinDTOList{Coins: convertToCoinDTO(coins)}, nil
}

// Run запуск сервера.
// Настраивает роутер и запускает HTTP-сервер на заданном адресе.
func (s *Server) Run() error { //TODO: вверху- старт с запуска сервера
	//s.router.Use(middleware.Logger)            // не нужен
	s.router.Get("/v1/getmaxrate", s.handleGetMax) // Регистрирует обработчик для получения максимального значения
	s.router.Get("/min", s.handleGetMin)           // Регистрирует обработчик для получения минимального значения
	s.router.Get("/avg", s.handleGetAverage)       // Регистрирует обработчик для получения среднего значения
	s.router.Get("/last", s.handleGetLastRate)     // Регистрирует обработчик для получения последнего значения
	log.Printf("Сервер запущен на порту %s", ":8080")
	return http.ListenAndServe(":8080", s.router) // Запускает сервер
}

// handleGetMax обработчик HTTP-запроса для получения максимального значения.
// Извлекает параметры запроса, вызывает метод GetMax и возвращает результат в формате JSON.
// TODO: которые не handle методы удалить кроме Run
func (s *Server) handleGetMax(rw http.ResponseWriter, req *http.Request) {
	titles := r.URL.Query()["title"]
	coins, err := s.GetMax(r.Context(), titles) //TODO: GetMax -это метод сервиса
	if err != nil {                             //TODO: здесь ещё проверить coin на nil
		http.Error(w, err.Error(), http.StatusInternalServerError) //TODO: здесь badrequest

		return
	}
	//TODO: здесь слайс dto-после проверки на наличие значений
	w.Header().Set("Content-Type", "application/json") //TODO: Header сначала добавить -потом обработать ошибку-потом write метод Header
	json.NewEncoder(w).Encode(coins)
}

// handleGetMin обработчик HTTP-запроса для получения минимального значения.
// Извлекает параметры запроса, вызывает метод GetMin и возвращает результат в формате JSON.
func (s *Server) handleGetMin(w http.ResponseWriter, r *http.Request) {
	titles := r.URL.Query()["title"]
	coins, err := s.GetMin(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coins)
}

// handleGetAverage обработчик HTTP-запроса для получения среднего значения.
// Извлекает параметры запроса, вызывает метод GetAverage и возвращает результат в формате JSON.
func (s *Server) handleGetAverage(w http.ResponseWriter, r *http.Request) {
	titles := r.URL.Query()["title"]
	coins, err := s.GetAverage(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coins)
}

// handleGetLastRate обработчик HTTP-запроса для получения последнего значения.
// Извлекает параметры запроса, вызывает метод GetLastRate и возвращает результат в формате JSON.
func (s *Server) handleGetLastRate(w http.ResponseWriter, r *http.Request) {
	titles := r.URL.Query()["title"]
	coins, err := s.GetLastRate(r.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coins)
}
