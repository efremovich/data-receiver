package entity

import "errors"

var (
	ErrObjectNotFound = errors.New("объект не найден")
	ErrTemporary      = errors.New("произошла временная ошибка при обработке запроса")
	ErrPermanent      = errors.New("запрос не может быть выполнен")
)
