package log

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// TraceID returns a traceid valuer.
func TraceID() Valuer {
	return func(ctx context.Context) interface{} {
		var nctx context.Context = ctx
		if c, ok := ctx.(*gin.Context); ok {
			nctx = c.Request.Context()
		}

		if span := trace.SpanContextFromContext(nctx); span.HasTraceID() {
			return span.TraceID().String()
		}
		return ""
	}
}
