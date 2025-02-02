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

// создать переменную ошибки- и далее использовать её в обёртке во всех случаях,
// добавляя соответсвующий комментарий

var SomeErr = errors.New("Ошибка:")

// конструктор

func NewCoin(title string, rate float64, date time.Time) (*Coin, error) {
	if title == "" {
		return nil, errors.Wrap(SomeErr, "пустое название криптовалюты")
	}
	if rate < 0 {
		return nil, errors.Wrap(SomeErr, "отрицательный курс криптовалюты")
	}
	return &Coin{Title: title, Rate: rate, Date: date}, nil
}
