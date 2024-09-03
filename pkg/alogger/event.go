package alogger

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"time"
)

// Event Структура события
type Event struct {
	CreatedAt  time.Time              `json:"created_at"`            // CreatedAt Дата создания события
	Level      Level                  `json:"level"`                 // Level Уровень события
	TraceId    string                 `json:"trace_id"`              // TraceId Идентификатор трейса, полученный из контекста
	Message    string                 `json:"message"`               // Message Сообщение события
	Caller     string                 `json:"caller"`                // Caller Строка вызова логирования
	ErrorCode  string                 `json:"error_code,omitempty"`  // ErrorCode Код ошибки (if srcError != nil && (err == AError || err == IError))
	SrcError   *customError           `json:"src_error,omitempty"`   // SrcError Исходная ошибка
	Attrs      map[string]interface{} `json:"attrs,omitempty"`       // Attrs Атрибуты события
	StackTrace []byte                 `json:"stack_trace,omitempty"` // StackTrace Стек вызова
}

func caller() string {
	callerSkip := 4
	_, file, line, ok := runtime.Caller(callerSkip)

	if !ok {
		return unknownData
	}

	return fmt.Sprintf("%s:%d", file, line)
}

func createEvent(msg string, level Level, traceId string) *Event {
	// Генерация caller
	caller := caller()

	// Создание события
	ev := Event{
		TraceId:   traceId,
		CreatedAt: time.Now(),
		Level:     level,
		Message:   msg,
		Caller:    caller,
	}

	return &ev
}

func (ev *Event) SetAttr(key string, value interface{}) *Event {
	if ev.Attrs == nil {
		ev.Attrs = make(map[string]interface{})
	}

	ev.Attrs[key] = value

	return ev
}

func (ev *Event) GetAttr(key string) (interface{}, bool) {
	if ev.Attrs == nil {
		return nil, false
	}

	v, ok := ev.Attrs[key]

	return v, ok
}

func (ev *Event) Wrap(err error) *Event {
	// Трансформация ошибки в кастомную
	customError := customError{
		srcError: err,
	}

	// Получение кода ошибки
	ev.ErrorCode = customError.Code()

	// Добавление исходной ошибки в событие
	ev.SrcError = &customError

	return ev
}

func (ev *Event) Unwrap() error {
	return ev.SrcError.srcError
}

// Flush Логирование события
func (ev *Event) Flush(w io.Writer, textFormat bool) error {
	var bufferForWrite []byte

	if textFormat {
		eventText, err := ev.textFormat()
		if err != nil {
			return err
		}

		bufferForWrite = eventText
	} else {
		eventJson, err := ev.jsonFormat()
		if err != nil {
			return err
		}

		bufferForWrite = eventJson
	}

	if _, err := w.Write(bufferForWrite); err != nil {
		return fmt.Errorf("[ALogger] Ошибка вывода события в поток: %w", err)
	}

	return nil
}

func (ev *Event) jsonFormat() ([]byte, error) {
	// Форматирование события в JSON
	eventJson, err := json.Marshal(ev)
	if err != nil {
		return nil, fmt.Errorf("[ALogger] Ошибка формирования события в JSON: %w", err)
	}

	eventJson = append(eventJson, []byte("\n")...)

	return eventJson, nil
}

func (ev *Event) textFormat() ([]byte, error) {
	var bufferForWrite []byte

	// Добавление события в буфер вывода
	format := "\n%s %s: %s\nTraceId: %s\nCaller: %s\n"
	bufferForWrite = append(bufferForWrite, fmt.Sprintf(format,
		ev.CreatedAt.Format("[02.01.2006 15:04:05.000]"),
		ev.Level,
		ev.Message,
		ev.TraceId,
		ev.Caller)...)

	// Добавление исходной ошибки в буфер вывода
	if ev.SrcError != nil {
		bufferForWrite = append(bufferForWrite, fmt.Sprintf("%s: %s\n", "SrcError", ev.SrcError.Error())...)
		bufferForWrite = append(bufferForWrite, fmt.Sprintf("%s: %s\n", "ErrorCode", ev.ErrorCode)...)
	}

	// Добавление атрибутов в буфер вывода
	if len(ev.Attrs) > 0 {
		bufferForWrite = append(bufferForWrite, "Attrs:\n"...)

		for k, v := range ev.Attrs {
			bufferForWrite = append(bufferForWrite, fmt.Sprintf("\t%s: %v\n", k, v)...)
		}
	}

	// Добавление cтека трейса в буфер вывода
	if len(ev.StackTrace) > 0 {
		bufferForWrite = append(bufferForWrite, fmt.Sprintf("%s: %s\n", "StackTrace", ev.StackTrace)...)
	}

	return bufferForWrite, nil
}
