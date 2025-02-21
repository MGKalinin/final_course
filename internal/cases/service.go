package cases

import (
	"context"
	cryptocompare "final_course/internal/adapters/storage/postgres"
	"final_course/internal/entities"
	"github.com/pkg/errors"
)

// Service содержит логику работы сервиса
// storage Storage
type Service struct {
	storage cryptocompare.Storage
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

// TO DO: здесь опциональные аргументы -здесь подстановка-передача

// GetCoins извлекает монеты с использованием слоя хранилища с опциями.
func (s *Service) GetCoins(ctx context.Context, titles []string, opts ...cryptocompare.Options) ([]entities.Coin, error) {
	return s.storage.Get(ctx, titles, opts...)
}
