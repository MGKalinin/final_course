package apiserver

import "context"

// RESTPort определяет интерфейс REST API
type RESTPort interface {
	GetCoins(ctx context.Context, titles []string) (*CoinsResponse, error)
	GetMaxRate(ctx context.Context, titles []string) (*CoinsResponse, error)
	GetMinRate(ctx context.Context, titles []string) (*CoinsResponse, error)
	GetAvgRate(ctx context.Context, titles []string) (*CoinsResponse, error)
}
