package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// TraceID returns a traceid valuer.
func TraceID() Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			return span.TraceID().String()
		}
		return ""
	}
}
