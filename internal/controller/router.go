package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/controller/middleware"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

func newRouter(grpcGateway *runtime.ServeMux, cfg config.Gateway, cardHandler func(*fiber.Ctx) error, metricsCollector metrics.Collector) *fiber.App {
	server := fiber.New()

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	server.All("/card/v1", middleware.CorsMiddleware(
		middleware.LoggerMiddleware(cardHandler),
	))

	server.Static("/receiver/swagger", cfg.PathToSwaggerDir)
	server.Static("/receiver/data-receiver/swagger", cfg.PathToSwaggerDir) // Swagger для локальной сборки.

	server.All("/metrics", middleware.CorsMiddleware(adaptor.HTTPHandler(metricsCollector.ServeHTTP())))

	server.All("/*", middleware.CorsMiddleware(
		middleware.LoggerMiddleware(adaptor.HTTPHandler(grpcGateway)),
	))

	return server
}
