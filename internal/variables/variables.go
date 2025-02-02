package variables

import "github.com/pkg/errors"

// создать переменную ошибки- и далее использовать её в обёртке во всех случаях,
// добавляя соответсвующий комментарий

var SomeErr = errors.New("Error:")
