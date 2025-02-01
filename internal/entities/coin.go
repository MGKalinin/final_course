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
		return nil, errors.Wrap(errors.New("Пустое название криптовалюты"), "NewCoin failed")
	}
	if rate < 0 {
		return nil, errors.Wrap(errors.New("Отрицательный курс криптовалюты"), "NewCoin failed")
	}
	return &Coin{Title: title, Rate: rate, Date: date}, nil
}

// сделать враппинг
