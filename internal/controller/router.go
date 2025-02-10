package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

func newRouter(grpcGateway *runtime.ServeMux, cfg config.Gateway,
	offerHandler func(*fiber.Ctx) error,
	stockHandler func(*fiber.Ctx) error,
	metricsCollector metrics.Collector) *fiber.App {
	server := fiber.New()

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	server.All("/feed/v1/offer", offerHandler)
	server.All("/feed/v1/stock", stockHandler)

	server.Static("/swagger", cfg.PathToSwaggerDir)
	server.Static("/data-receiver/swagger", cfg.PathToSwaggerDir) // Swagger для локальной сборки.
	server.All("/metrics", adaptor.HTTPHandler(metricsCollector.ServeHTTP()))
	server.All("/*", adaptor.HTTPHandler(grpcGateway))

	return server
}
