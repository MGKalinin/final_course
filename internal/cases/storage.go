package cases

import (
	"context"
	"final_course/internal/entities"
)

// Storage описывает интерфейс для работы с хранилищем данных
// TODO: Store кладёт в бд то что притащил  client
type Storage interface {
	Store(ctx context.Context, coins []entities.Coin) error
	Get(ctx context.Context, titles []string, opts ...interface{}) ([]entities.Coin, error)
}
