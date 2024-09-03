package alogger

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
)

// Проверка имплементации интерфейса
var (
	_ IALogger = (*ALogger)(nil)
	_ IEvent   = (*Event)(nil)
)

type Level int

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return unknownData
	}
}

func (l Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

type CtxKey string

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8

	// unknownData - это значение по умолчанию если нет данных
	unknownData = "unknown"
	// TraceIdKey - ключ для записи trace_id в контекст
	TraceIdKey = "trace_id"
	// DefaultStackTraceBuffer - длина буфера стектрейса по умолчанию
	DefaultStackTraceBuffer = 1024
)

var (
	// TODO формат времени который принимет кибана
	// timeFormatForKibana = "02.01.2006 15:04:05.000"

	// defaultConfig Дефолтная конфигурация ALogger
	defaultConfig = &Config{
		Output: os.Stdout,
		Level:  LevelInfo,
	}

	flushDriver func(aLogger *ALogger) error = nil
)

// IALogger Интерфейс логгера ALogger
type IALogger interface {
	// Flush Вывод событий
	Flush() error
	// Debugf Создание события уровня Debug
	Debugf(format string, a ...interface{})
	// Infof Создание события уровня Info
	Infof(format string, a ...interface{})
	// Warnf Создание события уровня Warn
	Warnf(format string, a ...interface{})
	// Errorf Создание события уровня Erro
	Errorf(format string, a ...interface{})
	// SetAttr Установка атрибута
	SetAttr(key string, value interface{}) *ALogger
	// SetAttr Установка атрибутов
	SetAttrs(attrs map[string]interface{}) *ALogger
}

// IEvent Интерфейс события
type IEvent interface {
	// SetAttr Установка атрибута
	SetAttr(key string, value interface{}) *Event
	// GetAttr Чтение атрибута
	GetAttr(key string) (interface{}, bool)
	// Wrap Добавление исходной ошибки
	Wrap(err error) *Event
	// Unwrap Получение исходной ошибки
	Unwrap() error
	// Flush Логирование события
	Flush(w io.Writer, textFormat bool) error
}

// Config Конфигурация ALogger
type Config struct {
	// Output Поток вывода
	Output io.Writer
	// Level Задаёт уровень логирования "debug","info","warn","error", по умолчанию - "info"
	Level Level
	// TextFormat Включает вывод в текстовом формате
	TextFormat bool
	// OnlyError Включает логирование только при ошибках
	OnlyError bool
}

// AError Минимальная реализация интерфейса ошибки AError
type AError interface {
	// Code Индивидуальный код ошибки
	Code() string
	// Error Реализация error, возвращает сообщение исходной ошибки
	Error() string
	// MarshalJSON Кастомный метод MarshalJSON
	MarshalJSON() ([]byte, error)
}

// customError Кастомная ошибка поддерживающая метод MarshalJSON для AError
type customError struct {
	srcError error
}

func (e *customError) Error() string {
	return e.srcError.Error()
}

func (e *customError) Code() string {
	var aerr AError
	if errors.As(e.srcError, &aerr) {
		return aerr.Code()
	} else {
		return ""
	}
}

func (e *customError) Unwrap() error {
	return e.srcError
}

func (e *customError) MarshalJSON() ([]byte, error) {
	if e.srcError == nil {
		return []byte("null"), nil
	}

	// Если ошибка srcError == AError требуется форматировать ее в JSON т.к. она поддерживает метод MarshalJSON.
	var aerr AError
	if errors.As(e.srcError, &aerr) {
		return json.Marshal(aerr)
	} else {
		return json.Marshal(e.srcError.Error())
	}
}

func getTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value(TraceIdKey).(string)
	if !ok || traceId == "" {
		return unknownData
	}

	return traceId
}
