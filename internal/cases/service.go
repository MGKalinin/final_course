package cases

import (
	"final_course/internal/entities"
	"github.com/pkg/errors"
)

// Service содержит логику работы сервиса

type Service struct {
	storage Storage // проверить на nil
	client  Client
}

// NewService конструктор - создает новый сервис
// проверка на nil
func NewService(storage Storage, client Client) (*Service, error) {
	if storage == nil {
		return nil, errors.Wrap(entities.SomeErr, "storage is nil")
	}
	if client == nil {
		return nil, errors.Wrap(entities.SomeErr, "client is nil")
	}
	return &Service{storage: storage, client: client}, nil
}
