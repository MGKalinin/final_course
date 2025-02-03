package adapters

import (
	"context"
	"encoding/json"
	"final_course/internal/entities"
	"fmt"
	"net/http"
	"strings"
)

// APIClient реализует интерфейс Client для получения данных о курсах валют
type APIClient struct {
	apiKey  string
	baseURL string
}

// NewAPIClient создает новый экземпляр APIClient
func NewAPIClient(apiKey string, baseURL string) *APIClient {
	return &APIClient{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// Get реализует метод интерфейса Client для получения данных о курсах валют
func (c *APIClient) Get(ctx context.Context, titles []string) ([]entities.Coin, error) {
	url := fmt.Sprintf("%s?titles=%s&apiKey=%s", c.baseURL, joinTitles(titles), c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var coins []entities.Coin
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return nil, err
	}
	fmt.Println(coins)
	return coins, nil
}

// joinTitles объединяет список названий валют в строку, разделенную запятыми
func joinTitles(titles []string) string {
	return strings.Join(titles, ",")
}
