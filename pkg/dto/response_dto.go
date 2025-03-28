package dto

import "time"

// DTO аннотации

// CoinDTO транспортная модель для Coin
// @Schema
type CoinDTO struct {
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

// CoinDTOList список CoinDTO
// @Schema
type CoinDTOList []CoinDTO
