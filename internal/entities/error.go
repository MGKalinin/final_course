package entities

import "github.com/pkg/errors"

// ErrorInvalidParams общая ошибка для невалидных параметров
// @Description Возникает при некорректных входных данных
// @Example "invalid params: example error message"
var ErrorInvalidParams = errors.New("invalid params:")
