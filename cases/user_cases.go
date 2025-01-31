package cases

import (
	"context"
	"final_course/adapters"
	"time"
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

// UpdateRates обновляет курсы валют каждые 5 минут

func (s *Service) UpdateRates(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			titles := []string{"BTC", "ETH"}
			coins, err := s.client.Get(ctx, titles)
			if err != nil {
				// Логирование ошибки (например, log.Println)
				continue
			}
			s.storage.Store(ctx, coins)
		case <-ctx.Done():
			return
		}
	}
}
