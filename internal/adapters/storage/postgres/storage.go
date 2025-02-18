//это postgress/storage.go

package cryptocompare

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

// Get извлекает монеты из базы данных по заданным заголовкам и опциям.
func (s *Storage) Get(ctx context.Context, titles []string, opts ...interface{}) ([]entities.Coin, error) {
	// TODO -реализовать сам паттерн-сюда притащить get max/ min/last/avergse -вычисление в бд(написать запросы в бд по соотвующим )-без поднятия бд-есть онлайн бд для проверки
}

//TODO придумать тип опция opts
