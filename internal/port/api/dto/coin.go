package dto

type CoinRequest struct {
	Symbol string `json:"symbol"`
}

type CoinResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}
