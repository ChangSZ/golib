package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetLogger returns global logger appliance as logger in current process.
func GetLoggerWithTrace() Logger {
	return With(GetLogger(), "trace.id", TraceID())
}

// Trace相关的函数, 需要ctx中存在opentelemetry的标准trace信息
func WithTrace(ctx context.Context) *Helper {
	return NewHelper(
		WithContext(ctx,
			With(GetLogger(),
				"caller", DefaultCaller,
				"trace.id", TraceID(),
			),
		),
	)
}

func SQLWithTrace(ctx context.Context) *Helper {
	return NewHelper(
		WithContext(ctx,
			With(GetLogger(),
				"trace.id", TraceID(),
			),
		),
	)
}

// timeEncoder 时间格式化函数
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05-07:00"))
}

type Config struct {
	FilePath string `toml:"filePath"`
	MaxDays  int    `toml:"maxDays"`
	LogLevel string `toml:"logLevel"`
	Std      bool   `toml:"Std"`
}

func Init(cfg Config) {
	writer, err := RotateDailyLog(cfg.FilePath, cfg.MaxDays)
	if err != nil {
		panic("创建日志文件失败")
	}

	var writeSyncer zapcore.WriteSyncer
	if cfg.Std {
		multiWriter := io.MultiWriter(os.Stdout, writer)
		writeSyncer = zapcore.AddSync(multiWriter)
	} else {
		writeSyncer = zapcore.AddSync(writer)
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zapcore.InfoLevel
		fmt.Printf("日志level(%v)设置不正确: %v, 已自动设置为: %v", cfg.LogLevel, err, zapcore.InfoLevel)
	}
	core := zapcore.NewCore(encoder, writeSyncer, level)
	SetLogger(NewZapLogger(zap.New(core)))
}
