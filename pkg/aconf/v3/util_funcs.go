package aconf

import (
	"reflect"
)

// isAPtrToStruct возвращает true, если переданный интерфейс
// является указателем на структуру
func isPtrToStruct(v interface{}) bool {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return false
	}
	return value.Elem().Kind() == reflect.Struct
}

// isPtrToMap проверяет, что значение переданного интерфейса
// является указателем на map
func isPtrToMap(v interface{}) bool {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return false
	}
	return value.Elem().Kind() == reflect.Map
}
