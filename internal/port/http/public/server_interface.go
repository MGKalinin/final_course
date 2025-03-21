// Файл public/server_interface.go
package public

import (
	"context"
	"final_course/pkg/dto"
)

//TODO: здесь имплементация методов Service- это вынести в отдельный файл

// ServerInterface определяет контракт HTTP-сервера
type ServerInterface interface {
	Run(addr string) error

	// Внутренние методы для обработки запросов (опционально)
	getMaxHandler(ctx context.Context, titles []string) (dto.CoinDTOList, error)
	getMinHandler(ctx context.Context, titles []string) (dto.CoinDTOList, error)
	getAverageHandler(ctx context.Context, titles []string) (dto.CoinDTOList, error)
	getLastRateHandler(ctx context.Context, titles []string) (dto.CoinDTOList, error)
}
