package cryptocompare

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	query.Set("fsyms", strings.Join(titles, ",")) // Добавлен параметр fsyms
	query.Set("tsyms", "USD")
	rawURL.RawQuery = query.Encode()
	//--------------------------------------
	// Добавьте это для логирования URL
	log.Printf("[DEBUG] Request URL: %s", rawURL.String())
	//--------------------------------------

	// Cоздаём HTTP GET запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL.String(), nil) // Используем rawURL.String() для создания запроса
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

	var rates map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return nil, fmt.Errorf("failed to decode response:%w", err)
	}

	// Преобразуем данные в слайс структур Coin
	var coins []entities.Coin
	for title, cost := range rates {
		coin, err := entities.NewCoin(title, cost["USD"], time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to create coin:%w", err)
		}
		coins = append(coins, *coin)
	}
	return coins, nil
}
