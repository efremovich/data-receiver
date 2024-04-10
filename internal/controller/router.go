package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	
	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/controller/middleware"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

func newRouter(grpcGateway *runtime.ServeMux, cfg config.Gateway, cmsHandler func(*fiber.Ctx) error, metricsCollector metrics.Collector) *fiber.App {
	server := fiber.New()

	server.All("/receiver/cms/v1", middleware.CorsMiddleware(
		middleware.LoggerMiddleware(cmsHandler),
	))

	server.Static("/receiver/swagger", cfg.PathToSwaggerDir)
	server.All("/receiver/metrics", middleware.CorsMiddleware(adaptor.HTTPHandler(metricsCollector.ServeHTTP())))

	server.All("/receiver/*", middleware.CorsMiddleware(
		middleware.LoggerMiddleware(adaptor.HTTPHandler(grpcGateway)),
	))

	return server
}
