package api

import "final_course/internal/cases"

// ServiceInterface описывает методы, которые должен реализовывать сервис
type ServiceInterface interface {
	GetCoins() ([]cases.Coin, error)
}
