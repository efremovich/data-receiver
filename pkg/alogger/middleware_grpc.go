package alogger

import (
	"context"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryTraceIdInterceptor - middleware для добавления trace_id в контекст унарных gRPC запросов
func UnaryTraceIdInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	traceId, md := getTraceIdWithMetadataFromGRPCContext(ctx)
	ctx = context.WithValue(ctx, TraceIdKey, traceId)
	ctx = metadata.NewIncomingContext(ctx, md)

	return handler(ctx, req)
}

// StreamTraceIdInterceptor - middleware для добавления trace_id в контекст потоковых gRPC запросов
func StreamTraceIdInterceptor(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	traceId, md := getTraceIdWithMetadataFromGRPCContext(stream.Context())
	ctx := context.WithValue(stream.Context(), TraceIdKey, traceId)
	ctx = metadata.NewIncomingContext(ctx, md)

	wrappedStream := &grpcmiddleware.WrappedServerStream{
		ServerStream:   stream,
		WrappedContext: metadata.NewIncomingContext(ctx, md),
	}

	return handler(srv, wrappedStream)
}
