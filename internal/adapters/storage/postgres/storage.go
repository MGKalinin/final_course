//это postgress/storage.go

package storage

import (
	"context"
	"final_course/internal/cases"
	"final_course/internal/entities"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

// Storage структура реализующая интерфейс Storage
type Storage struct {
	titles []string
	db     *pgxpool.Pool
}

// NewStorage конструктор для создания нового хранилища
func NewStorage(db *pgxpool.Pool, titles []string) (*Storage, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection pool cannot be nil")
	}
	return &Storage{
		titles: titles,
		db:     db,
	}, nil
}

// Store метод сохраняет монеты в бд
func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	//TODO:положить запрос в бд-посмотреть результат
	query := `INSERT INTO coins (title, rate, date)  
		      VALUES ($1, $2, $3)
		      ON CONFLICT (title) DO UPDATE
		      SET rate = EXCLUDED.rate, date = EXCLUDED.date;`

	// Выполняем запрос для каждой монеты
	for _, coin := range coins {
		_, err := s.db.Exec(ctx, query, coin.Title, coin.Rate, coin.Date)
		if err != nil {
			return fmt.Errorf("failed to store coin %s: %w", coin.Title, err)
		}
	}
	return nil
}

// Get метод извлекает монеты из базы данных
func (s *Storage) Get(ctx context.Context, titles []string, opts ...cases.Option) ([]entities.Coin, error) {

	return coins, nil
}
