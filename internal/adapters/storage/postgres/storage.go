//это postgress/storage.go

package storage

import (
	"context"
	"database/sql"
	"final_course/internal/cases"
	"final_course/internal/entities"
	"fmt"
	"strings"
)

// Storage структура реализующая интерфейс Storage
type Storage struct {
	titles []string
	// db     *sql.DB
}

// Store метод сохраняет монеты в бд
func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	// Реализация метода Store
	return nil
}

// Get метод извлекает монеты из базы данных
func (s *Storage) Get(ctx context.Context, titles []string, opts ...cases.Option) ([]entities.Coin, error) {

	return coins, nil
}
