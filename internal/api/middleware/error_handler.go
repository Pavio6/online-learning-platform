package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"online-learning-platform/internal/errors"
	"online-learning-platform/internal/logger"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 检查是否是AppError
			if appErr, ok := err.(*errors.AppError); ok {
				logger.WithError(appErr).Errorf("Request error: %s", appErr.Message)
				c.JSON(appErr.HTTPStatus(), gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return
			}

			// 其他错误
			logger.WithError(err).Error("Internal server error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    errors.ErrCodeInternal,
				"message": "内部服务器错误",
			})
		}
	}
}
