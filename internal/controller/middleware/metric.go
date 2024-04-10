package middleware

import (
	"context"

	"google.golang.org/grpc"

	"github.com/efremovich/data-receiver/pkg/metrics"
)

var ignoredMethodsForMetric = map[string]bool{
	"/marking.ServiceMarking/CheckHealth": true,
}

func MetricInterceptor(metricsCollector metrics.Collector) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !ignoredMethodsForMetric[info.FullMethod] {
			metricsCollector.AddAPIMethodUse(info.FullMethod)
		}

		return handler(ctx, req)
	}
}
