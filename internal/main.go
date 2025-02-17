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

	// Определяем, какие монеты будем запрашивать
	coinsToFetch := []string{"BTC", "ETH", "XRP"}
	//coinsToFetch := []string{}

	// Инициализируем клиент с базовым URL
	client, err := cryptocompare.NewClient(httpClient, "https://min-api.cryptocompare.com", coinsToFetch)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	// Делаем запрос для получения информации о криптовалютах
	coins, err := client.Get(ctx, nil) // nil для использования монет по умолчанию
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	// Выводим результат
	for _, coin := range coins {
		fmt.Printf("Title: %s, Rate: %f, Date: %s\n", coin.Title, coin.Rate, coin.Date)
	}
}

//package main
//
//import (
//"database/sql"
//"final_course/internal/cryptocompare"
//"log"
//)
//
//func main() {
//	db, err := sql.Open("postgres", "your-connection-string")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
//	storage := cryptocompare.NewStorage(cryptocompare.WithDB(db))
//
//	// Используйте storage для выполнения операций
//}
