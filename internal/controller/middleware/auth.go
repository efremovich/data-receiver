package middleware

import (
	"context"
	"strings"

	// data_receiver "github.com/efremovich/data-receiver/pkg/data-receiver-service"

	"github.com/gogo/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var ignoredMethodsForAuth = map[string]bool{
	// data_receiver.CardReceiver_CheckHealth_FullMethodName: true,
}

func AuthInterceptor(token string) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if ok := ignoredMethodsForAuth[info.FullMethod]; !ok {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "metadata not provided")
			}

			requestToken := md.Get("authorization")
			if len(requestToken) == 0 {
				return nil, status.Error(codes.Unauthenticated, "authorization token not provided")
			}

			if strings.TrimPrefix(requestToken[0], "Bearer ") != token {
				return nil, status.Error(codes.PermissionDenied, "invalid token")
			}
		}

		return handler(ctx, req)
	}
}
