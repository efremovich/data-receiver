package logger

import (
	"context"

	"github.com/efremovich/data-receiver/pkg/alogger"
)

type loggerKey string

const key loggerKey = "logger"

func GetLoggerFromContext(ctx context.Context) *alogger.ALogger {
	logger, ok := ctx.Value(key).(*alogger.ALogger)
	if ok {
		return logger
	}

	return alogger.NewALogger(ctx, nil)
}

func AddLoggerInContext(ctx context.Context, l *alogger.ALogger) context.Context {
	return context.WithValue(ctx, key, l)
}
