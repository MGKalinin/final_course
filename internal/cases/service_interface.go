// Файл cases/service_interface.go
package cases

import (
	"context"
	"final_course/internal/entities"
)

// ServiceInterface определяет контракт бизнес-логики
type ServiceInterface interface {
	FetchAndStoreCoins(ctx context.Context) error
	GetMaxRate(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMinRate(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAvgRate(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error)
}
