package application

import (
	_ "final_course/docs"
	"final_course/internal/adapters/externalclient/cryptocompare"
	"final_course/internal/adapters/storage/postgres"
	"final_course/internal/cases"
	"final_course/internal/port/http/public"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
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

// конфигурация
// LoadConfig грузим конфигурацию
//configs.LoadConfig()// раскомитить
//dbParams := viper.GetStringMapString("database")
//// Формируем DSN
//dsn := fmt.Sprintf("%s://%s:%s@localhost:%s/%s",
//dbParams["db_name"],
//dbParams["username"],
//dbParams["password"],
//dbParams["address"],
//dbParams["db_name"],
//)
//// Установка переменной окружения
////os.Setenv("DATABASE_URL", "postgres://maksimkalinin:password@localhost:5432/postgres")
//os.Setenv("DATABASE_URL", dsn)

// Создание контекста
ctx := context.Background()

// Определение монет для запроса
coinsToFetch := []string{"BTC", "ETH", "XRP"}

// конструктор client
// Инициализация клиента
//client, err := cryptocompare.NewClient("https://min-api.cryptocompare.com/data/pricemulti", coinsToFetch)
client, err := cryptocompare.NewClient(dbParams["client_address"], coinsToFetch)

if err != nil {
log.Fatalf("Error creating client: %v", err)
}
// конструктор storage

// Инициализация хранилища
storage, err := storage.NewStorage(ctx, os.Getenv("DATABASE_URL"))
if err != nil {
log.Fatalf("Error creating storage: %v", err)
}
// конструктор service

// Инициализация сервиса
service, err := cases.NewService(storage, client)
if err != nil {
log.Fatalf("Error creating service: %v", err)
}
// конструктор server

//Инициализация сервера
server, err := public.NewServer(service)
if err != nil {
log.Fatalf("Failed to create server: %v", err)
}

// Запуск сервера
server.Run()
// cron запрос в фоне горутина курсов с заданной периодичностью
