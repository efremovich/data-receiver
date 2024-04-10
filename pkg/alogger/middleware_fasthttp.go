package alogger

import "github.com/valyala/fasthttp"

// TraceIdMiddlewareFastHTTP - middleware для добавления trace_id в контекст
func TraceIdMiddlewareFastHTTP(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		traceId := getTraceIdFromHTTPRequestHeaders(&ctx.Request)

		if traceId != unknownData {
			ctx.SetUserValue(TraceIdKey, traceId)
			ctx.Response.Header.Add(TraceIdKey, traceId)
		}

		next(ctx)
	}
}
