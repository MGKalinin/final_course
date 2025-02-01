package entities

import (
	"errors"
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
	//прописать ошибки пустое поле, отрицательное значение
	if title == "" {
		return nil, errors.New("get errors") //обработка пакетом из github errors wrap
	}
	if rate < 0 {
		return nil, error
	}
	return &Coin{Title: title, Rate: rate, Date: date}, nil
}

// сделать враппинг
