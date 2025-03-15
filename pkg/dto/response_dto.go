package dto

import "time"

//type CoinDTO struct {
//Symbol string `json:"symbol"` //TO DO поля как у монет -здесь должен быть слайс структур coinsdto
//Price  float64 `json:"price"`
//}

type CoinDTO struct {
	Title string
	Rate  float64
	Date  time.Time
}

type CoinDTOList struct {
	Coins []CoinDTO
}
