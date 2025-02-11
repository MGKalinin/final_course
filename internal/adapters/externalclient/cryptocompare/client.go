package cryptocompare

import (
	"context"
	"encoding/json"
	"final_course/internal/entities"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

// Client структура реализующая интерфейс Client
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient конструктор, создаёт новый экземпляр Client
func NewClient(httpClient *http.Client, url string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    url,
	}
}

// Get реализует метод интерфейса Client
func (c *Client) Get(ctx context.Context, titles []string) ([]entities.Coin, error) {
	// базовый URL для запроса
	url := c.baseURL + "/data/pricemulti"

	// создаём HTTP GET запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "error creating request")
	}

	// Получаем объект url.Values
	query := req.URL.Query()

	// Устанавливаем query-параметры
	if len(titles) > 0 {
		query.Set("fsyms", strings.Join(titles, ","))
		query.Set("tsyms", "USD") // Предполагаем, что нам нужны курсы в долларах
	}

	// Обновляем URL запроса
	req.URL.RawQuery = query.Encode()

	// выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "failed to execute request")
	}
	defer resp.Body.Close()

	// проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "failed to get data")
	}

	// распарсить JSON ответ в слайс структур Coin
	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "failed to decode response")
	}

	// Преобразуем данные в слайс структур Coin
	var coins []entities.Coin
	for symbol, rates := range data {
		if rate, ok := rates["USD"]; ok {
			coin, err := entities.NewCoin(symbol, rate, time.Now())
			if err != nil {
				return nil, err
			}
			coins = append(coins, *coin)
		}
	}

	return coins, nil
}
