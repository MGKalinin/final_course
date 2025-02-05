package main

import (
	"context"
	"final_course/internal/adapters"
	"fmt"
	"log"
)

func main() {
	// Создаем контекст
	ctx := context.Background()

	// Инициализируем клиент с API ключом и базовым URL
	apiKey := "1c5c331d210b7d08b2efe2c0741139b5063317d646e84f3619dd69a25d79f5a5"
	client := adapters.NewClient(apiKey, "https://api.cryptocompare.com")

	// Делаем запрос для получения информации о криптовалюте BTH
	coins, err := client.Get(ctx, []string{"BTH"})
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	// Выводим результат
	for _, coin := range coins {
		fmt.Printf("Title: %s, Rate: %f, Date: %s\n", coin.Title, coin.Rate, coin.Date)
	}
}
