package cases

import (
	"context"
	"final_course/internal/entities"
)

// Storage описывает интерфейс для работы с хранилищем данных

type Storage interface {
	Store(ctx context.Context, coins []entities.Coin) error                                 //TODO: кладёт в бд то что притащил client
	Get(ctx context.Context, titles []string, opts ...interface{}) ([]entities.Coin, error) //TODO:-получение опции и подстановка в метод get-разобраться с добавлением опциональных аргументов-подставляю опцию на уровне сервиса
} //TODO:-get метод достаёт из базы данных

//TODO: 1 разобраться с опциональными аргументами; 2 подправить client ; 3 написать любой метод get на уровне storage
