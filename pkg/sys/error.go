package sys

import (
	"github.com/pkg/errors"

	"github.com/merynayr/PingerVK/pkg/sys/codes"
)

// Структура для обработчки http ошибок
type commonError struct {
	msg  string
	code codes.Code
}

// NewCommonError создаёт новвую ошибку
func NewCommonError(msg string, code codes.Code) *commonError {
	return &commonError{msg, code}
}

func (r *commonError) Error() string {
	return r.msg
}

func (r *commonError) Code() codes.Code {
	return r.code
}

// IsCommonError проверяет на соответствие ошибке
func IsCommonError(err error) bool {
	var ce *commonError
	return errors.As(err, &ce)
}

// GetCommonError получает ошбику
func GetCommonError(err error) *commonError {
	var ce *commonError
	if !errors.As(err, &ce) {
		return nil
	}

	return ce
}
