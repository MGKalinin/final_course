package cases

import (
	"final_course/internal/variables"
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
		return nil, errors.Wrap(variables.SomeErr, "storage is nil")
	}
	if client == nil {
		return nil, errors.Wrap(variables.SomeErr, "client is nil")
	}
	return &Service{storage: storage, client: client}, nil
}
