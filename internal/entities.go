package internal

import "time"

// Coin представляет информацию о криптовалюте
// Title - название криптовалюты
// Rate - текущий курс
// Date - дата обновления

type Coin struct {
	Title string    `json:"title"`
	Rate  float64   `json:"rate"`
	Date  time.Time `json:"date"`
}
