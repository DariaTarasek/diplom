package sharederrors

import "errors"

var (
	ErrInvalidValue = errors.New("некорректные данные")
)
