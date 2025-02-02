package entities

import (
	"github.com/pkg/errors"
	"time"
)

// Coin представляет информацию о криптовалюте
// Title - название криптовалюты
// Rate - текущий курс
// Date - дата обновления

type Coin struct {
	Title string
	Rate  float64
	Date  time.Time
}

// конструктор

func NewCoin(title string, rate float64, date time.Time) (*Coin, error) {
	if title == "" {
		return nil, errors.Wrap(ErrorEntity, "empty name of the cryptocurrency")
	}
	if rate < 0 {
		return nil, errors.Wrap(ErrorEntity, "negative cryptocurrency exchange rate")
	}
	return &Coin{Title: title, Rate: rate, Date: date}, nil
}
