package apiserver

import "time"

// CoinDTO struct DTO для монет
type CoinDTO struct {
	Title string    `json:"title"`
	Rate  float64   `json:"rate"`
	Date  time.Time `json:"date"`
}

type CoinsRequest struct {
	Titles []string `json:"titles"`
}

type CoinsResponse struct {
	Coins []CoinDTO `json:"coins"`
}
