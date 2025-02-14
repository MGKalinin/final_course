package cryptocompare

import (
	"context"
	"encoding/json"
	"final_course/internal/entities"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client структура реализующая интерфейс Client
type Client struct {
	httpClient *http.Client
	baseURL    string
	//TODO: добавить переменную для монет
}

// NewClient конструктор, создаёт новый экземпляр Client
func NewClient(httpClient *http.Client, url string) (*Client, error) {
	if httpClient == nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "httpClient cannot be nil")
	}
	if url == "" {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "url cannot be empty")
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    url,
	}, nil
}

// Get реализует метод интерфейса Client
func (c *Client) Get(ctx context.Context, titles []string) ([]entities.Coin, error) { //TODO: здесь в слайс приходят монеты-нет монет-парсит по дефолту монеты

	// Создаем объект url.Values
	query := url.Values{}

	//TODO: добавить проверку что нет монет-парсить по дефолту прописанные монеты(0-дефолт или
	//		//TODO: не 0-то что прописано в слайс)-завести переменную-но нужно завести конфиг
	// Устанавливаем query-параметры
	if len(titles) > 0 {
		query.Set("fsyms", strings.Join(titles, ","))
	} else {
		// Если titles пустой, используем монеты по умолчанию
		defaultCoins := []string{"BTC", "ETH", "LTC"} // Пример монет по умолчанию
		query.Set("fsyms", strings.Join(defaultCoins, ","))
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
		return nil, fmt.Errorf("failed to get data: %w", err)
	}

	// Распарсить JSON ответ в слайс структур Coin
	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
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
