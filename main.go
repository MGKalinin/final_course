package main

import (
	"context"
	"final_course/internal/adapters"
	"final_course/internal/cases"
	"final_course/internal/cases/"
	"fmt"
	"log"
)

func main() {
	apiKey := "1c5c331d210b7d08b2efe2c0741139b5063317d646e84f3619dd69a25d79f5a5"
	baseURL := "https://min-api.cryptocompare.com/data/pricemulti"

	client := adapters.NewAPIClient(apiKey, baseURL)

	ctx := context.Background()
	titles := []string{"BTC", "ETH"}

	coins, err := client.Get(ctx, titles)
	if err != nil {
		log.Fatalf("failed to get coins: %v", err)
	}

	for _, coin := range coins {
		fmt.Printf("Title: %s, Rate: %f\n", coin.Title, coin.Rate)
	}
}
