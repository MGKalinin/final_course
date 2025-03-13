package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"final_course/internal/adapters/externalclient/cryptocompare"
	"final_course/internal/adapters/storage/postgres"
	"final_course/internal/cases"
	"final_course/internal/port/api"
	"github.com/go-chi/chi/v5"
)

// TODO: создать дир порт, в ней ещё дир(см как реализовается рест апи)-создать интерфейс для связи между портом(точка входа в программу, моё приложение) и самим сервисом-мне нужен будет роутер из пакета chi.mux-(методы, структура) -на этом урове не использовать сущности коин-нужно data transfer object-ошибки из пакета http(ok/bad reqvest)
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

	// Создание обработчика REST API
	handler := api.NewAPIHandler(service)

	// Создание роутера
	r := chi.NewRouter()

	// Настройка маршрутов
	handler.SetupRoutes(r)

	// Запуск сервера
	addr := ":3000"
	log.Printf("Сервер запущен на порту %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
