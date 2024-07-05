# 使用方法

```go
// 我习惯用的
cfg := log.Config{
    FilePath: "/var/log",
    MaxDays:  7,
    LogLevel: "info",
}
log.Init(cfg)
log.Info("hello world!")


// 使用Std
logger := log.NewStdLogger(os.Stdout)
// fields & valuer
logger = log.With(logger,
    "service.name", "hellworld",
    "service.version", "v1.0.0",
    "ts", log.DefaultTimestamp,
    "caller", log.DefaultCaller,
)
logger.Log(log.LevelInfo, "key", "value")

// helper
helper := log.NewHelper(logger)
helper.Log(log.LevelInfo, "key", "value")
helper.Info("info message")
helper.Infof("info %s", "message")
helper.Infow("key", "value")
```