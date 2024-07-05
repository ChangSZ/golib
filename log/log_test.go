package log

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestInfo(_ *testing.T) {
	// 结构化
	cfg := Config{
		FilePath: "./tmp/log",
		MaxDays:  7,
		LogLevel: "info",
		Std:      true,
	}
	Init(cfg)
	Info("hello world!")

	parent := context.Background()
	tp := sdkTrace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("hello")
	ctx, _ := tracer.Start(parent, "1")
	WithTrace(ctx).Info("xxxxxxxxx")

	// 非结构化
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp)
	logger = With(logger, "caller", Caller(3))
	_ = logger.Log(LevelInfo, "key1", "value1")
}
