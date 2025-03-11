package cases

import (
	"context"
	"final_course/internal/entities"
)

// Storage описывает интерфейс для работы с хранилищем данных
type Storage interface {
	Store(ctx context.Context, coins []entities.Coin) error
	Get(ctx context.Context, titles []string, opts ...Option) ([]entities.Coin, error)
	GetAllTitles(ctx context.Context) ([]string, error) //TO DO: добавил метод GetAllTitles в Storage interface
}
