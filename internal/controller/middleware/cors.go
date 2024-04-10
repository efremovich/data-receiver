package middleware

import (
	"bytes"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CorsMiddleware(h func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response().Header.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		ctx.Response().Header.Set("Access-Control-Allow-Headers", "*")

		if bytes.Equal([]byte(ctx.Method()), []byte("OPTIONS")) {
			ctx.Response().SetStatusCode(http.StatusOK)
			return nil
		}

		return h(ctx)
	}
}
