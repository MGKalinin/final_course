//это postgress/storage.go

package storage

import (
	"context"
	"final_course/internal/cases"
	"final_course/internal/entities"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
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
	query := `INSERT INTO prices (title, rate, date)  
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

// Get метод извлекает монеты из бд
func (s *Storage) Get(ctx context.Context, titles []string, opts ...cases.Option) ([]entities.Coin, error) {
	options := &cases.Options{}
	for _, opt := range opts {
		opt(options)
	}

	var query string
	switch options.FuncType {
	case cases.Max:
		query = `SELECT title, MAX(rate) as rate, MAX(date) as date FROM prices WHERE title = ANY($1) GROUP BY title`
	case cases.Min:
		query = `SELECT title, MIN(rate) as rate, MIN(date) as date FROM prices WHERE title = ANY($1) GROUP BY title`
	case cases.Avg:
		query = `SELECT title, AVG(rate) as rate, MAX(date) as date FROM prices WHERE title = ANY($1) GROUP BY title`
	default:
		query = `SELECT title, rate, date FROM prices WHERE title = ANY($1)`
	}

	rows, err := s.db.Query(ctx, query, titles)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var coins []entities.Coin
	for rows.Next() {
		var coin entities.Coin
		if err := rows.Scan(&coin.Title, &coin.Rate, &coin.Date); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		coins = append(coins, coin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return coins, nil
}
