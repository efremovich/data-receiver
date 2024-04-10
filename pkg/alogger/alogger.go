package alogger

import (
	"context"
	"fmt"
	"os"
	"sync"
)

func SetDefaultConfig(cfg *Config) {
	defaultConfig = cfg
}

// SetFlushDriver Устанавливает кастомный метод вывода логов
func SetFlushDriver(f func(aLogger *ALogger) error) {
	flushDriver = f
}

// ALogger Структура логера
type ALogger struct {
	mu           sync.Mutex
	cfg          *Config                // Конфигурация логера
	generalAttrs map[string]interface{} // общие аттрибуты для всех событий
	events       []*Event               // Список событий
	traceId      string                 // Trace_id полученный из контекста
	hasError     bool                   // Признак наличия ошибок
}

func NewALogger(ctx context.Context, cfg *Config) *ALogger {
	if cfg == nil {
		cfg = defaultConfig
	}

	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	return &ALogger{
		cfg:     cfg,
		traceId: getTraceId(ctx),
	}
}

func (l *ALogger) appendEvent(ev *Event) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.events = append(l.events, ev)
}

func (l *ALogger) Debugf(format string, a ...interface{}) *Event {
	msg := fmt.Sprintf(format, a...)
	ev := createEvent(msg, LevelDebug, l.traceId)

	for key, value := range l.generalAttrs {
		ev.SetAttr(key, value)
	}

	l.appendEvent(ev)

	return ev
}

func (l *ALogger) Infof(format string, a ...interface{}) *Event {
	msg := fmt.Sprintf(format, a...)
	ev := createEvent(msg, LevelInfo, l.traceId)

	for key, value := range l.generalAttrs {
		ev.SetAttr(key, value)
	}

	l.appendEvent(ev)

	return ev
}

func (l *ALogger) Warnf(format string, a ...interface{}) *Event {
	msg := fmt.Sprintf(format, a...)
	ev := createEvent(msg, LevelWarn, l.traceId)

	for key, value := range l.generalAttrs {
		ev.SetAttr(key, value)
	}

	l.appendEvent(ev)

	return ev
}

func (l *ALogger) Errorf(format string, a ...interface{}) *Event {
	msg := fmt.Sprintf(format, a...)

	ev := createEvent(msg, LevelError, l.traceId)
	l.hasError = true

	for key, value := range l.generalAttrs {
		ev.SetAttr(key, value)
	}

	l.appendEvent(ev)

	return ev
}

func (l *ALogger) Flush() error {
	// Пользовательский вывод логов
	if flushDriver != nil {
		return flushDriver(l)
	}

	// Режим логирования только при ошибках
	if l.cfg.OnlyError && !l.hasError {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Логирование событий
	for _, ev := range l.events {
		// Пропуск событий с уровнем ниже установленного
		if ev.Level < l.cfg.Level {
			continue
		}

		if err := ev.Flush(l.cfg.Output, l.cfg.TextFormat); err != nil {
			return fmt.Errorf("[ALogger] Ошибка логирования события: %w", err)
		}
	}

	return nil
}

func (l *ALogger) SetAttr(key string, value interface{}) *ALogger {
	if l.generalAttrs == nil {
		l.generalAttrs = make(map[string]interface{})
	}

	l.generalAttrs[key] = value

	return l
}

func (l *ALogger) SetAttrs(attrs map[string]interface{}) *ALogger {
	l.generalAttrs = attrs

	return l
}

// DebugFromCtx Логирование события уровня Debug
func DebugFromCtx(ctx context.Context, msg string, err error, attrs map[string]interface{}, stack bool) {
	ev := createEvent(msg, LevelDebug, getTraceId(ctx))
	defer ev.Flush(defaultConfig.Output, defaultConfig.TextFormat)

	if err != nil {
		ev.Wrap(err)
	}

	for key, value := range attrs {
		ev.SetAttr(key, value)
	}

	if stack {
		ev.Stack()
	}
}

// InfoFromCtx Логирование события уровня Info
func InfoFromCtx(ctx context.Context, msg string, err error, attrs map[string]interface{}, stack bool) {
	ev := createEvent(msg, LevelInfo, getTraceId(ctx))
	defer ev.Flush(defaultConfig.Output, defaultConfig.TextFormat)

	if err != nil {
		ev.Wrap(err)
	}

	for key, value := range attrs {
		ev.SetAttr(key, value)
	}

	if stack {
		ev.Stack()
	}
}

// WarnFromCtx Логирование события уровня Warn
func WarnFromCtx(ctx context.Context, msg string, err error, attrs map[string]interface{}, stack bool) {
	ev := createEvent(msg, LevelWarn, getTraceId(ctx))
	defer ev.Flush(defaultConfig.Output, defaultConfig.TextFormat)

	if err != nil {
		ev.Wrap(err)
	}

	for key, value := range attrs {
		ev.SetAttr(key, value)
	}

	if stack {
		ev.Stack()
	}
}

// ErrorFromCtx Логирование события уровня Error
func ErrorFromCtx(ctx context.Context, msg string, err error, attrs map[string]interface{}, stack bool) {
	ev := createEvent(msg, LevelError, getTraceId(ctx))
	defer ev.Flush(defaultConfig.Output, defaultConfig.TextFormat)

	if err != nil {
		ev.Wrap(err)
	}

	for key, value := range attrs {
		ev.SetAttr(key, value)
	}

	if stack {
		ev.Stack()
	}
}
