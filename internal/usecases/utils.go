package usecases

import (
	"errors"
	"fmt"
)

// wrapErr определяет тип ошибки. Если ошибка == ErrObjectNotFound,
// то оборачивает в ErrPermanent, если любая другая ошибка, то оборачивает
// в ErrTemporary, если ошибка нулевая - возвращает nil.
func wrapErr(err error) error {
	var result error

	switch {
	case err == nil:
		result = nil
	case errors.Is(err, ErrObjectNotFound):
		result = fmt.Errorf("%w: %s", ErrPermanent, err.Error())
	default:
		result = fmt.Errorf("%w: %s", ErrTemporary, err.Error())
	}

	return result
}
