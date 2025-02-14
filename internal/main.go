package main

import (
	"context"
	"final_course/internal/adapters/externalclient/cryptocompare"
	"fmt"
	"log"
	"net/http"
)

// TO DO: написать метод get для client
// TODO: написать метод get для storage-надо чего-то прочитать....
func main() {
	// Создаем контекст
	ctx := context.Background()

	// Инициализируем HTTP клиент
	httpClient := &http.Client{}

	// Инициализируем клиент с базовым URL
	client, err := cryptocompare.NewClient(httpClient, "https://min-api.cryptocompare.com")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	// Делаем запрос для получения информации о криптовалютах BTC и ETH
	coins, err := client.Get(ctx, []string{"BTC", "ETH"}) //TODO: завести отдельную переменную-куда записать какие монеты парсить
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	// Выводим результат
	for _, coin := range coins {
		fmt.Printf("Title: %s, Rate: %f, Date: %s\n", coin.Title, coin.Rate, coin.Date)
	}
}
