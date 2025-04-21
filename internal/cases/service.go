package cases

import (
	"context"
	"final_course/internal/entities"
	"final_course/internal/port/http/public"
	"github.com/pkg/errors"
)

// Добавляем проверку реализации интерфейса
var _ public.Service = (*Service)(nil)

// Service содержит логику работы сервиса
type Service struct {
	storage Storage
	client  Client
}

// NewService конструктор - создает новый сервис
func NewService(storage Storage, client Client) (*Service, error) {
	if storage == nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "storage is nil")
	}
	if client == nil {
		return nil, errors.Wrap(entities.ErrorInvalidParams, "client is nil")
	}
	return &Service{storage: storage, client: client}, nil
}

// FetchAndStoreCoins метод вызывает метод бд для получения всех titles, затем использует клиент для получения данных о монетах и сохраняет их в бд
func (s *Service) FetchAndStoreCoins(ctx context.Context) error {
	// Получаем все titles из базы данных
	titles, err := s.storage.GetAllTitles(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get titles from database")
	}

	// Получаем данные о монетах с помощью клиента
	coins, err := s.client.Get(ctx, titles)
	if err != nil {
		return errors.Wrap(err, "failed to get coins from client")
	}

	// Сохраняем монеты в базе данных
	if err := s.storage.Store(ctx, coins); err != nil {
		return errors.Wrap(err, "failed to store coins in database")
	}

	return nil
}

type AggFunc int

const (
	_ AggFunc = iota
	Max
	Min
	Avg
)

type Options struct {
	FuncType AggFunc
}
type Option func(options *Options)

func (a AggFunc) String() string {
	return [...]string{"", "MAX", "MIN", "AVG"}[a]
}

// Вспомогательная функция для поиска отсутствующих монет
func findMissingTitles(requested, existing []string) []string {
	existingSet := make(map[string]struct{})
	for _, title := range existing {
		existingSet[title] = struct{}{}
	}

	var missing []string
	for _, title := range requested {
		if _, ok := existingSet[title]; !ok {
			missing = append(missing, title)
		}
	}
	return missing
}

// WithMaxFunc функция получения max значения
func WithMaxFunc() Option {
	return func(options *Options) {
		options.FuncType = Max
	}
}

// GetMaxRate возвращает максимальные значения, автоматически запрашивая отсутствующие монеты
func (s *Service) GetMaxRate(ctx context.Context, titles []string) ([]entities.Coin, error) {
	// Получаем все существующие монеты из БД
	existingTitles, err := s.storage.GetAllTitles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get existing titles")
	}

	// Определяем отсутствующие монеты
	missing := findMissingTitles(titles, existingTitles)

	// Запрашиваем и сохраняем отсутствующие
	if len(missing) > 0 {
		coins, err := s.client.Get(ctx, missing)
		if err != nil {
			return nil, errors.Wrap(err, "client failed to fetch missing coins")
		}
		if err := s.storage.Store(ctx, coins); err != nil {
			return nil, errors.Wrap(err, "failed to store new coins")
		}
	}

	// Возвращаем данные из хранилища
	return s.storage.Get(ctx, titles, WithMaxFunc())
}

// WithMinFunc функция получения min значения
func WithMinFunc() Option {
	return func(options *Options) {
		options.FuncType = Min
	}
}

// GetMinRate возвращает минимальные значения, автоматически запрашивая отсутствующие монеты
func (s *Service) GetMinRate(ctx context.Context, titles []string) ([]entities.Coin, error) {
	existingTitles, err := s.storage.GetAllTitles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get existing titles")
	}

	missing := findMissingTitles(titles, existingTitles)
	if len(missing) > 0 {
		coins, err := s.client.Get(ctx, missing)
		if err != nil {
			return nil, errors.Wrap(err, "client failed to fetch missing coins")
		}
		if err := s.storage.Store(ctx, coins); err != nil {
			return nil, errors.Wrap(err, "failed to store new coins")
		}
	}

	return s.storage.Get(ctx, titles, WithMinFunc())
}

// WithAvgFunc функция получения avg значения
func WithAvgFunc() Option {
	return func(options *Options) {
		options.FuncType = Avg
	}
}

// GetAvgRate возвращает средние значения, автоматически запрашивая отсутствующие монеты
func (s *Service) GetAvgRate(ctx context.Context, titles []string) ([]entities.Coin, error) {
	existingTitles, err := s.storage.GetAllTitles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get existing titles")
	}

	missing := findMissingTitles(titles, existingTitles)
	if len(missing) > 0 {
		coins, err := s.client.Get(ctx, missing)
		if err != nil {
			return nil, errors.Wrap(err, "client failed to fetch missing coins")
		}
		if err := s.storage.Store(ctx, coins); err != nil {
			return nil, errors.Wrap(err, "failed to store new coins")
		}
	}

	return s.storage.Get(ctx, titles, WithAvgFunc())
}

// GetLastRates возвращает последние значения, автоматически запрашивая отсутствующие монеты
func (s *Service) GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	existingTitles, err := s.storage.GetAllTitles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get existing titles")
	}

	missing := findMissingTitles(titles, existingTitles)
	if len(missing) > 0 {
		coins, err := s.client.Get(ctx, missing)
		if err != nil {
			return nil, errors.Wrap(err, "client failed to fetch missing coins")
		}
		if err := s.storage.Store(ctx, coins); err != nil {
			return nil, errors.Wrap(err, "failed to store new coins")
		}
	}

	// Вызов без опций для получения последних значений
	return s.storage.Get(ctx, titles)
}
