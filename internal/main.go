package main

import (
	"context"
	"final_course/internal/adapters/externalclient/cryptocompare"
	"final_course/internal/adapters/storage/postgres"
	"final_course/internal/cases"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Установите переменную окружения DATABASE_URL
	os.Setenv("DATABASE_URL", "postgres://maksimkalinin:password@localhost:5432/postgres")

	// Создаем контекст
	ctx := context.Background()

	// Определяем, какие монеты будем запрашивать
	coinsToFetch := []string{"BTC", "ETH", "XRP"}

	// Инициализируем клиент с базовым URL
	client, err := cryptocompare.NewClient("https://min-api.cryptocompare.com/data/pricemulti", coinsToFetch)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Инициализируем хранилище
	storage, err := storage.NewStorage(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error creating storage: %v", err)
	}

	// Инициализируем сервис
	service, err := cases.NewService(storage, client)
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	// Вызываем метод для получения и сохранения данных о монетах
	if err := service.FetchAndStoreCoins(ctx); err != nil {
		log.Fatalf("Error fetching and storing coins: %v", err)
	}

	// Делаем запрос для получения информации о криптовалютах
	coins, err := service.GetLastRates(ctx, coinsToFetch)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	// Устанавливаем текущую дату для новых записей
	currentDate := time.Now()
	for i := range coins {
		coins[i].Date = currentDate
	}

	// Сохраняем данные в базе данных
	if err := storage.Store(ctx, coins); err != nil {
		log.Fatalf("Error storing data: %v", err)
	}

	// Выводим результат
	fmt.Println("Fetched Coins:")
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
