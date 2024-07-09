package md

import (
	"net/http"
	"time"

	"github.com/ChangSZ/golib/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Rate 限流
func Rate(maxRequestsPerSecond int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limiter := rate.NewLimiter(rate.Every(time.Second*1), maxRequestsPerSecond)
		if !limiter.Allow() {
			operation := ctx.Request.URL.Path
			raw := ctx.Request.URL.RawQuery
			if raw != "" {
				operation = operation + "?" + raw
			}
			log.Log(log.LevelWarn,
				"Operation", operation,
				"Method", ctx.Request.Method,
				"StatusCode", http.StatusTooManyRequests,
				"Message", "该请求已被限流",
			)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
