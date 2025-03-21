package cases

import (
	"context"
	"final_course/internal/entities"
	"github.com/pkg/errors"
)

// Добавляем проверку реализации интерфейса
var _ ServiceInterface = (*Service)(nil)

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

// FetchAndStoreCoins метод вызывает метод бд для получения всех titles, затем использует клиент для получения данных о монетах и сохраняет их в бд
func (s *Service) FetchAndStoreCoins(ctx context.Context) error {
	// Получаем все titles из базы данных
	titles, err := s.storage.GetAllTitles(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get titles from database")
	}

	// Получаем данные о монетах с помощью клиента
	coins, err := s.client.Get(ctx, titles)
	if err != nil {
		return errors.Wrap(err, "failed to get coins from client")
	}

	// Сохраняем монеты в базе данных
	if err := s.storage.Store(ctx, coins); err != nil {
		return errors.Wrap(err, "failed to store coins in database")
	}

	return nil
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

// GetLastRates метод получения значений без опций
func (s *Service) GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coins, err := s.storage.Get(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "last value is missing")
	}
	return coins, nil
}
