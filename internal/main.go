package main

import (
	"context"
	"final_course/internal/adapters/externalclient/cryptocompare"
	storage "final_course/internal/adapters/storage/postgres"
	"final_course/internal/cases"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
)

func main() {
	// Установите переменную окружения DATABASE_URL
	os.Setenv("DATABASE_URL", "postgres://maksimkalinin:password@localhost:5432/postgres")

	// Создаем контекст
	ctx := context.Background()

	// Инициализируем HTTP клиент
	httpClient := &http.Client{}

	// Определяем, какие монеты будем запрашивать
	coinsToFetch := []string{"BTC", "ETH", "XRP"}
	//coinsToFetch := []string{}

	// Инициализируем клиент с базовым URL
	client, err := cryptocompare.NewClient(httpClient, "https://min-api.cryptocompare.com", coinsToFetch)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	// Инициализируем пул подключений к базе данных
	dbpool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Инициализируем хранилище
	storageInstance, err := storage.NewStorage(dbpool, coinsToFetch)
	if err != nil {
		log.Fatalf("Error creating storage: %v", err)
	}

	// Инициализируем сервис
	service, err := cases.NewService(storageInstance, client)
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	// Делаем запрос для получения информации о криптовалютах
	coins, err := client.Get(ctx, nil) // nil для использования монет по умолчанию
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	// Сохраняем данные в базу данных
	if err := storageInstance.Store(ctx, coins); err != nil {
		log.Fatalf("Error storing data: %v", err)
	}

	// Выводим результат
	for _, coin := range coins {
		fmt.Printf("Title: %s, Rate: %f, Date: %s\n", coin.Title, coin.Rate, coin.Date)
	}

	// Пример получения максимального значения
	maxCoins, err := service.GetMaxRate(ctx, coinsToFetch)
	if err != nil {
		log.Fatalf("Error getting max rate: %v", err)
	}
	fmt.Println("Max Rates:")
	for _, coin := range maxCoins {
		fmt.Printf("Title: %s, Rate: %f, Date: %s\n", coin.Title, coin.Rate, coin.Date)
	}
}
