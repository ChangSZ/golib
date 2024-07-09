# 推荐使用方式
```golang
package main

import (
	"github.com/ChangSZ/golib/gin/md"
	"github.com/ChangSZ/golib/log"
	"github.com/ChangSZ/golib/shutdown"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	eng := gin.Default()
	// 配置 CORS 中间件
	config := cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "UPDATE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding",
			"X-CSRF-Token", "Authorization", "X-Auth-Token", "X-Auth-UUID", "X-Auth-Openid",
			"referrer", "Authorization", "x-client-id", "x-client-version", "x-client-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	eng.Use(
		cors.New(config),
		md.Rate(100),
		md.Tracing("golib"),
		md.AccessLog(log.GetLoggerWithTrace()),
	)

	// 设置可信代理
	err := eng.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic(err)
	}

	log.Info("app Run...")
	if err := eng.Run(":8080"); err != nil {
		panic(err)
	}

	// 优雅关闭
	shutdown.NewHook().Close(
		func() {
			log.Info("shutdown...")
		},
	)
}

```

# 有没有感觉非常优雅！