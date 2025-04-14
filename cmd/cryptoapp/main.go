package main

import (
	"final_course/internal/application"
	"log"
)

// Общее описание API

// @title Crypto Currency API
// @version 1.0
// @description API для получения курсов криптовалют

// @contact.name Максим Калинин
// @contact.email ваш@email.com

// @host localhost:8080
// @BasePath /v1
func main() {
	app := application.NewApp()

	if err := app.Run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

// TODO: создать .ENV для docker compouse куда положить user pass базы данных туда где лежит docker compouse файл
//  логин , пароль вернуть
//  docker файлы не менять
