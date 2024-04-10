package aerrors

import (
	"fmt"
)

// ErrorId Идентификатор ошибки который является уникальным
type ErrorId int

func (e ErrorId) UserMessage() string {
	res := fmt.Sprintf("ID ошибки: %d.", e)
	msg, ok := userMessages[e]
	if !ok {
		return res
	}
	return res + " " + msg
}

func AppendToUserMessages(mgs map[ErrorId]string) {
	for k, v := range mgs {
		userMessages[k] = v
	}
}


// Идентификаторы ошибок
const (
	ErrorIdNoCode         ErrorId = 1_00_1_000

	// код ошибки - шести значное число.
	// например: 1_10_1_001
	// где первая цифра - система.
	// 1 - оператор
	// 2 - регистратор
	// 3 - астрал доверенность.
	// вторая и третья цифры - сервис
	// 00 - неизвестен (не оператор)
	// 01 - росеу
	// 11 - новая маркировка
	// четвертая цифра - критическая или нет
	// 0 - нет
	// 1 - да
	// три последние цифры - id ошибки.
)

// userMessages Дефолтные сообщения ошибок для пользователей
var userMessages = map[ErrorId]string{
	ErrorIdNoCode:         "Неизвестная ошибка - специалисты оповещены",
}
