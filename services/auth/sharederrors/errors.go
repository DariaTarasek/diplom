package sharederrors

import "errors"

var (
	ErrCodeInvalid     = errors.New("неверный код подтверждения")
	ErrCodeExpired     = errors.New("действие кода подтверждения истекло")
	ErrTooManyAttempts = errors.New("слишком много попыток ввода кода")
	ErrRateLimited     = errors.New("код уже отправлен")
	ErrPasswordInvalid = errors.New("неверный пароль")
)
