package cryptocompare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Выводим сырой JSON-ответ
	fmt.Println("Raw JSON response:")
	fmt.Println(string(body))

	// Используем json.NewDecoder для декодирования JSON-ответа
	var response struct {
		Response       string                        `json:"Response"`
		Message        string                        `json:"Message"`
		HasWarning     bool                          `json:"HasWarning"`
		Type           int                           `json:"Type"`
		RateLimit      map[string]interface{}        `json:"RateLimit"`
		Data           map[string]map[string]float64 `json:"Data"`
		ParamWithError string                        `json:"ParamWithError"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Response == "Error" {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	// Преобразуем данные в слайс структур Coin
	var coins []entities.Coin
	for title, rates := range response.Data {
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
