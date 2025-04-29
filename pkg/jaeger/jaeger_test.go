package jaeger

// TODO добавить негативные тесты

import (
	"context"
	"testing"
)

const (
	jaegerCollectorURL = "http://localhost:14268/api/traces"
)

func TestStart(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "Успешный запуск Jaeger",
			ctx:  context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jaeger := NewJaeger()
			if err := jaeger.Start("testService", jaegerCollectorURL); err != nil {
				t.Errorf("Start() error = %v", err)
			}
		})
	}
}

func TestStop(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "Успешная остановка Jaeger",
			ctx:  context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jaeger := NewJaeger()
			if err := jaeger.Start("testService", jaegerCollectorURL); err != nil {
				t.Errorf("Start() error = %v", err)
			}
			if err := jaeger.Stop(); err != nil {
				t.Errorf("Stop() error = %v", err)
			}
		})
	}
}

func TestStartSpan(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		spanName string
	}{
		{
			name:     "Создание спана",
			ctx:      context.Background(),
			spanName: "testSpan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			span, spanCtx := StartSpan(tt.ctx, tt.spanName)
			if span == nil {
				t.Errorf("StartSpan() didn't return a span")
			}

			if spanCtx == nil {
				t.Errorf("StartSpan() didn't return a context")
			}
		})
	}
}
