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
	db *pgxpool.Pool
}

// NewStorage конструктор для создания нового хранилища
func NewStorage(ctx context.Context, connString string) (*Storage, error) {
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		//TODO : обработать ошибку, заврапать
	}
	return &Storage{
		db: pool,
	}, nil
}

// Store метод сохраняет монеты в бд
func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	//TODO: prices заменить на coin_base
	query := `INSERT INTO prices (title, rate, date)  
		      VALUES ($1, $2, $3)`

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
	case cases.Max: //TODO: заменить на select distict
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
	defer rows.Close() //TODO: см что прилетит сюда

	var coins []entities.Coin
	for rows.Next() {
		var coin entities.Coin
		if err := rows.Scan(&coin.Title, &coin.Rate, &coin.Date); err != nil { //TODO: если нет агрегирующей функции то Date не нужна
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		coins = append(coins, coin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return coins, nil
}
