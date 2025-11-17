package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/logger"
)

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)
		status := c.Writer.Status()

		// 记录日志
		logger.WithFields(map[string]interface{}{
			"method":  method,
			"path":    path,
			"status":  status,
			"latency": latency,
			"ip":      c.ClientIP(),
		}).Info("HTTP Request")
	}
}
