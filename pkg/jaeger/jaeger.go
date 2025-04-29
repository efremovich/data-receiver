// Пакет оборачивает оригинальный jaeger для упрощения использования в проекте.
// Функционал будет добавляться из оргинального Jaerger при необходимости.
package jaeger

import (
	"context"
	"fmt"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	client "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

type IJaeger interface {
	// Start Запуск Jaeger. Позволяет использовать jaeger.StartSpan()
	Start(serviceName, collectorURL string) error
	// Stop правильная остановка Jaeger
	Stop() error
}

// Инициализация IJaeger.
func NewJaeger() IJaeger {
	return new(jaeger)
}

// jaeger реализация IJaerger.
type jaeger struct {
	closer io.Closer
}

func (j *jaeger) Start(serviceName, collectorURL string) error {
	// Создание конфигурации для Jaeger
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  client.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			CollectorEndpoint: collectorURL, // "http://jaeger-collector:14268/api/traces"
			LogSpans:          true,
		},
	}

	// Создание трейсера
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return fmt.Errorf("jaeger: ошибка инициализации трейсера: %w", err)
	}

	j.closer = closer

	// Установка трейсера в обертку opentracing
	opentracing.SetGlobalTracer(tracer)

	return nil
}

func (j *jaeger) Stop() error {
	err := j.closer.Close()
	if err != nil {
		return fmt.Errorf("jaeger: ошибка остановки jaerger: %w", err)
	}

	return nil
}

type Span opentracing.Span
type SpanContext opentracing.SpanContext

// StartSpan создание спана.
// Если jaeger не был проинициализирован - ошибки не будет как и записи в jaeger UI.
func StartSpan(ctx context.Context, nameSpan string) (Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, nameSpan)
}
