// Файл cases/service.go
package public

import (
	"context"
	"final_course/internal/entities"
)

// Service определяет контракт бизнес-логики
type Service interface {
	FetchAndStoreCoins(ctx context.Context) error
	GetMaxRate(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMinRate(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAvgRate(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error)
}
