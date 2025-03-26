package public

import (
	"encoding/json"
	"final_course/internal/entities"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Server структура реализующая интерфейс ServerInterface.
// Содержит поля для сервиса и роутера.
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
	//s.router.Get("/v1/min", s.GetMin)
	//s.router.Get("/v1/avg", s.GetAverage)
	//s.router.Get("/v1/last", s.GetLastRate)

	log.Printf("Server starting on :8080")
	http.ListenAndServe(":8080", s.router)
}

//TODO: здесь 4 метода всё, -это методы Server-в каждый метод объединить всё что ниже расписано-
// метод в итоге отдаёт json; всё расписано-определить порядок

//-----------------------------------------------------------------------------
//ctx := req.Context()
//coins, err := s.service.GetMaxRate(ctx, titles) //TODO:ctx:=req.Context() прописать в других методах
//TODO: все ошибки через http.Erorr и return

// GetMax обрабатывает GET-запросы к эндпоинту /v1/max, возвращая максимальные ставки для указанных монет.
func (s *Server) GetMax(rw http.ResponseWriter, req *http.Request) {
	// Извлекаем параметр titles из строки запроса
	// Это список названий монет, разделенных запятыми, например, "bitcoin,ethereum"
	titlesStr := req.URL.Query().Get("titles")
	if titlesStr == "" {
		// Если параметр отсутствует, возвращаем ошибку 400 (Bad Request)
		http.Error(rw, "missing 'titles' parameter", http.StatusBadRequest)
		return
	}

	// Разделяем строку titles на список и обрезаем пробелы
	// Например, "bitcoin, ethereum" преобразуется в ["bitcoin", "ethereum"]
	titles := strings.Split(titlesStr, ",")
	var validTitles []string
	for _, title := range titles {
		trimmed := strings.TrimSpace(title)
		if trimmed != "" {
			validTitles = append(validTitles, trimmed)
		}
	}

	// Проверяем, остались ли валидные названия после обработки
	// Если нет, возвращаем ошибку 400 (Bad Request)
	if len(validTitles) == 0 {
		http.Error(rw, "No valid titles provided.", http.StatusBadRequest)
		return
	}

	// Вызываем сервис для получения максимальных ставок для указанных монет
	// Сервис возвращает список entities.Coin с максимальными ставками
	ctx := req.Context()
	response, err := s.service.GetMaxRate(ctx, validTitles)
	if err != nil {
		// Если произошла ошибка в сервисе, возвращаем ошибку 500 (Internal Server Error)
		// Сообщение об ошибке передается клиенту
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type для JSON-ответа
	//rw.Header().Set("Content-Type", "application/json")
	rw.Header().Add("Content-Type", "application/json")

	// Кодируем и отправляем список монет в формате JSON
	// Предполагается, что entities.Coin имеет те же поля, что и dto.CoinDT,
	// поэтому прямое кодирование возможно
	json.NewEncoder(rw).Encode(response)
}

//-----------------------------------------------------------------------------
