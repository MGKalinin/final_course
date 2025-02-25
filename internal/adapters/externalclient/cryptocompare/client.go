package cryptocompare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"final_course/internal/entities"
	"github.com/pkg/errors"
)

// Client структура реализующая интерфейс Client
type Client struct {
	httpClient    http.Client
	baseURL       string
	defaultTitles []string // Переменная для хранения списка монет
}

// NewClient конструктор, создаёт новый экземпляр Client
func NewClient(url string, titles []string) (*Client, error) {
	cl := http.Client{}
	return &Client{
		httpClient:    cl,
		baseURL:       url,
		defaultTitles: titles,
	}, nil
}

// Get реализует метод интерфейса Client
func (c *Client) Get(ctx context.Context, titles []string) ([]entities.Coin, error) {
	rawURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, errors.Wrap(err, "error parse")
	}

	if len(titles) == 0 {
		titles = c.defaultTitles
	}
	query := rawURL.Query()
	query.Set("fsyms", "USDT")
	query.Set("tsyms", "USD")
	rawURL.RawQuery = query.Encode()
	// Cоздаём HTTP GET запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get data: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body) //TODO: это убрать
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data map[string]map[string]float64 //TODO: допилить с анмаршаллингом
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Выводим данные после анмаршалинга
	fmt.Printf("Unmarshalled data: %+v\n", data)

	// Преобразуем данные в слайс структур Coin
	var coins []entities.Coin
	for title, rates := range data {
		if rate, ok := rates["USD"]; ok {
			coin, err := entities.NewCoin(title, rate, time.Now())
			if err != nil {
				return nil, fmt.Errorf("failed to create coin for symbol %s: %w", title, err)
			}
			coins = append(coins, *coin)
		}
	}
	return coins, nil
}
