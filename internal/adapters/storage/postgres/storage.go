package storage

import (
	"context"
	"final_course/internal/cases"
	"final_course/internal/entities"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"log"
)

// Storage структура реализующая интерфейс Storage
type Storage struct {
	db *pgxpool.Pool
}

// NewStorage конструктор для создания нового хранилища
func NewStorage(ctx context.Context, connString string) (*Storage, error) {
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to the database")
	}
	return &Storage{
		db: pool,
	}, nil
}

// GetAllTitles метод извлекает все уникальные titles из бд
func (s *Storage) GetAllTitles(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT title FROM coin_base`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "query was not found")
	}
	defer rows.Close()

	var titles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		titles = append(titles, title)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return titles, nil
}

// Store метод сохраняет монеты в бд
func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := `INSERT INTO coin_base (title, rate, date)
		      VALUES ($1, $2, $3)`

	// Выполняем запрос для каждой монеты
	for _, coin := range coins {
		log.Printf("Storing coin: %s with rate: %f on date: %s", coin.Title, coin.Rate, coin.Date)
		_, err = tx.Exec(ctx, query, coin.Title, coin.Rate, coin.Date)
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
			if err := rows.Scan(&coin.Title, &coin.Rate, &coin.Date); err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}
		}
		coins = append(coins, coin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return coins, nil
}
