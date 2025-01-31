package cases

import (
	"final_course/adapters"
)

// Service содержит логику работы сервиса

type Service struct {
	storage adapters.Storage
	client  adapters.Client
}

// NewService создает новый сервис

func NewService(storage adapters.Storage, client adapters.Client) *Service {
	return &Service{storage: storage, client: client}
}
