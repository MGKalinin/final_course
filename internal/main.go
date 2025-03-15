package main

import (
	"context"
	"final_course/internal/adapters/externalclient/cryptocompare"
	"final_course/internal/adapters/storage/postgres"
	"final_course/internal/cases"
	"final_course/internal/port/http/public"
	"log"
	"os"
)

// TO DO: ОБЩЕЕ: создать дир порт, в ней ещё дир(см как реализовается рест апи)-создать интерфейс
//для связи между портом(точка входа в программу, моё приложение) и самим сервисом-мне нужен будет
//роутер из пакета chi.mux-(методы, структура) -на этом уровне не использовать сущности коин-нужно
//data transfer object-ошибки из пакета http(ok/bad reqvest)

func main() {
	// Установка переменной окружения
	os.Setenv("DATABASE_URL", "postgres://maksimkalinin:password@localhost:5432/postgres")

	// Создание контекста
	ctx := context.Background()

	// Определение монет для запроса
	coinsToFetch := []string{"BTC", "ETH", "XRP"}

	// Инициализация клиента
	client, err := cryptocompare.NewClient("https://min-api.cryptocompare.com/data/pricemulti", coinsToFetch)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Инициализация хранилища
	storage, err := storage.NewStorage(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error creating storage: %v", err)
	}

	// Инициализация сервиса
	service, err := cases.NewService(storage, client)
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	// Создание и запуск сервера
	server := public.NewServer(service)
	if err := server.Run(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
