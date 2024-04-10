package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.astralnalog.ru/utils/alogger"

	"github.com/gofiber/fiber/v2"
)

var ignoreEndpointsForLog = map[string]bool{
	"/receiver/health": true,
}

func LoggerMiddleware(next func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()
		err := next(ctx)
		endpoint := ctx.BaseURL() + ctx.Path()
		duration := time.Since(start)
		msg := fmt.Sprintf("[%s] %s | %s | %d", duration, ctx.Method(), endpoint, ctx.Response().StatusCode())

		if ctx.Response().StatusCode() != http.StatusOK {
			msg += " | body: " + string(ctx.Response().Body())
		}

		if !ignoreEndpointsForLog[ctx.Path()] {
			alogger.InfoFromCtx(context.Background(), msg, nil, nil, false)
		}

		return err
	}
}
