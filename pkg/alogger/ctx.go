package alogger

import (
	"context"
)

const key CtxKey = "logger"

func GetLoggerFromContext(ctx context.Context) *ALogger {
	logger, ok := ctx.Value(key).(*ALogger)
	if ok {
		return logger
	}

	return NewALogger(ctx, nil)
}

func AddLoggerInContext(ctx context.Context, l *ALogger) context.Context {
	return context.WithValue(ctx, key, l)
}
