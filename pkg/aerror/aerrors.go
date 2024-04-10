package aerrors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"
)

const (
	// TraceIdKey Ключ получение trace_id
	TraceIdKey = "trace_id"
	// unknownData Неизвестные данные
	unknownData    = "unknown"
	unknownTraceId = "000"
)

// AError Интерфейс ошибки
type AError interface {
	GetID() ErrorId
	// SetUserMsgf Установка сообщения для пользователя
	SetUserMsgf(string, ...interface{}) AError
	// SetAttr Установка атрибута ошибки
	SetAttr(key string, value interface{}) AError
	// GetAttr Получения атрибута ошибки
	GetAttr(key string) (interface{}, bool)
	// Code Индивидуальный код ошибки
	Code() string
	// Error Реализация error, возвращает сообщение исходной ошибки или сообщение разработчика
	Error() string
	// UserMessage Сообщение ошибки для пользователя
	UserMessage() string
	// DeveloperMessage Сообщение ошибки разработчика, если оно не заполнено возвращается сообщение исходной ошибки
	DeveloperMessage() string
	// IsCritical Является ли ошибка критической
	IsCritical() bool
	// Unwrap Возвращает исходную ошибку или nil
	Unwrap() error
	// IErrorMarshalJSON Конвертация в JSON аналогично MarshalJSON для IError
	IErrorMarshalJSON() ([]byte, error)
	// Is сравнивает только ErrorId. В противном случае сравнивается исходная ошибка
	Is(error) bool
}

type aError struct {
	// Id Индивидуальный код ошибки
	Id ErrorId `json:"id"`
	// CreatedAt дата создания
	CreatedAt time.Time `json:"created_at"`
	// TraceId уникальный идентификатор трассировки
	TraceId string `json:"trace_id"`
	// Сообщение для разработчика
	DeveloperMsg string `json:"developer_message"`
	// Сообщение для пользователя
	UserMsg string `json:"user_message"`
	// Caller строка вызова ошибки
	Caller string `json:"caller"`
	// Critical признак критической ошибки
	Critical bool `json:"critical"`
	// Attrs дополнительная информация для разработчиков
	Attrs map[string]interface{} `json:"attrs,omitempty"`
	// SrcError Исходная ошибка
	SrcError error `json:"src_error,omitempty"`
	// TODO: удалить как уйдем от IError. Временный признак трансформированной ошибки.
	IErrorCode string `json:"-"`
}

// New создание AError
func New(ctx context.Context, id ErrorId, err error, msgf string, a ...interface{}) AError {
	aE := createError(ctx, id, err, msgf, a...)

	return aE
}

// NewCritical создание AError со статусом critical
func NewCritical(ctx context.Context, id ErrorId, err error, msgf string, a ...interface{}) AError {
	aE := createError(ctx, id, err, msgf, a...)
	aE.Critical = true

	return aE
}

func (e *aError) GetID() ErrorId {
	return e.Id
}

func (e *aError) SetUserMsgf(msg string, args ...interface{}) AError {
	e.UserMsg = fmt.Sprintf(msg, args...)
	return e
}

func (e *aError) SetAttr(key string, value interface{}) AError {
	if e.Attrs == nil {
		e.Attrs = make(map[string]interface{})
	}

	e.Attrs[key] = value

	return e
}

func (e *aError) GetAttr(key string) (interface{}, bool) {
	if e.Attrs == nil {
		return "", false
	}

	v, ok := e.Attrs[key]

	return v, ok
}

func (e *aError) Code() string {
	// Если ошибка конвертирована из IError, то используется код IError
	if e.IErrorCode != "" {
		return e.IErrorCode
	}

	return fmt.Sprintf("%d", e.Id)

}

func (e *aError) Error() string {
	switch {
	case e.SrcError != nil:
		return e.SrcError.Error()
	default:
		return e.DeveloperMsg
	}
}

func (e *aError) DeveloperMessage() string {
	if e.DeveloperMsg == "" && e.SrcError != nil {
		return e.SrcError.Error()
	}

	return e.DeveloperMsg
}

func (e *aError) UserMessage() string {
	return e.UserMsg + fmt.Sprintf(". (%s_%s)", e.Code(), e.TraceId)
}

type AErrorMarshal struct {
	// Id Индивидуальный код ошибки
	Id ErrorId `json:"id"`
	// CreatedAt дата создания
	CreatedAt string `json:"created_at"`
	// TraceId уникальный идентификатор трассировки
	TraceId string `json:"trace_id"`
	// Сообщение для разработчика
	DeveloperMsg string `json:"developer_message"`
	// Сообщение для пользователя
	UserMsg string `json:"user_message"`
	// Caller строка вызова ошибки
	Caller string `json:"caller"`
	// Critical признак критической ошибки
	Critical bool `json:"critical"`
	// Attrs дополнительная информация для разработчиков
	Attrs map[string]interface{} `json:"attrs,omitempty"`
	// SrcError Исходная ошибка
	SrcError string `json:"src_error"`
	// Code уникальный код ошибки
	Code string `json:"code"`
}

func (e *aError) MarshalJSON() ([]byte, error) {
	// Создание структуры описания разработчика
	aErrorMarshal := AErrorMarshal{
		Id:           e.Id,
		CreatedAt:    e.CreatedAt.Format("2006-01-02 15:04:05.000"),
		TraceId:      e.TraceId,
		DeveloperMsg: e.DeveloperMsg,
		UserMsg:      e.UserMsg,
		Caller:       e.Caller,
		Critical:     e.Critical,
		Attrs:        e.Attrs,
		Code:         e.Code(),
	}

	// Сообщение оригинальной ошибки
	if e.SrcError != nil {
		aErrorMarshal.SrcError = e.SrcError.Error()
	}

	aErrorJson, err := json.Marshal(&aErrorMarshal)
	if err != nil {
		return nil, err
	}

	return aErrorJson, nil
}

func (e *aError) IsCritical() bool {
	return e.Critical
}

func (e *aError) Unwrap() error {
	return e.SrcError
}

// IErrorMarshalJSON - Реализация json.Marshal() аналогично IError
func (e *aError) IErrorMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message,omitempty"`
	}{Code: e.Code(), Message: e.UserMsg})
}

func (e *aError) Is(err error) bool {
	// Проверка ErrorId
	var aE *aError
	ok := errors.As(err, &aE)
	if ok {
		return aE.Id == e.Id
	}

	// Классический Is
	if e.SrcError != nil {
		return errors.Is(e.SrcError, err)
	}

	return false
}

// createError создание ошибки, получение из контекста trace_id
func createError(ctx context.Context, id ErrorId, err error, msgf string, a ...interface{}) *aError {
	// runtimeCallerLevel Уровень вызова runtime.Caller()
	runtimeCallerLevel := 2
	_, file, line, _ := runtime.Caller(runtimeCallerLevel)
	caller := fmt.Sprintf("%s:%d", file, line)

	// Получение trace_id
	traceId, ok := ctx.Value(TraceIdKey).(string)
	if !ok {
		traceId = unknownTraceId
	}

	e := aError{
		Id:           id,
		CreatedAt:    time.Now(),
		TraceId:      traceId,
		SrcError:     err,
		Caller:       caller,
		DeveloperMsg: fmt.Sprintf(msgf, a...),
		UserMsg:      userMessages[id],
	}

	return &e
}
