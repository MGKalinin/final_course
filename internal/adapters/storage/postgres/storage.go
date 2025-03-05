//это postgress/storage.go

package storage

import (
	"context"
	"final_course/internal/cases"
	"final_course/internal/entities"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"time"
)

// Storage структура реализующая интерфейс Storage
type Storage struct {
	db *pgxpool.Pool
}

// NewStorage конструктор для создания нового хранилища
func NewStorage(ctx context.Context, connString string) (*Storage, error) {
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		//TO DO : обработать ошибку, заврапать
		return nil, errors.Wrap(err, "failed to connect to the database")

	}
	return &Storage{
		db: pool,
	}, nil
}

// Store метод сохраняет монеты в бд
func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	//TO DO: prices заменить на coin_base
	query := `INSERT INTO coin_base (title, rate, date)  
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
	case cases.Max:
		query = `SELECT DISTINCT ON (title) title, rate, date FROM (SELECT title, rate, date FROM coin_base WHERE title = ANY($1) ORDER BY title, rate DESC, date DESC) AS subquery`
	case cases.Min:
		query = `SELECT DISTINCT ON (title) title, rate, date FROM (SELECT title, rate, date FROM coin_base WHERE title = ANY($1) ORDER BY title, rate ASC, date ASC) AS subquery`
	case cases.Avg:
		query = `SELECT title, AVG(rate) as rate, MAX(date) as date FROM coin_base WHERE title = ANY($1) GROUP BY title`
	default:
		query = `SELECT DISTINCT ON (title) title, rate, date FROM coin_base WHERE title = ANY($1) ORDER BY title, rate, date`
	}

	rows, err := s.db.Query(ctx, query, titles)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var coins []entities.Coin
	for rows.Next() {
		var coin entities.Coin
		if options.FuncType == cases.Avg {
			if err := rows.Scan(&coin.Title, &coin.Rate, &coin.Date); err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}
		} else {
			if err := rows.Scan(&coin.Title, &coin.Rate); err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}
			// Установите Date в нулевое значение, так как оно не используется
			coin.Date = time.Time{}
		}
		coins = append(coins, coin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return coins, nil
}
