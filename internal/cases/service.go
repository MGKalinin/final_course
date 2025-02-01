package cases

import "fmt"

// Service содержит логику работы сервиса

type Service struct {
	storage Storage //проверить на nil
	client  Client
}

// NewService конструктор - создает новый сервис
// проверка на nil
func NewService(storage Storage, client Client) (*Service, error) {
	if storage == nil {
		return nil, fmt.Errorf("NewService failed: storage is nil")
	}
	if client == nil {
		return nil, fmt.Errorf("NewService failed: client is nil")
	}
	return &Service{storage: storage, client: client}, nil
}
