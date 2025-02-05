package adapters

import (
	"context"
	"encoding/json"
	"final_course/internal/entities"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// структура реализующая интерфейс Client (ключ API & url адрес запроса)

type Client struct {
	apiKey  string
	baseURL string
}

// конструктор, создаёт новый экземпляр Client

func NewClient(apiKey string, url string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: url,
	}
}

// Get реализует метод интерфейса Client
func (c *Client) Get(ctx context.Context, titles []string) ([]entities.Coin, error) {
	// базовый URL для запроса
	url := c.baseURL + "/coins"

	// создаём HTTP GET запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(entities.ErrorEntity, "error creating request")
	}

	// Получаем объект url.Values
	query := req.URL.Query()

	// Устанавливаем query-параметры
	if len(titles) > 0 {
		query.Set("symbols", strings.Join(titles, ","))
	}

	// Обновляем URL запроса
	req.URL.RawQuery = query.Encode()

	// добавляем заголовок с API ключом
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// создаём HTTP клиент и выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(entities.ErrorEntity, "failed to execute request")
	}
	defer resp.Body.Close()

	// проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(entities.ErrorEntity, "failed to get data")
	}

	// распарсить JSON ответ в слайс структур Coin
	var coins []entities.Coin
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return nil, errors.Wrap(entities.ErrorEntity, "failed to decode response")
	}
	return coins, nil
}
