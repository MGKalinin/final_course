package cases

import (
	"context"
	"final_course/internal/entities"
	"github.com/pkg/errors"
)

// Service содержит логику работы сервиса
type Service struct {
	storage Storage
	client  Client
}

// NewService конструктор - создает новый сервис
// проверка на nil
func NewService(storage Storage, client Client) (*Service, error) {
	if storage == nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "storage is nil")
	}
	if client == nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "client is nil")
	}
	return &Service{storage: storage, client: client}, nil
}

type AggFunc int

const (
	_ AggFunc = iota
	Max
	Min
	Avg
)

type Options struct {
	FuncType AggFunc
}
type Option func(options *Options)

func (a AggFunc) String() string {
	return [...]string{"", "MAX", "MIN", "AVG"}[a]
}

func WithMaxFunc() Option {
	return func(options *Options) {
		options.FuncType = Max
	}
}

func (s *Service) GetMaxRate(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coins, err := s.storage.Get(ctx, titles, WithMaxFunc())
	if err != nil {
		return nil, err //TODO заврапать ошибку
	}
	return coins, nil
}

//TODO: дописать мин,сред....  TODO: нужен pgx 4 версия и pgx pull ; установить бд; нужна история установления соединения
