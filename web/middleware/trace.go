package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
	"go.opentelemetry.io/otel/trace"
	"web/config"
)

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//心跳检测不记录
		if ctx.Request.URL.Path == config.APICheckPath {
			ctx.Next()
			return
		}
		var parentCtx context.Context
		parentCtx = ctx
		if ctx.GetHeader("trace-id") != "" {
			traceId, _ := trace.TraceIDFromHex(ctx.GetHeader("trace-id"))
			spanId, _ := trace.SpanIDFromHex(ctx.GetHeader("span-id"))
			spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceId,
				SpanID:     spanId,
				TraceFlags: trace.FlagsSampled,
				TraceState: trace.TraceState{},
				Remote:     true,
			})
			parentCtx = trace.ContextWithSpanContext(ctx, spanCtx)
		}

		spanCtx, span := opentelemetry.StartSpanFromContext(parentCtx, nil, ctx.Request.URL.Path)
		ctx.Set(config.ContextFieldTraceId, span.SpanContext().TraceID().String())
		ctx.Set(config.ContextFieldParentContext, spanCtx)
		ctx.Next()

		//span.SetAttributes(attribute.String("request", ctx.Request.PostForm.Encode()))
		//span.SetAttributes(attribute.String("response", responseBody))
		span.End()
	}
}
