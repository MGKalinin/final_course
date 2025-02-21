//это postgress/storage.go

package storage

import (
	"context"
	"database/sql"
	"final_course/internal/entities"
	"fmt"
	"strings"
)

// Store сохраняет монеты в базе данных.
func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	// Реализация метода Store
	return nil
}

// Get извлекает монеты из базы данных по заданным заголовкам и опциям.// TODO: opts ...interface{}-заменить на опции
func (s *Storage) Get(ctx context.Context, titles []string, opts ...interface{}) ([]entities.Coin, error) {
	// TO DO -реализовать сам паттерн-сюда притащить get max/ min/last/avergse -вычисление в бд(написать запросы в бд по соотвующим )-без поднятия бд-есть онлайн бд для проверки
	option := GetOption(opts...)

	query := `SELECT * FROM coins WHERE title = ANY($1)`
	args := []interface{}{titles}

	switch option.CalculationType {
	case "max":
		query += " ORDER BY value DESC LIMIT 1"
	case "min":
		query += " ORDER BY value ASC LIMIT 1"
	case "last":
		query += " ORDER BY timestamp DESC LIMIT 1"
	case "average":
		query = "SELECT AVG(value) FROM coins WHERE title = ANY($1)"
	default:
		// Поведение по умолчанию, если опция не указана
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []entities.Coin
	for rows.Next() {
		var coin entities.Coin
		if err := rows.Scan(&coin.ID, &coin.Title, &coin.Value, &coin.Timestamp); err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	return coins, nil
}
