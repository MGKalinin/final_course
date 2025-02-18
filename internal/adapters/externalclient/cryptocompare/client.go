package cryptocompare

import (
	"context"
	"encoding/json"
	"final_course/internal/entities"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client структура реализующая интерфейс Client
type Client struct {
	httpClient *http.Client
	baseURL    string
	coins      []string // Переменная для хранения списка монет //TODO default.coin что-то типа
}

// NewClient конструктор, создаёт новый экземпляр Client
func NewClient(httpClient *http.Client, url string, coins []string) (*Client, error) { //TODO не coins- titles
	if httpClient == nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "httpClient cannot be nil")
	}
	if url == "" {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "url cannot be empty")
	}
	if len(coins) == 0 { //TODO здесь проверка что titles не равен 0-иначе по умолчанию
		coins = []string{"BTC", "ETH", "LTC"} // Монеты по умолчанию
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    url,
		coins:      coins,
	}, nil
}

// Get реализует метод интерфейса Client
func (c *Client) Get(ctx context.Context, titles []string) ([]entities.Coin, error) {
	// Создаем объект url.Values
	query := url.Values{}

	// Устанавливаем query-параметры
	if len(titles) > 0 {
		query.Set("fsyms", strings.Join(titles, ","))
	} else {
		query.Set("fsyms", strings.Join(c.coins, ","))
	}
	query.Set("tsyms", "USD") // Предполагаем, что нам нужны курсы в долларах

	// Cоздаём HTTP GET запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/data/pricemulti", nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}
	// Обновляем URL запроса
	req.URL.RawQuery = query.Encode()

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get data: status code %d, response: %s", resp.StatusCode, string(body))
	}
	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body) // Добавлено
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err) // Добавлено
	}

	// Логируем тело ответа для отладки
	//fmt.Println("Response body:", string(body))

	// Распарсить JSON ответ в слайс структур Coin
	var data map[string]map[string]float64
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Преобразуем данные в слайс структур Coin
	var coins []entities.Coin
	for symbol, rates := range data {
		if rate, ok := rates["USD"]; ok {
			coin, err := entities.NewCoin(symbol, rate, time.Now())
			if err != nil {
				return nil, fmt.Errorf("failed to create coin for symbol %s: %w", symbol, err)
			}
			coins = append(coins, *coin)
		}
	}
	return coins, nil
}
