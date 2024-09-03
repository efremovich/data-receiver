package alogger

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
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
	createdAt    time.Time
}

func NewALogger(ctx context.Context, cfg *Config) *ALogger {
	if cfg == nil {
		cfg = defaultConfig
	}

	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	return &ALogger{
		createdAt: time.Now(),
		cfg:       cfg,
		traceId:   getTraceId(ctx),
	}
}

func (l *ALogger) appendEvent(ev *Event) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.events = append(l.events, ev)
}

func (l *ALogger) Debugf(format string, a ...interface{}) {
	l.eventLog(LevelDebug, format, a...)
}

func (l *ALogger) Infof(format string, a ...interface{}) {
	l.eventLog(LevelInfo, format, a...)
}

func (l *ALogger) Warnf(format string, a ...interface{}) {
	l.eventLog(LevelWarn, format, a...)
}

func (l *ALogger) Errorf(format string, a ...interface{}) {
	l.eventLog(LevelError, format, a...)
}

func (l *ALogger) eventLog(level Level, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	ev := createEvent(msg, level, l.traceId)
	l.hasError = true

	for key, value := range l.generalAttrs {
		ev.SetAttr(key, value)
	}

	ev.SetAttr("logger_life_time", time.Since(l.createdAt).String())

	if ev.Level < l.cfg.Level {
		return
	}

	// Если включен режим логирования только при ошибках, добавим к массиву, иначе логируем сразу.
	if l.cfg.OnlyError {
		l.appendEvent(ev)

		return
	}

	if err := ev.Flush(l.cfg.Output, l.cfg.TextFormat); err != nil {
		fmt.Printf("[ALogger] Ошибка логирования события: %s\n", err.Error())
	}
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
func DebugFromCtx(ctx context.Context, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	ev := createEvent(msg, LevelDebug, getTraceId(ctx))
	defer ev.Flush(defaultConfig.Output, defaultConfig.TextFormat)
}

// InfoFromCtx Логирование события уровня Info
func InfoFromCtx(ctx context.Context, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	ev := createEvent(msg, LevelInfo, getTraceId(ctx))
	defer ev.Flush(defaultConfig.Output, defaultConfig.TextFormat)
}

// WarnFromCtx Логирование события уровня Warn
func WarnFromCtx(ctx context.Context, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	ev := createEvent(msg, LevelWarn, getTraceId(ctx))
	defer ev.Flush(defaultConfig.Output, defaultConfig.TextFormat)
}

// ErrorFromCtx Логирование события уровня Error
func ErrorFromCtx(ctx context.Context, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	ev := createEvent(msg, LevelWarn, getTraceId(ctx))
	defer ev.Flush(defaultConfig.Output, defaultConfig.TextFormat)
}
