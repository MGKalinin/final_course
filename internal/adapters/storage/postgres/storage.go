//это postgress/storage.go

package cryptocompare

import (
	"context"
	"database/sql"
	"final_course/internal/entities"
	"fmt"
	"strings"
)

//TODO: здесь метод get для storage? а что делает storage?....fmt print opts
//TODO:-получение опции и подстановка в метод get-разобраться с добавлением опциональных аргументов-подставляю опцию на уровне сервиса
//TODO:-get метод достаёт из базы данных
//TO DO: 1 разобраться с опциональными аргументами; 2 подправить client ; 3 написать любой метод get на уровне storage

// StorageOption определяет функциональную опцию для Storage.
type StorageOption func(*Storage)

// Storage представляет хранилище данных.
type Storage struct {
	db *sql.DB
}

// NewStorage создаёт новый экземпляр Storage с заданными опциями.
func NewStorage(opts ...StorageOption) *Storage {
	storage := &Storage{}
	for _, opt := range opts {
		opt(storage)
	}
	return storage
}

// WithDB устанавливает базу данных для Storage.
func WithDB(db *sql.DB) StorageOption {
	return func(s *Storage) {
		s.db = db
	}
}

// Store сохраняет монеты в базе данных.
func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	// Реализация метода Store
	// TODO: Store кладёт в бд то что притащил client
	return nil
}

// Get извлекает монеты из базы данных по заданным заголовкам и опциям.
func (s *Storage) Get(ctx context.Context, titles []string, opts ...interface{}) ([]entities.Coin, error) {
	// TODO: здесь метод get для storage? а что делает storage?....fmt print opts
	fmt.Println("Options:", opts)

	// TODO:-получение опции и подстановка в метод get
	var args []interface{}
	var whereClauses []string

	for i, title := range titles {
		whereClauses = append(whereClauses, fmt.Sprintf("title = $%d", i+1))
		args = append(args, title)
	}

	// TODO:-разобраться с добавлением опциональных аргументов
	for i, opt := range opts {
		switch v := opt.(type) {
		case string:
			whereClauses = append(whereClauses, fmt.Sprintf("additional_field = $%d", len(args)+1))
			args = append(args, v)
		case int:
			whereClauses = append(whereClauses, fmt.Sprintf("another_field = $%d", len(args)+1))
			args = append(args, v)
		default:
			return nil, fmt.Errorf("unsupported option type: %T", v)
		}
	}

	query := `SELECT id, title, price FROM coins WHERE ` + strings.Join(whereClauses, " OR ")
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []entities.Coin
	for rows.Next() {
		var coin entities.Coin
		if err := rows.Scan(&coin.ID, &coin.Title, &coin.Price); err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	return coins, nil
}
