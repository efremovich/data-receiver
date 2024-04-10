package alogger

import (
	"context"

	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/metadata"
)

func getTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value(TraceIdKey).(string)
	if !ok || traceId == "" {
		return unknownData
	}

	return traceId
}

// getTraceIdWithMetadataFromGRPCContext - Получить trace_id из контекста gRPC запроса, а так же набор метаданных,
// которые содержатся в этом запросе
func getTraceIdWithMetadataFromGRPCContext(ctx context.Context) (string, metadata.MD) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{TraceIdKey: unknownData})
		return unknownData, md
	}

	if len(md.Get(TraceIdKey)) == 0 {
		md.Set(TraceIdKey, unknownData)
	}

	return md.Get(TraceIdKey)[0], md
}

// getTraceIdFromHTTPRequestHeaders - Получить trace_id из заголовков http запроса.
// Если в http заголовках не было trace_id, то генерируется новый uuid
func getTraceIdFromHTTPRequestHeaders(req *fasthttp.Request) string {
	traceId := string(req.Header.Peek(TraceIdKey))
	if traceId == "" {
		traceId = unknownData
	}

	return traceId
}
