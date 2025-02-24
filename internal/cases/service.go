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

// WithMaxFunc функция получения max значения
func WithMaxFunc() Option {
	return func(options *Options) {
		options.FuncType = Max
	}
}

// GetMaxRate метод получения max значения
func (s *Service) GetMaxRate(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coins, err := s.storage.Get(ctx, titles, WithMaxFunc())
	if err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "maximum value is missing")
	}
	return coins, nil
}

// WithMinFunc функция получения min значения
func WithMinFunc() Option {
	return func(options *Options) {
		options.FuncType = Min
	}
}

// GetMinRate метод получения min значения
func (s *Service) GetMinRate(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coins, err := s.storage.Get(ctx, titles, WithMinFunc())
	if err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "minimum value is missing")
	}
	return coins, nil
}

// WithAvgFunc функция получения avg значения
func WithAvgFunc() Option {
	return func(options *Options) {
		options.FuncType = Avg
	}
}

// GetAvgRate метод получения avg значения
func (s *Service) GetAvgRate(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coins, err := s.storage.Get(ctx, titles, WithAvgFunc())
	if err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "average value is missing")
	}
	return coins, nil
}

//TODO: реализовать метод вызывает функцию без опций- из title
