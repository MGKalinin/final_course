package entities

import (
	"github.com/pkg/errors"
	"time"
)

//  Аннотации моделей

// Coin информация о криптовалюте
// @Schema
type Coin struct {
	// Название криптовалюты
	// @example BTC
	Title string

	// Текущий курс
	// @example 50000.0
	Rate float64

	// Дата обновления
	// @example 2023-10-01T12:00:00Z
	Date time.Time
}

// конструктор

func NewCoin(title string, rate float64, date time.Time) (*Coin, error) {
	if title == "" {
		return nil, errors.Wrap(ErrorInvalidParams, "empty name of the cryptocurrency")
	}
	if rate < 0 {
		return nil, errors.Wrap(ErrorInvalidParams, "negative cryptocurrency exchange rate")
	}
	return &Coin{Title: title, Rate: rate, Date: date}, nil
}
