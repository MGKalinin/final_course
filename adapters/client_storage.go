package adapters

import (
	"context"
	"final_course/internal"
)

// Storage описывает интерфейс для работы с хранилищем данных

type Storage interface {
	Store(ctx context.Context, coins []internal.Coin) error
	Get(ctx context.Context, titles []string, opts ...interface{}) ([]internal.Coin, error)
}

// Client описывает интерфейс для получения данных о курсах валют

type Client interface {
	Get(ctx context.Context, titles []string) ([]internal.Coin, error)
}
