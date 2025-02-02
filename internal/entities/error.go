package entities

import "github.com/pkg/errors"

// создать переменную ошибки- и далее использовать её в обёртке во всех случаях,
// добавляя соответсвующий комментарий

// TODO: Переименовать var Error в понятное для всех название

var ErrorEntity = errors.New("Error:")

// TODO: едем в клиента- ооп, структура...
