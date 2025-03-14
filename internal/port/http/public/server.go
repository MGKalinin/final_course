package public

import (
	"encoding/json"
	"final_course/pkg/dto"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CoinService interface {
	GetCoinPrice(symbol string) (float64, error)
}

//TODO: под реализацию методов в сервисе/кейсес создать интерфейс методов макс/мин...структура сервера, конструктор-сервер нужно запустить-здесь его запустить и здесь его методы
//TODO: ручки -пример на 24 строке -4 метода, 4 ручки

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

//TODO: см как сервис связан с портом-всё тоже самое
