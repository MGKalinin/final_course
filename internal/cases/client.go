package cases

import (
	"context"
	"final_course/internal/entities"
)

// здесь интерфейс клиента
// Client описывает интерфейс для получения данных о курсах валют

type Client interface {
	Get(ctx context.Context, titles []string) ([]entities.Coin, error)
}
