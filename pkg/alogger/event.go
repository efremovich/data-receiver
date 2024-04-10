package alogger

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"
)

// Event Структура события
type Event struct {
	CreatedAt   time.Time              `json:"created_at"`            // CreatedAt Дата создания события
	Level       Level                  `json:"level"`                 // Level Уровень события
	TraceId     string                 `json:"trace_id"`              // TraceId Идентификатор трейса, полученный из контекста
	Message     string                 `json:"message"`               // Message Сообщение события
	Caller      string                 `json:"caller"`                // Caller Строка вызова логирования
	PackageName string                 `json:"package_name"`          // PackageName Имя пакета
	ErrorCode   string                 `json:"error_code,omitempty"`  // ErrorCode Код ошибки (if srcError != nil && (err == AError || err == IError))
	SrcError    *customError           `json:"src_error,omitempty"`   // SrcError Исходная ошибка
	Attrs       map[string]interface{} `json:"attrs,omitempty"`       // Attrs Атрибуты события
	MetaInfo    MetaInfo               `json:"meta_info,omitempty"`   // MetaInfo Информация, выступающая в роли sub_trace_id
	StackTrace  []byte                 `json:"stack_trace,omitempty"` // StackTrace Стек вызова
}

func caller() (caller, packageName string) {
	callerSkip := 3
	pc, file, line, ok := runtime.Caller(callerSkip)

	if !ok {
		return unknownData, unknownData
	}

	// Получение строки вызова
	caller = fmt.Sprintf("%s:%d", file, line)

	// Получение packageName
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		fullName := fn.Name()
		sliceName := strings.Split(fullName, ".")

		if len(sliceName) > 0 {
			packageName = sliceName[0]
		}
	}

	if packageName == "" {
		packageName = unknownData
	}

	return caller, packageName
}

func createEvent(msg string, level Level, traceId string) *Event {
	// Генерация caller
	caller, packageName := caller()

	// Создание события
	ev := Event{
		TraceId:     traceId,
		CreatedAt:   time.Now(),
		Level:       level,
		Message:     msg,
		PackageName: packageName,
		Caller:      caller,
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

func (ev *Event) SetMetaInfo(key MetaInfoKey, value interface{}) *Event {
	switch key {
	case TPIdKey:
		if str, ok := value.(string); ok {
			ev.MetaInfo.TPId = str
		}
	case LSIdKey:
		if str, ok := value.(string); ok {
			ev.MetaInfo.LSId = str
		}
	case DocIdKey:
		if str, ok := value.(string); ok {
			ev.MetaInfo.DocId = str
		}
	case SenderIdKey:
		if str, ok := value.(string); ok {
			ev.MetaInfo.SenderId = str
		}
	case ReceiverIdKey:
		if str, ok := value.(string); ok {
			ev.MetaInfo.ReceiverId = str
		}
	}

	return ev
}

func (ev *Event) GetMetaInfo(key MetaInfoKey) (interface{}, bool) {
	switch key {
	case TPIdKey:
		if ev.MetaInfo.TPId != "" {
			return ev.MetaInfo.TPId, true
		}
	case LSIdKey:
		if ev.MetaInfo.LSId != "" {
			return ev.MetaInfo.LSId, true
		}
	case DocIdKey:
		if ev.MetaInfo.DocId != "" {
			return ev.MetaInfo.DocId, true
		}
	case SenderIdKey:
		if ev.MetaInfo.SenderId != "" {
			return ev.MetaInfo.SenderId, true
		}
	case ReceiverIdKey:
		if ev.MetaInfo.ReceiverId != "" {
			return ev.MetaInfo.ReceiverId, true
		}
	}

	return nil, false
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

func (ev *Event) Stack() *Event {
	stack := make([]byte, DefaultStackTraceBuffer)
	runtime.Stack(stack, true)

	ev.StackTrace = stack

	return ev
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
	format := "\n%s %s: %s\nTraceId: %s\nPackageName: %s\nCaller: %s\n"
	bufferForWrite = append(bufferForWrite, fmt.Sprintf(format,
		ev.CreatedAt.Format("[02.01.2006 15:04:05.000]"),
		ev.Level,
		ev.Message,
		ev.TraceId, ev.PackageName, ev.Caller)...)

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

	// Добавление мета информации в буфер вывода
	if !ev.MetaInfo.IsEmpty() {
		bufferForWrite = append(bufferForWrite, ev.MetaInfo.Bytes()...)
	}

	// Добавление cтека трейса в буфер вывода
	if len(ev.StackTrace) > 0 {
		bufferForWrite = append(bufferForWrite, fmt.Sprintf("%s: %s\n", "StackTrace", ev.StackTrace)...)
	}

	return bufferForWrite, nil
}
