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
	//TODO: добавить переменную для монет
}

// NewClient конструктор, создаёт новый экземпляр Client
func NewClient(httpClient *http.Client, url string) *Client { //TODO: конструктор должен возвращать ошибку
	return &Client{
		httpClient: httpClient,
		baseURL:    url,
	}
}

// Get реализует метод интерфейса Client
func (c *Client) Get(ctx context.Context, titles []string) ([]entities.Coin, error) { //TODO: здесь в слайс приходят монеты-нет монет-парсит по дефолту монеты
	// базовый URL для запроса
	url := c.baseURL + "/data/pricemulti"

	// создаём HTTP GET запрос с контекстом-2
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "error creating request") //TODO: здесь ошибка парсинга
	}

	// Получаем объект url.Values
	query := req.URL.Query()
	//TODO:сначала параметры-потом реквест-потом ...3-получение ответа
	// Устанавливаем query-параметры-1
	if len(titles) > 0 { //TODO: добавить проверку что нет монет-парсить по дефолту прописанные монеты(0-дефолт или не 0-то что прописано в слайс)-завести переменную-но нужно завести конфиг
		query.Set("fsyms", strings.Join(titles, ","))
		query.Set("tsyms", "USD") // Предполагаем, что нам нужны курсы в долларах
	}

	// Обновляем URL запроса
	req.URL.RawQuery = query.Encode()

	// выполняем запрос-3
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "failed to execute request") //TODO: тут другая ошибка-она не внутренняя уже
	}
	defer resp.Body.Close()

	// проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "failed to get data") //TODO: тут другая ошибка-она не внутренняя уже
	}

	// распарсить JSON ответ в слайс структур Coin
	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "failed to decode response") //TODO: тут другая ошибка-она не внутренняя уже-errorf
	}

	// Преобразуем данные в слайс структур Coin
	var coins []entities.Coin
	for symbol, rates := range data {
		if rate, ok := rates["USD"]; ok {
			coin, err := entities.NewCoin(symbol, rate, time.Now())
			if err != nil {
				return nil, err //TODO: ошибка чтоб понятно было что за ошибка-errof
			}
			coins = append(coins, *coin)
		}
	}

	return coins, nil
}
