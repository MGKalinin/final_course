package application

import (
	"context"
	"final_course/deploy/configs"
	"final_course/internal/adapters/externalclient/cryptocompare"
	storage "final_course/internal/adapters/storage/postgres"
	"final_course/internal/cases"
	"final_course/internal/port/http/public"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//TODO: 1.здесь создать все конструкторы ,запускаешь, обрабатываешь все ошибки ,запуск метода сервиса,
// запрос курсов в фоне(горутина)раз в n времени
// cron -запускает в фоне
// 2. в main запуск run
// 3. и всё это собрать в докере

//TODO: 4.перезапустить генерацию документации swag- будет перенесён main и изменилась структура (оставить
// @Failure 400, @Failure 404, @Failure 500 )
// убрать русскую анотацию

//эйч навыки собесы по golang

type App struct {
	server *public.Server
	cron   *cron.Cron
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configs.LoadConfig() // раскомитить
	dbParams := viper.GetStringMapString("database")
	// Формируем DSN
	dsn := fmt.Sprintf("%s://%s:%s@localhost:%s/%s",
		dbParams["db_name"],
		dbParams["username"],
		dbParams["password"],
		dbParams["address"],
		dbParams["db_name"],
	)
	// Установка переменной окружения
	os.Setenv("DATABASE_URL", dsn)

	// Определение монет для запроса
	coinsToFetch := []string{"BTC", "ETH", "XRP"}

	// Инициализация клиента
	client, err := cryptocompare.NewClient(dbParams["client_address"], coinsToFetch)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Инициализация хранилища
	storage, err := storage.NewStorage(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Инициализация сервиса
	service, err := cases.NewService(storage, client)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	// Инициализация сервера
	server, err := public.NewServer(service)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Настройка планировщика задач
	a.cron = cron.New()
	_, err = a.cron.AddFunc("@every 5m", func() {
		log.Println("[CRON] Запуск фонового обновления данных...")
		if err := service.FetchAndStoreCoins(ctx); err != nil {
			log.Printf("[CRON] Ошибка обновления: %v", err)
		} else {
			log.Println("[CRON] Данные успешно обновлены")
		}
	})
	if err != nil {
		log.Fatalf("Ошибка настройки расписания: %v", err)
	}
	a.cron.Start()
	log.Println("Фоновые задачи запланированы с интервалом 5 минут")

	// Настройка graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Получен сигнал завершения работы...")
		cancel()

		// Остановка планировщика
		a.cron.Stop()
		log.Println("Фоновые задачи остановлены")
	}()

	// Запуск HTTP-сервера
	log.Println("Запуск сервера на порту :8080")
	server.Run() // Блокирующий вызов

	return nil // Корректное завершение
}
