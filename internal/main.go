package main

import (
	"context"
	"final_course/configs"
	_ "final_course/docs"
	"final_course/internal/adapters/externalclient/cryptocompare"
	"final_course/internal/adapters/storage/postgres"
	"final_course/internal/cases"
	"final_course/internal/port/http/public"
	"github.com/spf13/viper"
	"log"
	"os"
)

// Общее описание API

// @title Crypto Currency API
// @version 1.0
// @description API для получения курсов криптовалют

// @contact.name Максим Калинин
// @contact.email ваш@email.com

// @host localhost:8080
// @BasePath /v1
func main() {
	// LoadConfig грузим конфигурацию
	configs.LoadConfig()
	dbParams := viper.GetStringMapString("database")
	// Установка переменной окружения
	//os.Setenv("DATABASE_URL", "postgres://maksimkalinin:password@localhost:5432/postgres")
	os.Setenv("DATABASE_URL", dbParams["db_name"]+"://"+dbParams["username"]+":"+dbParams["password"]+
		"@localhost:"+dbParams["address"]+"/"+dbParams["db_name"])

	// Создание контекста
	ctx := context.Background()

	// Определение монет для запроса
	coinsToFetch := []string{"BTC", "ETH", "XRP"}

	// Инициализация клиента
	//client, err := cryptocompare.NewClient("https://min-api.cryptocompare.com/data/pricemulti", coinsToFetch)
	client, err := cryptocompare.NewClient(dbParams["client_address"], coinsToFetch)

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

	//Инициализация сервера
	server, err := public.NewServer(service)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Запуск сервера
	server.Run()
}
